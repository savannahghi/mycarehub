package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	clinicalMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical/mock"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	smsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms/mock"
	twilioMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/silcomms"
	"github.com/segmentio/ksuid"
	pkgGorm "gorm.io/gorm"
)

func TestUseCasesUserImpl_Login_Unittest(t *testing.T) {
	ctx := context.Background()

	PIN := "1234"

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
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: true,
		},
		{
			name: "Happy case: Login pro",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: true,
		},
		{
			name: "Sad Case - Unable to create getstream token",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: true, // a user should still be able to log in
		},
		{
			name: "Happy Case - should not fail when CCC number is not found",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: true,
		},
		{
			name: "Sad Case - Unable to create getstream user PRO",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: true, // a user should still be able to log in
		},
		{
			name: "Sad Case - Unable to create getstream user CONSUMER",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: true, // a user should still be able to log in
		},
		{
			name: "Sad case - fail to get user profile by phonenumber",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad case - unable to get user PIN By User ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to create firebase token",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to authenticate token",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to update successful login time",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get client profile by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get staff profile by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get user roles by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get user permissions by user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - Fail to get chv user profile by chv user ID",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - failed to check if client has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - client has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourConsumer,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - failed to check if staff has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
				},
			},
			want1: false,
		},
		{
			name: "Sad Case - staff has pending pin reset request",
			args: args{
				ctx: ctx,
				input: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      PIN,
					Flavour:  feedlib.FlavourPro,
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Happy case: consumer login" {
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: time.Now(),
						ValidTo:   laterTime,
						IsValid:   true,
						Salt:      salt,
					}, nil
				}
			}
			if tt.name == "Happy case: Login pro" {
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: time.Now(),
						ValidTo:   laterTime,
						IsValid:   true,
						Salt:      salt,
					}, nil
				}
			}
			if tt.name == "Sad case - fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad Case - Fail to update successful login time" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user")
				}
			}

			if tt.name == "Sad case - unable to get user PIN By User ID" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get user PIN by user ID")
				}
			}

			if tt.name == "Sad Case - Fail to create firebase token" {
				fakeExtension.MockCreateFirebaseCustomTokenWithClaimsFn = func(ctx context.Context, uid string, claims map[string]interface{}) (string, error) {
					return "", fmt.Errorf("failed to create custom token")
				}
			}

			if tt.name == "Sad Case - Fail to authenticate token" {
				fakeExtension.MockAuthenticateCustomFirebaseTokenFn = func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
					return nil, fmt.Errorf("failed to authenticate token")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile by user ID" {
				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad Case - Fail to get staff profile by user ID" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}
			if tt.name == "Sad Case - Fail to get user roles by user ID" {
				fakeAuthority.MockGetUserRolesFn = func(ctx context.Context, userID string, organisationID string) ([]*domain.AuthorityRole, error) {
					return nil, fmt.Errorf("failed to get user role")
				}
			}
			if tt.name == "Sad Case - Fail to get user permissions by user ID" {
				fakeAuthority.MockGetUserPermissionsFn = func(ctx context.Context, userID string, organisationID string) ([]*domain.AuthorityPermission, error) {
					return nil, fmt.Errorf("failed to get user permission")
				}
			}
			if tt.name == "Sad Case - Fail to get chv user profile by chv user ID" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					userID := uuid.NewString()
					return &domain.User{
						ID: &userID,
					}, nil
				}

				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get chv profile")
				}
			}
			if tt.name == "Sad Case - Unable to create getstream token" {
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: time.Now(),
						ValidTo:   laterTime,
						IsValid:   true,
						Salt:      salt,
					}, nil
				}
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
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: time.Now(),
						ValidTo:   laterTime,
						IsValid:   true,
						Salt:      salt,
					}, nil
				}
				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc number identifier value")
				}
			}

			if tt.name == "Sad Case - Unable to create getstream user PRO" {
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: time.Now(),
						ValidTo:   laterTime,
						IsValid:   true,
						Salt:      salt,
					}, nil
				}
				fakeGetStream.MockCreateGetStreamUserFn = func(ctx context.Context, user *stream_chat.User) (*stream_chat.UpsertUserResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Unable to create getstream user CONSUMER" {
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: time.Now(),
						ValidTo:   laterTime,
						IsValid:   true,
						Salt:      salt,
					}, nil
				}
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
				fakeTwilio.MockSendSMSViaTwilioFn = func(ctx context.Context, phonenumber, message string) error {
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
				fakeDB.MockInvalidatePINFn = func(ctx context.Context, userID string) (bool, error) {
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
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return &silcomms.BulkSMSResponse{
						GUID: "123",
					}, nil
				}
			}
			if tt.name == "invalid: send in message error" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
				fakeDB.MockInvalidatePINFn = func(ctx context.Context, userID string) (bool, error) {
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
		ctx      context.Context
		username string
		flavour  feedlib.Flavour
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
				ctx:      ctx,
				username: gofakeit.Name(),
				flavour:  feedlib.FlavourConsumer,
			},
			want:    "111222",
			wantErr: false,
		},
		{
			name: "invalid: invalid flavour",
			args: args{
				ctx:      ctx,
				username: "0710000000",
				flavour:  feedlib.Flavour("Invalid_flavour"),
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by username",
			args: args{
				ctx:      ctx,
				username: gofakeit.Name(),
				flavour:  feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to check if user has pin",
			args: args{
				ctx:      ctx,
				username: gofakeit.Name(),
				flavour:  feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to generate and send OTP",
			args: args{
				ctx:      ctx,
				username: gofakeit.Name(),
				flavour:  feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to save otp",
			args: args{
				ctx:      ctx,
				username: gofakeit.Name(),
				flavour:  feedlib.FlavourConsumer,
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case - Invalid username" {
				fakeUser.MockRequestPINResetFn = func(ctx context.Context, username string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("invalid username")
				}
			}

			if tt.name == "invalid: invalid flavour" {
				fakeUser.MockRequestPINResetFn = func(ctx context.Context, username string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("invalid flavour defined")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by username" {
				fakeDB.MockGetUserProfileByUsernameFn = func(ctx context.Context, username string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad Case - Fail to check if user has pin" {
				fakeDB.MockCheckUserHasPinFn = func(ctx context.Context, userID string) (bool, error) {
					return false, fmt.Errorf("failed to check if user has pin")
				}
			}
			if tt.name == "Sad Case - Fail to generate and send OTP" {
				fakeOTP.MockGenerateAndSendOTPFn = func(ctx context.Context, username string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("failed to generate and send OTP")
				}
			}

			if tt.name == "Sad Case - Fail to save otp" {
				fakeOTP.MockGenerateAndSendOTPFn = func(ctx context.Context, username string, flavour feedlib.Flavour) (string, error) {
					return "111222", nil
				}
				fakeDB.MockSaveOTPFn = func(ctx context.Context, otpInput *domain.OTP) error {
					return fmt.Errorf("failed to save otp")
				}
			}

			got, err := us.RequestPINReset(tt.args.ctx, tt.args.username, tt.args.flavour)
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
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
				fakeDB.MockInvalidatePINFn = func(ctx context.Context, userID string) (bool, error) {
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Happy Case - Successfully verify pin" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					currentTime := time.Now()
					laterTime := currentTime.Add(time.Minute * 2005)
					salt, encryptedPin := utils.EncryptPIN("1234", nil)
					return &domain.UserPIN{
						UserID:    userID,
						HashedPIN: encryptedPin,
						ValidFrom: currentTime,
						ValidTo:   laterTime,
						IsValid:   false,
						Salt:      salt,
					}, nil
				}
			}

			if tt.name == "invalid: failed to get user pin by user id" {
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return nil, fmt.Errorf("failed to get user pin by user id")
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
			name: "Sad case: unable to check if username exists",
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
			name: "Sad case: fail if username exists",
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
		{
			name: "Sad case: unable to publish cms user to pubsub",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get logged in user id",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx:   context.Background(),
				input: payload,
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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: unable to register client" {
				fakeDB.MockRegisterClientFn = func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error")
				}
			}

			if tt.name == "Sad case: unable to check that facility exists" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("unable to check that facility exists")
				}
			}

			if tt.name == "Sad case: unable to retrieve facility by mfl code" {
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("unable to retrieve facility by mfl code")
				}
			}

			if tt.name == "Sad case: unable to invite user" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("failed to send sms")
				}
			}

			if tt.name == "Sad case: unable to create patient via pubsub" {

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("error notifying patient creation topic")
				}
			}

			if tt.name == "Sad case: unable to check if username exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, fmt.Errorf("unable to check if phone exists")
				}
			}

			if tt.name == "Sad case: unable to check identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("unable to check identifier exists")
				}
			}

			if tt.name == "Sad case: fail if identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: fail if username exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: fail if facility exists" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case: unable to publish cms user to pubsub" {

				fakePubsub.MockNotifyCreateCMSUserFn = func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
					return fmt.Errorf("unable to publish cms user to pubsub")
				}
			}

			if tt.name == "Sad case: unable to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user id")
				}
			}

			if tt.name == "Sad case: unable to get user profile by user id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get user profile by user id")
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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

	ctx := context.Background()
	input := []*dto.PatientRegistrationPayload{
		{
			MFLCode:   "1234",
			CCCNumber: "1234",
			Name:      "Jane Doe",
			ProgramID: uuid.NewString(),
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error fetching facility")
				}
			}

			if tt.name == "Sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case: cannot retrieve facility" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "Happy case: identifier already exists" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType, identifierValue string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: cannot create user" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
					return &domain.ClientProfile{
						ID: &id,
						DefaultFacility: &domain.Facility{
							ID: &id,
						},
					}, nil
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
					return &domain.ClientProfile{ID: &id, DefaultFacility: &domain.Facility{
						ID: &id,
					}}, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return nil
				}

				fakeDB.MockGetOrCreateContactFn = func(ctx context.Context, contact *domain.Contact) (*domain.Contact, error) {
					return nil, fmt.Errorf("error creating contact")
				}
			}

			if tt.name == "Sad case: cannot create next of kin" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
					return &domain.ClientProfile{ID: &id, DefaultFacility: &domain.Facility{
						ID: &id,
					}}, nil
				}

				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return nil
				}

				fakeDB.MockGetOrCreateNextOfKin = func(ctx context.Context, person *dto.NextOfKinPayload, clientID, contactID string) error {
					return fmt.Errorf("cannot create the next of kin")
				}
			}

			if tt.name == "Happy case: successfully create a client and next of kin" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
					return &domain.ClientProfile{ID: &id, DefaultFacility: &domain.Facility{
						ID: &id,
					}}, nil
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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error fetching facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility")
				}
			}

			if tt.name == "sad case: error retrieving clients with sync time" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("cannot retrieve clients")
				}
			}

			if tt.name == "sad case: error retrieving clients without sync time" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := uuid.NewString()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientsByParams = func(ctx context.Context, params gorm.Client, lastSyncTime *time.Time) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("cannot retrieve clients")
				}
			}

			if tt.name == "sad case: error retrieving client identifier" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
		{
			name: "Sad Case - unable to delete cms client via pub sub",
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
		{
			name: "Sad Case - unable to delete cms staff via pub sub",
			args: args{
				ctx: ctx,
				payload: &dto.PhoneInput{
					PhoneNumber: "",
					Flavour:     feedlib.FlavourPro,
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
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}

			if tt.name == "Sad Case - unable to get user profile by phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone number")
				}
			}

			if tt.name == "Sad Case - unable to get client profile" {
				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID string, programID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			if tt.name == "Sad Case - unable to delete user" {
				fakeDB.MockDeleteUserFn = func(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
					return fmt.Errorf("failed to delete user")
				}
			}

			if tt.name == "Sad Case - unable to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
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
			if tt.name == "Sad Case - unable to delete cms client via pub sub" {
				fakePubsub.MockNotifyDeleteCMSClientFn = func(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
					return fmt.Errorf("failed to delete cms client")
				}
			}

			if tt.name == "Sad Case - unable to delete cms staff via pub sub" {
				fakePubsub.MockNotifyDeleteCMSStaffFn = func(ctx context.Context, user *dto.DeleteCMSUserPayload) error {
					return fmt.Errorf("failed to delete cms staff")
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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

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
						DefaultFacility: &domain.Facility{
							ID:   &ID,
							Name: gofakeit.Name(),
						},
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
						DefaultFacility: &domain.Facility{
							ID:   &ID,
							Name: gofakeit.Name(),
						},
					}, nil
				}
				// update the client profile
				fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
						DefaultFacility: &domain.Facility{
							ID: &ID,
						},
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
						DefaultFacility: &domain.Facility{
							ID: &ID,
						},
					}, nil
				}
				// update the client profile
				fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &ID,
						DefaultFacility: &domain.Facility{
							ID: &ID,
						},
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
			name: "Sad Case - Unable to check username exists",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - username exists",
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
		{
			name: "Sad case: unable to publish cms staff to pubsub",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get logged in user id",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx:   context.Background(),
				input: *payload,
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
			fakeUser := mock.NewUserUseCaseMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

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
			if tt.name == "Sad Case - Unable to check username exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, fmt.Errorf("failed to check username exists")
				}
			}

			if tt.name == "Sad Case - username exists" {

				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad Case - unable to check facility exists" {

				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("failed to check facility exists")
				}
			}
			if tt.name == "Sad Case - raise error if facility does not exists" {

				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "Sad Case - unable to retrieve facility by MFL Code" {

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility by MFL Code")
				}
			}
			if tt.name == "Sad Case - unable to register staff" {

				fakeDB.MockRegisterStaffFn = func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to register staff")
				}
			}
			if tt.name == "Sad Case - unable to assign roles" {

				fakeDB.MockAssignRolesFn = func(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
					return false, fmt.Errorf("failed to assign roles")
				}
			}
			if tt.name == "Sad Case - unable to invite staff" {

				fakeUser.MockInviteUserFn = func(ctx context.Context, userID, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return false, fmt.Errorf("failed to invite user")
				}
			}
			if tt.name == "Sad case: unable to publish cms staff to pubsub" {
				fakePubsub.MockNotifyCreateCMSStaffFn = func(ctx context.Context, user *dto.PubsubCreateCMSStaffPayload) error {
					return fmt.Errorf("failed to publish cms staff to pubsub")
				}
			}
			if tt.name == "Sad case: unable to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user id")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user id" {

				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get user profile by user id")
				}
			}

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			_, err := us.RegisterStaff(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_SetStaffDefaultFacility(t *testing.T) {
	type args struct {
		ctx        context.Context
		staffID    string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: staff update default facility",
			args: args{
				ctx:        nil,
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to update default facility, invalid facility id",
			args: args{
				ctx:        nil,
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get staff profile by staff id",
			args: args{
				ctx:        nil,
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to update default facility, update error",
			args: args{
				ctx:        nil,
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: staff not assigned to facility",
			args: args{
				ctx:        nil,
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to retrieve current facility",
			args: args{
				ctx:        nil,
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: failed to get staff profile by staff id" {
				fakeDB.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}
			if tt.name == "Sad case: failed to update default facility, update error" {
				fakeDB.MockUpdateStaffFn = func(ctx context.Context, staff *domain.StaffProfile, updates map[string]interface{}) error {
					return fmt.Errorf("failed to update staff profile")
				}
			}

			if tt.name == "Sad case: failed to update default facility, invalid facility id" {
				fakeDB.MockGetStaffFacilitiesFn = func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get staff facilities")
				}
			}
			if tt.name == "Sad case: staff not assigned to facility" {
				fakeDB.MockGetStaffFacilitiesFn = func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, pkgGorm.ErrRecordNotFound
				}
			}

			if tt.name == "Sad case: failed to retrieve current facility" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := us.SetStaffDefaultFacility(tt.args.ctx, tt.args.staffID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetStaffDefaultFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_SetClientDefaultFacility(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: client update default facility",
			args: args{
				ctx:        nil,
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to update default facility, invalid facility id",
			args: args{
				ctx:        nil,
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get client profile by client",
			args: args{
				ctx:        nil,
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: client not assigned to facility",
			args: args{
				ctx:        nil,
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to update default facility, update error",
			args: args{
				ctx:        nil,
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to retrieve current facility",
			args: args{
				ctx:        nil,
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: failed to get client profile by client" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get  client profile")
				}
			}

			if tt.name == "Sad case: failed to update default facility, update error" {
				fakeDB.MockUpdateClientFn = func(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to update client profile")
				}
			}

			if tt.name == "Sad case: failed to update default facility, invalid facility id" {
				fakeDB.MockGetClientFacilitiesFn = func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get client facilities")
				}
			}

			if tt.name == "Sad case: client not assigned to facility" {
				fakeDB.MockGetClientFacilitiesFn = func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, pkgGorm.ErrRecordNotFound
				}
			}

			if tt.name == "Sad case: failed to retrieve current facility" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := us.SetClientDefaultFacility(tt.args.ctx, tt.args.clientID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetClientDefaultFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_AddFacilitiesToStaffProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		staffID    string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: assign facilities to staff",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: missing client id",
			args: args{
				ctx:        context.Background(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: missing facility id",
			args: args{
				ctx:     context.Background(),
				staffID: uuid.NewString(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to assign facilities to staff",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: failed to assign facilities to staff" {
				fakeDB.MockAddFacilitiesToStaffProfileFn = func(ctx context.Context, staffID string, facilities []string) error {
					return fmt.Errorf("failed to add facilities to staff profile")
				}
			}
			got, err := us.AddFacilitiesToStaffProfile(tt.args.ctx, tt.args.staffID, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.AddFacilitiesToStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.AddFacilitiesToStaffProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestUseCasesUserImpl_GetUserLinkedFacilities(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		userID     string
// 		pagination *dto.PaginationsInput
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []*domain.Facility
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy Case - Successfully get client linked facilities",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Happy Case - Successfully get staff linked facilities",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Sad Case - Invalid pagination input",
// 			args: args{
// 				ctx:        context.Background(),
// 				pagination: &dto.PaginationsInput{},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Fail to get staff linked facilities, missing user ID",
// 			args: args{
// 				ctx: context.Background(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Fail to get user profile by user id",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Fail to get client profile by user id",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Fail to get client facilities",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Fail to get staff profile by user id",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Fail to get staff facilities",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case - Invalid user type",
// 			args: args{
// 				ctx:    context.Background(),
// 				userID: uuid.NewString(),
// 				pagination: &dto.PaginationsInput{
// 					CurrentPage: 1,
// 					Limit:       10,
// 				},
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			fakeDB := pgMock.NewPostgresMock()
// 			fakeExtension := extensionMock.NewFakeExtension()
// 			fakeOTP := otpMock.NewOTPUseCaseMock()
// 			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
// 			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
// 			fakePubsub := pubsubMock.NewPubsubServiceMock()
// 			fakeClinical := clinicalMock.NewClinicalServiceMock()
// 			fakeSMS := smsMock.NewSMSServiceMock()
// 			fakeTwilio := twilioMock.NewTwilioServiceMock()

// 			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

// 			if tt.name == "Happy Case - Successfully get client linked facilities" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}
// 			}

// 			if tt.name == "Happy Case - Successfully get staff linked facilities" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}
// 			}

// 			if tt.name == "Sad Case - Fail to get user profile by user id" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return nil, fmt.Errorf("failed to get user profile by user ID")
// 				}
// 			}

// 			if tt.name == "Sad Case - Fail to get client profile by user id" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}

// 				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID string) (*domain.ClientProfile, error) {
// 					return nil, fmt.Errorf("failed to get client profile by user ID")
// 				}
// 			}

// 			if tt.name == "Sad Case - Fail to get client facilities" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}

// 				fakeDB.MockGetClientFacilitiesFn = func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
// 					return nil, nil, fmt.Errorf("failed to get client facilities")
// 				}
// 			}

// 			if tt.name == "Sad Case - Fail to get staff profile by user id" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}

// 				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
// 					return nil, fmt.Errorf("failed to get staff profile by user ID")
// 				}
// 			}

// 			if tt.name == "Sad Case - Fail to get staff facilities" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}

// 				fakeDB.MockGetStaffFacilitiesFn = func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
// 					return nil, nil, fmt.Errorf("failed to get staff facilities")
// 				}
// 			}

// 			if tt.name == "Sad Case - Invalid user type" {
// 				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
// 					return &domain.User{}, nil
// 				}
// 			}

// 			got, err := us.GetUserLinkedFacilities(tt.args.ctx, tt.args.userID, *tt.args.pagination)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("UseCasesUserImpl.GetUserLinkedFacilities() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got == nil && !tt.wantErr {
// 				t.Errorf("expected a response but got %v", got)
// 			}
// 		})
// 	}
// }

func TestUseCasesUserImpl_AddFacilitiesToClientProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: assign facilities to clients",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: missing client id",
			args: args{
				ctx:        context.Background(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: missing facility id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.NewString(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to assign facilities to clients",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: failed to assign facilities to clients" {
				fakeDB.MockAddFacilitiesToClientProfileFn = func(ctx context.Context, clientID string, facilities []string) error {
					return fmt.Errorf("error adding facilities to client profile")
				}
			}

			got, err := us.AddFacilitiesToClientProfile(tt.args.ctx, tt.args.clientID, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.AddFacilitiesToClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.AddFacilitiesToClientProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterCaregiver(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.CaregiverInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: register caregiver",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      false,
					AssignedClients: []dto.ClientCaregiverInput{
						{
							ClientID:      gofakeit.UUID(),
							CaregiverType: enums.CaregiverTypeFather,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: username exists",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      false,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: username check error",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      false,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid phone number",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     "+2547",
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      false,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: register caregiver error",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      false,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to send invite",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      true,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to assign client",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      true,
					AssignedClients: []dto.ClientCaregiverInput{
						{
							ClientID:      gofakeit.UUID(),
							CaregiverType: enums.CaregiverTypeFather,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get logged in user",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      true,
					AssignedClients: []dto.ClientCaregiverInput{
						{
							ClientID:      gofakeit.UUID(),
							CaregiverType: enums.CaregiverTypeFather,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get user profile of the logged in user",
			args: args{
				ctx: context.Background(),
				input: dto.CaregiverInput{
					Name:   gofakeit.Name(),
					Gender: enumutils.GenderMale,
					DateOfBirth: scalarutils.Date{
						Year:  10,
						Month: 10,
						Day:   10,
					},
					PhoneNumber:     gofakeit.Phone(),
					CaregiverNumber: gofakeit.SSN(),
					SendInvite:      true,
					AssignedClients: []dto.ClientCaregiverInput{
						{
							ClientID:      gofakeit.UUID(),
							CaregiverType: enums.CaregiverTypeFather,
						},
					},
				},
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "sad case: username check error" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, fmt.Errorf("failed to check phone number")
				}
			}

			if tt.name == "sad case: username exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "sad case: register caregiver error" {

				fakeDB.MockRegisterCaregiverFn = func(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
					return nil, fmt.Errorf("failed to register caregiver")
				}
			}

			if tt.name == "sad case: fail to send invite" {

				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: fail to assign client" {

				fakeDB.MockAddCaregiverToClientFn = func(ctx context.Context, clientCaregiver *domain.CaregiverClient) error {
					return fmt.Errorf("failed to assign caregiver")
				}
			}

			if tt.name == "sad case: unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "sad case: unable to get user profile of the logged in user" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			got, err := us.RegisterCaregiver(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}

		})
	}
}

func TestUseCasesUserImpl_RegisterClientAsCaregiver(t *testing.T) {

	type args struct {
		ctx             context.Context
		clientID        string
		caregiverNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create a caregiver",
			args: args{
				ctx:             nil,
				clientID:        gofakeit.UUID(),
				caregiverNumber: gofakeit.SSN(),
			},
			wantErr: false,
		},
		{
			name: "sad case: get client error",
			args: args{
				ctx:             nil,
				clientID:        gofakeit.UUID(),
				caregiverNumber: gofakeit.SSN(),
			},
			wantErr: true,
		},
		{
			name: "sad case: create caregiver error",
			args: args{
				ctx:             nil,
				clientID:        gofakeit.UUID(),
				caregiverNumber: gofakeit.SSN(),
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "sad case: get client error" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client")
				}
			}

			if tt.name == "sad case: create caregiver error" {
				fakeDB.MockCreateCaregiverFn = func(ctx context.Context, caregiver domain.Caregiver) (*domain.Caregiver, error) {
					return nil, fmt.Errorf("failed to create caregiver")
				}
			}

			got, err := us.RegisterClientAsCaregiver(tt.args.ctx, tt.args.clientID, tt.args.caregiverNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterClientAsCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_SearchCaregiverUser(t *testing.T) {
	type args struct {
		ctx             context.Context
		searchParameter string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search caregiver user",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: unable to search caregiver user" {
				fakeDB.MockSearchCaregiverUserFn = func(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
					return nil, fmt.Errorf("failed to search caregiver user")
				}
			}
			got, err := us.SearchCaregiverUser(tt.args.ctx, tt.args.searchParameter)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SearchCaregiverUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_RemoveFacilitiesFromClientProfile(t *testing.T) {

	type args struct {
		ctx        context.Context
		clientID   string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case: remove facilities from client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case: missing client id",
			args: args{
				ctx:        context.Background(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case: missing facility",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get client profile by client id",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case: failed to remove facilities from client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case: failed to get client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by client id")
				}
			}

			if tt.name == "Sad Case: failed to remove facilities from client profile" {
				fakeDB.MockRemoveFacilitiesFromClientProfileFn = func(ctx context.Context, clientID string, facilities []string) error {
					return fmt.Errorf("failed to remove facilities from client profile")
				}
			}

			got, err := us.RemoveFacilitiesFromClientProfile(tt.args.ctx, tt.args.clientID, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RemoveFacilitiesFromClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.RemoveFacilitiesFromClientProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_AssignCaregiver(t *testing.T) {
	ID := gofakeit.UUID()
	CaregiverID := gofakeit.UUID()
	CaregiverType := enums.CaregiverTypeFather

	type args struct {
		ctx   context.Context
		input dto.ClientCaregiverInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: add caregiver to client",
			args: args{
				ctx: nil,
				input: dto.ClientCaregiverInput{
					ClientID:      ID,
					CaregiverID:   CaregiverID,
					CaregiverType: CaregiverType,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing caregiver ID",
			args: args{
				ctx: nil,
				input: dto.ClientCaregiverInput{
					ClientID:      ID,
					CaregiverType: CaregiverType,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to add caregiver to client",
			args: args{
				ctx: nil,
				input: dto.ClientCaregiverInput{
					ClientID:      ID,
					CaregiverID:   CaregiverID,
					CaregiverType: CaregiverType,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to get logged in user",
			args: args{
				ctx: nil,
				input: dto.ClientCaregiverInput{
					ClientID:      ID,
					CaregiverID:   CaregiverID,
					CaregiverType: CaregiverType,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to get staff profile",
			args: args{
				ctx: nil,
				input: dto.ClientCaregiverInput{
					ClientID:      ID,
					CaregiverID:   CaregiverID,
					CaregiverType: CaregiverType,
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "sad case: unable to add caregiver to client" {
				fakeDB.MockAddCaregiverToClientFn = func(ctx context.Context, clientCaregiver *domain.CaregiverClient) error {
					return fmt.Errorf("failed to add caregiver to client")
				}
			}
			if tt.name == "sad case: unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}
			if tt.name == "sad case: unable to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, staffID string, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}
			got, err := us.AssignCaregiver(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.AddCaregiverToClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.AddCaregiverToClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_RemoveFacilitiesFromStaffProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		staffID    string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case: remove facilities from staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case: missing staff id",
			args: args{
				ctx:        context.Background(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case: missing facility",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get staff profile by staff id",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case: failed to remove facilities from staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case: failed to get staff profile by staff id" {
				fakeDB.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by staff id")
				}
			}

			if tt.name == "Sad Case: failed to remove facilities from staff profile" {
				fakeDB.MockRemoveFacilitiesFromStaffProfileFn = func(ctx context.Context, staffID string, facilities []string) error {
					return fmt.Errorf("failed to remove facilities from staff profile")
				}
			}

			got, err := us.RemoveFacilitiesFromStaffProfile(tt.args.ctx, tt.args.staffID, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RemoveFacilitiesFromStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.RemoveFacilitiesFromStaffProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_GetCaregiverManagedClients(t *testing.T) {
	type args struct {
		ctx         context.Context
		caregiverID string
		input       dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get managed clients",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				input: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
					Sort: dto.SortsInput{
						Direction: enums.SortDataTypeDesc,
						Field:     "id",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get managed clients",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				input: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
					Sort: dto.SortsInput{
						Direction: enums.SortDataTypeDesc,
						Field:     "id",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to validate pagination input",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad case: failed to get managed clients" {
				fakeDB.MockGetCaregiverManagedClientsFn = func(ctx context.Context, caregiverID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get managed clients")
				}
			}
			got, err := us.GetCaregiverManagedClients(tt.args.ctx, tt.args.caregiverID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetCaregiverManagedClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect error, got: %v", err)
			}
		})
	}
}

func TestUseCasesUserImpl_ListClientsCaregivers(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		pagination *dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.CaregiverProfileOutputPage
		wantErr bool
	}{
		{
			name: "Happy Case: list clients caregivers",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.NewString(),
				pagination: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: unable to list clients caregivers",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.NewString(),
				pagination: &dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case: unable to list clients caregivers" {
				fakeDB.MockListClientsCaregiversFn = func(ctx context.Context, clientID string, pagination *domain.Pagination) (*domain.ClientCaregivers, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("unable to list clients caregivers")
				}
			}
			got, err := us.ListClientsCaregivers(tt.args.ctx, tt.args.clientID, tt.args.pagination)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.ListClientsCaregivers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got %v", got)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_ConsentToAClientCaregiver(t *testing.T) {
	type args struct {
		ctx         context.Context
		clientID    string
		caregiverID string
		consent     bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case: client consent",
			args: args{
				ctx:         context.Background(),
				clientID:    uuid.NewString(),
				caregiverID: uuid.NewString(),
				consent:     true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case: client unable to consent",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.NewString(),
				consent:  true,
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case: client unable to consent" {
				fakeDB.MockUpdateCaregiverClientFn = func(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update caregiver client")
				}
			}
			got, err := us.ConsentToAClientCaregiver(tt.args.ctx, tt.args.clientID, tt.args.caregiverID, tt.args.consent)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.TestUseCasesUserImpl_ConsentToAClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.TestUseCasesUserImpl_ConsentToAClientCaregiver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_ConsentToManagingClient(t *testing.T) {
	type args struct {
		ctx         context.Context
		caregiverID string
		clientID    string
		consent     bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case: consent to managing client",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				clientID:    uuid.NewString(),
				consent:     true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case: unable to consent to managing client",
			args: args{
				ctx:         context.Background(),
				caregiverID: uuid.NewString(),
				clientID:    uuid.NewString(),
				consent:     true,
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case: unable to consent to managing client" {
				fakeDB.MockUpdateCaregiverClientFn = func(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error {
					return fmt.Errorf("unable to consent to managing client")
				}
			}
			got, err := us.ConsentToManagingClient(tt.args.ctx, tt.args.caregiverID, tt.args.clientID, tt.args.consent)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.ConsentToManagingClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.ConsentToManagingClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_FetchContactOrganisations(t *testing.T) {

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "sad case: invalid phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: "+123",
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "sad case: fail to find contacts",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "sad case: fail to find organisation",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "sad case: contact doesn't exist",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "happy case: find organisation",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "happy case: multiple similar organisations",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantCount: 1,
			wantErr:   false,
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
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakeGetStream, fakePubsub, fakeClinical, fakeSMS, fakeTwilio)

			if tt.name == "sad case: fail to find contacts" {
				fakeDB.MockFindContactsFn = func(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error) {
					return nil, fmt.Errorf("failed to find contact")
				}
			}

			if tt.name == "sad case: contact doesn't exist" {
				fakeDB.MockFindContactsFn = func(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error) {
					return []*domain.Contact{}, nil
				}
			}

			if tt.name == "sad case: fail to find organisation" {
				fakeDB.MockGetOrganisationFn = func(ctx context.Context, id string) (*domain.Organisation, error) {
					return nil, fmt.Errorf("failed to find organisation")
				}
			}

			if tt.name == "happy case: multiple similar organisations" {
				fakeDB.MockFindContactsFn = func(ctx context.Context, contactType, contactValue string) ([]*domain.Contact, error) {
					sameID := gofakeit.UUID()
					return []*domain.Contact{
						{
							OrganisationID: sameID,
						},
						{
							OrganisationID: sameID,
						},
						{
							OrganisationID: sameID,
						},
					}, nil
				}
			}

			got, err := us.FetchContactOrganisations(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.FetchContactOrganisations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.wantCount) {
				t.Errorf("PGInstance.FindContacts() = %v, want %v", got, tt.wantCount)
			}
		})
	}
}
