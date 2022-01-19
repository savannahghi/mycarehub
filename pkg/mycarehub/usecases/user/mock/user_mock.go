package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockLoginFn                         func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error)
	MockInviteUserFn                    func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
	MockSavePinFn                       func(ctx context.Context, input dto.PINInput) (bool, error)
	MockVerifyLoginPINFn                func(ctx context.Context, userID string, pin string) (bool, int, error)
	MockSetNickNameFn                   func(ctx context.Context, userID *string, nickname *string) (bool, error)
	MockRequestPINResetFn               func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error)
	MockResetPINFn                      func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockRefreshTokenFn                  func(ctx context.Context, userID string) (*domain.AuthCredentials, error)
	MockVerifyPINFn                     func(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error)
	MockGetClientCaregiverFn            func(ctx context.Context, clientID string) (*domain.Caregiver, error)
	MockCreateOrUpdateClientCaregiverFn func(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error)
}

// NewUserUseCaseMock creates in itializes create type mocks
func NewUserUseCaseMock() *UserUseCaseMock {
	var UUID = uuid.New().String()
	caregiver := &domain.Caregiver{
		ID:            UUID,
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		PhoneNumber:   gofakeit.Phone(),
		CaregiverType: enums.CaregiverTypeFather,
	}

	return &UserUseCaseMock{

		MockLoginFn: func(ctx context.Context, phoneNumber, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error) {
			ID := uuid.New().String()
			time := time.Now()
			return &domain.LoginResponse{
				Client: &domain.ClientProfile{
					ID: &ID,
					User: &domain.User{
						ID:               &ID,
						Username:         gofakeit.Username(),
						TermsAccepted:    true,
						Active:           true,
						NextAllowedLogin: &time,
						FailedLoginCount: 1,
					},
				},
				AuthCredentials: domain.AuthCredentials{
					RefreshToken: gofakeit.HipsterSentence(15),
					IDToken:      gofakeit.BeerAlcohol(),
					ExpiresIn:    gofakeit.BeerHop(),
				},
				Code:    1,
				Message: "Success",
			}, 1, nil
		},
		MockInviteUserFn: func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockSavePinFn: func(ctx context.Context, input dto.PINInput) (bool, error) {
			return true, nil
		},
		MockVerifyLoginPINFn: func(ctx context.Context, userID string, pin string) (bool, int, error) {
			return true, 0, nil
		},
		MockSetNickNameFn: func(ctx context.Context, userID, nickname *string) (bool, error) {
			return true, nil
		},
		MockRequestPINResetFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
			return "111222", nil
		},
		MockResetPINFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockRefreshTokenFn: func(ctx context.Context, userID string) (*domain.AuthCredentials, error) {
			return &domain.AuthCredentials{
				RefreshToken: uuid.New().String(),
				ExpiresIn:    "3600",
				IDToken:      uuid.New().String(),
			}, nil
		},
		MockVerifyPINFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
			return true, nil
		},
		MockGetClientCaregiverFn: func(ctx context.Context, clientID string) (*domain.Caregiver, error) {
			return caregiver, nil
		},
		MockCreateOrUpdateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
			return true, nil
		},
	}
}

// Login mocks the login functionality
func (f *UserUseCaseMock) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, int, error) {
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

// VerifyLoginPIN mocks the verify pin functionality
func (f *UserUseCaseMock) VerifyLoginPIN(ctx context.Context, userID string, pin string) (bool, int, error) {
	return f.MockVerifyLoginPINFn(ctx, userID, pin)
}

// SetNickName is used to mock the implementation ofsetting or changing the user's nickname
func (f *UserUseCaseMock) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	return f.MockSetNickNameFn(ctx, userID, nickname)
}

// RequestPINReset mocks the implementation of requesting pin reset
func (f *UserUseCaseMock) RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
	return f.MockRequestPINResetFn(ctx, phoneNumber, flavour)
}

// ResetPIN mocks the reset pin functionality
func (f *UserUseCaseMock) ResetPIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return f.MockResetPINFn(ctx, userID, flavour)
}

// RefreshToken mocks the implementation for refreshing a token
func (f *UserUseCaseMock) RefreshToken(ctx context.Context, userID string) (*domain.AuthCredentials, error) {
	return f.MockRefreshTokenFn(ctx, userID)
}

// VerifyPIN mocks the implementation for verifying a pin
func (f *UserUseCaseMock) VerifyPIN(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
	return f.MockVerifyPINFn(ctx, userID, flavour, pin)
}

// GetClientCaregiver mocks the implementation for getting the caregiver of a client
func (f *UserUseCaseMock) GetClientCaregiver(ctx context.Context, clientID string) (*domain.Caregiver, error) {
	return f.MockGetClientCaregiverFn(ctx, clientID)
}

// CreateOrUpdateClientCaregiver mocks the implementation for creating or updating the caregiver of a client
func (f *UserUseCaseMock) CreateOrUpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
	return f.MockCreateOrUpdateClientCaregiverFn(ctx, caregiverInput)
}
