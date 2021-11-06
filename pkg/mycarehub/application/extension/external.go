package extension

import (
	"context"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// ExternalMethodsExtension is an interface that represents methods that are
// called from external libraries. Adding this layer will help write unit tests
type ExternalMethodsExtension interface {
	CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error)
	AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
}

// External type implements external methods
type External struct {
	pinExt extension.PINExtension
}

// NewExternalMethodsImpl creates a new instance of the external methods
func NewExternalMethodsImpl() ExternalMethodsExtension {
	pinExtension := extension.NewPINExtensionImpl()
	return &External{
		pinExt: pinExtension,
	}
}

// CreateFirebaseCustomToken creates a custom auth token for the user with the
// indicated UID
func (e *External) CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error) {
	return firebasetools.CreateFirebaseCustomToken(ctx, uid)
}

// AuthenticateCustomFirebaseToken takes a custom Firebase auth token and tries to fetch an ID token
// If successful, a pointer to the ID token is returned
// Otherwise, an error is returned
func (e *External) AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
	return firebasetools.AuthenticateCustomFirebaseToken(customAuthToken)
}

// ComparePIN takes four arguments, the raw password, its generated salt, the encoded password,
// and a pointer to the Options struct, and returns a boolean value determining whether the password is the correct one or not.
// Passing `nil` as the last argument resorts to default options.
func (e *External) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return e.pinExt.ComparePIN(rawPwd, salt, encodedPwd, nil)
}
