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
	MockEncryptPINFn                      func(rawPwd string, options *extension.Options) (string, string)
	MockGenerateTempPINFn                 func(ctx context.Context) (string, error)
	MockSendSMSFn                         func(ctx context.Context, phoneNumbers []string, message string) error
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
		MockEncryptPINFn: func(rawPwd string, options *extension.Options) (string, string) {
			return uuid.New().String(), uuid.New().String()
		},
		MockGenerateTempPINFn: func(ctx context.Context) (string, error) {
			return uuid.New().String(), nil
		},
		MockSendSMSFn: func(ctx context.Context, phoneNumbers []string, message string) error {
			return nil
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

// EncryptPIN mocks the encrypt pin method
func (f *FakeExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return f.MockEncryptPINFn(rawPwd, options)
}

// GenerateTempPIN mocks the generate temporary pin method
func (f *FakeExtensionImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return f.MockGenerateTempPINFn(ctx)
}

// SendSMS mocks the send sms method
func (f *FakeExtensionImpl) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.MockSendSMSFn(ctx, phoneNumbers, message)
}
