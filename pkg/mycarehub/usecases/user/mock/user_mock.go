package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockGetUserProfileByPhoneNumberFn func(ctx context.Context, phoneNumber string) (*domain.User, error)
	MockSetUserPINFn                  func(ctx context.Context, input *dto.PinInput) (bool, error)
	MockLoginFn                       func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error)
}

// NewUserUseCaseMock creates in itializes create type mocks
func NewUserUseCaseMock() *UserUseCaseMock {
	return &UserUseCaseMock{
		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*domain.User, error) {
			id := uuid.New().String()
			return &domain.User{
				ID: &id,
			}, nil
		},

		MockSetUserPINFn: func(ctx context.Context, input *dto.PinInput) (bool, error) {
			return true, nil
		},

		MockLoginFn: func(ctx context.Context, phoneNumber, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error) {
			return &domain.AuthCredentials{
				User: &domain.User{
					Username: gofakeit.Username(),
				},
				RefreshToken: gofakeit.HipsterSentence(15),
				IDToken:      gofakeit.BeerAlcohol(),
				ExpiresIn:    gofakeit.BeerHop(),
			}, "", nil
		},
	}
}

// GetUserProfileByPhoneNumber mocks the implementation of fetching a user profile by phonenumber
func (f *UserUseCaseMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return f.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// SetUserPIN ...
func (f *UserUseCaseMock) SetUserPIN(ctx context.Context, input *dto.PinInput) (bool, error) {
	return f.MockSetUserPINFn(ctx, input)
}

// Login ...
func (f *UserUseCaseMock) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, string, error) {
	return f.MockLoginFn(ctx, phoneNumber, pin, flavour)
}
