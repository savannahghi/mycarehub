package testutils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/kevinburke/twilio-go"
	"github.com/mailgun/mailgun-go"
	"github.com/savannahghi/firebasetools"
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
)

var (
	surveysBaseURL   = serverutils.MustGetEnvVar("SURVEYS_BASE_URL")
	mailGunAPIKey    = serverutils.MustGetEnvVar("MAILGUN_API_KEY")
	mailGunDomain    = serverutils.MustGetEnvVar("MAILGUN_DOMAIN")
	twilioAccountSID = serverutils.MustGetEnvVar("TWILIO_ACCOUNT_SID")
	twilioAuthToken  = serverutils.MustGetEnvVar("TWILIO_ACCOUNT_AUTH_TOKEN")
)

const (
	twilioHTTPClientTimeoutSeconds = 10
)

// InitializeTestService sets up the structure that will be used by the usecase layer for
// integration tests
func InitializeTestService(ctx context.Context) (*usecases.MyCareHub, error) {
	fc := &firebasetools.FirebaseClient{}
	_, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate test repository: %v", err)
	}

	externalExt := externalExtension.NewExternalMethodsImpl()

	fcmService := fcm.NewService()

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)

	silCommsLib, err := silcomms.NewSILCommsLib()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate silcomms lib: %v", err)
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

	getStream := streamService.NewServiceGetStream(&stream.Client{})

	pubsub, err := pubsubmessaging.NewServicePubSubMessaging(externalExt, getStream, db, fcmService)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate pubsub service: %v", err)
	}

	clinicalClient := externalExtension.NewInterServiceClient("clinical", externalExt)
	clinicalService := clinical.NewServiceClinical(clinicalClient)

	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db, pubsub, externalExt)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, db, externalExt)
	contentUseCase := content.NewUseCasesContentImplementation(db, db, externalExt)

	mailClient := mailgun.NewMailgun(mailGunDomain, mailGunAPIKey)
	mailClient.SetAPIBase(mailgun.ApiBase)
	mailService := mail.NewServiceMail(mailClient)
	feedbackUsecase := feedback.NewUsecaseFeedback(db, db, mailService)

	notificationUseCase := notification.NewNotificationUseCaseImpl(fcmService, db, db, db, externalExt)
	appointmentUsecase := appointment.NewUseCaseAppointmentsImpl(externalExt, db, db, db, pubsub, notificationUseCase)
	communityUsecase := communities.NewUseCaseCommunitiesImpl(getStream, externalExt, db, db, pubsub, notificationUseCase, db)
	authorityUseCase := authority.NewUsecaseAuthority(db, db, externalExt, notificationUseCase)
	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt, otpUseCase, authorityUseCase, getStream, pubsub, clinicalService, smsService, twilioService)
	serviceRequestUseCase := servicerequest.NewUseCaseServiceRequestImpl(db, db, db, externalExt, userUsecase, notificationUseCase, smsService)
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
	programsUsecase := programs.NewUsecasePrograms(db, db, db, externalExt)

	organisationUsecase := organisation.NewUseCaseOrganisationImpl(db, db, db, externalExt)

	i := usecases.NewMyCareHubUseCase(
		userUsecase, termsUsecase, facilityUseCase,
		securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
		serviceRequestUseCase, authorityUseCase, communityUsecase, screeningToolsUsecases,
		appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
		programsUsecase,
		organisationUsecase,
	)
	return i, nil
}

// CheckIfCurrentDBIsLocal checks whether the database used to run the test is a test/local database. If not, the setup exits
func CheckIfCurrentDBIsLocal() bool {
	isLocal, err := strconv.ParseBool(os.Getenv("IS_LOCAL_DB"))
	if err != nil {
		return false
	}

	return isLocal
}
