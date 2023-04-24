package oauth

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
)

// OAuth2InternalGrantFactory creates an OAuth2 client credentials grant handler and registers
// an access token, refresh token and authorize code validator.
func OAuth2InternalGrantFactory(config fosite.Configurator, storage interface{}, strategy interface{}) interface{} {
	return &InternalGrantHandler{
		HandleHelper: &oauth2.HandleHelper{
			AccessTokenStrategy: strategy.(oauth2.AccessTokenStrategy),
			AccessTokenStorage:  storage.(oauth2.AccessTokenStorage),
			Config:              config,
		},
		RefreshTokenStrategy: strategy.(oauth2.RefreshTokenStrategy),
		RefreshTokenStorage:  storage.(oauth2.RefreshTokenStorage),
		Config:               config,
	}
}

type InternalGrantHandler struct {
	*oauth2.HandleHelper
	RefreshTokenStrategy oauth2.RefreshTokenStrategy
	RefreshTokenStorage  oauth2.RefreshTokenStorage
	Config               interface {
		fosite.AccessTokenLifespanProvider
	}
}

// PopulateTokenEndpointResponse is responsible for setting return values and should only be executed if
// the handler's HandleTokenEndpointRequest did not return ErrUnknownRequest.
func (i InternalGrantHandler) PopulateTokenEndpointResponse(ctx context.Context, request fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !i.CanHandleTokenEndpointRequest(ctx, request) {
		return fosite.ErrUnknownRequest
	}

	atLifespan := i.Config.GetAccessTokenLifespan(ctx)
	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(atLifespan))

	access, accessSignature, err := i.AccessTokenStrategy.GenerateAccessToken(ctx, request)
	if err != nil {
		return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
	}

	refresh, refreshSignature, err := i.RefreshTokenStrategy.GenerateRefreshToken(ctx, request)
	if err != nil {
		return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
	}

	if err = i.AccessTokenStorage.CreateAccessTokenSession(ctx, accessSignature, request.Sanitize([]string{})); err != nil {
		return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
	} else if refreshSignature != "" {
		if err = i.RefreshTokenStorage.CreateRefreshTokenSession(ctx, refreshSignature, request.Sanitize([]string{})); err != nil {
			return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
		}
	}

	responder.SetAccessToken(access)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(atLifespan)
	responder.SetExtra("refresh_token", refresh)

	return nil
}

// HandleTokenEndpointRequest handles an authorize request. If the handler is not responsible for handling
// the request, this method should return ErrUnknownRequest and otherwise handle the request.
func (i InternalGrantHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !i.CanHandleTokenEndpointRequest(ctx, request) {
		return fosite.ErrUnknownRequest
	}

	if !request.GetClient().GetGrantTypes().Has("internal") {
		return fosite.ErrUnauthorizedClient.WithHint("The OAuth 2.0 Client is not allowed to use authorization grant 'internal'.")
	}

	return nil
}

// CanSkipClientAuth indicates if client authentication can be skipped. By default it MUST be false, unless you are
// implementing extension grant type, which allows unauthenticated client. CanSkipClientAuth must be called
// before HandleTokenEndpointRequest to decide, if AccessRequester will contain authenticated client.
func (InternalGrantHandler) CanSkipClientAuth(ctx context.Context, request fosite.AccessRequester) bool {
	return true
}

// CanHandleRequest indicates, if TokenEndpointInternalGrantHandler can handle this request or not. If true,
// HandleTokenEndpointRequest can be called.
func (InternalGrantHandler) CanHandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) bool {
	return request.GetGrantTypes().ExactOne("internal") && request.GetClient().GetGrantTypes().Has("internal")
}
