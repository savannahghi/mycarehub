package testutils

import (
	"context"
	"fmt"
	"net/http"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/firebasetools"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	surveyInstance "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/screeningtools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/surveys"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/serverutils"
)

var (
	surveysBaseURL = serverutils.MustGetEnvVar("SURVEYS_BASE_URL")
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

	// add organization
	// createOrganization(pg)

	externalExt := externalExtension.NewExternalMethodsImpl()

	fcmService := fcm.NewService()

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)

	otpUseCase := otp.NewOTPUseCase(db, db, externalExt)
	getStream := streamService.NewServiceGetStream(&stream.Client{})

	pubsub, err := pubsubmessaging.NewServicePubSubMessaging(externalExt, getStream, db, fcmService)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate pubsub service: %v", err)
	}

	clinicalClient := externalExtension.NewInterServiceClient("clinical", externalExt)
	clinicalService := clinical.NewServiceClinical(clinicalClient)

	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db, pubsub)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, db, externalExt)
	contentUseCase := content.NewUseCasesContentImplementation(db, db, externalExt)
	feedbackUsecase := feedback.NewUsecaseFeedback(db, db, externalExt)

	notificationUseCase := notification.NewNotificationUseCaseImpl(fcmService, db, db, db, externalExt)
	appointmentUsecase := appointment.NewUseCaseAppointmentsImpl(externalExt, db, db, db, pubsub, notificationUseCase)
	communityUsecase := communities.NewUseCaseCommunitiesImpl(getStream, externalExt, db, db, pubsub, notificationUseCase, db)
	authorityUseCase := authority.NewUsecaseAuthority(db, db, externalExt, notificationUseCase)
	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt, otpUseCase, authorityUseCase, getStream, pubsub, clinicalService)
	serviceRequestUseCase := servicerequest.NewUseCaseServiceRequestImpl(db, db, db, externalExt, userUsecase, notificationUseCase)
	healthDiaryUseCase := healthdiary.NewUseCaseHealthDiaryImpl(db, db, db, serviceRequestUseCase)
	screeningToolsUsecases := screeningtools.NewUseCasesScreeningTools(db, db, db, externalExt)

	surveysClient := surveyInstance.ODKClient{
		BaseURL:    surveysBaseURL,
		HTTPClient: &http.Client{},
	}
	survey := surveyInstance.NewSurveysImpl(surveysClient)
	surveysUsecase := surveys.NewUsecaseSurveys(survey, db, db, db, notificationUseCase)
	metricsUsecase := metrics.NewUsecaseMetricsImpl(db)

	i := usecases.NewMyCareHubUseCase(
		userUsecase, termsUsecase, facilityUseCase,
		securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
		serviceRequestUseCase, authorityUseCase, communityUsecase, screeningToolsUsecases,
		appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase,
	)
	return i, nil
}
