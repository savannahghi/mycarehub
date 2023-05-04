package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ory/fosite"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
)

func (s Storage) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	session := request.GetSession().(*domain.Session)

	err = s.Create.CreateOrUpdateSession(ctx, session)
	if err != nil {
		return err
	}

	client := request.GetClient()

	data := &domain.RefreshToken{
		ID:                request.GetID(),
		Active:            true,
		Signature:         signature,
		RequestedAt:       request.GetRequestedAt(),
		ClientID:          client.GetID(),
		RequestedScopes:   request.GetRequestedScopes(),
		GrantedScopes:     request.GetGrantedScopes(),
		Form:              request.GetRequestForm(),
		SessionID:         session.ID,
		RequestedAudience: request.GetRequestedAudience(),
		GrantedAudience:   request.GetGrantedAudience(),
	}

	err = s.Create.CreateRefreshToken(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	refreshToken, err := s.Query.GetRefreshToken(ctx, domain.RefreshToken{Signature: signature})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	rq := &fosite.Request{
		ID:                refreshToken.ID,
		RequestedAt:       refreshToken.RequestedAt,
		Client:            refreshToken.Client,
		RequestedScope:    fosite.Arguments(refreshToken.RequestedScopes),
		GrantedScope:      fosite.Arguments(refreshToken.GrantedScopes),
		Form:              refreshToken.Form,
		Session:           &refreshToken.Session,
		RequestedAudience: fosite.Arguments(refreshToken.RequestedAudience),
		GrantedAudience:   fosite.Arguments(refreshToken.GrantedAudience),
	}

	if !refreshToken.Active {
		return rq, fosite.ErrInactiveToken
	}

	return rq, nil
}

func (s Storage) DeleteRefreshTokenSession(ctx context.Context, signature string) (err error) {
	return s.Delete.DeleteRefreshToken(ctx, signature)
}

// RevokeRefreshToken revokes a refresh token as specified in:
// https://tools.ietf.org/html/rfc7009#section-2.1
// If the particular
// token is a refresh token and the authorization server supports the
// revocation of access tokens, then the authorization server SHOULD
// also invalidate all access tokens based on the same authorization
// grant (see Implementation Note).
func (s Storage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	refreshToken, err := s.Query.GetRefreshToken(ctx, domain.RefreshToken{ID: requestID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err
	}

	if err := s.Update.UpdateRefreshToken(ctx, refreshToken, map[string]interface{}{"active": false}); err != nil {
		return fmt.Errorf("failed to invalidate authorization code: %w", err)
	}

	return nil
}

// RevokeRefreshTokenMaybeGracePeriod revokes a refresh token as specified in:
// https://tools.ietf.org/html/rfc7009#section-2.1
// If the particular
// token is a refresh token and the authorization server supports the
// revocation of access tokens, then the authorization server SHOULD
// also invalidate all access tokens based on the same authorization
// grant (see Implementation Note).
//
// If the Refresh Token grace period is greater than zero in configuration the token
// will have its expiration time set as UTCNow + GracePeriod.
func (s Storage) RevokeRefreshTokenMaybeGracePeriod(ctx context.Context, requestID string, signature string) error {
	return s.RevokeRefreshToken(ctx, requestID)
}
