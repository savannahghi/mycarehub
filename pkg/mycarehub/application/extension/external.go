package extension

import (
	"context"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

// ExternalMethodsExtension is an interface that represents methods that are
// called from external libraries. Adding this layer will help write unit tests
type ExternalMethodsExtension interface {
	CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error)
	AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	GenerateTempPIN(ctx context.Context) (string, error)
	EncryptPIN(rawPwd string, options *extension.Options) (string, string)
	SendSMS(ctx context.Context, phoneNumbers []string, message string) error
}

// External type implements external methods
type External struct {
	pinExt        extension.PINExtension
	engagementExt engagement.ServiceEngagement
}

// NewExternalMethodsImpl creates a new instance of the external methods
func NewExternalMethodsImpl() ExternalMethodsExtension {
	pinExtension := extension.NewPINExtensionImpl()

	var firebaseClient firebasetools.IFirebaseClient
	iscExt := extension.NewISCExtension()
	baseExt := extension.NewBaseExtensionImpl(firebaseClient)
	engagementExtension := engagement.NewServiceEngagementImpl(iscExt, baseExt)
	return &External{
		pinExt:        pinExtension,
		engagementExt: engagementExtension,
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

// GenerateTempPIN generates a temporary One Time PIN for a user
// The PIN will have 4 digits formatted as a string
func (o *External) GenerateTempPIN(ctx context.Context) (string, error) {
	return o.pinExt.GenerateTempPIN(ctx)
}

// EncryptPIN takes two arguments, a raw pin, and a pointer to an Options struct.
// In order to use default options, pass `nil` as the second argument.
// It returns the generated salt and encoded key for the user.
func (o *External) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return o.pinExt.EncryptPIN(rawPwd, nil)
}

// SendSMS does the actual delivery of messages to the provided phone numbers
func (o *External) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return o.engagementExt.SendSMS(ctx, phoneNumbers, message)
}
