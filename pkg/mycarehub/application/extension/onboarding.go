package extension

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// OnboardingLibraryExtension is an interface that represents methods in the
// `onboarding library`. Adding this layer will help write unit tests for the
// methods that depends on the onboarding library
type OnboardingLibraryExtension interface {
	EncryptPIN(rawPwd string, options *extension.Options) (string, string)
	ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	GenerateTempPIN(ctx context.Context) (string, error)
}

// OnboardingLibraryImpl represents onboarding library usecases
type OnboardingLibraryImpl struct {
	pinExt extension.PINExtension
}

// NewOnboardingLibImpl creates a new instance of the onboarding library extension
func NewOnboardingLibImpl() OnboardingLibraryExtension {
	pinExtension := extension.NewPINExtensionImpl()
	return &OnboardingLibraryImpl{
		pinExt: pinExtension,
	}
}

// EncryptPIN takes two arguments, a raw pin, and a pointer to an Options struct.
// In order to use default options, pass `nil` as the second argument.
// It returns the generated salt and encoded key for the user.
func (o *OnboardingLibraryImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return o.pinExt.EncryptPIN(rawPwd, nil)
}

// ComparePIN takes four arguments, the raw password, its generated salt, the encoded password,
// and a pointer to the Options struct, and returns a boolean value determining whether the password is the correct one or not.
// Passing `nil` as the last argument resorts to default options.
func (o *OnboardingLibraryImpl) ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool {
	return o.pinExt.ComparePIN(rawPwd, salt, encodedPwd, nil)
}

// GenerateTempPIN generates a temporary One Time PIN for a user
// The PIN will have 4 digits formatted as a string
func (o *OnboardingLibraryImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return o.pinExt.GenerateTempPIN(ctx)
}
