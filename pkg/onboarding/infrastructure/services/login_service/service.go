package loginservice

import (
	"context"
	"net/http"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// ServiceLogin represents external login service logic
type ServiceLogin interface {
	GetLoginFunc(ctx context.Context) http.HandlerFunc
	GetLogoutFunc(ctx context.Context) http.HandlerFunc
	GetRefreshFunc() http.HandlerFunc
	GetVerifyTokenFunc(ctx context.Context) http.HandlerFunc
}

// ServiceLoginImpl represents the service implementation object
type ServiceLoginImpl struct {
	baseExt extension.BaseExtension
}

// NewServiceLogin initializes a new login service
func NewServiceLogin(ext extension.BaseExtension) ServiceLogin {
	return &ServiceLoginImpl{
		baseExt: ext,
	}
}

// GetLoginFunc returns a function that can authenticate against both Slade 360 and Firebase
func (l ServiceLoginImpl) GetLoginFunc(ctx context.Context) http.HandlerFunc {
	return l.baseExt.GetLoginFunc(ctx)
}

// GetLogoutFunc logs the user out of Firebase
func (l ServiceLoginImpl) GetLogoutFunc(ctx context.Context) http.HandlerFunc {
	return l.baseExt.GetLogoutFunc(ctx)
}

// GetRefreshFunc is used to refresh OAuth tokens
func (l ServiceLoginImpl) GetRefreshFunc() http.HandlerFunc {
	return l.baseExt.GetRefreshFunc()
}

// GetVerifyTokenFunc confirms that an EDI access token (supplied) is valid.
// If it is valid, it exchanges it for a Firebase ID token.
func (l ServiceLoginImpl) GetVerifyTokenFunc(ctx context.Context) http.HandlerFunc {
	return l.baseExt.GetVerifyTokenFunc(ctx)
}
