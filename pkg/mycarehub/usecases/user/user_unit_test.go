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
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

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
	matrixMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	smsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms/mock"
	twilioMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/silcomms"
	"github.com/segmentio/ksuid"
	pkgGorm "gorm.io/gorm"
)

func TestMain(m *testing.M) {
	initialMatrixUserEnv := os.Getenv("MCH_MATRIX_USER")
	initialMatrixPasswordEnv := os.Getenv("MCH_MATRIX_PASSWORD")

	// set test envs
	os.Setenv("MCH_MATRIX_USER", "test user")
	os.Setenv("MCH_MATRIX_PASSWORD", "test pass")

	code := m.Run()

	// restore envs
	os.Setenv("MCH_MATRIX_USER", initialMatrixUserEnv)
	os.Setenv("MCH_MATRIX_PASSWORD", initialMatrixPasswordEnv)

	os.Exit(code)
}

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
			want1: false,
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
			name: "Happy Case - should not fail when CCC number is not found",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Happy case: consumer login" {
				currentTime := time.Now()
				laterTime := currentTime.Add(time.Minute * 2005)
				salt, encryptedPin := utils.EncryptPIN("1234", nil)
				fakeDB.MockGetUserPINByUserIDFn = func(ctx context.Context, userID string) (*domain.UserPIN, error) {
					return &domain.UserPIN{
						UserID:    "f3f8f8f8-f3f8-f3f8-f3f8-f3f8f8f8f8f8",
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
				fakeDB.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
					return nil, fmt.Errorf("failed to get client ccc number identifier value")
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
	validIntlPhone := "+32468799972"

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
			name: "Happy case: send sms to foreign phone number",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validIntlPhone,
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
		{
			name: "Sad case: unable to send sms to foreign phone number",
			args: args{
				ctx:         ctx,
				userID:      userID,
				phoneNumber: validIntlPhone,
				flavour:     validFlavour,
				reinvite:    false,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			if tt.name == "Sad case: unable to send sms to foreign phone number" {
				fakeTwilio.MockSendSMSViaTwilioFn = func(ctx context.Context, phonenumber, message string) error {
					return fmt.Errorf("unable to send SMS via twillio")
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			name: "Sad Case: nickname exists",
			args: args{
				ctx:      ctx,
				userID:   userID,
				nickname: gofakeit.Username(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			u := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Happy case: Successfully set nickname" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, nil
				}
			}
			if tt.name == "Sad Case: nickname exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case - no nickname" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
				fakeOTP.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTPResponse, error) {
					return nil, fmt.Errorf("failed to generate and send otp")
				}
			}

			if tt.name == "Sad Case - Fail to save otp" {
				fakeOTP.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTPResponse, error) {
					return &domain.OTPResponse{
						OTP:         "123456",
						PhoneNumber: interserviceclient.TestUserPhoneNumber,
					}, nil
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "1234",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid phone flavor",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					Username: gofakeit.Word(),
					Flavour:  "invalid",
					OTP:      "111222",
					PIN:      "1234",
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "abcd",
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "12345",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to get user profile by username",
			args: args{
				ctx: context.Background(),
				input: dto.UserResetPinInput{
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "1234",
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "1234",
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "1234",
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
					PIN:      "1234",
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
					Username: gofakeit.Word(),
					Flavour:  feedlib.FlavourConsumer,
					OTP:      "111222",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Happy Case - Successfully reset pin" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
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
			if tt.name == "invalid: failed to get user profile by username" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
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
				fakeDB.MockGetUserProfileByUsernameFn = func(ctx context.Context, username string) (*domain.User, error) {
					return nil, errors.New("failed to get user profile by phone")
				}
			}

			if tt.name == "invalid: failed to verify OTP" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
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
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
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
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID, flavour string) ([]*domain.SecurityQuestionResponse, error) {
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
		PhoneNumber: interserviceclient.TestUserPhoneNumber,
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
		{
			name: "Sad case: unable to register matrix user",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to check whether matrix user is an admin",
			args: args{
				ctx:   context.Background(),
				input: payload,
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get program by id",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("unable to check identifier exists")
				}
			}

			if tt.name == "Sad case: fail if identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
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
			if tt.name == "sad case: unable to check whether matrix user is an admin" {
				fakeMatrix.MockCheckIfUserIsAdminFn = func(ctx context.Context, auth *domain.MatrixAuth, userID string) (bool, error) {
					return false, fmt.Errorf("failed to check whether matrix user is an admin")
				}
			}
			if tt.name == "Sad case: unable to register matrix user" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("unable to register matrix user")
				}
			}
			if tt.name == "sad case: unable to get program by id" {
				fakeDB.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*domain.Program, error) {
					return nil, fmt.Errorf("error")
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

func TestUseCasesUserImpl_GetUserProfile(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()

	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()
	fakeMatrix := matrixMock.NewMatrixMock()

	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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

func TestUseCasesUserImpl_GetStaffProfile(t *testing.T) {

	type args struct {
		ctx       context.Context
		userID    string
		programID string
	}
	tests := []struct {
		name string
		args args

		wantErr bool
	}{
		{
			name: "happy case: get staff profile",
			args: args{
				ctx:       context.Background(),
				userID:    gofakeit.UUID(),
				programID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get staff profile",
			args: args{
				ctx:       context.Background(),
				userID:    gofakeit.UUID(),
				programID: gofakeit.UUID(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "sad case: failed to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
					return nil, errors.New("an error occurred")
				}
			}
			_, err := us.GetStaffProfile(tt.args.ctx, tt.args.userID, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
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
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()
	fakeMatrix := matrixMock.NewMatrixMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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

				fakeDB.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
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

				fakeDB.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
					return []*domain.Identifier{{Value: "123456"}}, nil
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

				fakeDB.MockGetClientIdentifiers = func(ctx context.Context, clientID string) ([]*domain.Identifier, error) {
					return []*domain.Identifier{{Value: "123456"}}, nil
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
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()
	fakeMatrix := matrixMock.NewMatrixMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			name: "Sad case: failed to search staff",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - failed to get logged in user id",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - failed to get user profile",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - failed to get staff profile",
			args: args{
				ctx:         ctx,
				staffNumber: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to search staff" {
				fakeDB.MockSearchStaffProfileFn = func(ctx context.Context, searchParameter string, programID *string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - failed to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - failed to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
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

func TestUseCasesUserImpl_SearchClientUser(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()
	fakeMatrix := matrixMock.NewMatrixMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			name: "Happy case: search client",
			args: args{
				ctx:             ctx,
				searchParameter: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to search client",
			args: args{
				ctx:             ctx,
				searchParameter: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Missing search parameter",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: failed to search client" {
				fakeDB.MockSearchClientProfileFn = func(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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

func TestUseCasesUserImpl_GetProgramClientProfileByIdentifier(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()
	fakeMatrix := matrixMock.NewMatrixMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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

		{
			name: "Sad Case - unable to get logged in user profile",
			args: args{
				ctx:       context.Background(),
				cccNumber: "123456789",
			},
			wantErr: true,
		},

		{
			name: "Sad Case - unable to get logged in user profile",
			args: args{
				ctx:       context.Background(),
				cccNumber: "123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad Case - unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("unable to get logged in user")
				}
			}

			if tt.name == "Sad Case - unable to get logged in user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("unable to get user profile")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile by ccc number" {
				fakeDB.MockGetProgramClientProfileByIdentifierFn = func(ctx context.Context, programID string, identifierType string, value string) (*domain.ClientProfile, error) {
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

func TestUseCasesUserImpl_TransferClientToFacility(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeOTP := otpMock.NewOTPUseCaseMock()
	fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeClinical := clinicalMock.NewClinicalServiceMock()
	fakeSMS := smsMock.NewSMSServiceMock()
	fakeTwilio := twilioMock.NewTwilioServiceMock()
	fakeMatrix := matrixMock.NewMatrixMock()
	us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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

func TestUseCasesUserImpl_RegisterStaffProfile(t *testing.T) {
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
			name: "Happy Case - register staff",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: false,
		},
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
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeUser := mock.NewUserUseCaseMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			if tt.name == "Sad Case - Unable to check identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("failed to check identifier exists")
				}
			}
			if tt.name == "Sad Case - identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
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
			if tt.name == "Sad Case - unable to invite staff" {

				fakeUser.MockInviteUserFn = func(ctx context.Context, userID, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return false, fmt.Errorf("failed to invite user")
				}
			}

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			_, err := us.RegisterStaffProfile(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
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
			name: "Happy Case - register staff",
			args: args{
				ctx:   ctx,
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
		{
			name: "sad case: unable to check whether matrix user is an admin",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to register matrix user",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

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
			if tt.name == "Sad case: unable to register matrix user" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("unable to register matrix user")
				}
			}
			if tt.name == "sad case: unable to check whether matrix user is an admin" {
				fakeMatrix.MockCheckIfUserIsAdminFn = func(ctx context.Context, auth *domain.MatrixAuth, userID string) (bool, error) {
					return false, fmt.Errorf("failed to check whether matrix user is an admin")
				}
			}

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			_, err := us.RegisterStaff(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterOrganisationAdmin(t *testing.T) {
	ctx := context.Background()

	payload := &dto.StaffRegistrationInput{
		Username:       gofakeit.Name(),
		Facility:       "1234",
		StaffName:      gofakeit.BeerName(),
		Gender:         enumutils.GenderMale,
		DateOfBirth:    scalarutils.Date{Year: 2000, Month: 2, Day: 20},
		PhoneNumber:    interserviceclient.TestUserPhoneNumber,
		IDNumber:       "123456789",
		StaffNumber:    "123456789",
		StaffRoles:     "Community Management",
		InviteStaff:    true,
		ProgramID:      gofakeit.UUID(),
		OrganisationID: gofakeit.UUID(),
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
			name: "Happy Case - register organisation admin",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get staff profile by id",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: true,
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
			name: "Sad case: unable to register staff",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: true,
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
			name: "Sad case: unable to get logged in user profile",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to register user in matrix",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()

			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			if tt.name == "Sad case: unable to get staff profile by id" {
				fakeDB.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*domain.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to register staff" {
				fakeDB.MockRegisterStaffFn = func(ctx context.Context, staffRegistrationPayload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to get logged in user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to register user in matrix" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			_, err := us.RegisterOrganisationAdmin(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterOrganisationAdmin() error = %v, wantErr %v", err, tt.wantErr)
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
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to update default facility, invalid facility id",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get staff profile by staff id",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to update default facility, update error",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: staff not assigned to facility",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to retrieve current facility",
			args: args{
				ctx:        context.Background(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to update default facility, invalid facility id",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get client profile by client",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: client not assigned to facility",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to update default facility, update error",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilityID: uuid.NewString(),
			},
			wantErr: true,
		},

		{
			name: "Sad case: failed to retrieve current facility",
			args: args{
				ctx:        context.Background(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			name: "Happy case: assign facility to staff",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy case: assign facilities to staff",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString(), gofakeit.UUID()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: no staff id provided",
			args: args{
				ctx: context.Background(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to retrieve facility",
			args: args{
				ctx:        context.Background(),
				staffID:    gofakeit.UUID(),
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
		{
			name: "Sad case: unable to get staff profile by staff id",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to send sms",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to retrieve facility" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}
			if tt.name == "Sad case: failed to assign facilities to staff" {
				fakeDB.MockAddFacilitiesToStaffProfileFn = func(ctx context.Context, staffID string, facilities []string) error {
					return fmt.Errorf("failed to add facilities to staff profile")
				}
			}
			if tt.name == "Sad case: unable to get staff profile by staff id" {
				fakeDB.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("unable to get staff profile")
				}
			}
			if tt.name == "Sad case: unable to send sms" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
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
			name: "Sad case: no client id provided",
			args: args{
				ctx:        context.Background(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to retrieve facility",
			args: args{
				ctx:        context.Background(),
				clientID:   gofakeit.UUID(),
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
		{
			name: "Sad case: unable to get client profile by client id",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to send sms",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to retrieve facility" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}
			if tt.name == "Sad case: failed to assign facilities to clients" {
				fakeDB.MockAddFacilitiesToClientProfileFn = func(ctx context.Context, clientID string, facilities []string) error {
					return fmt.Errorf("error adding facilities to client profile")
				}
			}
			if tt.name == "Sad case: unable to get client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("unable to get client profile")
				}
			}
			if tt.name == "Sad case: unable to send sms" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
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
					PhoneNumber:     interserviceclient.TestUserPhoneNumber,
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
		{
			name: "sad case: unable to check whether matrix user is an admin",
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
			name: "sad case: unable to register matrix user",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			if tt.name == "sad case: unable to check whether matrix user is an admin" {
				fakeMatrix.MockCheckIfUserIsAdminFn = func(ctx context.Context, auth *domain.MatrixAuth, userID string) (bool, error) {
					return false, fmt.Errorf("failed to check whether matrix user is an admin")
				}
			}
			if tt.name == "sad case: unable to register matrix user" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("failed to register matrix user")
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
				ctx:             context.Background(),
				clientID:        gofakeit.UUID(),
				caregiverNumber: gofakeit.SSN(),
			},
			wantErr: false,
		},
		{
			name: "sad case: get client error",
			args: args{
				ctx:             context.Background(),
				clientID:        gofakeit.UUID(),
				caregiverNumber: gofakeit.SSN(),
			},
			wantErr: true,
		},
		{
			name: "sad case: create caregiver error",
			args: args{
				ctx:             context.Background(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
		{
			name: "Sad case - failed to get logged in user id",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - failed to get user profile",
			args: args{
				ctx:             context.Background(),
				searchParameter: gofakeit.Name(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - failed to get staff profile",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case - failed to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - failed to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
			name: "sad case: unable to get user profile by logged in user id",
			args: args{
				ctx: context.Background(),
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
				ctx: context.Background(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			if tt.name == "sad case: unable to get user profile by logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by logged in user id")
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
		ctx    context.Context
		userID string
		input  dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get managed clients",
			args: args{
				ctx:    context.Background(),
				userID: uuid.NewString(),
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
				ctx:    context.Background(),
				userID: uuid.NewString(),
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
				ctx:    context.Background(),
				userID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to update user",
			args: args{
				ctx:    context.Background(),
				userID: uuid.NewString(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: failed to get managed clients" {
				fakeDB.MockGetCaregiverManagedClientsFn = func(ctx context.Context, userID string, pagination *domain.Pagination) ([]*domain.ManagedClient, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to get managed clients")
				}
			}

			if tt.name == "Sad case: failed to update user" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := us.GetCaregiverManagedClients(tt.args.ctx, tt.args.userID, tt.args.input)
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
		consent     enums.ConsentState
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
				consent:     enums.ConsentStateAccepted,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case: client unable to consent",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.NewString(),
				consent:  enums.ConsentStateAccepted,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
		consent     enums.ConsentState
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
				consent:     enums.ConsentStateAccepted,
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
				consent:     enums.ConsentStateAccepted,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

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

func TestUseCasesUserImpl_GetStaffFacilities(t *testing.T) {
	type args struct {
		ctx             context.Context
		staffID         string
		paginationInput dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.FacilityOutputPage
		wantErr bool
	}{
		{
			name: "Happy case: get staff facilities",
			args: args{
				ctx:     context.Background(),
				staffID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get staff facilities",
			args: args{
				ctx:     context.Background(),
				staffID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid pagination",
			args: args{
				ctx:     context.Background(),
				staffID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit: 10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx:     context.Background(),
				staffID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx:     context.Background(),
				staffID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to update user",
			args: args{
				ctx:     context.Background(),
				staffID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to get staff facilities" {
				fakeDB.MockGetStaffFacilitiesFn = func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: invalid pagination" {
				fakeDB.MockGetStaffFacilitiesFn = func(ctx context.Context, input dto.StaffFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("unable to find the loggedin user id")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}
			if tt.name == "Sad case: unable to update user" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			_, err := us.GetStaffFacilities(tt.args.ctx, tt.args.staffID, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetStaffFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterExistingUserAsClient(t *testing.T) {
	cccNumber := "1109410004"
	type args struct {
		ctx   context.Context
		input dto.ExistingUserClientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: register existing user as client",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: register existing user as client, user already has a client profile",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get client profiles",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: user has no client profile",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to check if client exists in program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: user already exists in program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to check if facility exist in program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: facility does not exist in program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to register existing user as client",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to create client in cms",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to create patient in fhir",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserClientInput{
					FacilityID:  uuid.NewString(),
					ClientTypes: []enums.ClientType{"PMTCT"},
					EnrollmentDate: scalarutils.Date{
						Year:  2020,
						Month: 1,
						Day:   1,
					},
					CCCNumber:    &cccNumber,
					Counselled:   true,
					InviteClient: true,
					UserID:       uuid.NewString(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "sad case: failed to get client profiles" {
				fakeDB.MockGetUserClientProfilesFn = func(ctx context.Context, userID string) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: user has no client profile" {
				fakeDB.MockGetUserClientProfilesFn = func(ctx context.Context, userID string) ([]*domain.ClientProfile, error) {
					return []*domain.ClientProfile{}, nil
				}
			}
			if tt.name == "sad case: failed to check if client exists in program" {
				fakeDB.MockCheckIfClientExistsInProgramFn = func(ctx context.Context, userID string, programID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: user already exists in program" {
				fakeDB.MockCheckIfClientExistsInProgramFn = func(ctx context.Context, userID string, programID string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "sad case: failed to get program" {
				fakeDB.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*domain.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to check if facility exist in program" {
				fakeDB.MockCheckIfFacilityExistsInProgramFn = func(ctx context.Context, programID string, facilityID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: facility does not exist in program" {
				fakeDB.MockCheckIfFacilityExistsInProgramFn = func(ctx context.Context, programID string, facilityID string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: failed to register existing user as client" {
				fakeDB.MockRegisterExistingUserAsClientFn = func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: failed to create client in cms" {
				fakePubsub.MockNotifyCreateCMSClientFn = func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to create patient in fhir" {
				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("an error occurred")
				}
			}

			_, err := us.RegisterExistingUserAsClient(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterExistingUserAsClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_GetClientFacilities(t *testing.T) {
	type args struct {
		ctx             context.Context
		clientID        string
		paginationInput dto.PaginationsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get client facilities",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get client facilities",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid pagination",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit: 10,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
				paginationInput: dto.PaginationsInput{
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to get client facilities" {
				fakeDB.MockGetClientFacilitiesFn = func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: invalid pagination" {
				fakeDB.MockGetClientFacilitiesFn = func(ctx context.Context, input dto.ClientFacilityInput, pagination *domain.Pagination) ([]*domain.Facility, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("unable to find the loggedin user id")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get logged in user")
				}
			}

			_, err := us.GetClientFacilities(tt.args.ctx, tt.args.clientID, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.GetClientFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_SetCaregiverCurrentClient(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: set caregiver current client",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get logged in user",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get client profile by client id",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get caregiver profile by user id",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to update caregiver profile",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get caregivers clients",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: client consent not accepted",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: caregiver does not manage client",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to update user",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "sad case: failed to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to get client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to get caregiver profile by user id" {
				fakeDB.MockGetCaregiverProfileByUserIDFn = func(ctx context.Context, userID string, clientID string) (*domain.CaregiverProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to update caregiver profile" {
				fakeDB.MockUpdateCaregiverFn = func(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to user profile by user id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to update user profile" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: failed to get caregivers clients" {
				fakeDB.MockGetCaregiversClientFn = func(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: client consent not accepted" {
				fakeDB.MockGetCaregiversClientFn = func(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error) {
					return []*domain.CaregiverClient{
						{
							ClientConsent: enums.ConsentStateRejected,
						},
					}, nil
				}
			}
			if tt.name == "sad case: caregiver does not manage client" {
				fakeDB.MockGetCaregiversClientFn = func(ctx context.Context, caregiverClient domain.CaregiverClient) ([]*domain.CaregiverClient, error) {
					return []*domain.CaregiverClient{}, nil
				}
			}
			if tt.name == "Sad case: unable to update user" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("error")
				}
			}

			got, err := us.SetCaregiverCurrentClient(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetCaregiverCurrentClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterExistingUserAsStaff(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.ExistingUserStaffInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: add new program to staff",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to check if staff is already registered in program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: staff user already registered in program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get staff profiles",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: user does not have staff profile",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get program",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get facility",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to check program facility exists",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: program facility does not exist",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to add program to staff profile",
			args: args{
				ctx: context.Background(),
				input: dto.ExistingUserStaffInput{
					UserID:      uuid.NewString(),
					ProgramID:   uuid.NewString(),
					FacilityID:  uuid.NewString(),
					StaffNumber: "123456789",
					StaffRoles:  "SYSTEM_ADMIN",
					InviteStaff: true,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to check if staff is already registered in program" {
				fakeDB.MockCheckStaffExistsInProgramFn = func(ctx context.Context, userID string, programID string) (bool, error) {
					return false, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: staff user already registered in program" {
				fakeDB.MockCheckStaffExistsInProgramFn = func(ctx context.Context, userID string, programID string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case: unable to get staff profiles" {
				fakeDB.MockGetUserStaffProfilesFn = func(ctx context.Context, userID string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: user does not have staff profile" {
				fakeDB.MockGetUserStaffProfilesFn = func(ctx context.Context, userID string) ([]*domain.StaffProfile, error) {
					return []*domain.StaffProfile{}, nil
				}
			}

			if tt.name == "Sad case: unable to get program" {
				fakeDB.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*domain.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to get facility" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to check program facility exists" {
				fakeDB.MockCheckIfFacilityExistsInProgramFn = func(ctx context.Context, programID string, facilityID string) (bool, error) {
					return false, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: program facility does not exist" {
				fakeDB.MockCheckIfFacilityExistsInProgramFn = func(ctx context.Context, programID string, facilityID string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case: unable to add program to staff profile" {
				fakeDB.MockRegisterExistingUserAsStaffFn = func(ctx context.Context, payload *domain.StaffRegistrationPayload) (*domain.StaffProfile, error) {
					return nil, errors.New("unable to register existing user as staff")
				}
			}
			_, err := us.RegisterExistingUserAsStaff(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterExistingUserAsStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_SetCaregiverCurrentFacility(t *testing.T) {
	type args struct {
		ctx         context.Context
		caregiverID string
		facilityID  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: set caregiver current facility",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
				facilityID:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get caregiver profile by id",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
				facilityID:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to retrieve facility",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
				facilityID:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to update caregiver profile",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
				facilityID:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get logged in user id",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
				facilityID:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get client profile by id",
			args: args{
				ctx:         context.Background(),
				caregiverID: gofakeit.UUID(),
				facilityID:  gofakeit.UUID(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("unable to get logged in user id")
				}
			}
			if tt.name == "Sad case: unable to get client profile by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, errors.New("unable to get logged in user profile by id")
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("unable to get logged in user profile by id")
				}
			}
			if tt.name == "sad case: failed to get caregiver profile by id" {
				fakeDB.MockGetCaregiverProfileByUserIDFn = func(ctx context.Context, userID string, organisationID string) (*domain.CaregiverProfile, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "sad case: failed to retrieve facility" {
				fakeDB.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "sad case: failed to update caregiver profile" {
				fakeDB.MockUpdateCaregiverFn = func(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := us.SetCaregiverCurrentFacility(tt.args.ctx, tt.args.caregiverID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.SetCaregiverCurrentFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_RegisterExistingUserAsCaregiver(t *testing.T) {
	type args struct {
		ctx             context.Context
		userID          string
		caregiverNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: register existing user as caregiver",
			args: args{
				ctx:             context.Background(),
				userID:          "user-id",
				caregiverNumber: "caregiver-number",
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to get logged in user id",
			args: args{
				ctx:             context.Background(),
				userID:          "user-id",
				caregiverNumber: "caregiver-number",
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to get logged in user profile",
			args: args{
				ctx:             context.Background(),
				userID:          "user-id",
				caregiverNumber: "caregiver-number",
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to register existing user as caregiver",
			args: args{
				ctx:             context.Background(),
				userID:          "user-id",
				caregiverNumber: "caregiver-number",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "sad case: unable to get logged in user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get logged in user profile" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "user-id", nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to register existing user as caregiver" {
				fakeDB.MockRegisterExistingUserAsCaregiverFn = func(ctx context.Context, input *domain.CaregiverRegistration) (*domain.CaregiverProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := us.RegisterExistingUserAsCaregiver(tt.args.ctx, tt.args.userID, tt.args.caregiverNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.RegisterExistingUserAsCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_UpdateUserProfile(t *testing.T) {
	UID := uuid.NewString()
	CCCNumber := "ccc-number"
	username := "username"
	phoneNumber := interserviceclient.TestUserPhoneNumber
	email := gofakeit.Email()
	type args struct {
		ctx         context.Context
		userID      string
		cccNumber   *string
		username    *string
		phoneNumber *string
		programID   string
		flavour     feedlib.Flavour
		email       *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "consumer happy case: update user profile",
			args: args{
				ctx:         context.Background(),
				userID:      uuid.NewString(),
				cccNumber:   &CCCNumber,
				username:    &username,
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourConsumer,
				email:       &email,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "pro happy case: update user profile",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourPro,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: unable to get logged in user",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to get staff profile of the currently logged in user",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to get client profile",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update client identifier",
			args: args{
				ctx:       context.Background(),
				cccNumber: &CCCNumber,
				programID: uuid.NewString(),
				flavour:   feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update client username",
			args: args{
				ctx:       context.Background(),
				username:  &username,
				programID: uuid.NewString(),
				flavour:   feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update client phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to verify phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update user",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to get user profile - pro",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update staff phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to verify staff phone number",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update user - pro",
			args: args{
				ctx:         context.Background(),
				phoneNumber: &phoneNumber,
				programID:   uuid.NewString(),
				flavour:     feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update username - pro",
			args: args{
				ctx:       context.Background(),
				username:  &username,
				programID: uuid.NewString(),
				flavour:   feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: invalid flavour defined",
			args: args{
				ctx:       context.Background(),
				username:  &username,
				programID: uuid.NewString(),
				flavour:   "invalid",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: unable to update email",
			args: args{
				ctx:       context.Background(),
				username:  &username,
				programID: uuid.NewString(),
				flavour:   feedlib.FlavourPro,
				email:     &email,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "sad case: unable to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get staff profile of the currently logged in user" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get client profile" {
				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID, programID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update client identifier" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
					return &domain.StaffProfile{
						ProgramID: programID,
					}, nil
				}

				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID, programID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID:        &UID,
						ProgramID: programID,
						User: &domain.User{
							ID:               &UID,
							CurrentProgramID: programID,
						},
					}, nil
				}
				fakeDB.MockUpdateClientIdentifierFn = func(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update client username" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update client phone number" {
				fakeDB.MockUpdateUserContactFn = func(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to verify phone number" {
				fakeOTP.MockVerifyPhoneNumberFn = func(ctx context.Context, phone string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update user" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to get user profile - pro" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.NewString(), nil
				}
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
					return &domain.StaffProfile{
						ID: &UID,
					}, nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update staff phone number" {
				fakeDB.MockUpdateUserContactFn = func(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to verify staff phone number" {
				fakeOTP.MockVerifyPhoneNumberFn = func(ctx context.Context, phone string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update user - pro" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update username - pro" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "sad case: unable to update email" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := us.UpdateUserProfile(tt.args.ctx, tt.args.userID, tt.args.cccNumber, tt.args.username, tt.args.phoneNumber, tt.args.programID, tt.args.flavour, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.UpdateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.UpdateUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_CreateSuperUser(t *testing.T) {
	ctx := context.Background()
	payload := &dto.StaffRegistrationInput{
		Facility:  "3232323",
		StaffName: gofakeit.BeerName(),
		Gender:    enumutils.GenderMale,
		DateOfBirth: scalarutils.Date{
			Year:  2000,
			Month: 2,
			Day:   20,
		},
		PhoneNumber: interserviceclient.TestUserPhoneNumber,
		IDNumber:    "54545444",
		StaffNumber: "12345545456789",
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
			name: "Sad Case - unable to invite staff",
			args: args{
				ctx:   ctx,
				input: *payload,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx:   context.Background(),
				input: *payload,
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to register user in matrix",
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			fakeUser := mock.NewUserUseCaseMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Happy Case - Register Staff" {
				fakeDB.MockCheckIfSuperUserExistsFn = func(ctx context.Context) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad Case - Unable to check username exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, fmt.Errorf("failed to check username exists")
				}
			}
			if tt.name == "Sad Case - Unable to check identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("failed to check identifier exists")
				}
			}
			if tt.name == "Sad Case - identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
					return true, nil
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
			if tt.name == "Sad Case - unable to invite staff" {

				fakeUser.MockInviteUserFn = func(ctx context.Context, userID, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
					return false, fmt.Errorf("failed to invite user")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user id" {

				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get user profile by user id")
				}
			}
			if tt.name == "Sad case: failed to register user in matrix" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := us.CreateSuperUser(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.CreateSuperUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did not expect nil, got %v", got)
			}
		})
	}
}

func TestUseCasesUserImpl_CheckIdentifierExists(t *testing.T) {
	type args struct {
		ctx             context.Context
		identifierType  enums.UserIdentifierType
		identifierValue string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check if identifier exists",
			args: args{
				ctx:             context.Background(),
				identifierType:  enums.UserIdentifierTypeCCC,
				identifierValue: gofakeit.UUID(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: unable to check identifier exists",
			args: args{
				ctx:             context.Background(),
				identifierType:  enums.UserIdentifierTypeCCC,
				identifierValue: gofakeit.UUID(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to check identifier exists" {
				fakeDB.MockCheckIdentifierExists = func(ctx context.Context, identifierType enums.UserIdentifierType, identifierValue string) (bool, error) {
					return false, fmt.Errorf("unable to check identifier exists")
				}
			}
			got, err := us.CheckIdentifierExists(tt.args.ctx, tt.args.identifierType, tt.args.identifierValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.CheckIdentifierExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.CheckIdentifierExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_CheckIfPhoneExists(t *testing.T) {
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check if phone exists",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Sad case: invalid phone",
			args: args{
				ctx:         context.Background(),
				phoneNumber: "invalid",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to check if phone exists",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: failed to check if phone exists" {
				fakeDB.MockCheckPhoneExistsFn = func(ctx context.Context, phone string) (bool, error) {
					return false, fmt.Errorf("unable to check if phone exists")
				}
			}
			got, err := us.CheckIfPhoneExists(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.CheckIfPhoneExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.CheckIfPhoneExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_UpdateOrganisationAdminPermission(t *testing.T) {
	type args struct {
		ctx                 context.Context
		staffID             string
		isOrganisationAdmin bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: update is organisation admin",
			args: args{
				ctx:                 context.Background(),
				staffID:             gofakeit.UUID(),
				isOrganisationAdmin: false,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: failed get staff profile",
			args: args{
				ctx:                 context.Background(),
				staffID:             gofakeit.UUID(),
				isOrganisationAdmin: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed update staff profile",
			args: args{
				ctx:                 context.Background(),
				staffID:             gofakeit.UUID(),
				isOrganisationAdmin: false,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: failed get staff profile" {
				fakeDB.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed update staff profile" {
				fakeDB.MockUpdateStaffFn = func(ctx context.Context, staff *domain.StaffProfile, updates map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := us.UpdateOrganisationAdminPermission(tt.args.ctx, tt.args.staffID, tt.args.isOrganisationAdmin)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.UpdateOrganisationAdminPermission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.UpdateOrganisationAdminPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_NotifyNewFacilityAdded(t *testing.T) {
	type args struct {
		ctx                context.Context
		assignedFacilities []string
		userProfile        *domain.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: notify new single facility added",
			args: args{
				ctx:                context.Background(),
				assignedFacilities: []string{gofakeit.UUID()},
				userProfile: &domain.User{
					Username: gofakeit.BeerName(),
					Contacts: &domain.Contact{
						ContactValue: interserviceclient.TestUserPhoneNumber,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: notify new multiple facility added",
			args: args{
				ctx:                context.Background(),
				assignedFacilities: []string{gofakeit.UUID(), gofakeit.UUID()},
				userProfile: &domain.User{
					Username: gofakeit.BeerName(),
					Contacts: &domain.Contact{
						ContactValue: interserviceclient.TestUserPhoneNumber,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to send SMS",
			args: args{
				ctx:                context.Background(),
				assignedFacilities: []string{gofakeit.UUID(), gofakeit.UUID()},
				userProfile: &domain.User{
					Username: gofakeit.BeerName(),
					Contacts: &domain.Contact{
						ContactValue: interserviceclient.TestUserPhoneNumber,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to send SMS" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if err := us.NotifyNewFacilityAdded(tt.args.ctx, tt.args.assignedFacilities, tt.args.userProfile); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.NotifyNewFacilityAdded() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesUserImpl_DeleteClientProfile(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx      context.Context
		clientID string
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
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable to get client profile",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get client profiles",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get staff profiles",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to get contact by user ID",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - unable to delete FHIR patient profile",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable to delete cms client via pub sub",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - unable to deactivate matrix user",
			args: args{
				ctx:      ctx,
				clientID: gofakeit.UUID(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Happy Case - Successfully delete client" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
					input := dto.BasicUserInput{
						Username: gofakeit.Word(),
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

			if tt.name == "Sad Case - unable to get client profile" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad Case - unable to get client profiles" {
				fakeDB.MockGetUserClientProfilesFn = func(ctx context.Context, userID string) ([]*domain.ClientProfile, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad Case - unable to get staff profiles" {
				fakeDB.MockGetUserStaffProfilesFn = func(ctx context.Context, userID string) ([]*domain.StaffProfile, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad Case - unable to get contact by user ID" {
				fakeDB.MockGetContactByUserIDFn = func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
					return nil, errors.New("an error occurred")
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

			if tt.name == "Sad Case - unable to deactivate matrix user" {
				fakeDB.MockGetUserClientProfilesFn = func(ctx context.Context, userID string) ([]*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return []*domain.ClientProfile{{
						ID:     &id,
						UserID: id,
					}}, nil
				}

				fakeDB.MockGetUserStaffProfilesFn = func(ctx context.Context, userID string) ([]*domain.StaffProfile, error) {
					return []*domain.StaffProfile{}, nil
				}

				fakeDB.MockSearchCaregiverUserFn = func(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
					return []*domain.CaregiverProfile{}, nil
				}

				fakeMatrix.MockDeactivateUserFn = func(ctx context.Context, userID string, auth *domain.MatrixAuth) error {
					return fmt.Errorf("failed to deactivate matrix user")
				}
			}

			got, err := us.DeleteClientProfile(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.DeleteClientProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesUserImpl.DeleteClientProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_ClientSignUp(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *dto.ClientSelfSignUp
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: client sign up",
			args: args{ctx: context.Background(),
				input: &dto.ClientSelfSignUp{
					Username:    "test",
					ClientName:  "test client",
					Gender:      "MALE",
					DateOfBirth: scalarutils.Date{},
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: client sign up unable to sign up",
			args: args{ctx: context.Background(),
				input: &dto.ClientSelfSignUp{
					Username:    "test",
					ClientName:  "test client",
					Gender:      "MALE",
					DateOfBirth: scalarutils.Date{},
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: client sign up unable to sign up" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, fmt.Errorf("error")
				}
			}

			_, err := us.ClientSignUp(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.ClientSignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesUserImpl_Register(t *testing.T) {
	now, err := scalarutils.NewDate(time.Now().Day(), int(time.Now().Month()), time.Now().Year())
	if err != nil {
		t.Errorf("unable to setup date")
		return
	}

	ID := gofakeit.UUID()

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
		PhoneNumber: interserviceclient.TestUserPhoneNumber,
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
		ctx            context.Context
		payload        *dto.SignUpPayload
		selfRegistered bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: register self registering client",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: false,
		},
		{
			name: "Happy case: register invited clients",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: false,
		},
		{
			name: "Sad case: self onboarded: unable to check facility exist by identifier",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: true,
		},
		{
			name: "Sad case: self onboarded: unable to retrieve facility by identifier",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: true,
		},
		{
			name: "Sad case: self onboarded: unable register client",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: true,
		},
		{
			name: "Sad case: self onboarded: unable to create pub sub patient",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: false,
		},
		{
			name: "Sad case: self onboarded: unable to create cms client",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: false,
		},
		{
			name: "Sad case: self onboarded: unable to invite user",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: true,
		},
		{
			name: "Sad case: self onboarded: unable to notify matrix registration",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: true,
		},
		{
			name: "Sad case: self onboarded: unable to create cms client",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: &dto.ClientRegistrationInput{
						Username:       "test1234",
						Facility:       "1234",
						ClientName:     "test 1234",
						Gender:         "MALE",
						DateOfBirth:    *now,
						PhoneNumber:    interserviceclient.TestUserPhoneNumber,
						EnrollmentDate: *now,
					},
				},
				selfRegistered: true,
			},
			wantErr: false,
		},
		{
			name: "Sad case: hcw invited: unable to register client",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: hcw invited: unable to create patient",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: false,
		},
		{
			name: "Sad case: hcw invited: unable to create cms client",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: false,
		},
		{
			name: "Sad case: hcw invited: unable to invite user",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: hcw invited: unable to register matrix user",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: hcw invited: unable to check if username exists",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: hcw invited: fail if username exist",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: false,
			},
			wantErr: true,
		},
		{
			name: "Sad case: self onboarded: fail if facility doesn't exist",
			args: args{
				ctx: context.Background(),
				payload: &dto.SignUpPayload{
					ClientInput: payload,
					UserProfile: &domain.User{
						ID:                    &ID,
						CurrentOrganizationID: ID,
						CurrentProgramID:      ID,
					},
					UserProgram: &domain.Program{
						ID: ID,
						Organisation: domain.Organisation{
							ID: ID,
						},
					},
					Facility: &domain.Facility{
						ID: &ID,
					},
					Matrix: &domain.MatrixAuth{},
				},
				selfRegistered: true,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()
			us := user.NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: self onboarded: unable to check facility exist by identifier" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: self onboarded: unable to retrieve facility by identifier" {
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: self onboarded: unable register client" {
				fakeDB.MockRegisterClientFn = func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: self onboarded: unable to create cms client" {
				fakePubsub.MockNotifyCreateCMSClientFn = func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
					return fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: self onboarded: unable to create pub sub patient" {
				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: self onboarded: unable to invite user" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: self onboarded: unable to notify matrix registration" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to register client" {
				fakeDB.MockRegisterClientFn = func(ctx context.Context, payload *domain.ClientRegistrationPayload) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to create cms client" {
				fakePubsub.MockNotifyCreateCMSClientFn = func(ctx context.Context, user *dto.PubsubCreateCMSClientPayload) error {
					return fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to invite user" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("failed to send sms")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to register matrix user" {
				fakePubsub.MockNotifyRegisterMatrixUserFn = func(ctx context.Context, payload *dto.MatrixUserRegistrationPayload) error {
					return fmt.Errorf("unable to register matrix user")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to create patient" {
				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to create cms client" {
				fakePubsub.MockNotifyCreatePatientFn = func(ctx context.Context, client *dto.PatientCreationOutput) error {
					return fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: hcw invited: unable to check if username exists" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return false, fmt.Errorf("error")
				}
			}
			if tt.name == "Sad case: hcw invited: fail if username exist" {
				fakeDB.MockCheckIfUsernameExistsFn = func(ctx context.Context, username string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad case: self onboarded: fail if facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			_, err := us.Register(tt.args.ctx, tt.args.payload, tt.args.selfRegistered)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesUserImpl.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
