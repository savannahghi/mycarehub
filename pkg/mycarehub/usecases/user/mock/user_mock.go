package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockLoginFn       func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, int, error)
	MockInviteUserFn  func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
	MockSavePinFn     func(ctx context.Context, input dto.PINInput) (bool, error)
	MockVerifyPINFn   func(ctx context.Context, userID string, pin string) (bool, error)
	MockSetNickNameFn func(ctx context.Context, userID *string, nickname *string) (bool, error)
}

// NewUserUseCaseMock creates in itializes create type mocks
func NewUserUseCaseMock() *UserUseCaseMock {
	return &UserUseCaseMock{

		MockLoginFn: func(ctx context.Context, phoneNumber, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, int, error) {
			ID := uuid.New().String()
			time := time.Now()
			return &domain.AuthCredentials{
				User: &domain.User{
					ID:               &ID,
					Username:         gofakeit.Username(),
					TermsAccepted:    true,
					Active:           true,
					NextAllowedLogin: &time,
					FailedLoginCount: 1,
				},
				RefreshToken: gofakeit.HipsterSentence(15),
				IDToken:      gofakeit.BeerAlcohol(),
				ExpiresIn:    gofakeit.BeerHop(),
			}, 1, nil
		},
		MockInviteUserFn: func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockSavePinFn: func(ctx context.Context, input dto.PINInput) (bool, error) {
			return true, nil
		},
		MockVerifyPINFn: func(ctx context.Context, userID string, pin string) (bool, error) {
			return true, nil
		},
		MockSetNickNameFn: func(ctx context.Context, userID, nickname *string) (bool, error) {
			return true, nil
		},
	}
}

// Login mocks the login functionality
func (f *UserUseCaseMock) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.AuthCredentials, int, error) {
	return f.MockLoginFn(ctx, phoneNumber, pin, flavour)
}

// InviteUser mocks the invite functionality
func (f *UserUseCaseMock) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	return f.MockInviteUserFn(ctx, userID, phoneNumber, flavour)
}

// SavePin mocks the save pin functionality
func (f *UserUseCaseMock) SavePin(ctx context.Context, input dto.PINInput) (bool, error) {
	return f.MockSavePinFn(ctx, input)
}

// VerifyPIN mocks the verify pin functionality
func (f *UserUseCaseMock) VerifyPIN(ctx context.Context, userID string, pin string) (bool, error) {
	return f.MockVerifyPINFn(ctx, userID, pin)
}

// SetNickName is used to mock the implementation ofsetting or changing the user's nickname
func (f *UserUseCaseMock) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	return f.MockSetNickNameFn(ctx, userID, nickname)
}
