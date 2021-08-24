package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/savannahghi/onboarding/pkg/onboarding/usecases/ussd"
	"gitlab.slade360emr.com/go/apiclient"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/chargemaster"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"

	loginservice "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/login_service"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	crmExt "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/crm"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/graph"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/rest"
	adminSrv "github.com/savannahghi/onboarding/pkg/onboarding/usecases/admin"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

const (
	mbBytes              = 1048576
	serverTimeoutSeconds = 120
	engagementService    = "engagement"
	ediService           = "edi"
)

// AllowedOrigins is list of CORS origins allowed to interact with
// this service
var AllowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
}
var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type",
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}
	fsc, err := firebaseApp.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to initialize Firestore: %s", err)
	}

	fbc, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}

	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}

	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	var repo repository.OnboardingRepository

	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		firestoreExtension := fb.NewFirestoreClientExtension(fsc)
		repo = fb.NewFirebaseRepository(firestoreExtension, fbc)
	}

	// Initialize base (common) extension
	baseExt := extension.NewBaseExtensionImpl(fc)

	// Initialize ISC clients
	engagementClient := utils.NewInterServiceClient(engagementService, baseExt)
	ediClient := utils.NewInterServiceClient(ediService, baseExt)

	// Initialize new instance of our infrastructure services
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	engage := engagement.NewServiceEngagementImpl(engagementClient, baseExt)
	edi := edi.NewEdiService(ediClient, repo)
	mes := messaging.NewServiceMessagingImpl(baseExt)
	pinExt := extension.NewPINExtensionImpl()

	// hubspot usecases
	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(ctx, hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr)
	crmExt := crmExt.NewCrmService(hubspotUsecases)
	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		baseExt,
		crmExt,
		edi,
		repo,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new pubsub messaging service: %w", err)
	}

	// Initialize the usecases
	profile := usecases.NewProfileUseCase(repo, baseExt, engage, pubSub, crmExt)
	supplier := usecases.NewSupplierUseCases(repo, profile, chrg, engage, mes, baseExt, pubSub)
	login := usecases.NewLoginUseCases(repo, profile, baseExt, pinExt)
	survey := usecases.NewSurveyUseCases(repo, baseExt)
	userpin := usecases.NewUserPinUseCase(repo, profile, baseExt, pinExt, engage)
	su := usecases.NewSignUpUseCases(repo, profile, userpin, supplier, baseExt, engage, pubSub, edi)
	nhif := usecases.NewNHIFUseCases(repo, profile, baseExt, engage)
	aitUssd := ussd.NewUssdUsecases(repo, baseExt, profile, userpin, su, pinExt, pubSub, crmExt)
	sms := usecases.NewSMSUsecase(repo, baseExt)
	role := usecases.NewRoleUseCases(repo, baseExt)
	adminSrv := adminSrv.NewService(baseExt)

	i, err := interactor.NewOnboardingInteractor(
		profile, su, supplier, login, survey,
		userpin, chrg, engage, mes, nhif, pubSub,
		sms, aitUssd, edi, adminSrv, crmExt,
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}

	h := rest.NewHandlersInterfaces(i)
	loginService := loginservice.NewServiceLogin(baseExt)

	r := mux.NewRouter() // gorilla mux
	r.Use(otelmux.Middleware(serverutils.MetricsCollectorService("onboarding")))
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(serverutils.RequestDebugMiddleware())

	// Add Middleware that records the metrics for HTTP routes
	r.Use(serverutils.CustomHTTPRequestMetricsMiddleware())

	//USSD routes
	r.Path("/ait_ussd").
		Methods(http.MethodPost, http.MethodOptions).
		HandlerFunc(h.IncomingUSSDHandler())

	// Unauthenticated routes
	r.Path("/optout").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(h.OptOut())
	r.Path("/switch_flagged_features").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(
		h.SwitchFlaggedFeaturesHandler(),
	)

	// login service routes
	r.Path("/login").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(
		loginService.GetLoginFunc(ctx),
	)
	r.Path("/logout").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(
		loginService.GetLogoutFunc(ctx),
	)
	r.Path("/refresh").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(
		loginService.GetRefreshFunc(),
	)
	r.Path("/verify_access_token").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(
		loginService.GetVerifyTokenFunc(ctx),
	)

	r.Path("/pubsub").Methods(
		http.MethodPost).
		HandlerFunc(pubSub.ReceivePubSubPushMessages)

	// misc routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	// Admin service polling
	r.Path("/poll_services").Methods(http.MethodGet).HandlerFunc(h.PollServices())

	// signup routes
	r.Path("/verify_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.VerifySignUpPhoneNumber())
	r.Path("/create_user_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateUserWithPhoneNumber())
	r.Path("/user_recovery_phonenumbers").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.UserRecoveryPhoneNumbers())
	r.Path("/set_primary_phonenumber").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SetPrimaryPhoneNumber())

	// LoginByPhone routes
	r.Path("/login_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginByPhone())
	r.Path("/login_anonymous").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginAnonymous())
	r.Path("/refresh_token").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RefreshToken())

	// PIN Routes
	r.Path("/reset_pin").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.ResetPin())

	r.Path("/request_pin_reset").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RequestPINReset())

	//OTP routes
	r.Path("/send_otp").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SendOTP())

	r.Path("/send_retry_otp").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SendRetryOTP())

	r.Path("/remove_user").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveUserByPhoneNumber())

	r.Path("/add_admin_permissions").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddAdminPermsToUser())

	r.Path("/remove_admin_permissions").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveAdminPermsToUser())

	r.Path("/incoming_ait_messages").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.IncomingATSMS())

	// Authenticated routes
	rs := r.PathPrefix("/roles").Subrouter()
	rs.Use(apiclient.AuthenticationMiddleware(firebaseApp))
	rs.Path("/create_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateRole())
	rs.Path("/assign_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AssignRole())
	rs.Path("/remove_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveRoleByName())

	rs.Path("/add_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddRoleToUser())

	rs.Path("/remove_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveRoleToUser())

	// Interservice Authenticated routes
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(interserviceclient.InterServiceAuthenticationMiddleware())
	isc.Path("/supplier").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.FindSupplierByUID())
	isc.Path("/user_profile").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.GetUserProfileByUID())
	isc.Path("/retrieve_user_profile").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.GetUserProfileByPhoneOrEmail())
	isc.Path("/update_covers").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.UpdateCovers())
	isc.Path("/contactdetails/{attribute}/").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.ProfileAttributes())
	isc.Path("/check_permission").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CheckHasPermission())

	// Interservice Authenticated routes
	// The reason for the below endpoints to be used for interservice communication
	// is to allow for the creation and deletion of internal `test` users that can be used
	// to run tests in other services that require an authenticated user.
	// These endpoint have been used in the `Base` lib to create and delete the test users
	iscTesting := r.PathPrefix("/testing").Subrouter()
	iscTesting.Use(interserviceclient.InterServiceAuthenticationMiddleware())
	iscTesting.Path("/verify_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.VerifySignUpPhoneNumber())
	iscTesting.Path("/create_user_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateUserWithPhoneNumber())
	iscTesting.Path("/login_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginByPhone())
	iscTesting.Path("/remove_user").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveUserByPhoneNumber())
	iscTesting.Path("/register_push_token").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RegisterPushToken())
	iscTesting.Path("/add_admin_permissions").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddAdminPermsToUser())
	iscTesting.Path("/add_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddRoleToUser())
	iscTesting.Path("/remove_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveRoleToUser())
	iscTesting.Path("/update_user_profile").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.UpdateUserProfile())
	iscTesting.Path("/create_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateRole())
	iscTesting.Path("/assign_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AssignRole())
	iscTesting.Path("/remove_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveRoleByName())

	// Authenticated routes
	authR := r.Path("/graphql").Subrouter()
	authR.Use(apiclient.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, i))

	return r, nil
}

// PrepareServer starts up a server
func PrepareServer(ctx context.Context, port int, allowedOrigins []string) *http.Server {
	// start up the router
	r, err := Router(ctx)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	// start the server
	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)
	h = handlers.CORS(
		handlers.AllowedHeaders(allowedHeaders),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)
	h = handlers.CombinedLoggingHandler(os.Stdout, h)
	h = handlers.ContentTypeHandler(h, "application/json", "application/x-www-form-urlencoded")
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: serverTimeoutSeconds * time.Second,
		ReadTimeout:  serverTimeoutSeconds * time.Second,
	}
	log.Infof("Server running at port %v", addr)
	return srv
}

//HealthStatusCheck endpoint to check if the server is working.
func HealthStatusCheck(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	service *interactor.Interactor,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, service)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	server := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	}
}
