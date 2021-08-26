package usecases_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"

	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"

	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"

	mockRepo "github.com/savannahghi/onboarding/pkg/onboarding/repository/mock"

	extMock "github.com/savannahghi/onboarding/pkg/onboarding/application/extension/mock"
	engagementMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement/mock"

	messagingMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging/mock"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	pubsubmessagingMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub/mock"
	adminSrv "github.com/savannahghi/onboarding/pkg/onboarding/usecases/admin"
)

const (
	otpService        = "otp"
	engagementService = "engagement"
)

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	envOriginalValue := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")
	emailOriginalValue := os.Getenv("SAVANNAH_ADMIN_EMAIL")
	os.Setenv("SAVANNAH_ADMIN_EMAIL", "test@bewell.co.ke")
	debugEnvValue := os.Getenv("DEBUG")
	os.Setenv("DEBUG", "true")
	os.Setenv("REPOSITORY", "firebase")
	collectionEnvValue := os.Getenv("ROOT_COLLECTION_SUFFIX")
	// !NOTE!
	// Under no circumstances should you remove this env var when testing
	// You risk purging important collections, like our prod collections
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("onboarding_ci_%v", time.Now().Unix()))

	ctx := context.Background()

	fsc, fbc := InitializeTestFirebaseClient(ctx)
	if fsc == nil {
		log.Panicf("failed to initialize test FireStore client")
	}
	if fbc == nil {
		log.Panicf("failed to initialize test FireBase client")
	}

	purgeRecords := func() {
		if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
			r := fb.Repository{}
			collections := []string{
				r.GetCustomerProfileCollectionName(),
				r.GetPINsCollectionName(),
				r.GetUserProfileCollectionName(),
				r.GetSupplierProfileCollectionName(),
				r.GetSurveyCollectionName(),
				r.GetCommunicationsSettingsCollectionName(),
				r.GetCustomerProfileCollectionName(),
				r.GetExperimentParticipantCollectionName(),
				r.GetKCYProcessCollectionName(),
				r.GetNHIFDetailsCollectionName(),
				r.GetProfileNudgesCollectionName(),
				r.GetRolesCollectionName(),
			}
			for _, collection := range collections {
				ref := fsc.Collection(collection)
				firebasetools.DeleteCollection(ctx, fsc, ref, 10)
			}
		}

	}

	// try clean up first
	purgeRecords()

	// do clean up
	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	// restore environment variables to original values
	os.Setenv(envOriginalValue, "ENVIRONMENT")
	os.Setenv(emailOriginalValue, "SAVANNAH_ADMIN_EMAIL")
	os.Setenv("DEBUG", debugEnvValue)
	os.Setenv("ROOT_COLLECTION_SUFFIX", collectionEnvValue)

	os.Exit(code)
}

func InitializeTestFirebaseClient(ctx context.Context) (*firestore.Client, *auth.Client) {
	fc := firebasetools.FirebaseClient{}
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
	engage := engagement.NewServiceEngagementImpl(engagementClient, ext)

	ps, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		ext,
		repo,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize new pubsub messaging service: %w", err)
	}
	mes := messaging.NewServiceMessagingImpl(ext)
	pinExt := extension.NewPINExtensionImpl()
	profile := usecases.NewProfileUseCase(repo, ext, engage, ps)

	supplier := usecases.NewSupplierUseCases(repo, profile, engage, mes, ext, ps)
	login := usecases.NewLoginUseCases(repo, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(repo, ext)
	userpin := usecases.NewUserPinUseCase(repo, profile, ext, pinExt, engage)
	su := usecases.NewSignUpUseCases(repo, profile, userpin, supplier, ext, engage, ps)
	nhif := usecases.NewNHIFUseCases(repo, profile, ext, engage)

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
	}, nil
}

func generateTestOTP(t *testing.T, phone string) (*profileutils.OtpResponse, error) {
	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	testAppID := uuid.New().String()
	return s.Engagement.GenerateAndSendOTP(ctx, phone, &testAppID)
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
	phone := interserviceclient.TestUserPhoneNumber
	flavour := feedlib.FlavourConsumer
	pin := interserviceclient.TestUserPin
	testAppID := uuid.New().String()
	otp, err := s.Signup.VerifyPhoneNumber(ctx, phone, &testAppID)
	if err != nil {
		if strings.Contains(err.Error(), exceptions.CheckPhoneNumberExistError().Error()) {
			logInCreds, err := s.Login.LoginByPhone(
				ctx,
				phone,
				interserviceclient.TestUserPin,
				flavour,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to log in test user: %v", err)
			}

			return &auth.Token{
				UID: logInCreds.Auth.UID,
			}, nil
		}

		return nil, fmt.Errorf("failed to check if test phone exists: %v", err)
	}

	u, err := s.Signup.CreateUserByPhone(
		ctx,
		&dto.SignUpInput{
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

	return &auth.Token{
		UID: u.Auth.UID,
	}, nil
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
		firebasetools.AuthTokenContextKey,
		auth,
	)
	return authenticatedContext, auth, nil
}

func TestGetTestAuthenticatedContext(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	if ctx == nil {
		t.Errorf("nil context")
		return
	}
	if auth == nil {
		t.Errorf("nil auth data")
		return
	}
}

func TestLoginUseCasesImpl_LoginByPhone(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	flavour := feedlib.FlavourConsumer
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx     context.Context
		phone   string
		PIN     string
		flavour feedlib.Flavour
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
				phone:   interserviceclient.TestUserPhoneNumber,
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "sad case: wrong pin number supplied",
			args: args{
				ctx:     ctx,
				phone:   interserviceclient.TestUserPhoneNumber,
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
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect phone number",
			args: args{
				ctx:     ctx,
				phone:   "+2541234",
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: true,
		},
		{
			name: "sad case: incorrect flavour",
			args: args{
				ctx:     ctx,
				phone:   interserviceclient.TestUserPhoneNumber,
				PIN:     interserviceclient.TestUserPin,
				flavour: "not-a-correct-flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResponse, err := s.Login.LoginByPhone(
				tt.args.ctx,
				tt.args.phone,
				tt.args.PIN,
				tt.args.flavour,
			)
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
var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakePinExt extMock.PINExtensionImpl
var fakeEngagementSvs engagementMock.FakeServiceEngagement
var fakeMessagingSvc messagingMock.FakeServiceMessaging
var fakePubSub pubsubmessagingMock.FakeServicePubSub

// InitializeFakeOnboaridingInteractor represents a fakeonboarding interactor
func InitializeFakeOnboardingInteractor() (*interactor.Interactor, error) {
	var r repository.OnboardingRepository = &fakeRepo
	var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs
	var messagingSvc messaging.ServiceMessaging = &fakeMessagingSvc
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt
	var ps pubsubmessaging.ServicePubSub = &fakePubSub

	profile := usecases.NewProfileUseCase(r, ext, engagementSvc, ps)
	login := usecases.NewLoginUseCases(r, profile, ext, pinExt)
	survey := usecases.NewSurveyUseCases(r, ext)
	supplier := usecases.NewSupplierUseCases(
		r, profile, engagementSvc, messagingSvc, ext, ps,
	)
	userpin := usecases.NewUserPinUseCase(r, profile, ext, pinExt, engagementSvc)
	su := usecases.NewSignUpUseCases(r, profile, userpin, supplier, ext, engagementSvc, ps)
	nhif := usecases.NewNHIFUseCases(r, profile, ext, engagementSvc)
	adminSrv := adminSrv.NewService(ext)
	role := usecases.NewRoleUseCases(r, ext)

	i, err := interactor.NewOnboardingInteractor(
		profile, su, supplier, login,
		survey, userpin,
		engagementSvc, messagingSvc, nhif, ps,
		adminSrv, role,
	)
	if err != nil {
		return nil, fmt.Errorf("can't instantiate service : %w", err)
	}
	return i, nil

}
