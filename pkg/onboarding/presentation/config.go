package presentation

import (
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"time"

	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
)

const (
	serverTimeoutSeconds = 120
)

// AllowedOrigins is list of CORS origins allowed to interact with
// this service
var AllowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
	"http://localhost:8085",
	"http://localhost:8082",
}
var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type",
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	// fc := &firebasetools.FirebaseClient{}
	// firebaseApp, err := fc.InitFirebase()
	// if err != nil {
	// 	return nil, err
	// }

	// infra, err := infrastructure.NewInfrastructureInteractor()
	// if err != nil {
	// 	return nil, err
	// }
	// onboardingUsecases := interactor.NewUseCasesInteractor(infra)
	// if onboardingUsecases == nil {
	// 	return nil, fmt.Errorf("can't instantiate onboarding usecases: %w", err)
	// }

	// i := interactor.NewUseCasesInteractor(infra)

	r := mux.NewRouter() // gorilla mux

	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(serverutils.RequestDebugMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(serverutils.HealthStatusCheck)

	// Authenticated routes
	authR := r.Path("/graphql").Subrouter()
	// authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, &interactor.Interactor{}))

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

// GQLHandler sets up a GraphQL resolver
func GQLHandler(ctx context.Context,
	services *interactor.Interactor,
) http.HandlerFunc {
	resolver, err := graph.NewResolver(ctx, services)
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
