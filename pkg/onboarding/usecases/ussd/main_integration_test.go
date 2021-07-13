package ussd_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/edi"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"

	erp "gitlab.slade360emr.com/go/commontools/accounting/pkg/usecases"
	erpMock "gitlab.slade360emr.com/go/commontools/accounting/pkg/usecases/mock"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases/ussd"

	mockRepo "gitlab.slade360emr.com/go/profile/pkg/onboarding/repository/mock"

	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	chargemasterMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster/mock"
	ediMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/edi/mock"
	engagementMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement/mock"

	messagingMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging/mock"
	pubsubmessagingMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub/mock"
)

const (
	otpService        = "otp"
	engagementService = "engagement"
	ediService        = "edi"
)

func InitializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	fc := base.FirebaseClient{}
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

	ext := extension.NewBaseExtensionImpl(&base.FirebaseClient{})

	// Initialize ISC clients
	engagementClient := utils.NewInterServiceClient(engagementService, ext)
	ediClient := utils.NewInterServiceClient(ediService, ext)

	erp := erp.NewAccounting()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	crm := hubspot.NewHubSpotService()
	edi := edi.NewEdiService(ediClient, repo)
	ps, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		ext,
		erp,
		crm,
		edi,
		repo,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new pubsub messaging service: %w", err)
	}
	engage := engagement.NewServiceEngagementImpl(engagementClient, ext, ps)
	mes := messaging.NewServiceMessagingImpl(ext)
	pinExt := extension.NewPINExtensionImpl()
	profile := usecases.NewProfileUseCase(repo, ext, engage, ps)

	supplier := usecases.NewSupplierUseCases(repo, profile, erp, chrg, engage, mes, ext, ps)
	login := usecases.NewLoginUseCases(repo, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(repo, ext)
	userpin := usecases.NewUserPinUseCase(repo, profile, ext, pinExt, engage)
	su := usecases.NewSignUpUseCases(repo, profile, userpin, supplier, ext, engage, ps, edi)
	nhif := usecases.NewNHIFUseCases(repo, profile, ext, engage)
	sms := usecases.NewSMSUsecase(repo, ext)
	aitUssd := ussd.NewUssdUsecases(repo, ext, profile, userpin, su, pinExt, ps)

	return &interactor.Interactor{
		Onboarding:   profile,
		Signup:       su,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engagement:   engage,
		NHIF:         nhif,
		PubSub:       ps,
		SMS:          sms,
		AITUSSD:      aitUssd,
		CRM:          crm,
	}, nil
}

func TestImpl_HandleResponseFromUSSDGateway(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize service")
		return
	}

	sessionID := uuid.New().String()
	unregisteredPhoneNumber := "0723456756"
	registeredPhoneNumber := base.TestUserPhoneNumber

	unregisteredValidPayload := &dto.SessionDetails{
		SessionID:   sessionID,
		PhoneNumber: &unregisteredPhoneNumber,
	}

	registeredValidPayload := &dto.SessionDetails{
		SessionID:   sessionID,
		PhoneNumber: &registeredPhoneNumber,
	}

	invalidPayload := &dto.SessionDetails{
		SessionID:   "",
		PhoneNumber: &registeredPhoneNumber,
	}

	type args struct {
		ctx     context.Context
		payload *dto.SessionDetails
	}
	tests := []struct {
		name     string
		args     args
		response string
	}{
		{
			name: "Happy case ):_Success case_Unregistered_user",
			args: args{
				ctx:     ctx,
				payload: unregisteredValidPayload,
			},
			response: "CON Welcome to Be.Well\r\n" +
				"1. Register\r\n" +
				"2. Opt Out\r\n",
		},
		{
			name: "Happy case ):_Success case_Registered_user",
			args: args{
				ctx:     ctx,
				payload: registeredValidPayload,
			},
			response: "CON Welcome to Be.Well.Please enter\r\n" +
				"your PIN to continue(enter 00 if\r\n" +
				"you forgot your PIN)\r\n",
		},
		{
			name: "SAD case ):Fail case_invalid_sessionID",
			args: args{
				ctx:     ctx,
				payload: invalidPayload,
			},
			response: "END Something went wrong. Please try again.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotresp := u.AITUSSD.HandleResponseFromUSSDGateway(tt.args.ctx, tt.args.payload); gotresp != tt.response {
				t.Errorf("Impl.HandleResponseFromUSSDGateway() = %v, want %v", gotresp, tt.response)
			}
		})
	}
}

var fakeRepo mockRepo.FakeOnboardingRepository
var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakePinExt extMock.PINExtensionImpl
var fakeEngagementSvs engagementMock.FakeServiceEngagement
var fakeMessagingSvc messagingMock.FakeServiceMessaging
var fakeEPRSvc erpMock.FakeServiceCommonTools
var fakeChargeMasterSvc chargemasterMock.FakeServiceChargeMaster
var fakePubSub pubsubmessagingMock.FakeServicePubSub
var fakeEDISvc ediMock.FakeServiceEDI

// InitializeFakeOnboaridingInteractor represents a fakeonboarding interactor
func InitializeFakeOnboaridingInteractor() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo
	var erpSvc erp.AccountingUsecase = &fakeEPRSvc
	var chargemasterSvc chargemaster.ServiceChargeMaster = &fakeChargeMasterSvc
	var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs
	var messagingSvc messaging.ServiceMessaging = &fakeMessagingSvc
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt
	var ps pubsubmessaging.ServicePubSub = &fakePubSub
	var ediSvc edi.ServiceEdi = &fakeEDISvc

	profile := usecases.NewProfileUseCase(r, ext, engagementSvc, ps)
	login := usecases.NewLoginUseCases(r, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(r, ext)
	supplier := usecases.NewSupplierUseCases(
		r, profile, erpSvc, chargemasterSvc, engagementSvc, messagingSvc, ext, ps,
	)
	userpin := usecases.NewUserPinUseCase(r, profile, ext, pinExt, engagementSvc)
	crm := hubspot.NewHubSpotService()
	su := usecases.NewSignUpUseCases(r, profile, userpin, supplier, ext, engagementSvc, ps, ediSvc)
	nhif := usecases.NewNHIFUseCases(r, profile, ext, engagementSvc)
	sms := usecases.NewSMSUsecase(r, ext)
	admin := usecases.NewAdminUseCases(r, engagementSvc, ext, userpin)
	agent := usecases.NewAgentUseCases(r, engagementSvc, ext, userpin)
	aitUssd := ussd.NewUssdUsecases(r, ext, profile, userpin, su, pinExt, ps)

	i, err := interactor.NewOnboardingInteractor(
		r, profile, su, supplier, login,
		survey, userpin, erpSvc, chargemasterSvc,
		engagementSvc, messagingSvc, nhif, ps, sms, aitUssd, crm, agent, admin, ediSvc,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}
	return i, nil

}
