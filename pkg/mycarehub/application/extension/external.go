package extension

import (
	"context"

	engagementInfra "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	engagementOTP "github.com/savannahghi/engagementcore/pkg/engagement/usecases/otp"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

const (
	engagementService = "engagement"
)

// ExternalMethodsExtension is an interface that represents methods that are
// called from external libraries. Adding this layer will help write unit tests
type ExternalMethodsExtension interface {
	CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error)
	AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	EncryptPIN(rawPwd string, options *extension.Options) (string, string)
	GenerateTempPIN(ctx context.Context) (string, error)
	SendInviteSMS(ctx context.Context, phoneNumbers []string, message string) error
	GenerateAndSendOTP(ctx context.Context, phoneNumber string) (string, error)
	GenerateOTP(ctx context.Context) (string, error)
}

// External type implements external methods
type External struct {
	pinExt        extension.PINExtension
	engagementExt engagement.ServiceEngagement
	otpExtension  engagementOTP.ImplOTP
}

// NewExternalMethodsImpl creates a new instance of the external methods
func NewExternalMethodsImpl() ExternalMethodsExtension {
	var firebaseClient firebasetools.IFirebaseClient

	pinExtension := extension.NewPINExtensionImpl()
	baseExt := extension.NewBaseExtensionImpl(firebaseClient)
	engagementISC := utils.NewInterServiceClient(engagementService, baseExt)
	engagementExtension := engagement.NewServiceEngagementImpl(engagementISC, baseExt)
	otpExt := engagementOTP.NewOTP(engagementInfra.NewInteractor())
	return &External{
		pinExt:        pinExtension,
		engagementExt: engagementExtension,
		otpExtension:  *otpExt,
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

// EncryptPIN takes two arguments, a raw pin, and a pointer to an Options struct.
// In order to use default options, pass `nil` as the second argument.
// It returns the generated salt and encoded key for the user.
func (e *External) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return e.pinExt.EncryptPIN(rawPwd, nil)
}

// GenerateTempPIN generates a temporary One Time PIN for a user
// The PIN will have 4 digits formatted as a string
func (e *External) GenerateTempPIN(ctx context.Context) (string, error) {
	return e.pinExt.GenerateTempPIN(ctx)
}

// SendInviteSMS does the actual delivery of messages to the provided phone numbers
func (e *External) SendInviteSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return e.engagementExt.SendSMS(ctx, phoneNumbers, message)
}

// GenerateAndSendOTP generates a new OTP and sends it to the provided phone number
func (e *External) GenerateAndSendOTP(ctx context.Context, phoneNumber string) (string, error) {
	return e.otpExtension.GenerateAndSendOTP(ctx, phoneNumber, nil)
}

// GenerateOTP generates an OTP
func (e *External) GenerateOTP(ctx context.Context) (string, error) {
	return e.otpExtension.GenerateOTP(ctx)
}
