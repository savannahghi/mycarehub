package otp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	smsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms/mock"
	twilioMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/silcomms"
	"github.com/segmentio/ksuid"
)

func TestUseCaseOTPImpl_VerifyPhoneNumber(t *testing.T) {
	ctx := context.Background()

	phone := interserviceclient.TestUserPhoneNumber
	badPhone := ksuid.New().String()
	veryBadPhone := gofakeit.HipsterSentence(200)

	type args struct {
		ctx     context.Context
		phone   string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.OtpResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - bad flavour",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: "feedlib.FlavourPro",
			},
			wantErr: true,
		},
		{
			name: "Sad case - bad phone",
			args: args{
				ctx:     ctx,
				phone:   badPhone,
				flavour: "feedlib.FlavourPro",
			},
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - very bad phone and invalid flavour",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: "feedlib.FlavourPro",
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get profile by phone",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get profile by invalid phone number",
			args: args{
				ctx:     ctx,
				phone:   badPhone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get profile with very bad phone number",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get profile with very bad phone number and invalid flavor",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: "feedlib.FlavourPro",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to save otp",
			args: args{
				ctx:     ctx,
				phone:   interserviceclient.TestUserPhoneNumber,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - fail to send SMS",
			args: args{
				ctx:     ctx,
				phone:   interserviceclient.TestUserPhoneNumber,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Nonexistent phone",
			args: args{
				ctx:     ctx,
				phone:   interserviceclient.TestUserPhoneNumber,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			o := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case - Fail to save otp" {
				fakeDB.MockSaveOTPFn = func(ctx context.Context, otpInput *domain.OTP) error {
					return fmt.Errorf("failed to save otp")
				}
			}

			if tt.name == "Sad Case - fail to send SMS" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Nonexistent phone" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad flavour" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad phone" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - very bad phone and invalid flavour" {
				fakeDB.MockCheckIfPhoneNumberExistsFn = func(ctx context.Context, phone string, isOptedIn bool, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to get profile by phone" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to get profile by invalid phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to get profile with very bad phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to get profile with very bad phone number and invalid flavor" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := o.VerifyPhoneNumber(tt.args.ctx, tt.args.phone, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.VerifyPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_VerifyOTP_Unittest(t *testing.T) {
	ctx := context.Background()

	flavour := feedlib.FlavourConsumer

	validOTPPayload := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload2 := &dto.VerifyOTPInput{
		PhoneNumber: "",
		OTP:         uuid.New().String(),
		Flavour:     flavour,
	}
	invalidOTPPayload3 := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         "",
		Flavour:     flavour,
	}
	invalidOTPPayload4 := &dto.VerifyOTPInput{
		PhoneNumber: uuid.New().String(),
		OTP:         uuid.New().String(),
		Flavour:     "flavour",
	}
	invalidOTPPayload5 := &dto.VerifyOTPInput{
		PhoneNumber: "otpInput.PhoneNumber",
		OTP:         "otpInput.OTP",
		Flavour:     "flavour",
	}
	invalidOTPPayload6 := &dto.VerifyOTPInput{
		PhoneNumber: gofakeit.HipsterParagraph(1, 10, 100, ""),
		OTP:         gofakeit.HipsterParagraph(1, 10, 100, ""),
		Flavour:     "gofakeit.HipsterParagraph(300, 10, 100)",
	}

	type args struct {
		ctx     context.Context
		payload *dto.VerifyOTPInput
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
				payload: validOTPPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no phone",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload2,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no otp",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload3,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad flavour",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload4,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - bad inputs",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload5,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - extreme bad inputs",
			args: args{
				ctx:     ctx,
				payload: invalidOTPPayload6,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			o := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension, fakeSMS, fakeTwilio)

			if tt.name == "Sad case - no user ID" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no phone" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no otp" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad flavour" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad inputs" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - extreme bad inputs" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := o.VerifyOTP(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.VerifyOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.VerifyOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseOTPImpl_GenerateRetryOTP(t *testing.T) {
	ctx := context.Background()

	validPayload := &dto.SendRetryOTPPayload{
		Username: gofakeit.Name(),
		Flavour:  feedlib.FlavourConsumer,
	}

	invalidPayload := &dto.SendRetryOTPPayload{
		Username: "",
		Flavour:  feedlib.FlavourConsumer,
	}

	type args struct {
		ctx     context.Context
		payload *dto.SendRetryOTPPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				payload: validPayload,
			},
			wantErr: false,
		},
		{
			name: "Sad case - unable to get phone",
			args: args{
				ctx:     ctx,
				payload: invalidPayload,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to get user profile",
			args: args{
				ctx:     ctx,
				payload: validPayload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to save otp",
			args: args{
				ctx:     ctx,
				payload: validPayload,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - unable to send SMS",
			args: args{
				ctx:     ctx,
				payload: validPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			o := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension, fakeSMS, fakeTwilio)

			if tt.name == "Sad case - unable to get phone" {
				fakeDB.MockGetContactByUserIDFn = func(ctx context.Context, userID *string, contactType string) (*domain.Contact, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to get user profile" {
				fakeDB.MockGetUserProfileByUsernameFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case - Fail to save otp" {
				fakeDB.MockSaveOTPFn = func(ctx context.Context, otpInput *domain.OTP) error {
					return fmt.Errorf("failed to save otp")
				}
			}
			if tt.name == "Sad Case - unable to send SMS" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := o.GenerateRetryOTP(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.GenerateRetryOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseOTPImpl_SendOTP(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		phoneNumber string
		code        string
		message     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully send an otp to kenyan number",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				code:        "111222",
				message:     gofakeit.HipsterSentence(5),
			},
			want:    "111222",
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to send an otp to kenyan number",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				code:        "111222",
				message:     gofakeit.HipsterSentence(5),
			},
			wantErr: true,
		},
		{
			name: "Happy Case - Successfully send an otp to foreign number",
			args: args{
				ctx:         ctx,
				phoneNumber: "+14049370053",
				code:        "111222",
				message:     gofakeit.HipsterSentence(5),
			},
			want:    "111222",
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to send an otp to foreign number",
			args: args{
				ctx:         ctx,
				phoneNumber: "+14049370053",
				code:        "111222",
				message:     gofakeit.HipsterSentence(5),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()

			o := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension, fakeSMS, fakeTwilio)

			if tt.name == "Sad Case - Fail to send an otp to kenyan number" {
				fakeSMS.MockSendSMSFn = func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad Case - Fail to send an otp to foreign number" {
				fakeTwilio.MockSendSMSViaTwilioFn = func(ctx context.Context, phonenumber, message string) error {
					return fmt.Errorf("failed to send sms")
				}
			}

			got, err := o.SendOTP(tt.args.ctx, tt.args.phoneNumber, tt.args.code, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.SendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseOTPImpl.SendOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
