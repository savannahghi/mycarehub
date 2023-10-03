package mock

import "github.com/savannahghi/authutils"

// AuthClientMock mocks the slade 360's authentication client implementations
type AuthClientMock struct {
	MockAuthenticateFn func() (*authutils.OAUTHResponse, error)
}

// NewAuthClientMock is the constructor that initializes the client mocks
func NewAuthClientMock() *AuthClientMock {
	return &AuthClientMock{
		MockAuthenticateFn: func() (*authutils.OAUTHResponse, error) {
			return &authutils.OAUTHResponse{
				Scope:        "",
				ExpiresIn:    3600,
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				TokenType:    "Bearer",
			}, nil
		},
	}
}

// Authenticate mocks the implementation of client's authenticate method
func (a *AuthClientMock) Authenticate() (*authutils.OAUTHResponse, error) {
	return a.MockAuthenticateFn()
}
