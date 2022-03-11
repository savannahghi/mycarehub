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
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	loginservice "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/login"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
	internalRest "github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/rest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/faq"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/screeningtools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
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

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate repository in resolver: %v", err)
	}

	externalExt := externalExtension.NewExternalMethodsImpl()
	loginsvc := loginservice.NewServiceLoginImpl(externalExt)
	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db)

	otpUseCase := otp.NewOTPUseCase(db, db, externalExt)

	authorityUseCase := authority.NewUsecaseAuthority(db, db, externalExt)

	getStream := streamService.NewServiceGetStream()
	// Initialize user usecase
	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt, otpUseCase, authorityUseCase, getStream)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, db, externalExt)

	contentUseCase := content.NewUseCasesContentImplementation(db, db, externalExt)

	feedbackUsecase := feedback.NewUsecaseFeedback(db, externalExt)

	faq := faq.NewUsecaseFAQ(db)

	serviceRequestUseCase := servicerequest.NewUseCaseServiceRequestImpl(db, db, db)

	communitiesUseCase := communities.NewUseCaseCommunitiesImpl(getStream, externalExt, db, db)

	healthDiaryUseCase := healthdiary.NewUseCaseHealthDiaryImpl(db, db, serviceRequestUseCase)
	screeningToolsUsecases := screeningtools.NewUseCasesScreeningTools(db, db, db)

	useCase := usecases.NewMyCareHubUseCase(
		userUsecase, termsUsecase, facilityUseCase,
		securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
		faq, serviceRequestUseCase, authorityUseCase, communitiesUseCase, screeningToolsUsecases,
	)

	internalHandlers := internalRest.NewMyCareHubHandlersInterfaces(*useCase)

	r := mux.NewRouter() // gorilla mux
	r.Use(otelmux.Middleware(serverutils.MetricsCollectorService("mycarehub")))
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
	// openSourcePresentation.SharedUnauthenticatedRoutes(h, r)
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	r.Path("/login_by_phone").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.LoginByPhone())

	r.Path("/refresh_token").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.RefreshToken())

	r.Path("/refresh_getstream_token").Methods(
		http.MethodPost,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.RefreshGetStreamToken())

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

	// This endpoint will be used by external services to get a token that will be used to
	// authenticate against our APIs
	r.Path("/login").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(loginsvc.Login(ctx))

	// KenyaEMR routes
	r.Path("/register_patient").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.RegisterKenyaEMRPatients())

	r.Path("/health_diary").Methods(
		http.MethodGet,
		http.MethodOptions,
	).HandlerFunc(internalHandlers.GetClientHealthDiaryEntries())

	r.Path("/service_request").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.GetServiceRequests())

	r.Path("/patients").Methods(
		http.MethodOptions,
		http.MethodGet,
	).HandlerFunc(internalHandlers.RegisteredFacilityPatients())

	// Graphql route
	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Methods(
		http.MethodPost,
		http.MethodGet,
	).HandlerFunc(GQLHandler(ctx, *useCase))

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
