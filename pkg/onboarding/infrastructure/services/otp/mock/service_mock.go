package mock

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceOTP is an `OTP` service mock .
type FakeServiceOTP struct {
	GenerateAndSendOTPFn func(ctx context.Context, phone string) (*base.OtpResponse, error)
	SendRetryOTPFn       func(ctx context.Context, msisdn string, retryStep int) (*base.OtpResponse, error)
	VerifyOTPFn          func(ctx context.Context, phone, OTP string) (bool, error)
	VerifyEmailOTPFn     func(ctx context.Context, email, OTP string) (bool, error)
}

// GenerateAndSendOTP ...
func (f *FakeServiceOTP) GenerateAndSendOTP(ctx context.Context, phone string) (*base.OtpResponse, error) {
	return f.GenerateAndSendOTPFn(ctx, phone)
}

// SendRetryOTP ...
func (f *FakeServiceOTP) SendRetryOTP(ctx context.Context, msisdn string, retryStep int) (*base.OtpResponse, error) {
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
