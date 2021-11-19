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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
)

func TestUseCaseOTPImpl_GenerateAndSendOTP(t *testing.T) {
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
			name: "Happy Case - Successfully generate and send otp",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				flavour:     feedlib.FlavourConsumer,
			},
			want:    "111222",
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to generate and send otp",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by phone number",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
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
			name: "Sad Case - Fail to save OTP",
			args: args{
				ctx:         ctx,
				phoneNumber: "0710000000",
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeOTP := mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad Case - Fail to generate and send otp" {
				fakeExtension.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string) (string, error) {
					return "", fmt.Errorf("failed to generate and send otp")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "invalid: invalid flavour" {
				fakeOTP.MockGenerateAndSendOTPFn = func(
					ctx context.Context,
					userID string,
					phoneNumber string,
					flavour feedlib.Flavour,
				) (string, error) {
					return "", fmt.Errorf("invalid flavour")
				}
			}

			if tt.name == "Sad Case - Fail to save OTP" {
				fakeDB.MockSaveOTPFn = func(ctx context.Context, otpInput *domain.OTP) error {
					return fmt.Errorf("failed to save user pin")
				}
			}

			got, err := otp.GenerateAndSendOTP(tt.args.ctx, tt.args.phoneNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.GenerateAndSendOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseOTPImpl.GenerateAndSendOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseOTPImpl_GenerateOTP(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Happy Case - Successfully generate otp",
			args:    args{ctx: ctx},
			want:    "111222",
			wantErr: false,
		},
		{
			name:    "Sad Case - Fail to generate otp",
			args:    args{ctx: ctx},
			want:    "111222",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeOTP := mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad Case - Fail to generate otp" {
				fakeOTP.MockGenerateAndSendOTPFn = func(
					ctx context.Context,
					userID string,
					phoneNumber string,
					flavour feedlib.Flavour,
				) (string, error) {
					return "", fmt.Errorf("failed to generate otp")
				}
			}

			got, err := otp.GenerateOTP(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseOTPImpl.GenerateOTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
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
			name: "Sad case - unable to send otp",
			args: args{
				ctx:     ctx,
				phone:   phone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to send otp with invalid phone number",
			args: args{
				ctx:     ctx,
				phone:   badPhone,
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case - unable to send otp with very bad phone number",
			args: args{
				ctx:     ctx,
				phone:   veryBadPhone,
				flavour: feedlib.FlavourPro,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewOTPUseCaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)

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
			if tt.name == "Sad case - unable to send otp" {
				fakeExtension.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string) (string, error) {
					return "", fmt.Errorf("failed to generate and send otp")
				}
			}
			if tt.name == "Sad case - unable to send otp with invalid phone number" {
				fakeExtension.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string) (string, error) {
					return "", fmt.Errorf("failed to generate and send otp")
				}
			}
			if tt.name == "Sad case - unable to send otp with very bad phone number" {
				fakeExtension.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string) (string, error) {
					return "", fmt.Errorf("failed to generate and send otp")
				}
			}
			_, err := otp.VerifyPhoneNumber(tt.args.ctx, tt.args.phone, tt.args.flavour)
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
			otp := otp.NewOTPUseCase(fakeDB, fakeDB, fakeExtension)

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

			got, err := otp.VerifyOTP(tt.args.ctx, tt.args.payload)
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
