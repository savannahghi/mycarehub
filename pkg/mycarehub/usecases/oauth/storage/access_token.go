package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ory/fosite"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
)

func (s Storage) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) (err error) {
	session := request.GetSession().(*domain.Session)

	err = s.Create.CreateOrUpdateSession(ctx, session)
	if err != nil {
		return err
	}

	client := request.GetClient()

	data := &domain.AccessToken{
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

	err = s.Create.CreateAccessToken(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (request fosite.Requester, err error) {
	accessToken, err := s.Query.GetAccessToken(ctx, domain.AccessToken{Signature: signature})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	rq := &fosite.Request{
		ID:                accessToken.ID,
		RequestedAt:       accessToken.RequestedAt,
		Client:            accessToken.Client,
		RequestedScope:    fosite.Arguments(accessToken.RequestedScopes),
		GrantedScope:      fosite.Arguments(accessToken.GrantedScopes),
		Form:              accessToken.Form,
		Session:           &accessToken.Session,
		RequestedAudience: fosite.Arguments(accessToken.RequestedAudience),
		GrantedAudience:   fosite.Arguments(accessToken.GrantedAudience),
	}

	return rq, nil
}

func (s Storage) DeleteAccessTokenSession(ctx context.Context, signature string) (err error) {
	return s.Delete.DeleteAccessToken(ctx, signature)
}

func (s Storage) RevokeAccessToken(ctx context.Context, requestID string) error {
	accessToken, err := s.Query.GetAccessToken(ctx, domain.AccessToken{ID: requestID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err
	}

	if err := s.Update.UpdateAccessToken(ctx, accessToken, map[string]interface{}{"active": false}); err != nil {
		return fmt.Errorf("failed to invalidate authorization code: %w", err)
	}

	return nil
}
