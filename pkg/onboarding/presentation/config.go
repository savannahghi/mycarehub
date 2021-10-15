package presentation

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/savannahghi/firebasetools"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	infra "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	postgres "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
	internalRest "github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/rest"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	metrics "github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/metric"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/staff"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	osinfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	openSourcePresentation "github.com/savannahghi/onboarding/pkg/onboarding/presentation"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/rest"
	osusecases "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
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
	fc := &firebasetools.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	// Initialize base (common) extension
	baseExt := extension.NewBaseExtensionImpl(fc)

	// Initialize ISC clients

	pinExt := extension.NewPINExtensionImpl()

	// Initialize new instances of the infrastructure services
	// Initialize new open source interactors
	infrastructure := osinfra.NewInfrastructureInteractor()

	openSourceUsecases := osusecases.NewUsecasesInteractor(infrastructure, baseExt, pinExt)

	// initialize internal infrastructure
	infra := infra.NewInteractor()

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(infra)

	//Initialize metric usecases
	metricsUsecase := metrics.NewMetricUsecase(infra)

	//Initialize staff usecases
	staffUsecase := staff.NewUsecasesStaffProfileImpl(infra)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate repository in resolver: %v", err)
	}

	db := postgres.NewOnboardingDb(pg, pg, pg)

	// Initialize the interactor
	i, err := interactor.NewOnboardingInteractor(
		infrastructure,
		*db,
		openSourceUsecases,
		facilityUseCase,
		metricsUsecase,
		staffUsecase,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}

	h := rest.NewHandlersInterfaces(infrastructure, openSourceUsecases)
	internalHandlers := internalRest.NewOnboardingHandlersInterfaces(infra, *i)

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

	// Shared unauthenticated routes
	openSourcePresentation.SharedUnauthenticatedRoutes(h, r)
	// Shared authenticated ISC routes
	openSourcePresentation.SharedAuthenticatedISCRoutes(h, r)
	// Shared authenticated routes
	openSourcePresentation.SharedAuthenticatedRoutes(h, r)

	// Onboarding service rest routes
	r.Path("/collect_metrics").
		Methods(http.MethodPost, http.MethodOptions).
		HandlerFunc(internalHandlers.CollectMetricsHandler())

	// Graphql route
	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
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
