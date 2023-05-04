package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/ory/fosite"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
)

// CreateAuthorizeCodeSession stores the authorization request for a given authorization code.
func (s Storage) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) (err error) {
	session := request.GetSession().(*domain.Session)

	err = s.Create.CreateOrUpdateSession(ctx, session)
	if err != nil {
		return err
	}

	client := request.GetClient()

	data := &domain.AuthorizationCode{
		ID:                request.GetID(),
		Active:            true,
		Code:              code,
		RequestedAt:       request.GetRequestedAt(),
		ClientID:          client.GetID(),
		RequestedScopes:   request.GetRequestedScopes(),
		GrantedScopes:     request.GetGrantedScopes(),
		Form:              request.GetRequestForm(),
		SessionID:         session.ID,
		RequestedAudience: request.GetRequestedAudience(),
		GrantedAudience:   request.GetGrantedAudience(),
	}

	err = s.Create.CreateAuthorizationCode(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

// GetAuthorizeCodeSession hydrates the session based on the given code and returns the authorization request.
// If the authorization code has been invalidated with `InvalidateAuthorizeCodeSession`, this
// method should return the ErrInvalidatedAuthorizeCode error.
//
// Make sure to also return the fosite.Requester value when returning the fosite.ErrInvalidatedAuthorizeCode error!
func (s Storage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (request fosite.Requester, err error) {
	authorizationCode, err := s.Query.GetAuthorizationCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	rq := &fosite.Request{
		ID:                authorizationCode.ID,
		RequestedAt:       authorizationCode.RequestedAt,
		Client:            authorizationCode.Client,
		RequestedScope:    fosite.Arguments(authorizationCode.RequestedScopes),
		GrantedScope:      fosite.Arguments(authorizationCode.GrantedScopes),
		Form:              authorizationCode.Form,
		Session:           &authorizationCode.Session,
		RequestedAudience: fosite.Arguments(authorizationCode.RequestedAudience),
		GrantedAudience:   fosite.Arguments(authorizationCode.GrantedAudience),
	}

	if !authorizationCode.Active {
		return rq, fosite.ErrInvalidatedAuthorizeCode
	}

	return rq, nil
}

// InvalidateAuthorizeCodeSession is called when an authorize code is being used. The state of the authorization
// code should be set to invalid and consecutive requests to GetAuthorizeCodeSession should return the
// ErrInvalidatedAuthorizeCode error.
func (s Storage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) (err error) {
	authorizationCode, err := s.Query.GetAuthorizationCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fosite.ErrNotFound
		}
		return err
	}

	if err := s.Update.UpdateAuthorizationCode(ctx, authorizationCode, map[string]interface{}{"active": false}); err != nil {
		return fmt.Errorf("failed to invalidate authorization code: %w", err)
	}

	return nil
}
