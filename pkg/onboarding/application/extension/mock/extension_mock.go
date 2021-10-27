package mock

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// FakeOnboardingLibraryExtensionImpl is a fake representation of the onboarding library
type FakeOnboardingLibraryExtensionImpl struct {
	EncryptPINFn      func(rawPwd string, options *extension.Options) (string, string)
	ComparePINFn      func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	GenerateTempPINFn func(ctx context.Context) (string, error)
	SendSMSFn         func(ctx context.Context, phoneNumbers []string, message string) error
}

// EncryptPIN mocks the encrypt pin method
func (f *FakeOnboardingLibraryExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return f.EncryptPINFn(rawPwd, options)
}

// ComparePIN mocks the compare pin method
func (f *FakeOnboardingLibraryExtensionImpl) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return f.ComparePINFn(rawPwd, salt, encodedPwd, options)
}

// GenerateTempPIN mocks the generate temporary pin method
func (f *FakeOnboardingLibraryExtensionImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return f.GenerateTempPINFn(ctx)
}

// SendSMS mocks the send sms method
func (f *FakeOnboardingLibraryExtensionImpl) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.SendSMSFn(ctx, phoneNumbers, message)
}
