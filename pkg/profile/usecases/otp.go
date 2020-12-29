package usecases

import "context"

// OTPUseCases represent the business logic required for management of OTP
type OTPUseCases interface {
	SendRetryOTP(ctx context.Context, msisdn string, retryStep int) (string, error)
	GenerateAndSendOTP(phone string) (string, error)
	// TODO consider moving this to OTP service or making an isc for it
	VerifyEmailOTP(ctx context.Context, msisdn, otp, flavour string) (bool, error)
}
