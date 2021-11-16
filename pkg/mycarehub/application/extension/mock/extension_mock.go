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
	MockGenerateTempPINFn                 func(ctx context.Context) (string, error)
	MockEncryptPINFn                      func(rawPwd string, options *extension.Options) (string, string)
	MockSendInviteSMSFn                   func(ctx context.Context, phoneNumbers []string, message string) error
	MockGenerateAndSendOTPFn              func(ctx context.Context, phoneNumber string) (string, error)
	MockGenerateOTPFn                     func(ctx context.Context) (string, error)
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
		MockGenerateTempPINFn: func(ctx context.Context) (string, error) {
			return uuid.New().String(), nil
		},
		MockEncryptPINFn: func(rawPwd string, options *extension.Options) (string, string) {
			return uuid.New().String(), uuid.New().String()
		},
		MockSendInviteSMSFn: func(ctx context.Context, phoneNumbers []string, message string) error {
			return nil
		},
		MockGenerateAndSendOTPFn: func(ctx context.Context, phoneNumber string) (string, error) {
			return "111222", nil
		},
		MockGenerateOTPFn: func(ctx context.Context) (string, error) {
			return "111222", nil
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

// GenerateTempPIN mocks the generate temp pin method
func (f *FakeExtensionImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return f.MockGenerateTempPINFn(ctx)
}

// EncryptPIN mocks the encrypt pin method
func (f *FakeExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return f.MockEncryptPINFn(rawPwd, options)
}

// SendInviteSMS mocks the send sms method
func (f *FakeExtensionImpl) SendInviteSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.MockSendInviteSMSFn(ctx, phoneNumbers, message)
}

// GenerateAndSendOTP mocks the generate and send OTP method
func (f *FakeExtensionImpl) GenerateAndSendOTP(ctx context.Context, phoneNumber string) (string, error) {
	return f.MockGenerateAndSendOTPFn(ctx, phoneNumber)
}

// GenerateOTP mocks the GenerateOTP implementation
func (f *FakeExtensionImpl) GenerateOTP(ctx context.Context) (string, error) {
	return f.MockGenerateOTPFn(ctx)
}
