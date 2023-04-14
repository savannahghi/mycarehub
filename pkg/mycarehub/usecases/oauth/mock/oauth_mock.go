package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/ory/fosite"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// OauthUseCaseMock mocks the implementation of oauth usecase
type OauthUseCaseMock struct {
	MockCreateOauthClientFn func(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error)
	MockFositeProviderFn    func() fosite.OAuth2Provider
}

// NewOauthUseCaseMock initializes a new instance mock of the oauth usecase
func NewOauthUseCaseMock() *OauthUseCaseMock {

	return &OauthUseCaseMock{
		MockCreateOauthClientFn: func(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error) {
			return &domain.OauthClient{
				ID: gofakeit.UUID(),
			}, nil
		},
		MockFositeProviderFn: func() fosite.OAuth2Provider {
			return nil
		},
	}
}

// CreateOauthClient is the resolver for the createOauthClient field.
func (u *OauthUseCaseMock) CreateOauthClient(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error) {
	return u.MockCreateOauthClientFn(ctx, input)
}

func (u *OauthUseCaseMock) FositeProvider() fosite.OAuth2Provider {
	return u.MockFositeProviderFn()
}
