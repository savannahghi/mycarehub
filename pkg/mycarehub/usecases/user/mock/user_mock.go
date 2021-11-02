package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockGetUserProfileByPhoneNumberFn func(ctx context.Context, phoneNumber string) (*domain.User, error)
	SetUserPINFn                      func(ctx context.Context, input *dto.PinInput) (bool, error)
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

		SetUserPINFn: func(ctx context.Context, input *dto.PinInput) (bool, error) {
			return true, nil
		},
	}
}

// GetUserProfileByPhoneNumber mocks the implementation of fetching a user profile by phonenumber
func (f *UserUseCaseMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return f.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// SetUserPIN ...
func (f *UserUseCaseMock) SetUserPIN(ctx context.Context, input *dto.PinInput) (bool, error) {
	return f.SetUserPINFn(ctx, input)
}
