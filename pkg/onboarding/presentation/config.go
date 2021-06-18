package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	loginservice "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/login_service"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/rest"
)

const (
	mbBytes              = 1048576
	serverTimeoutSeconds = 120
	engagementService    = "engagement"
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
	fc := &base.FirebaseClient{}
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

	projectID, err := base.GetEnvVar(base.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			base.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}

	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	var repo repository.OnboardingRepository

	if base.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		firestoreExtension := fb.NewFirestoreClientExtension(fsc)
		repo = fb.NewFirebaseRepository(firestoreExtension, fbc)
	}

	// Initialize base (common) extension
	baseExt := extension.NewBaseExtensionImpl(fc)

	// Initialize ISC clients
	engagementClient := utils.NewInterServiceClient(engagementService, baseExt)

	// Initialize new instance of our infrastructure services
	erp := erp.NewERPService(repo)
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	crm := hubspot.NewHubSpotService()
	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		baseExt,
		erp,
		crm,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new pubsub messaging service: %w", err)
	}
	engage := engagement.NewServiceEngagementImpl(engagementClient, baseExt, pubSub)
	mes := messaging.NewServiceMessagingImpl(baseExt)
	pinExt := extension.NewPINExtensionImpl()
	aitUssd := usecases.NewUssdUsecases(repo, baseExt)

	// Initialize the usecases
	profile := usecases.NewProfileUseCase(repo, baseExt, engage, pubSub)
	supplier := usecases.NewSupplierUseCases(repo, profile, erp, chrg, engage, mes, baseExt, pubSub)
	login := usecases.NewLoginUseCases(repo, profile, baseExt, pinExt)
	survey := usecases.NewSurveyUseCases(repo, baseExt)
	userpin := usecases.NewUserPinUseCase(repo, profile, baseExt, pinExt, engage)
	su := usecases.NewSignUpUseCases(repo, profile, userpin, supplier, baseExt, engage, pubSub)
	nhif := usecases.NewNHIFUseCases(repo, profile, baseExt, engage)
	sms := usecases.NewSMSUsecase(repo, baseExt)
	agent := usecases.NewAgentUseCases(repo, profile, engage, mes, baseExt)

	i, err := interactor.NewOnboardingInteractor(
		repo, profile, su, supplier, login, survey,
		userpin, erp, chrg, engage, mes, nhif, pubSub, sms, aitUssd, crm, agent,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}

	h := rest.NewHandlersInterfaces(i)
	loginService := loginservice.NewServiceLogin(baseExt)

	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(base.RequestDebugMiddleware())

	// Add Middleware that records the metrics for HTTP routes
	r.Use(base.CustomHTTPRequestMetricsMiddleware())

	//USSD routes
	r.Path("/ait_ussd").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(h.IncomingUSSDHandler(ctx))
	r.Path("/ait_end_note_ussd").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(h.USSDEndNotificationHandler(ctx))

	// Unauthenticated routes

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

	// signup routes
	r.Path("/verify_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.VerifySignUpPhoneNumber(ctx))
	r.Path("/create_user_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateUserWithPhoneNumber(ctx))
	r.Path("/user_recovery_phonenumbers").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.UserRecoveryPhoneNumbers(ctx))
	r.Path("/set_primary_phonenumber").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SetPrimaryPhoneNumber(ctx))

	// LoginByPhone routes
	r.Path("/login_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginByPhone(ctx))
	r.Path("/login_anonymous").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginAnonymous(ctx))
	r.Path("/refresh_token").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RefreshToken(ctx))

	// PIN Routes
	r.Path("/reset_pin").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.ResetPin(ctx))

	r.Path("/request_pin_reset").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RequestPINReset(ctx))

	//OTP routes
	r.Path("/send_otp").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SendOTP(ctx))

	r.Path("/send_retry_otp").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SendRetryOTP(ctx))

	r.Path("/remove_user").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveUserByPhoneNumber(ctx))

	r.Path("/add_admin_permissions").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddAdminPermsToUser(ctx))

	r.Path("/remove_admin_permissions").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveAdminPermsToUser(ctx))

	r.Path("/incoming_ait_messages").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.IncomingATSMS(ctx))

	// Authenticated routes
	rs := r.PathPrefix("/roles").Subrouter()
	rs.Use(base.AuthenticationMiddleware(firebaseApp))
	rs.Path("/add_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddRoleToUser(ctx))

	rs.Path("/remove_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveRoleToUser(ctx))

	// Interservice Authenticated routes
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(base.InterServiceAuthenticationMiddleware())
	isc.Path("/supplier").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.FindSupplierByUID(ctx))
	isc.Path("/user_profile").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.GetUserProfileByUID(ctx))
	isc.Path("/update_covers").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.UpdateCovers(ctx))
	isc.Path("/contactdetails/{attribute}/").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.ProfileAttributes(ctx))

	// Interservice Authenticated routes
	// The reason for the below endpoints to be used for interservice communication
	// is to allow for the creation and deletion of internal `test` users that can be used
	// to run tests in other services that require an authenticated user.
	// These endpoint have been used in the `Base` lib to create and delete the test users
	iscTesting := r.PathPrefix("/testing").Subrouter()
	iscTesting.Use(base.InterServiceAuthenticationMiddleware())
	iscTesting.Path("/verify_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.VerifySignUpPhoneNumber(ctx))
	iscTesting.Path("/create_user_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateUserWithPhoneNumber(ctx))
	iscTesting.Path("/login_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginByPhone(ctx))
	iscTesting.Path("/remove_user").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveUserByPhoneNumber(ctx))
	iscTesting.Path("/register_push_token").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RegisterPushToken(ctx))
	iscTesting.Path("/add_admin_permissions").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddAdminPermsToUser(ctx))
	iscTesting.Path("/add_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.AddRoleToUser(ctx))
	iscTesting.Path("/remove_user_role").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveRoleToUser(ctx))
	iscTesting.Path("/update_user_profile").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.UpdateUserProfile(ctx))

	// Authenticated routes
	authR := r.Path("/graphql").Subrouter()
	authR.Use(base.AuthenticationMiddleware(firebaseApp))
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
		base.LogStartupError(ctx, err)
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
		base.LogStartupError(ctx, err)
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
