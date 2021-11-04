package user_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	user "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	utilsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func createPINInputData(phone string, PIN string, ConfirmedPIN string, flavour feedlib.Flavour) *dto.PinInput {
	return &dto.PinInput{
		PhoneNumber:  phone,
		PIN:          PIN,
		ConfirmedPin: ConfirmedPIN,
		Flavour:      flavour,
	}
}

func TestUseCasesUserImpl_SetUserPIN_Unittest(t *testing.T) {
	ctx := context.Background()

	phone := "+2547100000000"
	PIN := "1234"
	ConfirmedPin := "1234"
	flavour := feedlib.FlavourConsumer

	validPINInput := createPINInputData(phone, PIN, ConfirmedPin, flavour)

	noPhoneInput := createPINInputData("", PIN, ConfirmedPin, flavour)

	invalidPINLength := createPINInputData(phone, "123", "123", flavour)

	pinNotDigits := createPINInputData(phone, "page", ConfirmedPin, flavour)

	pinMismatch := createPINInputData(phone, "4321", ConfirmedPin, flavour)

	type args struct {
		ctx   context.Context
		input *dto.PinInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: validPINInput,
			},
			wantErr: false,
		},

		{
			name: "Sad case - no phone",
			args: args{
				ctx:   ctx,
				input: noPhoneInput,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid PIN length",
			args: args{
				ctx:   ctx,
				input: invalidPINLength,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin not digits",
			args: args{
				ctx:   ctx,
				input: pinNotDigits,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin mismatch",
			args: args{
				ctx:   ctx,
				input: pinMismatch,
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to save PIN",
			args: args{
				ctx:   ctx,
				input: validPINInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUser := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeOnboardingLibraryExtension()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case - no phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to set user pin")
				}
			}
			if tt.name == "Sad case - invalid PIN length" {
				fakeUser.MockSetUserPINFn = func(ctx context.Context, input *dto.PinInput) (bool, error) {
					return false, fmt.Errorf("failed to set user pin")
				}
			}
			if tt.name == "Sad case - pin not digits" {
				fakeUser.MockSetUserPINFn = func(ctx context.Context, input *dto.PinInput) (bool, error) {
					return false, fmt.Errorf("failed to set user pin")
				}
			}

			if tt.name == "Sad case - pin mismatch" {
				fakeExtension.MockComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad case - fail to save PIN" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("failed to save PIN")
				}
			}

			_, err := u.SetUserPIN(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUseCasesUserImpl_Login_Unittest(t *testing.T) {
	ctx := context.Background()

	phoneNumber := "+2547100000000"
	PIN := "1234"
	flavour := feedlib.FlavourConsumer

	type args struct {
		ctx         context.Context
		phoneNumber string
		pin         string
		flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to get user profile by phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get user PIN By User ID",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254710000000",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - pin mismatch",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254710000000",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create firebase token",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to authenticate token",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeOnboardingLibraryExtension()
			fakeUtils := utilsMock.NewUtilsMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case - no phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad case - fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad case - unable to get user PIN By User ID" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get user PIN by user ID")
				}
			}

			if tt.name == "Sad case - check PIN expiry" {
				fakeUtils.MockCheckPINExpiryFn = func(currentTime time.Time, pinData *domain.UserPIN) bool {
					return false
				}
			}

			if tt.name == "Sad case - pin mismatch" {
				fakeExtension.MockComparePINFn = func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad Case - Fail to create firebase token" {
				fakeExtension.MockCreateFirebaseCustomTokenFn = func(ctx context.Context, uid string) (string, error) {
					return "", fmt.Errorf("failed to create custom token")
				}
			}

			if tt.name == "Sad Case - Fail to authenticate token" {
				fakeExtension.MockAuthenticateCustomFirebaseTokenFn = func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
					return nil, fmt.Errorf("failed to authenticate token")
				}
			}

			_, _, err := u.Login(tt.args.ctx, tt.args.phoneNumber, tt.args.pin, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
