package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// LoginUseCases  represents all the business logic involved in logging in a user and managing their authorization credentials.
type LoginUseCases interface {
	LoginByPhone(ctx context.Context, phone, pin, flavour string) (*domain.UserResponse, error)
	RefreshToken(ctx context.Context, token string) (*domain.AuthCredentialResponse, error)
}
