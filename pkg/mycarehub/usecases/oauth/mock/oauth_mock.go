package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// OauthUseCaseMock mocks the implementation of oauth usecase
type OauthUseCaseMock struct {
	MockCreateOauthClientFn func(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error)
}

// NewOauthUseCaseMock initializes a new instance mock of the oauth usecase
func NewOauthUseCaseMock() *OauthUseCaseMock {

	return &OauthUseCaseMock{
		MockCreateOauthClientFn: func(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error) {
			return &domain.OauthClient{
				ID: gofakeit.UUID(),
			}, nil
		},
	}
}

// CreateOauthClient is the resolver for the createOauthClient field.
func (u *OauthUseCaseMock) CreateOauthClient(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error) {
	return u.MockCreateOauthClientFn(ctx, input)
}
