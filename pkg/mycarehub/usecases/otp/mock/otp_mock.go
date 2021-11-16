package mock

import (
	"context"

	"github.com/savannahghi/feedlib"
)

// OTPUseCaseMock mocks the implementation of OTP usecase
type OTPUseCaseMock struct {
	MockGenerateAndSendOTPFn func(
		ctx context.Context,
		userID string,
		phoneNumber string,
		flavour feedlib.Flavour,
	) (string, error)
}

// NewOTPUseCaseMock initializes a new instance mock of the OTP usecase
func NewOTPUseCaseMock() *OTPUseCaseMock {
	return &OTPUseCaseMock{
		MockGenerateAndSendOTPFn: func(
			ctx context.Context,
			userID string,
			phoneNumber string,
			flavour feedlib.Flavour,
		) (string, error) {
			return "111222", nil
		},
	}
}

// GenerateAndSendOTP mocks the generate and send OTP method
func (o *OTPUseCaseMock) GenerateAndSendOTP(
	ctx context.Context,
	userID string,
	phoneNumber string,
	flavour feedlib.Flavour,
) (string, error) {
	return o.MockGenerateAndSendOTPFn(ctx, userID, phoneNumber, flavour)
}
