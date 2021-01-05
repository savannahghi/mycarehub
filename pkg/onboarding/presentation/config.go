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

	i, err := interactor.NewOnboardingInteractor()
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

	// LoginByPhone routes
	r.Path("/login_by_phone").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.LoginByPhone(ctx))
	r.Path("/refresh_token").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RefreshToken(ctx))

	// PIN Routes
	r.Path("/change_pin").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.ChangePin(ctx))

	r.Path("/request_pin_reset").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.RequestPINReset(ctx))

	//OTP routes
	r.Path("/send_retry_otp").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.SendRetryOTP(ctx))

	// Interservice Authenticated routes
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(base.InterServiceAuthenticationMiddleware())
	isc.Path("/supplier").Methods(
		http.MethodPost,
		http.MethodOptions).
		HandlerFunc(h.FindSupplierByUID(ctx))

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
