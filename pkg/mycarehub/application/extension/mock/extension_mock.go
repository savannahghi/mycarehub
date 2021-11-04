package mock

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// FakeOnboardingLibraryExtensionImpl is a fake representation of the onboarding library
type FakeOnboardingLibraryExtensionImpl struct {
	MockEncryptPINFn      func(rawPwd string, options *extension.Options) (string, string)
	MockComparePINFn      func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	MockGenerateTempPINFn func(ctx context.Context) (string, error)
	MockSendSMSFn         func(ctx context.Context, phoneNumbers []string, message string) error
}

// NewFakeOnboardingLibraryExtension initializes a new onboarding library mocks
func NewFakeOnboardingLibraryExtension() *FakeOnboardingLibraryExtensionImpl {
	return &FakeOnboardingLibraryExtensionImpl{
		MockEncryptPINFn: func(rawPwd string, options *extension.Options) (string, string) {
			return "hash", "salt"
		},

		MockComparePINFn: func(rawPwd, salt, encodedPwd string, options *extension.Options) bool {
			return true
		},

		MockGenerateTempPINFn: func(ctx context.Context) (string, error) {
			return "temppin", nil
		},

		MockSendSMSFn: func(ctx context.Context, phoneNumbers []string, message string) error {
			return nil
		},
	}
}

// EncryptPIN mocks the encrypt pin method
func (f *FakeOnboardingLibraryExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return f.MockEncryptPINFn(rawPwd, options)
}

// ComparePIN mocks the compare pin method
func (f *FakeOnboardingLibraryExtensionImpl) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return f.MockComparePINFn(rawPwd, salt, encodedPwd, options)
}

// GenerateTempPIN mocks the generate temporary pin method
func (f *FakeOnboardingLibraryExtensionImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return f.MockGenerateTempPINFn(ctx)
}

// SendSMS mocks the send sms method
func (f *FakeOnboardingLibraryExtensionImpl) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.MockSendSMSFn(ctx, phoneNumbers, message)
}
