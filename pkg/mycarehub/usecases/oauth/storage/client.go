package storage

import (
	"context"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"

	"github.com/ory/fosite"
)

// ClientAssertionJWTValid returns an error if the JTI is known or the DB check failed
// and nil if the JTI is not known.
func (s Storage) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	clientJWT, err := s.Query.GetClientJWT(ctx, jti)
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
func (s Storage) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	client, err := s.Query.GetOauthClient(ctx, id)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// SetClientAssertionJWT marks a JTI as known for the given expiry time.
// Before inserting the new JTI, it will clean up any existing JTIs that have expired as those tokens cannot be replayed due to the expiry.
func (s Storage) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	_, err := s.Query.GetValidClientJWT(ctx, jti)
	if err != nil {
		return err
	}

	jwt := &domain.OauthClientJWT{
		Active:    true,
		JTI:       jti,
		ExpiresAt: exp,
	}

	err = s.Create.CreateOauthClientJWT(ctx, jwt)
	if err != nil {
		return err
	}

	return nil
}
