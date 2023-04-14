package domain

import (
	"net/url"
	"time"
)

// AccessToken represents an oauth2 access token
type AccessToken struct {
	ID        string
	Active    bool
	Signature string

	RequestedAt       time.Time
	RequestedScopes   []string
	GrantedScopes     []string
	Form              url.Values
	RequestedAudience []string
	GrantedAudience   []string

	ClientID  string
	Client    OauthClient
	SessionID string
	Session   Session
}

// AuthorizationCode represents an Oauth2 authorization code
type AuthorizationCode struct {
	ID     string
	Active bool
	Code   string

	RequestedAt       time.Time
	RequestedScopes   []string
	GrantedScopes     []string
	Form              url.Values
	RequestedAudience []string
	GrantedAudience   []string

	SessionID string
	Session   Session
	ClientID  string
	Client    OauthClient
}

// OauthClientJWT
type OauthClientJWT struct {
	ID        string
	Active    bool
	JTI       string
	ExpiresAt time.Time
}

// PKCE
type PKCE struct {
	ID        string
	Active    bool
	Signature string

	RequestedAt       time.Time
	RequestedScopes   []string
	GrantedScopes     []string
	Form              url.Values
	RequestedAudience []string
	GrantedAudience   []string

	SessionID string
	Session   Session
	ClientID  string
	Client    OauthClient
}

// RefreshToken
type RefreshToken struct {
	ID        string
	Active    bool
	Signature string

	RequestedAt       time.Time
	RequestedScopes   []string
	GrantedScopes     []string
	Form              url.Values
	RequestedAudience []string
	GrantedAudience   []string

	ClientID  string
	Client    OauthClient
	SessionID string
	Session   Session
}
