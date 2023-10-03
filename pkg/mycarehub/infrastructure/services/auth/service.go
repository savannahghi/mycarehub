package auth

import (
	"context"

	"github.com/savannahghi/authutils"
)

// IServiceAuth holds the method used to communicate with SIL's apiclient
type IServiceAuth interface {
	AuthenticateWithSlade360(ctx context.Context) (*authutils.OAUTHResponse, error)
}

// ISILAuthServerClient defines the method used to initialize the client used to make request to SIL auth server
type ISILAuthServerClient interface {
	Authenticate() (*authutils.OAUTHResponse, error)
}

// SILAuthServiceImpl initializes auth client
type SILAuthServiceImpl struct {
	client ISILAuthServerClient
}

// NewAuthService is the constructor which initializes auth service implementation.
func NewAuthService(client ISILAuthServerClient) *SILAuthServiceImpl {
	return &SILAuthServiceImpl{
		client: client,
	}
}

// AuthenticateWithSlade360 is used to authenticate mch service with slade 360 auth server
func (sa *SILAuthServiceImpl) AuthenticateWithSlade360(ctx context.Context) (*authutils.OAUTHResponse, error) {
	return sa.client.Authenticate()
}
