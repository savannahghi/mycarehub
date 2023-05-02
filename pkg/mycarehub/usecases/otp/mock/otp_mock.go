package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/profileutils"
)

// OTPUseCaseMock mocks the implementation of OTP usecase
type OTPUseCaseMock struct {
	MockGenerateAndSendOTPFn func(
		ctx context.Context,
		phoneNumber string,
		flavour feedlib.Flavour,
	) (*domain.OTPResponse, error)
	MockVerifyPhoneNumberFn func(ctx context.Context, username string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error)
	MockGenerateRetryOTPFn  func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
	MockSendOTPFn           func(ctx context.Context, phoneNumber string, code string, message string) (string, error)
	MockVerifyOTP           func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
}

// NewOTPUseCaseMock initializes a new instance mock of the OTP usecase
func NewOTPUseCaseMock() *OTPUseCaseMock {
	return &OTPUseCaseMock{
		MockGenerateAndSendOTPFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.OTPResponse, error) {
			return &domain.OTPResponse{
				OTP:         "111222",
				PhoneNumber: interserviceclient.TestUserPhoneNumber,
			}, nil
		},
		MockVerifyPhoneNumberFn: func(ctx context.Context, username string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
			return &profileutils.OtpResponse{
				OTP: "111222",
			}, nil
		},
		MockGenerateRetryOTPFn: func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
			return "test-OTP", nil
		},
		MockSendOTPFn: func(ctx context.Context, phoneNumber string, code string, message string) (string, error) {
			return "111222", nil
		},
		MockVerifyOTP: func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
			return true, nil
		},
	}
}

// GenerateAndSendOTP mocks the generate and send OTP method
func (o *OTPUseCaseMock) GenerateAndSendOTP(
	ctx context.Context,
	phoneNumber string,
	flavour feedlib.Flavour,
) (*domain.OTPResponse, error) {
	return o.MockGenerateAndSendOTPFn(ctx, phoneNumber, flavour)
}

// VerifyPhoneNumber mock the implementtation of phone verification
func (o *OTPUseCaseMock) VerifyPhoneNumber(ctx context.Context, username string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
	return o.MockVerifyPhoneNumberFn(ctx, username, flavour)
}

// GenerateRetryOTP mock the implementtation of generating a retry OTP
func (o *OTPUseCaseMock) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	return o.MockGenerateRetryOTPFn(ctx, payload)
}

// SendOTP mocks the implementation of sending an OTP
func (o *OTPUseCaseMock) SendOTP(ctx context.Context, phoneNumber string, code string, message string) (string, error) {
	return o.MockSendOTPFn(ctx, phoneNumber, code, message)
}

// VerifyOTP mocks the implementation of verifying an OTP
func (o *OTPUseCaseMock) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return o.MockVerifyOTP(ctx, payload)
}
