package oauth

import (
	"context"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"golang.org/x/crypto/bcrypt"

	"github.com/ory/fosite"
)

// ClientAssertionJWTValid returns an error if the JTI is known or the DB check failed
// and nil if the JTI is not known.
func (u UseCasesOauthImpl) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	clientJWT, err := u.Query.GetClientJWT(ctx, jti)
	if err != nil {
		return err
	}

	if clientJWT.ExpiresAt.After(time.Now()) {
		return fosite.ErrJTIKnown
	}

	return nil
}

// GetClient loads the client by its ID or returns an error
// if the client does not exist or another error occurred.
func (u UseCasesOauthImpl) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	client, err := u.Query.GetOauthClient(ctx, id)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// SetClientAssertionJWT marks a JTI as known for the given expiry time.
// Before inserting the new JTI, it will clean up any existing JTIs that have expired as those tokens cannot be replayed due to the expiry.
func (u UseCasesOauthImpl) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	_, err := u.Query.GetValidClientJWT(ctx, jti)
	if err != nil {
		return err
	}

	jwt := &domain.OauthClientJWT{
		Active:    true,
		JTI:       jti,
		ExpiresAt: exp,
	}

	err = u.Create.CreateOauthClientJWT(ctx, jwt)
	if err != nil {
		return err
	}

	return nil
}

// CreateOauthClient is the resolver for the createOauthClient field.
func (u UseCasesOauthImpl) CreateOauthClient(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error) {
	secret, err := bcrypt.GenerateFromPassword([]byte(input.Secret), fosite.DefaultBCryptWorkFactor)
	if err != nil {
		return nil, err
	}

	client := &domain.OauthClient{
		Name:   input.Name,
		Secret: string(secret),
	}

	err = u.Create.CreateOauthClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ListOauthClients is the resolver for the listOauthClients field.
func (u UseCasesOauthImpl) ListOauthClients(ctx context.Context) ([]*domain.OauthClient, error) {
	return nil, nil
}
