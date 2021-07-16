package mock

import (
	"context"

	"github.com/savannahghi/profileutils"
)

// FakeServiceOTP is an `OTP` service mock .
type FakeServiceOTP struct {
	GenerateAndSendOTPFn func(ctx context.Context, phone string) (*profileutils.OtpResponse, error)
	SendRetryOTPFn       func(ctx context.Context, msisdn string, retryStep int) (*profileutils.OtpResponse, error)
	VerifyOTPFn          func(ctx context.Context, phone, OTP string) (bool, error)
	VerifyEmailOTPFn     func(ctx context.Context, email, OTP string) (bool, error)
}

// GenerateAndSendOTP ...
func (f *FakeServiceOTP) GenerateAndSendOTP(ctx context.Context, phone string) (*profileutils.OtpResponse, error) {
	return f.GenerateAndSendOTPFn(ctx, phone)
}

// SendRetryOTP ...
func (f *FakeServiceOTP) SendRetryOTP(ctx context.Context, msisdn string, retryStep int) (*profileutils.OtpResponse, error) {
	return f.SendRetryOTPFn(ctx, msisdn, retryStep)
}

// VerifyOTP ...
func (f *FakeServiceOTP) VerifyOTP(ctx context.Context, phone, OTP string) (bool, error) {
	return f.VerifyOTPFn(ctx, phone, OTP)
}

// VerifyEmailOTP ...
func (f *FakeServiceOTP) VerifyEmailOTP(ctx context.Context, phone, OTP string) (bool, error) {
	return f.VerifyEmailOTPFn(ctx, phone, OTP)
}
