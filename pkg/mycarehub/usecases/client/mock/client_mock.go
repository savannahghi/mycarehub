package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ClientUsecaseMock initializes all the Client Usecase methods
type ClientUsecaseMock struct {
	MockRegisterClientFn func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error)
	MockSavePinFn        func(ctx context.Context, pinData *domain.UserPIN) (bool, error)
}

// NewClientUsecaseMock creates in itializes create type mocks
func NewClientUsecaseMock() *ClientUsecaseMock {
	ID := uuid.New().String()

	testTime := time.Now()

	clientProfile := &domain.ClientUserProfile{
		User: &domain.User{
			ID:                  &ID,
			FirstName:           gofakeit.FirstName(),
			LastName:            gofakeit.LastName(),
			Username:            gofakeit.Username(),
			MiddleName:          gofakeit.BeerAlcohol(),
			DisplayName:         gofakeit.BeerHop(),
			Gender:              enumutils.GenderMale,
			Active:              true,
			LastSuccessfulLogin: &testTime,
			LastFailedLogin:     &testTime,
			NextAllowedLogin:    &testTime,
			TermsAccepted:       true,
			AcceptedTermsID:     ID,
		},
		Client: &domain.ClientProfile{
			ID:             &ID,
			UserID:         &ID,
			ClientType:     enums.ClientTypeOvc,
			HealthRecordID: &ID,
		},
	}

	return &ClientUsecaseMock{
		MockRegisterClientFn: func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
			return clientProfile, nil
		},
		MockSavePinFn: func(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
			return true, nil
		},
	}
}

// RegisterClient mocks the implementation of `gorm's` RegisterClient method
func (f *ClientUsecaseMock) RegisterClient(
	ctx context.Context,
	userInput *dto.UserInput,
	clientInput *dto.ClientProfileInput,
) (*domain.ClientUserProfile, error) {
	return f.MockRegisterClientFn(ctx, userInput, clientInput)
}

// SavePin mocks the save pin implementation
func (f *ClientUsecaseMock) SavePin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	return f.MockSavePinFn(ctx, pinData)
}
