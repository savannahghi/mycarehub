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
	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kevinburke/twilio-go"
	"github.com/mailgun/mailgun-go"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	loginservice "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/login"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/mail"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	serviceSMS "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms"
	surveyInstance "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
	serviceTwilio "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
	internalRest "github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/rest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	appointment "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/appointments"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/metrics"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/programs"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/screeningtools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/surveys"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/serverutils"
	"github.com/savannahghi/silcomms"
	log "github.com/sirupsen/logrus"
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

var (
	getStreamAPIKey    = serverutils.MustGetEnvVar("GET_STREAM_KEY")
	getStreamAPISecret = serverutils.MustGetEnvVar("GET_STREAM_SECRET")

	// surveys
	surveysBaseURL = serverutils.MustGetEnvVar("SURVEYS_BASE_URL")

	mailGunAPIKey = serverutils.MustGetEnvVar("MAILGUN_API_KEY")
	mailGunDomain = serverutils.MustGetEnvVar("MAILGUN_DOMAIN")

	twilioAccountSID = serverutils.MustGetEnvVar("TWILIO_ACCOUNT_SID")
	twilioAuthToken  = serverutils.MustGetEnvVar("TWILIO_ACCOUNT_AUTH_TOKEN")

	clinicalDepsName = "clinical"
)

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

	fcmService := fcm.NewService()
	streamClient, err := stream.NewClient(getStreamAPIKey, getStreamAPISecret)
	if err != nil {
		log.Fatalf("failed to start getstream client: %v", err)
	}

	silCommsLib, err := silcomms.NewSILCommsLib()
	if err != nil {
		log.Fatalf("failed to start silcomms client: %v", err)
	}
	smsService := serviceSMS.NewServiceSMS(silCommsLib)

	// Twilio
	httpClient := &http.Client{
		Timeout: time.Second * twilioHTTPClientTimeoutSeconds,
	}
	twilioClient := twilio.NewClient(twilioAccountSID, twilioAuthToken, httpClient)
	twilioMessageObj := twilioClient.Messages
	twilioService := serviceTwilio.NewServiceTwilio(twilioMessageObj)

	otpUseCase := otp.NewOTPUseCase(db, db, externalExt, smsService, twilioService)

	getStream := streamService.NewServiceGetStream(streamClient)

	pubSub, err := pubsubmessaging.NewServicePubSubMessaging(externalExt, getStream, db, fcmService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize pubsub messaging service: %w", err)
	}

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db, pubSub)

	// Initialize user usecase
	notificationUseCase := notification.NewNotificationUseCaseImpl(fcmService, db, db, db, externalExt)

	authorityUseCase := authority.NewUsecaseAuthority(db, db, externalExt, notificationUseCase)

	clinicalClient := externalExtension.NewInterServiceClient(clinicalDepsName, externalExt)
	clinicalService := clinical.NewServiceClinical(clinicalClient)

	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt, otpUseCase, authorityUseCase, getStream, pubSub, clinicalService, smsService, twilioService)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, db, externalExt)

	contentUseCase := content.NewUseCasesContentImplementation(db, db, externalExt)

	mailClient := mailgun.NewMailgun(mailGunDomain, mailGunAPIKey)
	mailClient.SetAPIBase(mailgun.ApiBase)
	mailService := mail.NewServiceMail(mailClient)

	feedbackUsecase := feedback.NewUsecaseFeedback(db, db, mailService)

	serviceRequestUseCase := servicerequest.NewUseCaseServiceRequestImpl(db, db, db, externalExt, userUsecase, notificationUseCase, smsService)

	communitiesUseCase := communities.NewUseCaseCommunitiesImpl(getStream, externalExt, db, db, pubSub, notificationUseCase, db)

	appointmentUsecase := appointment.NewUseCaseAppointmentsImpl(externalExt, db, db, db, pubSub, notificationUseCase)

	healthDiaryUseCase := healthdiary.NewUseCaseHealthDiaryImpl(db, db, db, serviceRequestUseCase)

	screeningToolsUsecases := screeningtools.NewUseCasesScreeningTools(db, db, db, externalExt)

	surveysClient := surveyInstance.ODKClient{
		BaseURL:    surveysBaseURL,
		HTTPClient: &http.Client{},
	}
	survey := surveyInstance.NewSurveysImpl(surveysClient)
	surveysUsecase := surveys.NewUsecaseSurveys(survey, db, db, db, notificationUseCase, serviceRequestUseCase)

	metricsUsecase := metrics.NewUsecaseMetricsImpl(db)
	questionnaireUsecase := questionnaires.NewUseCaseQuestionnaire(db, db, db, db)
	programsUsecase := programs.NewUsecasePrograms(db, db)

	organisationUsecase := organisation.NewUseCaseOrganisationImpl(db)

	useCase := usecases.NewMyCareHubUseCase(
		userUsecase, termsUsecase, facilityUseCase,
		securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
		serviceRequestUseCase, authorityUseCase, communitiesUseCase, screeningToolsUsecases,
		appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
		programsUsecase,
		organisationUsecase,
	)

	internalHandlers := internalRest.NewMyCareHubHandlersInterfaces(*useCase)

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

	r.Use(DefaultOrganisationMiddleware())

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

	r.Path("/pubsub").Methods(http.MethodPost).HandlerFunc(pubSub.ReceivePubSubPushMessages)

	// This endpoint will be used by external services to get a token that will be used to
	// authenticate against our APIs
	r.Path("/login").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(loginsvc.Login(ctx))

	// This endpoint is a webhook listener. Getstream events --> `Push events` will be published
	// to this endpoint. It is mainly used for implementing a notification system for myCareHub professional app
	r.Path("/getstream_webhook").Methods(
		http.MethodPost,
	).HandlerFunc(internalHandlers.ReceiveGetstreamEvents())

	// KenyaEMR routes. These endpoints are authenticated and are used for integrations
	// between myCareHub and KenyaEMR
	kenyaEMR := r.PathPrefix("/kenya-emr").Subrouter()
	kenyaEMR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))

	kenyaEMR.Path("/register_patient").Methods(
		http.MethodOptions,
		http.MethodPost,
	).HandlerFunc(internalHandlers.RegisterKenyaEMRPatients())

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

	// Graphql route
	authR := r.Path("/graphql").Subrouter()
	authR.Use(firebasetools.AuthenticationMiddleware(firebaseApp))
	authR.Use(OrganisationMiddleware())
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
