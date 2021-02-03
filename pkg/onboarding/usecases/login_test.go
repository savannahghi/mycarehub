package usecases_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"

	otpMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp/mock"
	mockRepo "gitlab.slade360emr.com/go/profile/pkg/onboarding/repository/mock"

	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	chargemasterMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster/mock"

	engagementMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement/mock"

	erpMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp/mock"

	mailgunMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun/mock"

	messagingMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging/mock"
)

const (
	otpService        = "otp"
	mailgunService    = "mailgun"
	engagementService = "engagement"
)

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	envOriginalValue := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")
	debugEnvValue := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "true")
	collectionEnvValue := os.Getenv("ROOT_COLLECTION_SUFFIX")
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("onboarding_ci_%v", time.Now().Unix()))

	ctx := context.Background()
	r := database.Repository{} // They are nil
	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}

	purgeRecords := func() {
		collections := []string{
			r.GetCustomerProfileCollectionName(),
			r.GetPINsCollectionName(),
			r.GetUserProfileCollectionName(),
			r.GetSupplierProfileCollectionName(),
			r.GetSurveyCollectionName(),
			r.GetKCYProcessCollectionName(),
		}
		for _, collection := range collections {
			ref := fsc.Collection(collection)
			base.DeleteCollection(ctx, fsc, ref, 10)
		}
	}

	// try clean up first
	purgeRecords()

	// do clean up
	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	// restore environment varibles to original values
	os.Setenv(envOriginalValue, "ENVIRONMENT")
	os.Setenv("DEBUG", debugEnvValue)
	os.Setenv("ROOT_COLLECTION_SUFFIX", collectionEnvValue)

	os.Exit(code)
}

func InitializeTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	fc := base.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Panicf("unable to initialize Firebase: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Panicf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up tests: %s", err)
	}
	return fsc, fbc
}

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

	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)
	if err != nil {
		return nil, err
	}

	ext := extension.NewBaseExtensionImpl()

	// Initialize ISC clients
	otpClient := utils.NewInterServiceClient(otpService, ext)
	mailgunClient := utils.NewInterServiceClient(mailgunService, ext)
	engagementClient := utils.NewInterServiceClient(engagementService, ext)

	erp := erp.NewERPService()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	engage := engagement.NewServiceEngagementImpl(engagementClient)
	mg := mailgun.NewServiceMailgunImpl(mailgunClient)
	mes := messaging.NewServiceMessagingImpl(ext)
	pinExt := extension.NewPINExtensionImpl()
	otp := otp.NewOTPService(otpClient, ext)
	profile := usecases.NewProfileUseCase(fr, otp, ext, engage)
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes, ext)
	login := usecases.NewLoginUseCases(fr, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(fr, ext)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile, ext, pinExt)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier, otp, ext)
	nhif := usecases.NewNHIFUseCases(fr, profile, ext)

	return &interactor.Interactor{
		Onboarding:   profile,
		Signup:       su,
		Otp:          otp,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engagement:   engage,
		NHIF:         nhif,
	}, nil
}

func generateTestOTP(t *testing.T, phone string) (*base.OtpResponse, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	return s.Otp.GenerateAndSendOTP(ctx, phone)
}

// CreateTestUserByPhone creates a user that is to be used in
// running of our test cases.
// If the test user already exists then they are logged in
// to get their auth credentials
func CreateOrLoginTestUserByPhone(t *testing.T) (*auth.Token, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service")
	}
	phone := base.TestUserPhoneNumber
	flavour := base.FlavourConsumer
	pin := base.TestUserPin
	exists, err := s.Onboarding.CheckPhoneExists(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check if test phone exists: %v", err)
	}
	if !exists {
		otp, err := generateTestOTP(t, phone)
		if err != nil {
			return nil, fmt.Errorf("failed to generate test OTP: %v", err)
		}

		u, err := s.Signup.CreateUserByPhone(
			ctx,
			&resources.SignUpInput{
				PhoneNumber: &phone,
				PIN:         &pin,
				Flavour:     flavour,
				OTP:         &otp.OTP,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create a test user: %v", err)
		}
		if u == nil {
			return nil, fmt.Errorf("nil test user response")
		}
		authCred := &auth.Token{
			UID: u.Auth.UID,
		} // We add the test user UID to the expected auth.Token
		return authCred, nil
	}
	logInCreds, err := s.Login.LoginByPhone(
		ctx,
		phone,
		base.TestUserPin,
		flavour,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to log in test user: %v", err)
	}
	authCred := &auth.Token{
		UID: logInCreds.Auth.UID,
	}
	return authCred, nil
}

// TestAuthenticatedContext returns a logged in context, useful for test purposes
func GetTestAuthenticatedContext(t *testing.T) (context.Context, *auth.Token, error) {
	ctx := context.Background()
	auth, err := CreateOrLoginTestUserByPhone(t)
	if err != nil {
		return nil, nil, err
	}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		auth,
	)
	return authenticatedContext, auth, nil
}
func TestLoginUseCasesImpl_LoginByPhone(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	flavour := base.FlavourConsumer
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx     context.Context
		phone   string
		PIN     string
		flavour base.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: valid login",
			args: args{
				ctx:     ctx,
				phone:   base.TestUserPhoneNumber,
				PIN:     base.TestUserPin,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "sad case: wrong pin number supplied",
			args: args{
				ctx:     ctx,
				phone:   base.TestUserPhoneNumber,
				PIN:     "4567",
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: user profile without a primary phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2547900900", // not a primary phone number
				PIN:     base.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2541234",
				PIN:     base.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect flavour",
			args: args{
				ctx:     ctx,
				phone:   base.TestUserPhoneNumber,
				PIN:     base.TestUserPin,
				flavour: "not-a-correct-flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResponse, err := s.Login.LoginByPhone(tt.args.ctx, tt.args.phone, tt.args.PIN, tt.args.flavour)
			if tt.wantErr && authResponse != nil {
				t.Errorf("expected nil auth response but got %v, since the error %v occurred",
					authResponse,
					err,
				)
				return
			}

			if !tt.wantErr && authResponse == nil {
				t.Errorf("expected an auth response but got nil, since no error occurred")
				return
			}
		})
	}
}

var fakeRepo mockRepo.FakeOnboardingRepository
var fakeOtp otpMock.FakeServiceOTP
var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakePinExt extMock.PINExtensionImpl
var fakeMailgunSvc mailgunMock.FakeServiceMailgun
var fakeEngagementSvs engagementMock.FakeServiceEngagement
var fakeMessagingSvc messagingMock.FakeServiceMessaging
var fakeEPRSvc erpMock.FakeServiceERP
var fakeChargeMasterSvc chargemasterMock.FakeServiceChargeMaster

// InitializeFakeOnboaridingInteractor represents a fakeonboarding interactor
func InitializeFakeOnboaridingInteractor() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo
	var otpSvc otp.ServiceOTP = &fakeOtp
	var erpSvc erp.ServiceERP = &fakeEPRSvc
	var chargemasterSvc chargemaster.ServiceChargeMaster = &fakeChargeMasterSvc
	var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs
	var mailgunSvc mailgun.ServiceMailgun = &fakeMailgunSvc
	var messagingSvc messaging.ServiceMessaging = &fakeMessagingSvc
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt

	profile := usecases.NewProfileUseCase(r, otpSvc, ext, engagementSvc)
	login := usecases.NewLoginUseCases(r, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(r, ext)
	supplier := usecases.NewSupplierUseCases(
		r, profile, erpSvc, chargemasterSvc, engagementSvc, mailgunSvc, messagingSvc, ext,
	)
	userpin := usecases.NewUserPinUseCase(r, otpSvc, profile, ext, pinExt)
	su := usecases.NewSignUpUseCases(r, profile, userpin, supplier, otpSvc, ext)
	nhif := usecases.NewNHIFUseCases(r, profile, ext)

	i, err := interactor.NewOnboardingInteractor(
		r, profile, su, otpSvc, supplier, login,
		survey, userpin, erpSvc, chargemasterSvc, engagementSvc, mailgunSvc, messagingSvc, nhif,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}
	return i, nil

}
