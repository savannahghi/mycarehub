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
	clinicalMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical/mock"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
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
		ctx   context.Context
		input *dto.LoginInput
	}
	tests := []struct {
		name  string
		args  args
		want1 bool
	}{
		{
			name: "Happy case: consumer login",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: true,
		},
		{
			name: "Happy case: Login pro",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: true,
		},
		{
			name: "Sad Case - Unable to create getstream token",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: true, // a user should still be able to log in
		},
		{
			name: "Happy Case - should not fail when CCC number is not found",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want1: true,
		},
		{
			name: "Sad Case - Unable to create getstream user PRO",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: true, // a user should still be able to log in
		},
		{
			name: "Sad Case - Unable to create getstream user CONSUMER",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want1: true, // a user should still be able to log in
		},
		{
			name: "Sad case - fail to get user profile by phonenumber",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad case - unable to get user PIN By User ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad case - pin mismatch",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to create firebase token",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to authenticate token",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to update successful login time",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get client profile by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     flavour,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get staff profile by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get user roles by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get user permissions by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get chv user profile by chv user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - failed to check if client has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - client has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - failed to check if staff has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - staff has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					PhoneNumber: &phoneNumber,
					PIN:         &PIN,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "Sad case - fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad Case - Fail to update successful login time" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user")
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

			if tt.name == "Sad Case - failed to check if staff has pending pin reset request" {
				fakeDB.MockCheckIfStaffHasUnresolvedServiceRequestsFn = func(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
					return false, fmt.Errorf("failed to check if staff has pending pin reset request")
				}
			}

			if tt.name == "Sad Case - staff has pending pin reset request" {
				fakeDB.MockCheckIfStaffHasUnresolvedServiceRequestsFn = func(ctx context.Context, staffID string, serviceRequestType string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Happy Case - should not fail when CCC number is not found" {
				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc number identifier value")
				}
			}

			if tt.name == "Sad Case - Unable to create getstream user PRO" {
				fakeGetStream.MockCreateGetStreamUserFn = func(ctx context.Context, user *stream_chat.User) (*stream_chat.UpsertUserResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Unable to create getstream user CONSUMER" {
				fakeGetStream.MockCreateGetStreamUserFn = func(ctx context.Context, user *stream_chat.User) (*stream_chat.UpsertUserResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, got1 := u.Login(tt.args.ctx, tt.args.input)
			if got1 != tt.want1 {
				t.Errorf("UseCasesUserImpl.Login() got1 = %v, want %v", got1, tt.want1)
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
		reinvite    bool
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
				reinvite:    false,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Happy case - Send invite via twilio",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
				reinvite:    true,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case - Fail to Send invite via twilio",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validPhone,
				flavour:     validFlavour,
				reinvite:    true,
			},
			wantErr: true,
			want:    false,
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
		{
			name: "Sad Case - Fail to update user",
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "valid: valid phone number" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Happy case - Send invite via twilio" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case - Fail to Send invite via twilio" {
				fakeExtension.MockSendSMSViaTwilioFn = func(ctx context.Context, phonenumber, message string) error {
					return fmt.Errorf("failed to send sms")
				}
			}

			if tt.name == "valid: valid phone number without country code prefix" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "valid: valid flavour" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return false, fmt.Errorf("phone number is invalid")
				}
			}

			if tt.name == "valid: valid phone number without country code prefix" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "invalid: invalid flavour" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
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
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "invalid: get invite link error" {
				fakeUserMock.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
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

			if tt.name == "Sad Case - Fail to update user" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user")
				}
			}

			got, err := us.InviteUser(tt.args.ctx, tt.args.userID, tt.args.phoneNumber, tt.args.flavour, tt.args.reinvite)
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

			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			name: "Happy case: Successfully set nickname",
			args: args{
				ctx:      ctx,
				userID:   userID,
				nickname: gofakeit.Username(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: unable to set nickname",
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
		{
			name: "Sad Case: failed to update user profile",
			args: args{
				ctx:      ctx,
				userID:   userID,
				nickname: gofakeit.Username(),
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "Happy case: Successfully set nickname" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, nil
				}

				fakeDB.MockSetNickNameFn = func(ctx context.Context, userID, nickname *string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: unable to set nickname" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, nil
				}
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
			if tt.name == "Sad Case: failed to update user profile" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, nil
				}
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
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
			name: "Sad Case - Fail to generate and send OTP",
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			if tt.name == "Sad Case - Fail to generate and send OTP" {
				fakeOTP.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("failed to generate and send OTP")
				}
			}

			if tt.name == "Sad Case - Fail to save otp" {
				fakeOTP.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
					return "111222", nil
				}
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeUser := mock.NewUserUseCaseMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
				flavour: "invalid flavour",
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			name: "valid: Get client caregiver",
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
			name: "invalid: failed to get client caregiver",
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
		{
			name: "invalid: failed to get client caregiver",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "valid: no caregiver assigned" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					ID := uuid.New().String()
					return &domain.ClientProfile{
						ID: &ID,
					}, nil
				}
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
			if tt.name == "invalid: failed to get client caregiver" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					ID := uuid.New().String()
					return &domain.ClientProfile{
						ID:                      &ID,
						UserID:                  uuid.New().String(),
						TreatmentEnrollmentDate: &time.Time{},
						CaregiverID:             &ID,
					}, nil
				}

				fakeDB.MockGetClientCaregiverFn = func(ctx context.Context, clientID string) (*domain.Caregiver, error) {
					return nil, fmt.Errorf("failed to get client caregiver")
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						CaregiverID: nil,
					}, nil
				}
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeUser := mock.NewUserUseCaseMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

	payload := &dto.ClientRegistrationInput{
		Facility:    "123456789",
		ClientTypes: []enums.ClientType{"PMTCT"},
		ClientName:  gofakeit.BeerName(),
		Gender:      enumutils.GenderMale,
		DateOfBirth: scalarutils.Date{
			Year:  2000,
			Month: 01,
			Day:   02,
		},
		PhoneNumber: gofakeit.PhoneFormatted(),
		EnrollmentDate: scalarutils.Date{
			Year:  2000,
			Month: 01,
			Day:   02,
		},
		CCCNumber:    "123456789",
		Counselled:   true,
		InviteClient: true,
	}

	ID := uuid.New().String()
	phone := &domain.Contact{
		ID:           &ID,
		ContactType:  "PHONE",
		ContactValue: interserviceclient.TestUserPhoneNumber,
		Active:       true,
		OptedIn:      false,
		UserID:       &ID,
		Flavour:      feedlib.FlavourConsumer,
	}
	ccc := domain.Identifier{
		ID:                  "123456789",
		IdentifierType:      "CCC",
		IdentifierValue:     payload.CCCNumber,
		IdentifierUse:       "OFFICIAL",
		Description:         "CCC Number, Primary Identifier",
		IsPrimaryIdentifier: true,
	}
	facility := &domain.Facility{
		ID:                 &ID,
		Name:               gofakeit.Name(),
		Code:               20,
		Phone:              interserviceclient.TestUserPhoneNumber,
		Active:             true,
		County:             gofakeit.Name(),
		Description:        gofakeit.BeerAlcohol(),
		FHIROrganisationID: ID,
	}

	userProfile := &domain.User{
		ID:               &ID,
		Username:         gofakeit.Name(),
		Name:             gofakeit.Name(),
		Active:           true,
		TermsAccepted:    true,
		Gender:           enumutils.GenderMale,
		FailedLoginCount: 3,
		Contacts:         phone,
	}

	clientProfile := &domain.ClientProfile{
		ID:                      &ID,
		User:                    userProfile,
		Active:                  false,
		ClientTypes:             []enums.ClientType{},
		UserID:                  ID,
		TreatmentEnrollmentDate: &time.Time{},
		FHIRPatientID:           &ID,
		HealthRecordID:          &ID,
		TreatmentBuddy:          "",
		ClientCounselled:        true,
		OrganisationID:          ID,
		FacilityID:              ID,
		FacilityName:            facility.Name,
		CHVUserID:               &ID,
		CHVUserName:             "name",
		CaregiverID:             &ID,
		CCCNumber:               "123456789",
	}

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
			name: "Happy case: successfully register client",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register client",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to check that facility exists",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to retrieve facility by mfl code",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create user",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create client",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to invite user",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create patient via pubsub",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to check if phone exists",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to check identifier exists",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail if identifier exists",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail if phone exists",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: fail if facility exists",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case: successfully register client" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "Sad case: unable to create user" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return nil, errors.New("error")
				}
			}
			if tt.name == "Sad case: unable to register client" {
				fakeDB.MockRegisterClientFn = func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: unable to check that facility exists" {
				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &ccc, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("unable to check that facility exists")
				}
			}
			if tt.name == "Sad case: unable to retrieve facility by mfl code" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("unable to retrieve facility by mfl code")
				}
			}
			if tt.name == "Sad case: unable to create client" {
				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return phone, nil
				}
				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &ccc, nil
				}
				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("unable to create client")
				}
			}

			if tt.name == "Sad case: unable to invite user" {
				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return phone, nil
				}
				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &ccc, nil
				}
				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					return clientProfile, nil
				}
				fakeUser.MockInviteUserFn = func(ctx context.Context, userID, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return false, fmt.Errorf("unable to invite user")
				}
			}
			if tt.name == "Sad case: unable to create patient via pubsub" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}

				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return phone, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &ccc, nil
				}

				fakeDB.MockRegisterClientFn = func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
					return clientProfile, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("error notifying patient creation topic")
				}
			}
			if tt.name == "Sad case: unable to check if phone exists" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("unable to check if phone exists")
				}
			}
			if tt.name == "Sad case: unable to check identifier exists" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("unable to check identifier exists")
				}
			}
			if tt.name == "Sad case: fail if identifier exists" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad case: fail if phone exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad case: fail if facility exists" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
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
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()

	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			ClientType:  enums.ClientTypeKenyaEMR,
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
			name: "Sad case: cannot retrieve facility",
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
			name: "Happy case: identifier already exists",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Sad case: cannot create user",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot normalize phone number",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot get or create contact",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot create identifier",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot create client",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: cannot publish to patient topic",
			args: args{
				ctx:   ctx,
				input: input,
			},
			want:    nil,
			wantErr: false,
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

			if tt.name == "Sad case: cannot retrieve facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "Happy case: identifier already exists" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: cannot create user" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return nil, fmt.Errorf("cannot create user")
				}
			}

			if tt.name == "Sad case: cannot get or create contact" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return nil, fmt.Errorf("cannot get or create contact")
				}
			}

			if tt.name == "Sad case: cannot create identifier" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return nil, fmt.Errorf("cannot create identifier")
				}

			}

			if tt.name == "Sad case: cannot create client" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &domain.Identifier{ID: gofakeit.UUID()}, nil
				}

				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("cannot create client")
				}

			}

			if tt.name == "Sad case: cannot publish to patient topic" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &domain.Identifier{ID: gofakeit.UUID()}, nil
				}

				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{ID: &id}, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("error notifying patient creation topic")
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateNextOfKin = func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
					return nil
				}

			}

			if tt.name == "Sad case: cannot create next of kin contact" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &domain.Identifier{ID: gofakeit.UUID()}, nil
				}

				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{ID: &id}, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return nil, fmt.Errorf("error creating contact")
				}
			}

			if tt.name == "Sad case: cannot create next of kin" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &domain.Identifier{ID: gofakeit.UUID()}, nil
				}

				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{ID: &id}, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return nil
				}

				fakeDB.MockGetOrCreateNextOfKin = func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
					return fmt.Errorf("cannot create the next of kin")
				}
			}

			if tt.name == "Happy case: successfully create a client and next of kin" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					id := gofakeit.UUID()
					return &domain.User{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockCreateIdentifierFn = func(ctx context.Context, identifier domain.Identifier) (*domain.Identifier, error) {
					return &domain.Identifier{ID: gofakeit.UUID()}, nil
				}

				fakeDB.MockCreateClientFn = func(ctx context.Context, client domain.ClientProfile, contactID, identifierID string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{ID: &id}, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					id := gofakeit.UUID()
					return &domain.Contact{ID: &id}, nil
				}

				fakeDB.MockGetOrCreateNextOfKin = func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
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
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
			want: nil,
			// Shouldnt throw an error, it will accumulate errors and report them to sentry
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
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

func TestUseCasesUserImpl_SearchStaffByStaffNumber(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

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
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockSearchStaffProfileFn = func(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no staffID" {
				fakeDB.MockSearchStaffProfileFn = func(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := us.SearchStaffUser(tt.args.ctx, tt.args.staffNumber)
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
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

	type args struct {
		ctx             context.Context
		searchParameter string
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
				ctx:             ctx,
				searchParameter: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:             ctx,
				searchParameter: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty CCC number",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockSearchClientProfileFn = func(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty CCC number" {
				fakeDB.MockSearchClientProfileFn = func(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := us.SearchClientUser(tt.args.ctx, tt.args.searchParameter)
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
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully withdraw consent",
			args: args{
				ctx:         ctx,
				phoneNumber: gofakeit.Phone(),
				flavour:     feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to purge user details",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "Sad Case - Fail to purge user details" {
				fakeDB.MockDeleteUserFn = func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
					return fmt.Errorf("failed to purge user details")
				}
			}

			got, err := us.Consent(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
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

func TestUseCasesUserImpl_RegisterPushToken(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register a push token",
			args: args{
				ctx:   ctx,
				token: "valid token",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - invalid token length",
			args: args{
				ctx:   ctx,
				token: "123",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get logged in user",
			args: args{
				ctx:   ctx,
				token: "valid token",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user profile",
			args: args{
				ctx:   ctx,
				token: "valid token",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "Sad Case - Fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "Sad Case - Fail to update user profile" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user profile")
				}
			}

			got, err := us.RegisterPushToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterPushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.RegisterPushToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_GetClientProfileByCCCNumber(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

	type args struct {
		ctx       context.Context
		cccNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get client profile by ccc number",
			args: args{
				ctx:       context.Background(),
				cccNumber: "123456789",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get client profile by ccc number",
			args: args{
				ctx:       context.Background(),
				cccNumber: "123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad Case - Fail to get client profile by ccc number" {
				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, cccNumber string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by ccc number")
				}
			}

			got, err := us.GetClientProfileByCCCNumber(tt.args.ctx, tt.args.cccNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetClientProfileByCCCNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesUserImpl.GetClientProfileByCCCNumber() = %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_DeleteUser(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		payload *dto.PhoneInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete client",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully delete staff",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
					Flavour:     feedlib.FlavourPro,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable to get user profile by phone number",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get client profile",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to delete user",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get staff profile",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourPro,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to delete staff user",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourPro,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to delete getstream user",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourPro,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to delete getstream user - Consumer",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to delete FHIR patient profile",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, otp, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

			if tt.name == "Happy Case - Successfully delete client" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					input := dto.PhoneInput{
						PhoneNumber: interserviceclient.TestUserPhoneNumber,
					}

					payload, err := json.Marshal(input)
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

			if tt.name == "Sad Case - unable to get user profile by phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad Case - unable to get client profile" {
				fakeDB.MockGetClientProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad Case - unable to delete user" {
				fakeDB.MockDeleteUserFn = func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
					return fmt.Errorf("failed to delete user")
				}
			}

			if tt.name == "Sad Case - unable to get staff profile" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}

			if tt.name == "Sad Case - unable to delete staff user" {
				fakeDB.MockDeleteUserFn = func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
					return fmt.Errorf("failed to delete user")
				}
			}

			if tt.name == "Sad Case - unable to delete getstream user" {
				fakeGetStream.MockDeleteUsersFn = func(ctx context.Context, userIDs []string, options stream_chat.DeleteUserOptions) (*stream_chat.AsyncTaskResponse, error) {
					return nil, fmt.Errorf("failed to delete getstream user")
				}
			}

			if tt.name == "Sad Case - unable to delete getstream user - Consumer" {
				fakeGetStream.MockDeleteUsersFn = func(ctx context.Context, userIDs []string, options stream_chat.DeleteUserOptions) (*stream_chat.AsyncTaskResponse, error) {
					return nil, fmt.Errorf("failed to delete getstream user")
				}
			}

			if tt.name == "Sad Case - unable to delete FHIR patient profile" {
				fakeClinical.MockDeleteFHIRPatientByPhoneFn = func(ctx context.Context, phoneNumber string) error {
					return fmt.Errorf("failed to delete FHIR patient profile")
				}
			}

			got, err := us.DeleteUser(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_DeleteStreamUser(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete getstream user",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Unable delete getstream user",
			args: args{
				ctx: ctx,
				id:  "",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - Unable delete getstream user" {
				fakeGetStream.MockDeleteUsersFn = func(ctx context.Context, userIDs []string, options stream_chat.DeleteUserOptions) (*stream_chat.AsyncTaskResponse, error) {
					return nil, fmt.Errorf("failed to delete getstream user")
				}
			}

			err := us.DeleteStreamUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.DeleteStreamUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_TransferClientToFacility(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

	type args struct {
		ctx        context.Context
		clientID   *string
		facilityID *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully transfer client to facility",
			args: args{
				ctx:        ctx,
				clientID:   &ID,
				facilityID: &ID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Missing client ID or facility ID",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to get client profile by clientID",
			args: args{
				ctx:        ctx,
				clientID:   &ID,
				facilityID: &ID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to update client",
			args: args{
				ctx:        ctx,
				clientID:   &ID,
				facilityID: &ID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to get client service requests",
			args: args{
				ctx:        ctx,
				clientID:   &ID,
				facilityID: &ID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to update client service requests",
			args: args{
				ctx:        ctx,
				clientID:   &ID,
				facilityID: &ID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad Case - Failed to get client profile by clientID" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by clientID")
				}
			}
			if tt.name == "Sad Case - Failed to update client" {
				// get the client profile
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
					}, nil
				}
				fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to update client")
				}
			}

			if tt.name == "Sad Case - Failed to get client service requests" {
				// get the client profile
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
					}, nil
				}
				// update the client profile
				fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
					}, nil
				}

				fakeDB.MockGetClientServiceRequestsFn = func(ctx context.Context, requestType string, status string, clientID string, facilityID string) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get client service requests")
				}
			}

			if tt.name == "Sad Case - Failed to update client service requests" {
				// get the client profile
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
					}, nil
				}
				// update the client profile
				fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
					}, nil
				}
				// get Service Requests
				fakeDB.MockGetClientServiceRequestsFn = func(ctx context.Context, requestType string, status string, clientID string, facilityID string) ([]*domain.ServiceRequest, error) {
					return []*domain.ServiceRequest{{ID: ID}}, nil
				}

				fakeDB.MockUpdateClientServiceRequestFn = func(ctx context.Context, clientServiceRequest *domain.ServiceRequest, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update client service requests")
				}
			}

			got, err := us.TransferClientToFacility(tt.args.ctx, tt.args.clientID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.TransferClientToFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.TransferClientToFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterStaff(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeUser := mock.NewUserUseCaseMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical)

	ID := "123"

	payload := &dto.StaffRegistrationInput{
		Facility:  "1234",
		StaffName: gofakeit.BeerName(),
		Gender:    enumutils.GenderMale,
		DateOfBirth: scalarutils.Date{
			Year:  2000,
			Month: 2,
			Day:   20,
		},
		PhoneNumber: interserviceclient.TestUserPhoneNumber,
		IDNumber:    "123456789",
		StaffNumber: "123456789",
		StaffRoles:  "Community Management",
		InviteStaff: true,
	}

	phone := &domain.Contact{
		ID:           &ID,
		ContactType:  "PHONE",
		ContactValue: interserviceclient.TestUserPhoneNumber,
		Active:       true,
		OptedIn:      false,
		UserID:       &ID,
		Flavour:      feedlib.FlavourConsumer,
	}

	userProfile := &domain.User{
		ID:               &ID,
		Username:         gofakeit.Name(),
		Name:             gofakeit.Name(),
		Active:           true,
		TermsAccepted:    true,
		Gender:           enumutils.GenderMale,
		FailedLoginCount: 3,
		Contacts:         phone,
	}

	facility := &domain.Facility{
		ID:                 &ID,
		Name:               gofakeit.Name(),
		Code:               20,
		Phone:              interserviceclient.TestUserPhoneNumber,
		Active:             true,
		County:             gofakeit.Name(),
		Description:        gofakeit.BeerAlcohol(),
		FHIROrganisationID: ID,
	}

	staffProfile := &domain.StaffProfile{
		ID:                  &ID,
		User:                userProfile,
		UserID:              *userProfile.ID,
		Active:              true,
		StaffNumber:         "1234",
		DefaultFacilityID:   *facility.ID,
		DefaultFacilityName: facility.Name,
	}

	type args struct {
		ctx   context.Context
		input dto.StaffRegistrationInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.StaffRegistrationOutput
		wantErr bool
	}{
		{
			name: "Sad Case - Unable to check identifier exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - identifier exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to check phone exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - phone exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Register Staff",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to create user",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to check facility exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - raise error if facility does not exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to retrieve facility by MFL Code",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to register staff",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to assign roles",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to invite staff",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - Unable to check identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("failed to check identifier exists")
				}
			}
			if tt.name == "Sad Case - identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad Case - Unable to check phone exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to check phone exists")
				}
			}
			if tt.name == "Sad Case - phone exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Happy Case - Register Staff" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRegisterStaffFn = func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return staffProfile, nil
				}
				fakeDB.MockAssignRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad Case - unable to create user" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return nil, fmt.Errorf("failed to create user")
				}
			}
			if tt.name == "Sad Case - unable to check facility exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("failed to check facility exists")
				}
			}
			if tt.name == "Sad Case - raise error if facility does not exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "Sad Case - unable to retrieve facility by MFL Code" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility by MFL Code")
				}
			}
			if tt.name == "Sad Case - unable to register staff" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return facility, nil
				}
				fakeDB.MockRegisterStaffFn = func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to register staff")
				}
			}
			if tt.name == "Sad Case - unable to assign roles" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return facility, nil
				}
				fakeDB.MockRegisterStaffFn = func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return staffProfile, nil
				}
				fakeDB.MockAssignRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return false, fmt.Errorf("failed to assign roles")
				}
			}
			if tt.name == "Sad Case - unable to invite staff" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, nil
				}
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
				fakeDB.MockCreateUserFn = func(ctx context.Context, user domain.User) (*domain.User, error) {
					return userProfile, nil
				}
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return facility, nil
				}
				fakeDB.MockRegisterStaffFn = func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return staffProfile, nil
				}
				fakeDB.MockAssignRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return true, nil
				}
				fakeUser.MockInviteUserFn = func(ctx context.Context, userID, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return false, fmt.Errorf("failed to invite user")
				}
			}
			_, err := us.RegisterStaff(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
