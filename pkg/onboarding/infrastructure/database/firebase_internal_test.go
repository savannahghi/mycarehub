package database

import (
	"context"
	"fmt"
	"log"
	"testing"

	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

func InitializeTestService(ctx context.Context) (*interactor.Interactor, error) {
	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		return nil, err
	}
	ext := extension.NewBaseExtensionImpl()
	otp := otp.NewOTPService(fr, ext)
	profile := usecases.NewProfileUseCase(fr, otp, ext)
	erp := erp.NewERPService(fr)
	chrg := chargemaster.NewChargeMasterUseCasesImpl(fr)
	engage := engagement.NewServiceEngagementImpl(fr)
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes, ext)
	login := usecases.NewLoginUseCases(fr, profile, ext)
	survey := usecases.NewSurveyUseCases(fr, ext)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile, ext)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier, otp, ext)

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
	}, nil
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
	exists, err := s.Signup.CheckPhoneExists(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check if test phone exists: %v", err)
	}
	if !exists {
		otp, err := s.Otp.GenerateAndSendOTP(ctx, phone)
		log.Println("The otp is:", otp)
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

func Test_checkIdentifierExists(t *testing.T) {

	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	fr, err := NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create new Firebase Repository: %v", err)
		return
	}

	userProfile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	if err != nil {
		t.Errorf("failed to get a user profile")
		return
	}
	type args struct {
		profile *base.UserProfile
		UID     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Happy Case - Check if identifier exists",
			args: args{
				profile: userProfile,
				UID:     auth.UID,
			},
			want: true,
		},
		{
			name: "Sad Case - invalid UID",
			args: args{
				profile: userProfile,
				UID:     "invalid",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIdentifierExists(tt.args.profile, tt.args.UID); got != tt.want {
				t.Errorf("checkIdentifierExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
