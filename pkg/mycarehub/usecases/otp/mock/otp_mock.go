package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/profileutils"
)

// OTPUseCaseMock mocks the implementation of OTP usecase
type OTPUseCaseMock struct {
	MockGenerateAndSendOTPFn func(
		ctx context.Context,
		phoneNumber string,
		flavour feedlib.Flavour,
	) (string, error)
	MockVerifyPhoneNumberFn func(ctx context.Context, phone *string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error)
	MockGenerateOTPFn       func(ctx context.Context) (string, error)
	MockGenerateRetryOTPFn  func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
}

// NewOTPUseCaseMock initializes a new instance mock of the OTP usecase
func NewOTPUseCaseMock() *OTPUseCaseMock {
	return &OTPUseCaseMock{
		MockGenerateAndSendOTPFn: func(
			ctx context.Context,
			phoneNumber string,
			flavour feedlib.Flavour,
		) (string, error) {
			return "111222", nil
		},
		MockVerifyPhoneNumberFn: func(ctx context.Context, phone *string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
			return &profileutils.OtpResponse{
				OTP: "1234",
			}, nil
		},
		MockGenerateOTPFn: func(ctx context.Context) (string, error) {
			return "111222", nil
		},
		MockGenerateRetryOTPFn: func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
			return "test-OTP", nil
		},
	}
}

// GenerateAndSendOTP mocks the generate and send OTP method
func (o *OTPUseCaseMock) GenerateAndSendOTP(
	ctx context.Context,
	phoneNumber string,
	flavour feedlib.Flavour,
) (string, error) {
	return o.MockGenerateAndSendOTPFn(ctx, phoneNumber, flavour)
}

// VerifyPhoneNumber mock the implementtation of phone verification
func (o *OTPUseCaseMock) VerifyPhoneNumber(ctx context.Context, phone *string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
	return o.MockVerifyPhoneNumberFn(ctx, phone, flavour)
}

// GenerateOTP mocks the generate otp method
func (o *OTPUseCaseMock) GenerateOTP(ctx context.Context) (string, error) {
	return o.MockGenerateOTPFn(ctx)
}

// GenerateRetryOTP mock the implementtation of generating a retry OTP
func (o *OTPUseCaseMock) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	return o.MockGenerateRetryOTPFn(ctx, payload)
}
