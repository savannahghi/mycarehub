package wire

import (
	"fmt"
	"net/http"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/wire"
	"github.com/kevinburke/twilio-go"
	"github.com/mailgun/mailgun-go"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/mail"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	serviceSMS "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms"
	surveyInstance "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
	serviceTwilio "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio"
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
	twilioHTTPClientTimeoutSeconds = 10
)

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

// ProviderUseCases adds data that will be used to inject data to other usecases
func ProviderUseCases() (*usecases.MyCareHub, error) {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate repository in resolver: %v", err)
	}

	externalExt := externalExtension.NewExternalMethodsImpl()
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
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db, pubSub, externalExt)

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
	surveysUsecase := surveys.NewUsecaseSurveys(survey, db, db, db, notificationUseCase, serviceRequestUseCase, externalExt)

	metricsUsecase := metrics.NewUsecaseMetricsImpl(db)
	questionnaireUsecase := questionnaires.NewUseCaseQuestionnaire(db, db, db, db, externalExt)
	programsUsecase := programs.NewUsecasePrograms(db, db, db, externalExt, getStream, pubSub)

	organisationUsecase := organisation.NewUseCaseOrganisationImpl(db, db, db, externalExt, pubSub)

	useCase := usecases.NewMyCareHubUseCase(
		userUsecase, termsUsecase, facilityUseCase,
		securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
		serviceRequestUseCase, authorityUseCase, communitiesUseCase, screeningToolsUsecases,
		appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
		programsUsecase,
		organisationUsecase, pubSub,
	)

	return useCase, nil
}

// WireSet contains a set of functions that carry data to be injected
var WireSet = wire.NewSet(ProviderUseCases)
