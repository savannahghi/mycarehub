package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
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
	otpService           = "otp"
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
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)

	// Initialize ISC clients
	otpClient := utils.NewInterServiceClient(otpService)
	engagementClient := utils.NewInterServiceClient(engagementService)

	// Initialize new instance of our infrastructure services
	erp := erp.NewERPService()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	engage := engagement.NewServiceEngagementImpl(engagementClient)
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	baseExt := extension.NewBaseExtensionImpl()
	pinExt := extension.NewPINExtensionImpl()
	otp := otp.NewOTPService(otpClient, baseExt)

	// Initialize the usecases
	profile := usecases.NewProfileUseCase(fr, otp, baseExt, engage)
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes, baseExt)
	login := usecases.NewLoginUseCases(fr, profile, baseExt, pinExt)
	survey := usecases.NewSurveyUseCases(fr, baseExt)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile, baseExt, pinExt)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier, otp, baseExt)

	i, err := interactor.NewOnboardingInteractor(
		fr, profile, su, otp, supplier, login, survey, userpin, erp, chrg, engage, mg, mes,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}

	h := rest.NewHandlersInterfaces(i)

	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(base.RequestDebugMiddleware())

	// Unauthenticated routes

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

	// todo(dexter) : restore this after shikami fixes the iOS pipeline. Not my favorite thing but it had to be done
	// env := os.Getenv(base.Environment)
	// if env == base.TestingEnv || env == base.StagingEnv {
	// 	r.Path("/remove_user").Methods(
	// 		http.MethodPost,
	// 		http.MethodOptions).
	// 		HandlerFunc(h.RemoveUserByPhoneNumber(ctx))
	// }

	// todo(dexter): remove this after shikami fixes the iOS pipeline. This is temporarily
	r.Path("/remove_user").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RemoveUserByPhoneNumber(ctx))

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

	// The reason for the below endpoints to be used for interservice communication
	// is to allow for the creation of internal `test` users that can be used
	// to run tests in other services that require an authenticated user.
	// The creation of the users occurs in `Base` which is a central util accessible by
	// all the other services
	isc.Path("/verify_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.VerifySignUpPhoneNumber(ctx))
	isc.Path("/create_user_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.CreateUserWithPhoneNumber(ctx))
	isc.Path("/login_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginByPhone(ctx))

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
	h = handlers.ContentTypeHandler(h, "application/json")
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
