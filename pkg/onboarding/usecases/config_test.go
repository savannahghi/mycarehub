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
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
)

// NewInterServiceClient initializes an external service in the correct environment given its name
func InitializeTestService(
	ctx context.Context,
) (infrastructure.Interactor, error) {
	infra := infrastructure.NewInteractor()
	return infra, nil

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
				r.GetPINsCollectionName(),
				r.GetUserProfileCollectionName(),
				r.GetSurveyCollectionName(),
				r.GetCommunicationsSettingsCollectionName(),
				r.GetExperimentParticipantCollectionName(),
				r.GetProfileNudgesCollectionName(),
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

	// restore environment varibles to original values
	os.Setenv(envOriginalValue, "ENVIRONMENT")
	os.Setenv(emailOriginalValue, "SAVANNAH_ADMIN_EMAIL")
	os.Setenv("DEBUG", debugEnvValue)
	os.Setenv("ROOT_COLLECTION_SUFFIX", collectionEnvValue)

	os.Exit(code)
}

// CreateTestUserByPhone creates a user that is to be used in
// running of our test cases.
// If the test user already exists then they are logged in
// to get their auth credentials
func CreateOrLoginTestUserByPhone(t *testing.T) (*auth.Token, error) {
	ctx := context.Background()
	i, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service")
	}
	phone := "+254720000000"
	flavour := feedlib.FlavourConsumer
	pin := interserviceclient.TestUserPin
	testAppID := uuid.New().String()
	otp, err := i.VerifyPhoneNumber(ctx, phone, &testAppID)
	if err != nil {
		if strings.Contains(err.Error(), exceptions.CheckPhoneNumberExistError().Error()) {
			logInCreds, err := i.LoginByPhone(
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

	u, err := i.CreateUserByPhone(
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

func generateTestOTP(t *testing.T, phone string) (*profileutils.OtpResponse, error) {
	ctx := context.Background()
	i, err := InitializeTestService(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service: %v", err)
	}
	testAppID := uuid.New().String()
	return i.GenerateAndSendOTP(ctx, phone, &testAppID)
}
