package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/ory/fosite"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/oauth"
)

// OauthUseCaseMock mocks the implementation of oauth usecase
type OauthUseCaseMock struct {
	MockCreateOauthClientFn      func(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error)
	MockFositeProviderFn         func() fosite.OAuth2Provider
	MockGenerateUserAuthTokensFn func(ctx context.Context, userID string) (*oauth.AuthTokens, error)
	MockRefreshAutTokenFn        func(ctx context.Context, refreshToken string) (*oauth.AuthTokens, error)
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
		MockGenerateUserAuthTokensFn: func(ctx context.Context, userID string) (*oauth.AuthTokens, error) {
			return &oauth.AuthTokens{
				AccessToken:  "access",
				ExpiresIn:    3600,
				RefreshToken: "refresh",
			}, nil
		},
		MockRefreshAutTokenFn: func(ctx context.Context, refreshToken string) (*oauth.AuthTokens, error) {
			return &oauth.AuthTokens{
				AccessToken:  "access",
				ExpiresIn:    3600,
				RefreshToken: "refresh",
			}, nil
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

func (u *OauthUseCaseMock) GenerateUserAuthTokens(ctx context.Context, userID string) (*oauth.AuthTokens, error) {
	return u.MockGenerateUserAuthTokensFn(ctx, userID)
}

// RefreshAutToken mocks the implementation of RefreshAutToken method
func (u *OauthUseCaseMock) RefreshAutToken(ctx context.Context, refreshToken string) (*oauth.AuthTokens, error) {
	return u.MockRefreshAutTokenFn(ctx, refreshToken)
}
