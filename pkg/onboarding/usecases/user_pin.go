package usecases

import "context"

// UserPINUseCases represents all the business logic that touch on user PIN Management
type UserPINUseCases interface {
	SetUserPin(ctx context.Context, msisdn string, pin string) (bool, error)
	ChangeUserPin(ctx context.Context, msisdn string, pin string, otp string) (bool, error)
	RequestPinReset(ctx context.Context, msisdn string) (string, error)
}
