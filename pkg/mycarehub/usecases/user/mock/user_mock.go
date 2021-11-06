package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockLoginFn func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error)
}

// NewUserUseCaseMock creates in itializes create type mocks
func NewUserUseCaseMock() *UserUseCaseMock {
	return &UserUseCaseMock{
		MockLoginFn: func(ctx context.Context, phoneNumber, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error) {
			ID := uuid.New().String()
			return &domain.AuthCredentials{
				User: &domain.User{
					ID:       &ID,
					Username: gofakeit.Username(),
				},
				RefreshToken: gofakeit.HipsterSentence(15),
				IDToken:      gofakeit.BeerAlcohol(),
				ExpiresIn:    gofakeit.BeerHop(),
			}, "", nil
		},
	}
}

// Login mocks the login functionality
func (f *UserUseCaseMock) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error) {
	return f.MockLoginFn(ctx, phoneNumber, pin, flavour)
}
