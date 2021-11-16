package otp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
)

func TestUseCaseOTPImpl_GenerateAndSendOTP(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		userID      string
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
				userID:      "1234",
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
				userID:      "1234",
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
			otp := otp.NewOTPUseCase(fakeDB, fakeExtension)

			if tt.name == "Happy Case - Successfully generate and send otp" {
				fakeOTP.MockGenerateAndSendOTPFn = func(
					ctx context.Context,
					userID string,
					phoneNumber string,
					flavour feedlib.Flavour,
				) (string, error) {
					return "111222", nil
				}
			}

			if tt.name == "Sad Case - Fail to generate and send otp" {
				fakeExtension.MockGenerateAndSendOTPFn = func(ctx context.Context, phoneNumber string) (string, error) {
					return "", fmt.Errorf("failed to generate and send otp")
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

			got, err := otp.GenerateAndSendOTP(tt.args.ctx, tt.args.userID, tt.args.phoneNumber, tt.args.flavour)
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
			otp := otp.NewOTPUseCase(fakeDB, fakeExtension)

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
