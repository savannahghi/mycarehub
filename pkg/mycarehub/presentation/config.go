package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	loginservice "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/login"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
	internalRest "github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/rest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"

	injector "github.com/savannahghi/mycarehub/wire"
)

const (
	serverTimeoutSeconds           = 120
	twilioHTTPClientTimeoutSeconds = 10
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

	externalExt := externalExtension.NewExternalMethodsImpl()
	loginsvc := loginservice.NewServiceLoginImpl(externalExt)

	useCases, err := injector.InitializeUseCases(ctx)
	if err != nil {
		return nil, err
	}

	internalHandlers := internalRest.NewMyCareHubHandlersInterfaces(*useCases)

	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(serverutils.RequestDebugMiddleware())

	// Add Middleware that records the metrics for HTTP routes
	r.Use(serverutils.CustomHTTPRequestMetricsMiddleware())

	oauth2Routes := r.PathPrefix("/oauth").Subrouter()

	oauth2Routes.Path("/authorize").Methods(
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	).HandlerFunc(internalHandlers.AuthorizeHandler())

	oauth2Routes.Path("/token").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.TokenHandler())

	oauth2Routes.Path("/revoke").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.RevokeHandler())

	oauth2Routes.Path("/introspect").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.IntrospectionHandler())

	// Shared unauthenticated routes
	// openSourcePresentation.SharedUnauthenticatedRoutes(h, r)
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	r.Path("/login_by_phone").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.LoginByPhone())

	r.Path("/contact_organisations").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.FetchContactOrganisations())

	r.Path("/organisations").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.Organisations())

	r.Path("/refresh_token").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.RefreshToken())

	r.Path("/verify_security_questions").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.VerifySecurityQuestions())

	r.Path("/verify_phone").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.VerifyPhone())

	r.Path("/verify_otp").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.VerifyOTP())

	r.Path("/send_otp").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.SendOTP())

	// PIN routes
	r.Path("/request_pin_reset").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.RequestPINReset())

	r.Path("/reset_pin").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.ResetPIN())

	r.Path("/send_retry_otp").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.SendRetryOTP())

	r.Path("/get_user_responded_security_questions").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.GetUserRespondedSecurityQuestions())

	r.Path("/service-requests").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.CreatePinResetServiceRequest())

	r.Path("/facilities").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.SyncFacilities())

	r.Path("/delete-user").Methods(
		http.MethodOptions,
		http.MethodDelete,
	).HandlerFunc(internalHandlers.DeleteUser())

	r.Path("/pubsub").Methods(http.MethodPost).HandlerFunc(useCases.Pubsub.ReceivePubSubPushMessages)

	// This endpoint will be used by external services to get a token that will be used to
	// authenticate against our APIs
	r.Path("/login").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(loginsvc.Login(ctx))

	// KenyaEMR routes. These endpoints are authenticated and are used for integrations
	// between myCareHub and KenyaEMR
	kenyaEMR := r.PathPrefix("/kenya-emr").Subrouter()
	kenyaEMR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))

	kenyaEMR.Path("/health_diary").Methods(
		http.MethodGet,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.GetClientHealthDiaryEntries())

	kenyaEMR.Path("/service_request").Methods(
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	).HandlerFunc(internalHandlers.ServiceRequests())

	kenyaEMR.Path("/patients").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.RegisteredFacilityPatients())

	kenyaEMR.Path("/appointments").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.CreateOrUpdateKenyaEMRAppointments())

	kenyaEMR.Path("/observations").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.AddPatientsRecords())
	kenyaEMR.Path("/appointment-service-request").Methods(
		http.MethodOptions,
		http.MethodGet,
		http.MethodPost,
	).HandlerFunc(internalHandlers.AppointmentsServiceRequests())

	// ISC routes. These are inter-service route
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(interserviceclient.InterServiceAuthenticationMiddleware())
	isc.Path("/user-profile/{id}").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.GetUserProfile())
	isc.Path("/add-fhir-id").Methods(
		http.MethodOptions,
		http.MethodPatch,
	).HandlerFunc(internalHandlers.AddClientFHIRID())

	isc.Path("/facilities").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.AddFacilityFHIRID())

	isc.Path("/program").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.UpdateProgramTenantID())

	// Graphql route
	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, *useCases))

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

// HealthStatusCheck endpoint to check if the server is working.
func HealthStatusCheck(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}
}

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	usecase usecases.MyCareHub,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, usecase)
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
