package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
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
