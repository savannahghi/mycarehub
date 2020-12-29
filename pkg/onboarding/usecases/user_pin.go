package usecases

import "context"

// UserPINUseCases represents all the business logic that touch on user PIN Management
type UserPINUseCases interface {
	SetUserPIN(ctx context.Context, phone string, pin string) (bool, error)
	ChangeUserPIN(ctx context.Context, phone string, pin string, otp string) (bool, error)
	ResetPIN(ctx context.Context, phone string) (string, error)
}
