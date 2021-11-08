package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// FakeExtensionImpl mocks the external calls logic
type FakeExtensionImpl struct {
	MockComparePINFn                      func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	MockCreateFirebaseCustomTokenFn       func(ctx context.Context, uid string) (string, error)
	MockAuthenticateCustomFirebaseTokenFn func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
}

// NewFakeExtension initializes a new instance of the external calls mock
func NewFakeExtension() *FakeExtensionImpl {
	return &FakeExtensionImpl{
		MockComparePINFn: func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
			return true
		},

		MockCreateFirebaseCustomTokenFn: func(ctx context.Context, uid string) (string, error) {
			return uuid.New().String(), nil
		},

		MockAuthenticateCustomFirebaseTokenFn: func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
			return &firebasetools.FirebaseUserTokens{
				IDToken:      uuid.New().String(),
				RefreshToken: uuid.NewString(),
				ExpiresIn:    "1000",
			}, nil
		},
	}
}

// ComparePIN mocks the compare pin method
func (f *FakeExtensionImpl) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return f.MockComparePINFn(rawPwd, salt, encodedPwd, options)
}

// CreateFirebaseCustomToken mocks the create firebase custom token method
func (f *FakeExtensionImpl) CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error) {
	return f.MockCreateFirebaseCustomTokenFn(ctx, uid)
}

// AuthenticateCustomFirebaseToken mocks the authenticate custom firebase token method
func (f *FakeExtensionImpl) AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error) {
	return f.MockAuthenticateCustomFirebaseTokenFn(customAuthToken)
}
