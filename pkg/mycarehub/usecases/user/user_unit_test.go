package user_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	openSourceDto "github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/scalarutils"
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
		{
			name: "invalid: invalid flavour",
			args: args{
				ctx:         ctx,
				phoneNumber: "0710000000",
				pin:         PIN,
				flavour:     feedlib.Flavour("Invalid_flavour"),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update successful login time",
			args: args{
				ctx:         ctx,
				phoneNumber: "0710000000",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile by user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: "0710000000",
				pin:         PIN,
				flavour:     flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff profile by user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeUserMock := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Sad case - no phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad case - fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad case - unable to get user PIN By User ID" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
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
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "invalid: invalid flavour" {
				fakeUserMock.MockLoginFn = func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error) {
					return nil, 2, fmt.Errorf("invalid flavour defined")
				}
			}

			if tt.name == "Sad Case - Fail to update successful login time" {
				fakeDB.MockUpdateUserLastSuccessfulLoginTimeFn = func(ctx context.Context, userID string) error {
					return fmt.Errorf("failed to update last successfult login time")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile by user ID" {
				fakeDB.MockGetClientProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad Case - Fail to get staff profile by user ID" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
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
		ID:       &userID,
		Name:     "Test User",
		Username: "testuser",
		Active:   true,
		Gender:   enumutils.GenderFemale,
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
				flavour:     validFlavour,
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
		{
			name: "Sad Case - Fail to invalidate pin",
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
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

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

			if tt.name == "Sad Case - Fail to invalidate pin" {
				fakeDB.MockInvalidatePINFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to invalidate pin")
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
				fakeExtension.MockSendSMSFn = func(ctx context.Context, phoneNumbers string, message string, from enumutils.SenderID) (*openSourceDto.SendMessageResponse, error) {
					return &openSourceDto.SendMessageResponse{
						SMSMessageData: &openSourceDto.SMS{
							Recipients: []openSourceDto.Recipient{
								{
									Number: interserviceclient.TestUserPhoneNumber,
								},
							},
						},
					}, nil
				}
			}
			if tt.name == "invalid: send in message error" {
				fakeExtension.MockSendInviteSMSFn = func(ctx context.Context, phoneNumber, message string) error {
					return fmt.Errorf("failed to send SMS")
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
	invalidInput := ""
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
		{
			name: "Sad Case - Empty input",
			args: args{
				ctx: ctx,
				input: dto.PINInput{
					UserID:     &invalidInput,
					PIN:        &invalidInput,
					ConfirmPIN: &invalidInput,
					Flavour:    flavour,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to invalidate pin",
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
			name: "Sad Case - Fail to save pin",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()

			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

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
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad Case - Empty input" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("empty input")
				}
			}

			if tt.name == "Sad Case - Fail to save pin" {
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, fmt.Errorf("empty input")
				}
			}

			if tt.name == "Sad Case - Fail to invalidate pin" {
				fakeDB.MockInvalidatePINFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to invalidate pin")
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

func TestUseCasesUserImpl_VerifyLoginPIN(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		pin     string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully verify pin",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user pin",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "3456",
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Pin mismatch",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by ID",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user login count",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update last failed login time",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update next allowed login time",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid flavour",
			args: args{
				ctx:     ctx,
				userID:  "12345",
				pin:     "1234",
				flavour: "Invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUserMock := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Happy Case - Successfully verify pin" {
				fakeUserMock.MockVerifyLoginPINFn = func(ctx context.Context, userID string, pin string, flavour feedlib.Flavour) (bool, int, error) {
					return true, 0, nil
				}
			}

			if tt.name == "Sad Case - Fail to get user pin" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get a pin")
				}
			}

			if tt.name == "Sad Case - Invalid flavour" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get a pin")
				}
			}

			if tt.name == "Sad Case - Pin mismatch" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by ID" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user ID")
				}
			}

			if tt.name == "Sad Case - Fail to update user login count" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}

				fakeDB.MockUpdateUserFailedLoginCountFn = func(ctx context.Context, userID string, failedLoginAttempts int) error {
					return fmt.Errorf("failed to update user failed login count")
				}
			}

			if tt.name == "Sad Case - Fail to update last failed login time" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}

				fakeDB.MockUpdateUserLastFailedLoginTimeFn = func(ctx context.Context, userID string) error {
					return fmt.Errorf("failed to update user failed login count")
				}
			}

			if tt.name == "Sad Case - Fail to update next allowed login time" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}

				fakeDB.MockUpdateUserNextAllowedLoginTimeFn = func(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
					return fmt.Errorf("failed to update user failed login count")
				}
			}

			got, _, err := u.VerifyLoginPIN(tt.args.ctx, tt.args.userID, tt.args.pin, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.VerifyLoginPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.VerifyLoginPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_SetNickName(t *testing.T) {
	ctx := context.Background()

	userID := ksuid.New().String()
	nickname := gofakeit.BeerName()

	type args struct {
		ctx      context.Context
		userID   *string
		nickname *string
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
				ctx:      ctx,
				userID:   &userID,
				nickname: &nickname,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				userID:   &userID,
				nickname: &nickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:      ctx,
				userID:   nil,
				nickname: &nickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no nickname",
			args: args{
				ctx:      ctx,
				userID:   &userID,
				nickname: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Both userID and nickname nil",
			args: args{
				ctx:      ctx,
				userID:   nil,
				nickname: nil,
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
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Sad case" {
				fakeDB.MockSetNickNameFn = func(ctx context.Context, userID, nickname *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeDB.MockSetNickNameFn = func(ctx context.Context, userID, nickname *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no nickname" {
				fakeDB.MockSetNickNameFn = func(ctx context.Context, userID, nickname *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Both userID and nickname nil" {
				fakeDB.MockSetNickNameFn = func(ctx context.Context, userID, nickname *string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := u.SetNickName(tt.args.ctx, tt.args.userID, tt.args.nickname)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetNickName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.SetNickName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_RequestPINReset(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		phoneNumber string
		flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully request pin reset",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				flavour:     feedlib.FlavourConsumer,
			},
			want:    "111222",
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: "0732313",
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid: invalid flavour",
			args: args{
				ctx:         ctx,
				phoneNumber: "0710000000",
				flavour:     feedlib.Flavour("Invalid_flavour"),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to check if user has pin",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to save otp",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUser := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Sad Case - Invalid phonenumber" {
				fakeUser.MockRequestPINResetFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("invalid phonenumber")
				}
			}

			if tt.name == "invalid: invalid flavour" {
				fakeUser.MockRequestPINResetFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("invalid flavour defined")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad Case - Fail to check if user has pin" {
				fakeDB.MockCheckUserHasPinFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to check if user has pin")
				}
			}

			if tt.name == "Sad Case - Fail to save otp" {
				fakeDB.MockSaveOTPFn = func(ctx context.Context, otpInput *domain.OTP) error {
					return fmt.Errorf("failed to save otp")
				}
			}

			got, err := us.RequestPINReset(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RequestPINReset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.RequestPINReset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_ResetPIN(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.UserResetPinInput
	}
	tests := []struct {
		name string

		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully reset pin",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid phone number",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: "str",
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid phone flavor",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     "invalid",
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: string pin",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "abcd",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid pin length",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "12345",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to get user profile by phone",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to verify OTP",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to get security question responses",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: security question responses are incorrect",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to invalidate pin",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to save pin",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
					PIN:         "1234",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid reset pin input",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "111222",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUser := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Happy Case - Successfully reset pin" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
			}
			if tt.name == "invalid: failed to get user profile by phone" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, errors.New("failed to get user profile by phone")
				}
			}

			if tt.name == "invalid: failed to verify OTP" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "invalid: failed to get security question responses" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return nil, errors.New("failed to get user security question responses")
				}
			}

			if tt.name == "invalid: security question responses are incorrect" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  false,
						},
					}, nil
				}
			}

			if tt.name == "invalid: failed to invalidate pin" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
				fakeDB.MockInvalidatePINFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, errors.New("failed to invalidate pin")
				}
			}

			if tt.name == "invalid: failed to save pin" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
				fakeDB.MockSavePinFn = func(ctx context.Context, pin *domain.UserPIN) (bool, error) {
					return false, errors.New("failed to save pin")
				}
			}

			if tt.name == "invalid: invalid reset pin input" {
				fakeUser.MockResetPINFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := us.ResetPIN(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.ResetPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.ResetPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_RefreshToken(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully refresh a token",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create firebase custom token",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to authenticate firebase custom token",
			args: args{
				ctx:   ctx,
				token: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Sad Case - Fail to create firebase custom token" {
				fakeExtension.MockCreateFirebaseCustomTokenFn = func(ctx context.Context, uid string) (string, error) {
					return "", fmt.Errorf("failed to create firebase custom token")
				}
			}

			if tt.name == "Sad Case - Fail to authenticate firebase custom token" {
				fakeExtension.MockAuthenticateCustomFirebaseTokenFn = func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
					return nil, fmt.Errorf("failed to authenticate custom token")
				}
			}

			got, err := us.RefreshToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_VerifyPIN(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
		pin     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully verify pin",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
				pin:     "1234",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: missing user id",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
				pin:     "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: missing pin",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: missing flavour",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
				pin:    "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: invalid flavour",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.Flavour("invalid flavour"),
				pin:     "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: failed to get user pin by user id",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
				pin:     "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: pin mismatch",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
				pin:     "1234",
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "invalid: failed to compare pin",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
				pin:     "1234",
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "invalid: failed to get user pin by user id" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get user pin by user id")
				}
			}

			if tt.name == "invalid: pin mismatch" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: gofakeit.UUID(),
						ValidFrom: time.Now().AddDate(0, 0, -1),
						ValidTo:   time.Now().AddDate(0, 0, 10),
						Flavour:   feedlib.FlavourConsumer,
						IsValid:   true,
						Salt:      gofakeit.UUID(),
					}, nil
				}
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			if tt.name == "invalid: failed to compare pin" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
			}

			got, err := us.VerifyPIN(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.VerifyPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.VerifyPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_GetClientCaregiver(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Caregiver
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "valid: no caregiver assigned",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "invalid: missing client id",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "invalid: failed to get client by id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "valid: no caregiver assigned" {
				fakeDB.MockGetClientCaregiverFn = func(ctx context.Context, clientID string) (*domain.Caregiver, error) {
					return &domain.Caregiver{}, nil
				}
			}
			if tt.name == "invalid: failed to get client caregiver" {
				fakeDB.MockGetClientCaregiverFn = func(ctx context.Context, clientID string) (*domain.Caregiver, error) {
					return nil, fmt.Errorf("failed to get client caregiver")
				}
			}

			if tt.name == "invalid: failed to get client by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client by id")
				}
			}
			_, err := us.GetClientCaregiver(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_CreateOrUpdateClientCaregiver(t *testing.T) {
	type args struct {
		ctx            context.Context
		caregiverInput *dto.CaregiverInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: update if client id was provided",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					ClientID:      uuid.New().String(),
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid: create if no client id was provided",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid phone number",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					ClientID:      uuid.New().String(),
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.BeerAlcohol(),
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid phone caregiver type",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					ClientID:      uuid.New().String(),
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: "invalid",
				},
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "invalid: failed to get client by id",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					ClientID:      uuid.New().String(),
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to update caregiver",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					ClientID:      uuid.New().String(),
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			want:    false,
			wantErr: true,
		},

		{
			name: "invalid: failed to create caregiver",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: invalid caregiver type",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					FirstName:     gofakeit.FirstName(),
					LastName:      gofakeit.LastName(),
					PhoneNumber:   gofakeit.Phone(),
					CaregiverType: "enums.CaregiverTypeFather",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "invalid: failed to get client by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client by id")
				}
			}

			if tt.name == "invalid: failed to update caregiver" {
				ID := uuid.New().String()
				client := &domain.ClientProfile{
					ID:          &ID,
					CaregiverID: &ID,
				}
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return client, nil
				}

				fakeDB.MockUpdateClientCaregiverFn = func(ctx context.Context, caregiver *dto.CaregiverInput) error {
					return fmt.Errorf("failed to update caregiver")
				}
			}

			if tt.name == "invalid: failed to create caregiver" {
				fakeDB.MockCreateClientCaregiverFn = func(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
					return fmt.Errorf("failed to create caregiver")
				}
			}

			got, err := us.CreateOrUpdateClientCaregiver(tt.args.ctx, tt.args.caregiverInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.CreateOrUpdateClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.CreateOrUpdateClientCaregiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_CompleteOnboardingTour(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
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
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and flavour",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Sad case - no userID" {
				fakeDB.MockUpdateUserPinChangeRequiredStatusFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case - no userID and flavour" {
				fakeDB.MockUpdateUserPinChangeRequiredStatusFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}

			}

			got, err := us.CompleteOnboardingTour(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.CompleteOnboardingTour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.CompleteOnboardingTour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterClient(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *dto.ClientRegistrationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register client",
			args: args{
				ctx: context.Background(),
				input: &dto.ClientRegistrationInput{
					Facility: "Test Facility",
					DateOfBirth: scalarutils.Date{
						Year:  1990,
						Month: 3,
						Day:   12,
					},
					ClientName:  gofakeit.FirstName(),
					Gender:      enumutils.GenderFemale,
					PhoneNumber: "+254700000000",
					CCCNumber:   "5432",
					Counselled:  true,
					EnrollmentDate: scalarutils.Date{
						Year:  1990,
						Month: 3,
						Day:   12,
					},
					ClientType: enums.ClientTypeDreams,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to make request",
			args: args{
				ctx: context.Background(),
				input: &dto.ClientRegistrationInput{
					Facility: "Kanairo",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - User not authorized",
			args: args{
				ctx: context.Background(),
				input: &dto.ClientRegistrationInput{
					Facility: "Test Facility",
					DateOfBirth: scalarutils.Date{
						Year:  1990,
						Month: 3,
						Day:   12,
					},
					ClientName:  gofakeit.FirstName(),
					Gender:      enumutils.GenderFemale,
					PhoneNumber: "+254700000000",
					CCCNumber:   "5432",
					Counselled:  true,
					EnrollmentDate: scalarutils.Date{
						Year:  1990,
						Month: 3,
						Day:   12,
					},
					ClientType: enums.ClientTypeDreams,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority)

			if tt.name == "Sad Case - User not authorized" {
				fakeAuthority.MockCheckUserPermissionFn = func(ctx context.Context, permission enums.PermissionType) error {
					return fmt.Errorf("user not authorized")
				}
			}

			got, err := us.RegisterClient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got nil")
				return
			}
		})
	}
}
