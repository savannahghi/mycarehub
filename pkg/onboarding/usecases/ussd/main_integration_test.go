package ussd_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"

	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases/ussd"

	mockCrm "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/crm/mock"
	mockRepo "github.com/savannahghi/onboarding/pkg/onboarding/repository/mock"

	extMock "github.com/savannahghi/onboarding/pkg/onboarding/application/extension/mock"
	ediMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi/mock"
	engagementMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement/mock"

	crmExt "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/crm"
	messagingMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging/mock"
	pubsubmessagingMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub/mock"
	adminSrv "github.com/savannahghi/onboarding/pkg/onboarding/usecases/admin"
	hubspotRepo "gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/database/fs"
	hubspotUsecases "gitlab.slade360emr.com/go/commontools/crm/pkg/usecases"
)

const (
	otpService        = "otp"
	engagementService = "engagement"
	ediService        = "edi"
)

func InitializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	fc := firebasetools.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Fatalf("unable to initialize Firestore for the Feed: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}

	var repo repository.OnboardingRepository

	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		firestoreExtension := fb.NewFirestoreClientExtension(fsc)
		repo = fb.NewFirebaseRepository(firestoreExtension, fbc)
	}

	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		return nil, fmt.Errorf(
			"can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err,
		)
	}
	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize pubsub client: %w", err)
	}

	ext := extension.NewBaseExtensionImpl(&firebasetools.FirebaseClient{})

	// Initialize ISC clients
	engagementClient := utils.NewInterServiceClient(engagementService, ext)
	ediClient := utils.NewInterServiceClient(ediService, ext)

	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(context.Background(), hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr)
	crmExt := crmExt.NewCrmService(hubspotUsecases)
	engage := engagement.NewServiceEngagementImpl(engagementClient, ext)
	edi := edi.NewEdiService(ediClient, repo)
	ps, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		ext,
		crmExt,
		edi,
		repo,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new pubsub messaging service: %w", err)
	}
	mes := messaging.NewServiceMessagingImpl(ext)
	pinExt := extension.NewPINExtensionImpl()
	profile := usecases.NewProfileUseCase(repo, ext, engage, ps, crmExt)

	supplier := usecases.NewSupplierUseCases(repo, profile, engage, mes, ext, ps)
	login := usecases.NewLoginUseCases(repo, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(repo, ext)
	userpin := usecases.NewUserPinUseCase(repo, profile, ext, pinExt, engage)
	su := usecases.NewSignUpUseCases(repo, profile, userpin, supplier, ext, engage, ps, edi)
	nhif := usecases.NewNHIFUseCases(repo, profile, ext, engage)
	sms := usecases.NewSMSUsecase(repo, ext)

	aitUssd := ussd.NewUssdUsecases(repo, ext, profile, userpin, su, pinExt, ps, crmExt)

	return &interactor.Interactor{
		Onboarding: profile,
		Signup:     su,
		Supplier:   supplier,
		Login:      login,
		Survey:     survey,
		UserPIN:    userpin,
		Engagement: engage,
		NHIF:       nhif,
		PubSub:     ps,
		SMS:        sms,
		AITUSSD:    aitUssd,
		CrmExt:     crmExt,
	}, nil
}

// InitializeFakeOnboaridingInteractor represents a fakeonboarding interactor
func InitializeFakeOnboardingInteractor() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo

	var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs
	var messagingSvc messaging.ServiceMessaging = &fakeMessagingSvc
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt
	var ps pubsubmessaging.ServicePubSub = &fakePubSub
	var ediSvc edi.ServiceEdi = &fakeEDISvc

	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(context.Background(), hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr)
	crmExt := crmExt.NewCrmService(hubspotUsecases)
	profile := usecases.NewProfileUseCase(r, ext, engagementSvc, ps, crmExt)
	login := usecases.NewLoginUseCases(r, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(r, ext)
	supplier := usecases.NewSupplierUseCases(
		r, profile, engagementSvc, messagingSvc, ext, ps,
	)
	userpin := usecases.NewUserPinUseCase(r, profile, ext, pinExt, engagementSvc)
	su := usecases.NewSignUpUseCases(r, profile, userpin, supplier, ext, engagementSvc, ps, ediSvc)
	nhif := usecases.NewNHIFUseCases(r, profile, ext, engagementSvc)
	sms := usecases.NewSMSUsecase(r, ext)
	role := usecases.NewRoleUseCases(r, ext)

	aitUssd := ussd.NewUssdUsecases(r, ext, profile, userpin, su, pinExt, ps, crmExt)
	adminSrv := adminSrv.NewService(ext)

	i, err := interactor.NewOnboardingInteractor(
		profile, su, supplier, login,
		survey, userpin,
		engagementSvc, messagingSvc, nhif, ps, sms, aitUssd, ediSvc, adminSrv, crmExt,
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}
	return i, nil

}

func TestImpl_HandleResponseFromUSSDGateway(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	unregisteredPhoneNumber := "0723456756"
	registeredPhoneNumber := interserviceclient.TestUserPhoneNumber

	unregisteredValidPayload := &dto.SessionDetails{
		SessionID:   uuid.New().String(),
		PhoneNumber: &unregisteredPhoneNumber,
	}

	// registeredValidPayload := &dto.SessionDetails{
	// 	SessionID:   uuid.New().String(),
	// 	PhoneNumber: &registeredPhoneNumber,
	// }

	invalidPayload := &dto.SessionDetails{
		SessionID:   "",
		PhoneNumber: &registeredPhoneNumber,
	}

	type args struct {
		ctx     context.Context
		payload *dto.SessionDetails
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy case ):_Success case_Unregistered_user",
			args: args{
				ctx:     ctx,
				payload: unregisteredValidPayload,
			},
			want: "CON Welcome to Be.Well\r\n" +
				"1. Register\r\n" +
				"2. Opt Out\r\n",
		},
		// {
		// 	name: "Happy case ):_Success case_Registered_user",
		// 	args: args{
		// 		ctx:     ctx,
		// 		payload: registeredValidPayload,
		// 	},
		// 	want: "CON Welcome to Be.Well.Please enter\r\n" +
		// 		"your PIN to continue(enter 00 if\r\n" +
		// 		"you forgot your PIN)\r\n",
		// },
		{
			name: "SAD case ):Fail case_invalid_sessionID",
			args: args{
				ctx:     ctx,
				payload: invalidPayload,
			},
			want: "END Something went wrong. Please try again.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ussd := u
			got := ussd.AITUSSD.HandleResponseFromUSSDGateway(tt.args.ctx, tt.args.payload)
			if got != tt.want {
				t.Errorf("Impl.HandleResponseFromUSSDGateway() = %v, want %v", got, tt.want)
			}
		})
	}
}

var fakeRepo mockRepo.FakeOnboardingRepository
var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakePinExt extMock.PINExtensionImpl
var fakeEngagementSvs engagementMock.FakeServiceEngagement
var fakeMessagingSvc messagingMock.FakeServiceMessaging
var fakePubSub pubsubmessagingMock.FakeServicePubSub
var fakeEDISvc ediMock.FakeServiceEDI
var fakeCrm mockCrm.FakeServiceCrm

// InitializeFakeUSSDTestService represents a fakeussd interactor
func InitializeFakeUSSDTestService() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo
	var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs
	var messagingSvc messaging.ServiceMessaging = &fakeMessagingSvc
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt
	var ps pubsubmessaging.ServicePubSub = &fakePubSub
	var ediSvc edi.ServiceEdi = &fakeEDISvc
	var crmSvc crmExt.ServiceCrm = &fakeCrm

	hubspotService := hubspot.NewHubSpotService()
	hubspotfr, err := hubspotRepo.NewHubSpotFirebaseRepository(context.Background(), hubspotService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hubspot crm repository: %w", err)
	}
	hubspotUsecases := hubspotUsecases.NewHubSpotUsecases(hubspotfr)
	crmExt := crmExt.NewCrmService(hubspotUsecases)
	profile := usecases.NewProfileUseCase(r, ext, engagementSvc, ps, crmExt)
	login := usecases.NewLoginUseCases(r, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(r, ext)
	supplier := usecases.NewSupplierUseCases(
		r, profile, engagementSvc, messagingSvc, ext, ps,
	)
	userpin := usecases.NewUserPinUseCase(r, profile, ext, pinExt, engagementSvc)
	su := usecases.NewSignUpUseCases(r, profile, userpin, supplier, ext, engagementSvc, ps, ediSvc)
	nhif := usecases.NewNHIFUseCases(r, profile, ext, engagementSvc)
	sms := usecases.NewSMSUsecase(r, ext)
	role := usecases.NewRoleUseCases(r, ext)
	aitUssd := ussd.NewUssdUsecases(r, ext, profile, userpin, su, pinExt, ps, crmSvc)
	adminSrv := adminSrv.NewService(ext)

	i, err := interactor.NewOnboardingInteractor(
		profile, su, supplier, login,
		survey, userpin,
		engagementSvc, messagingSvc, nhif, ps, sms, aitUssd, ediSvc, adminSrv, crmExt,
		role,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}
	return i, nil

}
