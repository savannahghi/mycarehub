package user_test

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/segmentio/ksuid"
)

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
		{
			name: "Sad case - un-normalized phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "0710000000",
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
			fakeExtension := extensionMock.NewFakeExtension()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)

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

			if tt.name == "Sad case - un-normalized phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
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

func TestUnit_InviteUser(t *testing.T) {
	ctx := context.Background()

	invalidPhone := "invalid"
	validPhone := "+2547100000000"
	validPhoneWithNoCountryCode := "07100000000"

	validFlavour := feedlib.FlavourConsumer
	invalidFlavour := "INVALID_FLAVOUR"

	userID := ksuid.New().String()

	userOutput := &domain.User{
		ID:          &userID,
		DisplayName: "Test User",
		Username:    "testuser",
		FirstName:   "Test",
		LastName:    "User",
		MiddleName:  "",
		Active:      true,
		Gender:      enumutils.GenderFemale,
	}

	type args struct {
		ctx         context.Context
		userID      string
		phoneNumber string
		flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: valid phone number",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "valid: valid phone number without country code prefix",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhoneWithNoCountryCode,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: phone number is invalid",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: invalidPhone,
				flavour:     validFlavour,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "valid: valid flavour",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: invalid flavour",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     feedlib.Flavour(invalidFlavour),
			},
			wantErr: true,
			want:    false,
		},

		{
			name: "valid: fetched user by user ID",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: failed to get user by user ID",
			args: args{
				ctx:         ctx,
				userID:      "12345",
				phoneNumber: validPhone,
				flavour:     feedlib.Flavour(invalidFlavour),
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "valid: generated a temporary PIN successfully",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: failed go generate temporary pin",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "valid: saved temporary pin successfully",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: failed to save temporary pin",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "valid: get invite link success",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: get invite link error",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     feedlib.Flavour(invalidFlavour),
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "valid: send invite message success",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: send in message error",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()

			fakeUserMock := mock.NewUserUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "valid: valid phone number" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "valid: valid phone number without country code prefix" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "valid: valid flavour" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("phone number is invalid")
				}
			}

			if tt.name == "valid: valid phone number without country code prefix" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid: invalid flavour" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("flavour is invalid")
				}
			}

			if tt.name == "valid: fetched user by user ID" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return userOutput, nil
				}
			}
			if tt.name == "invalid: failed to get user by user ID" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user by user ID")
				}
			}

			if tt.name == "valid: generated a temporary PIN successfully" {
				fakeExtension.MockGenerateTempPINFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
			}

			if tt.name == "invalid: failed go generate temporary pin" {
				fakeExtension.MockGenerateTempPINFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to generate temporary pin")
				}
			}

			if tt.name == "valid: saved temporary pin successfully" {
				fakeDB.MockSaveTemporaryUserPinFn = func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: failed to save temporary pin" {
				fakeDB.MockSaveTemporaryUserPinFn = func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("failed to get user by user ID")
				}
			}

			if tt.name == "valid: get invite link success" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: get invite link error" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to get invite link")
				}
			}
			if tt.name == "valid: send invite message success" {
				fakeExtension.MockSendInviteSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return nil
				}
			}
			if tt.name == "invalid: send in message error" {
				fakeExtension.MockSendInviteSMSFn = func(ctx context.Context, phoneNumbers []string, message string) error {
					return fmt.Errorf("failed to send sms")
				}
			}

			got, err := us.InviteUser(tt.args.ctx, tt.args.userID, tt.args.phoneNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.InviteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.InviteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_SetUserPIN(t *testing.T) {
	ctx := context.Background()
	UserID := ksuid.New().String()
	PIN := "1234"
	longPIN := "12345"
	shortPIN := "123"
	tooLongPIN := strconv.Itoa(int(math.Pow(10, 6)))
	invalidPINString := "invalid"
	flavour := feedlib.FlavourConsumer

	nonMatchedPin := "0000"

	type args struct {
		ctx   context.Context
		input dto.PINInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: set user pin successfully",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &PIN,
					ConfirmPIN: &PIN,
					Flavour:    flavour,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: user not found",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &PIN,
					ConfirmPIN: &PIN,
					Flavour:    flavour,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: pin is not integer explicit",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &invalidPINString,
					ConfirmPIN: &invalidPINString,
					Flavour:    flavour,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: pin length long",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &longPIN,
					ConfirmPIN: &longPIN,
					Flavour:    flavour,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: pin too long pin",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &tooLongPIN,
					ConfirmPIN: &tooLongPIN,
					Flavour:    flavour,
				},
			},

			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: pin length short",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &shortPIN,
					ConfirmPIN: &shortPIN,
					Flavour:    flavour,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: confirm pin mismatch",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &UserID,
					PIN:        &PIN,
					ConfirmPIN: &nonMatchedPin,
					Flavour:    flavour,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()

			// fakeUserMock := mock.NewUserUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "invalid: user not found" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user by user ID")
				}
			}

			if tt.name == "invalid: pin is not integer explicit" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("only number input allowed")
				}
			}
			if tt.name == "invalid: pin length long" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("pin length id longer than 4")
				}
			}

			if tt.name == "invalid: pin too long pin" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("pin length is too long")
				}
			}

			if tt.name == "invalid: pin length short" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("pin length is too short")
				}
			}

			if tt.name == "invalid: confirm pin mismatch" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("confirm pin does not mach the pin")
				}
			}

			got, err := us.SetUserPIN(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.SetUserPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}
