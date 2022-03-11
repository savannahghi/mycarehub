package loginservice

import (
	"context"
	"net/http"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
)

// ServiceLogin represents an external login service login
type ServiceLogin interface {
	Login(ctx context.Context) http.HandlerFunc
}

// ServiceLoginImpl represents the service implementation object
type ServiceLoginImpl struct {
	ExternalExtension extension.ExternalMethodsExtension
}

// NewServiceLoginImpl initializes a new login service
func NewServiceLoginImpl(ext extension.ExternalMethodsExtension) ServiceLogin {
	return ServiceLoginImpl{
		ExternalExtension: ext,
	}
}

// Login returns a function that authenticates against firebase and returns a valid
// token that can be used to access our APIs
func (s ServiceLoginImpl) Login(ctx context.Context) http.HandlerFunc {
	return s.ExternalExtension.Login(ctx)
}
