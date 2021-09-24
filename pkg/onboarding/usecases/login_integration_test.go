package usecases_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	baseExt "github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	libInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	libOnboardingUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

func InitializeTestService(ctx context.Context) (*usecases.UseCaseLoginImpl, *usecases.UseCaseSignUpImpl, error) {
	i := libInfra.NewInfrastructureInteractor()
	var fc firebasetools.IFirebaseClient
	baseExtension := baseExt.NewBaseExtensionImpl(fc)
	pinExtension := baseExt.NewPINExtensionImpl()
	profile := libOnboardingUsecase.NewProfileUseCase(i, baseExtension)
	libUsecase := libOnboardingUsecase.NewLoginUseCases(i, profile, baseExtension, pinExtension)
	userPinUseCase := libOnboardingUsecase.NewUserPinUseCase(i, profile, baseExtension, pinExtension)
	libSignUp := libOnboardingUsecase.NewSignUpUseCases(i, profile, userPinUseCase, baseExtension)
	loginUsecase := usecases.NewUseCaseLogin(libUsecase)
	signUpUsecase := usecases.NewSignUpUseCase(libSignUp)
	return loginUsecase, signUpUsecase, nil
}

// CreateTestUserByPhone creates a user that is to be used in
// running of our test cases.
// If the test user already exists then they are logged in
// to get their auth credentials
func CreateOrLoginTestUserByPhone(t *testing.T) (*auth.Token, error) {
	ctx := context.Background()
	l, s, err := InitializeTestService(ctx)
	fmt.Printf("43: THE ERROR IS: %v\n", err)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize test service")
	}
	phone := "+254720000000"
	flavour := feedlib.FlavourConsumer
	pin := interserviceclient.TestUserPin
	testAppID := uuid.New().String()
	otp, err := s.LibSignUp.VerifyPhoneNumber(ctx, phone, &testAppID)
	fmt.Printf("52: THE OTP IS: %v\n", otp)
	fmt.Printf("53: THE ERROR IS: %v\n", err)
	if err != nil {
		if strings.Contains(err.Error(), exceptions.CheckPhoneNumberExistError().Error()) {
			logInCreds, err := l.LibLogin.LoginByPhone(
				ctx,
				phone,
				interserviceclient.TestUserPin,
				flavour,
			)
			fmt.Printf("62: THE ERROR IS: %v\n", err)
			if err != nil {
				return nil, fmt.Errorf("failed to log in test user: %v", err)
			}

			return &auth.Token{
				UID: logInCreds.Auth.UID,
			}, nil
		}

		return nil, fmt.Errorf("failed to check if test phone exists: %v", err)
	}

	u, err := s.LibSignUp.CreateUserByPhone(
		ctx,
		&dto.SignUpInput{
			PhoneNumber: &phone,
			PIN:         &pin,
			Flavour:     flavour,
			OTP:         &otp.OTP,
		},
	)
	fmt.Printf("84: THE ERROR IS: %v\n", err)
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
	fmt.Printf("99: THE ERROR IS: %v\n", err)
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
	fmt.Printf("THE ERROR IS: %v\n\n", err)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	flavour := feedlib.FlavourConsumer
	phone := "+254720000000"
	l, _, err := InitializeTestService(ctx)
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
				phone:   phone,
				PIN:     interserviceclient.TestUserPin,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "sad case: wrong pin number supplied",
			args: args{
				ctx:     ctx,
				phone:   phone,
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
				phone:   phone,
				PIN:     interserviceclient.TestUserPin,
				flavour: "not-a-correct-flavour",
			},
			// TODO: // Return this to true
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authResponse, err := l.LibLogin.LoginByPhone(
				tt.args.ctx,
				tt.args.phone,
				tt.args.PIN,
				tt.args.flavour,
			)
			fmt.Printf("214: THE ERROR IS: %v\n", err)
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
