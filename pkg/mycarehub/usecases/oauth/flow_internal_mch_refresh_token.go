package oauth

import (
	"context"
	"errors"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
)

// OAuth2InternalRefreshFactory creates an OAuth2 client credentials refresh token handler and registers
// an access token, refresh token and authorize code validator.
func OAuth2InternalRefreshFactory(config fosite.Configurator, storage interface{}, strategy interface{}) interface{} {
	return &InternalRefreshHandler{
		AccessTokenStrategy:    strategy.(oauth2.AccessTokenStrategy),
		RefreshTokenStrategy:   strategy.(oauth2.RefreshTokenStrategy),
		TokenRevocationStorage: storage.(oauth2.TokenRevocationStorage),
		Config:                 config,
	}
}

type InternalRefreshHandler struct {
	AccessTokenStrategy    oauth2.AccessTokenStrategy
	RefreshTokenStrategy   oauth2.RefreshTokenStrategy
	TokenRevocationStorage oauth2.TokenRevocationStorage
	Config                 interface {
		fosite.AccessTokenLifespanProvider
		fosite.RefreshTokenLifespanProvider
		fosite.ScopeStrategyProvider
		fosite.AudienceStrategyProvider
		fosite.RefreshTokenScopesProvider
	}
}

// PopulateTokenEndpointResponse is responsible for setting return values and should only be executed if
// the handler's HandleTokenEndpointRequest did not return ErrUnknownRequest.
func (i InternalRefreshHandler) PopulateTokenEndpointResponse(ctx context.Context, request fosite.AccessRequester, responder fosite.AccessResponder) error {
	if !i.CanHandleTokenEndpointRequest(ctx, request) {
		return fosite.ErrUnknownRequest
	}

	err := i.HandleTokenEndpointRequest(ctx, request)
	if err != nil {
		return err
	}

	accessToken, accessSignature, err := i.AccessTokenStrategy.GenerateAccessToken(ctx, request)
	if err != nil {
		return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
	}

	refreshToken, refreshSignature, err := i.RefreshTokenStrategy.GenerateRefreshToken(ctx, request)
	if err != nil {
		return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
	}

	signature := i.RefreshTokenStrategy.RefreshTokenSignature(ctx, request.GetRequestForm().Get("refresh_token"))

	ts, err := i.TokenRevocationStorage.GetRefreshTokenSession(ctx, signature, nil)
	if err != nil {
		return err
	} else if err := i.TokenRevocationStorage.RevokeRefreshToken(ctx, ts.GetID()); err != nil {
		return err
	}

	storeReq := request.Sanitize([]string{})
	storeReq.SetID(ts.GetID())

	if err = i.TokenRevocationStorage.CreateAccessTokenSession(ctx, accessSignature, storeReq); err != nil {
		return err
	}

	if err = i.TokenRevocationStorage.CreateRefreshTokenSession(ctx, refreshSignature, storeReq); err != nil {
		return err
	}

	atLifespan := i.Config.GetAccessTokenLifespan(ctx)
	responder.SetAccessToken(accessToken)
	responder.SetTokenType("bearer")
	responder.SetExpiresIn(atLifespan)
	responder.SetScopes(request.GetGrantedScopes())
	responder.SetExtra("refresh_token", refreshToken)

	return nil
}

// HandleTokenEndpointRequest handles an authorize request. If the handler is not responsible for handling
// the request, this method should return ErrUnknownRequest and otherwise handle the request.
func (i InternalRefreshHandler) HandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) error {
	if !i.CanHandleTokenEndpointRequest(ctx, request) {
		return fosite.ErrUnknownRequest
	}

	if !request.GetClient().GetGrantTypes().Has("internal_refresh_token") {
		return fosite.ErrUnauthorizedClient.WithHint("The OAuth 2.0 Client is not allowed to use authorization grant 'internal_refresh_token'.")
	}

	refresh := request.GetRequestForm().Get("refresh_token")
	signature := i.RefreshTokenStrategy.RefreshTokenSignature(ctx, refresh)

	originalRequest, err := i.TokenRevocationStorage.GetRefreshTokenSession(ctx, signature, request.GetSession())
	if err != nil {
		if errors.Is(err, fosite.ErrInactiveToken) {
			return fosite.ErrInactiveToken.WithWrap(err).WithDebug(err.Error())
		} else if errors.Is(err, fosite.ErrNotFound) {
			return fosite.ErrInvalidGrant.WithWrap(err).WithDebugf("The refresh token has not been found: %s", err.Error())
		} else {
			return fosite.ErrServerError.WithWrap(err).WithDebug(err.Error())
		}
	}

	if err := i.RefreshTokenStrategy.ValidateRefreshToken(ctx, originalRequest, refresh); err != nil {
		if errors.Is(err, fosite.ErrTokenExpired) {
			return fosite.ErrInvalidGrant.WithWrap(err).WithDebug(err.Error())
		}

		return fosite.ErrInvalidRequest.WithWrap(err).WithDebug(err.Error())
	}

	// The authorization server MUST ... and ensure that the refresh token was issued to the authenticated client
	if originalRequest.GetClient().GetID() != request.GetClient().GetID() {
		return fosite.ErrInvalidGrant.WithHint("The OAuth 2.0 Client ID from this request does not match the ID during the initial token issuance.")
	}

	request.SetSession(originalRequest.GetSession().Clone())
	request.SetRequestedScopes(originalRequest.GetRequestedScopes())
	request.SetRequestedAudience(originalRequest.GetRequestedAudience())

	atLifespan := i.Config.GetAccessTokenLifespan(ctx)
	request.GetSession().SetExpiresAt(fosite.AccessToken, time.Now().UTC().Add(atLifespan).Round(time.Second))

	rtLifespan := i.Config.GetRefreshTokenLifespan(ctx)
	if rtLifespan > -1 {
		request.GetSession().SetExpiresAt(fosite.RefreshToken, time.Now().UTC().Add(rtLifespan).Round(time.Second))
	}

	return nil
}

// CanSkipClientAuth indicates if client authentication can be skipped. By default it MUST be false, unless you are
// implementing extension grant type, which allows unauthenticated client. CanSkipClientAuth must be called
// before HandleTokenEndpointRequest to decide, if AccessRequester will contain authenticated client.
func (InternalRefreshHandler) CanSkipClientAuth(ctx context.Context, request fosite.AccessRequester) bool {
	return true
}

// CanHandleRequest indicates, if TokenEndpointInternalRefreshHandler can handle this request or not. If true,
// HandleTokenEndpointRequest can be called.
func (InternalRefreshHandler) CanHandleTokenEndpointRequest(ctx context.Context, request fosite.AccessRequester) bool {
	return request.GetGrantTypes().ExactOne("internal_refresh_token") && request.GetClient().GetGrantTypes().Has("internal_refresh_token")
}
