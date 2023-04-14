package domain

import "github.com/ory/fosite"

// OauthClient represents an application that is authorized to access a user's resources on a server using the OAuth2 protocol.
type OauthClient struct {
	// unique identifier of the OAuth2 client
	ID string
	// human-readable name of the OAuth2 client
	Name string
	// indicates whether the OAuth2 client is active or not
	Active bool
	// secret used by the client to authenticate with the authorization server
	// uses bcrypt by default
	Secret string
	// list of previously used secrets that were rotated out
	RotatedSecrets []string
	// indicates whether the client is a public or confidential client
	Public bool
	// list of valid URIs to redirect the user after authorization
	RedirectURIs []string
	// list of scopes the client is authorized to request
	Scopes []string
	// list of intended audiences for the access token
	Audience []string
	// a list of OAuth2 grant types that the client is authorized to use when requesting access tokens
	// e.g ["authorization_code", "refresh_token"]
	Grants []string
	// a list of OAuth2 response types that the client is authorized to use when requesting authorization
	// e.g  ["code", "token"]
	ResponseTypes []string
	// the authentication method that the client uses to authenticate with the auth server when requesting tokens.
	// e.g  "client_secret_basic"
	TokenEndpointAuthMethod string
}

// GetID returns the client ID.
func (c OauthClient) GetID() string {
	return c.ID
}

// GetHashedSecret returns the hashed secret as it is stored in the store.
func (c OauthClient) GetHashedSecret() []byte {
	return []byte(c.Secret)
}

func (c OauthClient) GetRotatedHashes() [][]byte {
	var secrets [][]byte

	for _, secret := range c.RotatedSecrets {
		secrets = append(secrets, []byte(secret))
	}

	return secrets
}

// GetRedirectURIs returns the client's allowed redirect URIs.
func (c OauthClient) GetRedirectURIs() []string {
	return c.RedirectURIs
}

// GetGrantTypes returns the client's allowed grant types.
func (c OauthClient) GetGrantTypes() fosite.Arguments {
	return c.Grants
}

// GetResponseTypes returns the client's allowed response types.
// All allowed combinations of response types have to be listed, each combination having
// response types of the combination separated by a space.
func (c OauthClient) GetResponseTypes() fosite.Arguments {
	return c.ResponseTypes
}

// GetScopes returns the scopes this client is allowed to request.
func (c OauthClient) GetScopes() fosite.Arguments {
	return c.Scopes
}

// IsPublic returns true, if this client is marked as public.
func (c OauthClient) IsPublic() bool {
	return c.Public
}

// GetAudience returns the allowed audience(s) for this client.
func (c OauthClient) GetAudience() fosite.Arguments {
	return c.Audience
}
