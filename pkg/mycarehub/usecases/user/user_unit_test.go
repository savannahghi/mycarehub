package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
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
		{
			name: "Sad Case - Fail to get user roles by user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user permissions by user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get chv user profile by chv user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to create getstream user",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to create getstream token",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - failed to check if client has pending pin reset request",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - client has pending pin reset request",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get client ccc number identifier value",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				pin:         PIN,
				flavour:     feedlib.FlavourConsumer,
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
				fakeUserMock.MockLoginFn = func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, error) {
					return nil, fmt.Errorf("invalid flavour defined")
				}
			}

			if tt.name == "Sad Case - Fail to update successful login time" {
				fakeDB.MockUpdateUserProfileAfterLoginSuccessFn = func(ctx context.Context, userID string) error {
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
			if tt.name == "Sad Case - Fail to get user roles by user ID" {
				fakeAuthority.MockGetUserRolesFn = func(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get user role")
				}
			}
			if tt.name == "Sad Case - Fail to get user permissions by user ID" {
				fakeAuthority.MockGetUserPermissionsFn = func(ctx context.Context, userID string) ([]*domain.AuthorityPermission, error) {
					return nil, fmt.Errorf("failed to get user permission")
				}
			}
			if tt.name == "Sad Case - Fail to get chv user profile by chv user ID" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					userID := uuid.NewString()
					return &domain.User{
						ID:      &userID,
						Flavour: feedlib.FlavourConsumer,
					}, nil
				}

				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get chv profile")
				}
			}
			if tt.name == "Sad Case - Unable to create getstream user" {
				fakeGetStream.MockCreateGetStreamUserFn = func(ctx context.Context, user *stream_chat.User) (*stream_chat.UpsertUserResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - Unable to create getstream token" {
				fakeGetStream.MockCreateGetStreamUserTokenFn = func(ctx context.Context, userID string) (string, error) {
					return "", fmt.Errorf("failed to create getstream token")
				}
			}

			if tt.name == "Sad Case - failed to check if client has pending pin reset request" {
				fakeDB.MockCheckIfClientHasUnresolvedServiceRequestsFn = func(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
					return false, fmt.Errorf("failed to check if client has pending pin reset request")
				}
			}

			if tt.name == "Sad Case - client has pending pin reset request" {
				fakeDB.MockCheckIfClientHasUnresolvedServiceRequestsFn = func(ctx context.Context, clientID string, serviceRequestType string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad Case - unable to get client ccc number identifier value" {
				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc number identifier value")
				}
			}

			_, err := u.Login(tt.args.ctx, tt.args.phoneNumber, tt.args.pin, tt.args.flavour)
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
		{
			name: "Sad Case - Fail to save pin update required",
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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

			if tt.name == "Sad Case - Fail to save pin update required" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to save pin")
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
	UUID := uuid.New().String()
	now := time.Now()
	past := now.Add(-time.Hour * 1000000)
	userProfile := &domain.User{
		ID:                     &UUID,
		Username:               gofakeit.Name(),
		UserType:               enums.ClientUser,
		Name:                   gofakeit.Name(),
		Gender:                 enumutils.GenderMale,
		Active:                 true,
		LastSuccessfulLogin:    &now,
		LastFailedLogin:        &now,
		FailedLoginCount:       0,
		NextAllowedLogin:       &now,
		PinChangeRequired:      false,
		HasSetPin:              false,
		HasSetSecurityQuestion: false,
		IsPhoneVerified:        false,
		TermsAccepted:          false,
		AcceptedTermsID:        0,
		Flavour:                "CONSUMER",
		Suspended:              false,
		Avatar:                 "",
		DateOfBirth:            &past,
	}

	type args struct {
		ctx         context.Context
		userProfile *domain.User
		pin         string
		flavour     feedlib.Flavour
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
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user pin",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "3456",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Pin mismatch",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user login count",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update last failed login time",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update next allowed login time",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid flavour",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     "Invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get the updated user profile after updating last failed login time and login count",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - failed to get the updated user profile after updating next allowed login time",
			args: args{
				ctx:         ctx,
				userProfile: userProfile,
				pin:         "1234",
				flavour:     feedlib.FlavourConsumer,
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Happy Case - Successfully verify pin" {
				fakeUserMock.MockVerifyLoginPINFn = func(ctx context.Context, userProfile *domain.User, pin string, flavour feedlib.Flavour) (bool, error) {
					return true, nil
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

			if tt.name == "Sad Case - failed to get the updated user profile after updating last failed login time and login count" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user ID")
				}
			}

			if tt.name == "Sad Case - failed to get the updated user profile after updating next allowed login time" {
				fakeExtension.MockComparePINFn = func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
					return false
				}
				fakeDB.MockUpdateUserNextAllowedLoginTimeFn = func(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
					return nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user ID")
				}
			}

			got, err := u.VerifyLoginPIN(tt.args.ctx, tt.args.userProfile, tt.args.pin, tt.args.flavour)
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
		userID   string
		nickname string
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
				userID:   userID,
				nickname: gofakeit.Username(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				userID:   userID,
				nickname: nickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:      ctx,
				userID:   "",
				nickname: nickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no nickname",
			args: args{
				ctx:      ctx,
				userID:   userID,
				nickname: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Both userID and nickname nil",
			args: args{
				ctx:      ctx,
				userID:   "",
				nickname: "",
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Happy case" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, nil
				}

				fakeDB.MockSetNickNameFn = func(ctx context.Context, userID, nickname *string) (bool, error) {
					return true, nil
				}
			}

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
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
		{
			name: "invalid: failed update pin update required status",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUser := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
				fakeUser.MockResetPINFn = func(ctx context.Context, input dto.UserResetPinInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "invalid: failed update pin update required status" {
				fakeDB.MockUpdateUserPinUpdateRequiredStatusFn = func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
					return fmt.Errorf("failed to update pin update required status")
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Sad case - no userID" {
				fakeDB.MockCompleteOnboardingTourFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case - no userID and flavour" {
				fakeDB.MockCompleteOnboardingTourFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
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
					ClientType:   enums.ClientTypeDreams,
					InviteClient: true,
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
		// TODO: Restore after aligning with frontend
		// {
		// 	name: "Sad Case - User not authorized",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		input: &dto.ClientRegistrationInput{
		// 			Facility: "Test Facility",
		// 			DateOfBirth: scalarutils.Date{
		// 				Year:  1990,
		// 				Month: 3,
		// 				Day:   12,
		// 			},
		// 			ClientName:  gofakeit.FirstName(),
		// 			Gender:      enumutils.GenderFemale,
		// 			PhoneNumber: "+254700000000",
		// 			CCCNumber:   "5432",
		// 			Counselled:  true,
		// 			EnrollmentDate: scalarutils.Date{
		// 				Year:  1990,
		// 				Month: 3,
		// 				Day:   12,
		// 			},
		// 			ClientType: enums.ClientTypeDreams,
		// 		},
		// 	},
		// 	wantErr: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Happy Case - Successfully register client" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					registrationOutput := dto.ClientRegistrationOutput{
						ID: uuid.New().String(),
					}

					payload, err := json.Marshal(registrationOutput)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to make request" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

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

func TestUseCasesUserImpl_AddClientFHIRID(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

	type args struct {
		ctx   context.Context
		input dto.ClientFHIRPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: success updating client profile fhir id",
			args: args{
				ctx: context.Background(),
				input: dto.ClientFHIRPayload{
					ClientID: gofakeit.UUID(),
					FHIRID:   gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error retrieving client profile",
			args: args{
				ctx: context.Background(),
				input: dto.ClientFHIRPayload{
					ClientID: gofakeit.UUID(),
					FHIRID:   gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error updating client profile",
			args: args{
				ctx: context.Background(),
				input: dto.ClientFHIRPayload{
					ClientID: gofakeit.UUID(),
					FHIRID:   gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.name == "sad case: error retrieving client profile" {
			fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
				return nil, fmt.Errorf("error retrieving client profile")
			}
		}

		if tt.name == "sad case: error updating client profile" {
			fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
				return nil, fmt.Errorf("error updating client profile")
			}

		}
		t.Run(tt.name, func(t *testing.T) {
			if err := us.AddClientFHIRID(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.AddClientFHIRID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesUserImpl_GetUserProfile(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "happy case: get user profile",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := us.GetUserProfile(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected user profile to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected user profile not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_RefreshGetStreamToken(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully refresh token",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to refresh token",
			args: args{
				ctx: context.Background(),
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Sad Case - Fail to refresh token" {
				fakeGetStream.MockCreateGetStreamUserTokenFn = func(ctx context.Context, userID string) (string, error) {
					return "", fmt.Errorf("failed to generate token")
				}
			}

			got, err := us.RefreshGetStreamToken(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RefreshGetStreamToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Errorf("expected a response but got %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterKenyaEMRPatients(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

	ctx := context.Background()
	input := []*dto.PatientRegistrationPayload{
		{
			MFLCode:   "1234",
			CCCNumber: "1234",
			Name:      "Jane Doe",
			DateOfBirth: scalarutils.Date{
				Year:  2000,
				Month: 12,
				Day:   25,
			},
			ClientType:  "KenyaEMR",
			PhoneNumber: gofakeit.Phone(),
			EnrollmentDate: scalarutils.Date{
				Year:  2020,
				Month: 12,
				Day:   25,
			},
			BirthDateEstimated: false,
			Gender:             "male",
			Counselled:         false,
			NextOfKin: dto.NextOfKinPayload{
				Name:         "John Doe",
				Contact:      gofakeit.Phone(),
				Relationship: "spouse",
			},
		},
	}

	type args struct {
		ctx   context.Context
		input []*dto.PatientRegistrationPayload
	}
	tests := []struct {
		name    string
		args    args
		want    []*dto.ClientRegistrationOutput
		wantErr bool
	}{
		{
			name: "Sad case: check facility error",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: facility doesn't exist",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: check identifier exists error",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: identifier already exists",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot register client",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot create next of kin contact",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot create next of kin",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Happy case: successfully create a client and next of kin",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: check facility error" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("error fetching facility")
				}
			}

			if tt.name == "Sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case: check identifier exists error" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return &domain.Facility{}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("error checking for identifier")
				}
			}

			if tt.name == "Sad case: identifier already exists" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return &domain.Facility{}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: cannot register client" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return &domain.Facility{}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("create client request fails")
				}
			}

			if tt.name == "Sad case: cannot create next of kin contact" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return &domain.Facility{}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					registrationOutput := dto.ClientRegistrationOutput{
						ID: uuid.New().String(),
					}

					payload, err := json.Marshal(registrationOutput)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}

				fakeDB.MockCreateContact = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return nil, fmt.Errorf("error creating contact")
				}
			}

			if tt.name == "Sad case: cannot create next of kin" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return &domain.Facility{}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateContact = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					registrationOutput := dto.ClientRegistrationOutput{
						ID: uuid.New().String(),
					}

					payload, err := json.Marshal(registrationOutput)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}

				fakeDB.MockCreateNextOfKin = func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
					return fmt.Errorf("cannot create the next of kin")
				}
			}

			if tt.name == "Happy case: successfully create a client and next of kin" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return &domain.Facility{}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					registrationOutput := dto.ClientRegistrationOutput{
						ID: uuid.New().String(),
					}

					payload, err := json.Marshal(registrationOutput)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}

				fakeDB.MockCreateContact = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateNextOfKin = func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
					return nil
				}
			}

			_, err := us.RegisterKenyaEMRPatients(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterKenyaEMRPatients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestUseCasesUserImpl_RegisteredFacilityPatients(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

	ctx := context.Background()
	syncTime := time.Now()

	type args struct {
		ctx   context.Context
		input dto.PatientSyncPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.PatientSyncResponse
		wantErr bool
	}{
		{
			name: "sad case: error checking facility",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: &syncTime,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: facility doesn't exist",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: &syncTime,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: error retrieving facility",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: &syncTime,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: error retrieving clients with sync time",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: &syncTime,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: error retrieving clients without sync time",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: error retrieving client identifier",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case: success retrieving new clients without sync time",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: nil,
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "happy case: success retrieving new clients with sync time",
			args: args{
				ctx: ctx,
				input: dto.PatientSyncPayload{
					MFLCode:  0,
					SyncTime: &syncTime,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: error checking facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("error fetching facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility")
				}
			}

			if tt.name == "sad case: error retrieving clients with sync time" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("cannot retrieve clients")
				}
			}

			if tt.name == "sad case: error retrieving clients without sync time" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("cannot retrieve clients")
				}
			}

			if tt.name == "sad case: error retrieving client identifier" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					id := uuid.NewString()
					return []*domain.ClientProfile{
						{
							ID: &id,
						},
					}, nil
				}

				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return nil, fmt.Errorf("cannot get ccc identifier")
				}
			}

			if tt.name == "happy case: success retrieving new clients with sync time" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					id := uuid.NewString()
					return []*domain.ClientProfile{
						{
							ID: &id,
						},
					}, nil
				}

				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return &domain.Identifier{IdentifierValue: "123456"}, nil
				}
			}

			if tt.name == "happy case: success retrieving new clients without sync time" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					id := uuid.NewString()
					return []*domain.ClientProfile{
						{
							ID: &id,
						},
					}, nil
				}

				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return &domain.Identifier{IdentifierValue: "123456"}, nil
				}
			}

			got, err := us.RegisteredFacilityPatients(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisteredFacilityPatients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected community to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected community not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterStaff(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.StaffRegistrationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register staff",
			args: args{
				ctx: context.Background(),
				input: dto.StaffRegistrationInput{
					Facility:  "Test Facility",
					StaffName: gofakeit.Name(),
					Gender:    enumutils.GenderFemale,
					DateOfBirth: scalarutils.Date{
						Year:  1990,
						Month: 3,
						Day:   12,
					},
					PhoneNumber: "+254700000000",
					IDNumber:    "1234567890",
					StaffNumber: "MS-01",
					StaffRoles:  "CONTENT_MANAGEMENT",
					InviteStaff: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - invalid ID number",
			args: args{
				ctx: context.Background(),
				input: dto.StaffRegistrationInput{
					Facility:  "Test Facility",
					StaffName: gofakeit.Name(),
					Gender:    enumutils.GenderFemale,
					DateOfBirth: scalarutils.Date{
						Year:  1990,
						Month: 3,
						Day:   12,
					},
					PhoneNumber: "+254700000000",
					IDNumber:    "s1234567890",
					StaffNumber: "MS-01",
					StaffRoles:  "CONTENT_MANAGEMENT",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to make request",
			args: args{
				ctx: context.Background(),
				input: dto.StaffRegistrationInput{
					Facility: "non existent",
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Happy Case - Successfully register staff" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					registrationOutput := dto.StaffRegistrationOutput{
						ID: uuid.New().String(),
					}

					payload, err := json.Marshal(registrationOutput)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to make request" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make a request")
				}
			}

			got, err := us.RegisterStaff(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got nil")
				return
			}
		})
	}
}

func TestUseCasesUserImpl_SearchStaffByStaffNumber(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

	type args struct {
		ctx         context.Context
		staffNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - no staffID",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockSearchStaffProfileByStaffNumberFn = func(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no staffID" {
				fakeDB.MockSearchStaffProfileByStaffNumberFn = func(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := us.SearchStaffByStaffNumber(tt.args.ctx, tt.args.staffNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SearchStaffByStaffNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected staff profiles to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected staff profiles not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_SearchClientByCCCNumber(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

	type args struct {
		ctx       context.Context
		CCCNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				CCCNumber: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				CCCNumber: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty CCC number",
			args: args{
				ctx:       ctx,
				CCCNumber: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockSearchClientProfilesByCCCNumberFn = func(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty CCC number" {
				fakeDB.MockSearchClientProfilesByCCCNumberFn = func(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := us.SearchClientsByCCCNumber(tt.args.ctx, tt.args.CCCNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetClientByCCCNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got nil")
				return
			}
		})
	}
}

func TestUseCasesUserImpl_Consent(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		phoneNumber string
		flavour     feedlib.Flavour
		active      bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully offer consent",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
				flavour:     feedlib.FlavourConsumer,
				active:      true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully withdraw consent",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
				flavour:     feedlib.FlavourConsumer,
				active:      false,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get user profile by phone",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				flavour:     "",
				active:      true,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user status",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				flavour:     "",
				active:      true,
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub)

			if tt.name == "Sad Case - Fail to get user profile by phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phonenumber")
				}
			}

			if tt.name == "Sad Case - Fail to update user status" {
				fakeDB.MockUpdateUserActiveStatusFn = func(ctx context.Context, userID string, flavour feedlib.Flavour, active bool) error {
					return fmt.Errorf("failed to update user active status")
				}
			}

			got, err := us.Consent(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour, tt.args.active)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.Consent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.Consent() = %v, want %v", got, tt.want)
			}
		})
	}
}
