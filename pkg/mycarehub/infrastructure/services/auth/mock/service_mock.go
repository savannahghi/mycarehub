package mock

import (
	"context"

	"github.com/savannahghi/authutils"
)

// AuthServiceMock mocks the slade 360's authentication service implementations
type AuthServiceMock struct {
	MockAuthenticateWithSlade360Fn func(ctx context.Context) (*authutils.OAUTHResponse, error)
}

// NewAuthServiceMock is the constructor that initializes the service mocks
func NewAuthServiceMock() *AuthServiceMock {
	return &AuthServiceMock{
		MockAuthenticateWithSlade360Fn: func(ctx context.Context) (*authutils.OAUTHResponse, error) {
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

// AuthenticateWithSlade360 mocks the implementation of service's authenticate method
func (a *AuthServiceMock) AuthenticateWithSlade360(ctx context.Context) (*authutils.OAUTHResponse, error) {
	return a.MockAuthenticateWithSlade360Fn(ctx)
}
