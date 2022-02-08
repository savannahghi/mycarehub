package extension

import (
	"context"
	"fmt"

	openSourceDto "github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	engagementInfra "github.com/savannahghi/engagementcore/pkg/engagement/infrastructure"
	engagementEmail "github.com/savannahghi/engagementcore/pkg/engagement/usecases/mail"
	engagementOTP "github.com/savannahghi/engagementcore/pkg/engagement/usecases/otp"
	engagementSMS "github.com/savannahghi/engagementcore/pkg/engagement/usecases/sms"
	engagementTwilio "github.com/savannahghi/engagementcore/pkg/engagement/usecases/twilio"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/serverutils"
)

// ExternalMethodsExtension is an interface that represents methods that are
// called from external libraries. Adding this layer will help write unit tests
type ExternalMethodsExtension interface {
	CreateFirebaseCustomToken(ctx context.Context, uid string) (string, error)
	AuthenticateCustomFirebaseToken(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	ComparePIN(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	EncryptPIN(rawPwd string, options *extension.Options) (string, string)
	GenerateTempPIN(ctx context.Context) (string, error)
	SendSMS(ctx context.Context, phoneNumbers string, message string, from enumutils.SenderID) (*openSourceDto.SendMessageResponse, error)
	GenerateAndSendOTP(ctx context.Context, phoneNumber string) (string, error)
	GenerateOTP(ctx context.Context) (string, error)
	GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
	SendSMSViaTwilio(ctx context.Context, phonenumber, message string) error
	SendInviteSMS(ctx context.Context, phoneNumber, message string) error
	SendFeedback(ctx context.Context, subject, feedbackMessage string) (bool, error)
	GetLoggedInUserUID(ctx context.Context) (string, error)
}

// External type implements external methods
type External struct {
	pinExt          extension.PINExtension
	otpExtension    engagementOTP.ImplOTP
	twilioExtension engagementTwilio.ImplTwilio
	smsExtension    engagementSMS.UsecaseSMS
	emailExtension  engagementEmail.UsecaseMail
}

// NewExternalMethodsImpl creates a new instance of the external methods
func NewExternalMethodsImpl() ExternalMethodsExtension {
	pinExtension := extension.NewPINExtensionImpl()
	otpExt := engagementOTP.NewOTP(engagementInfra.NewInteractor())
	twilioExt := engagementTwilio.NewImplTwilio(engagementInfra.NewInteractor())
	smsExt := engagementSMS.NewSMS(engagementInfra.NewInteractor())
	emailExt := engagementEmail.NewMail(engagementInfra.NewInteractor())
	return &External{
		pinExt:          pinExtension,
		otpExtension:    *otpExt,
		twilioExtension: *twilioExt,
		smsExtension:    smsExt,
		emailExtension:  emailExt,
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

// SendSMS does the actual delivery of messages to the provided phone numbers
func (e *External) SendSMS(ctx context.Context, phoneNumbers string, message string, from enumutils.SenderID) (*openSourceDto.SendMessageResponse, error) {
	return e.smsExtension.Send(ctx, message, phoneNumbers, from)
}

// GenerateAndSendOTP generates a new OTP and sends it to the provided phone number
func (e *External) GenerateAndSendOTP(ctx context.Context, phoneNumber string) (string, error) {
	return e.otpExtension.GenerateAndSendOTP(ctx, phoneNumber, nil)
}

// GenerateOTP generates an OTP
func (e *External) GenerateOTP(ctx context.Context) (string, error) {
	return e.otpExtension.GenerateOTP(ctx)
}

// GenerateRetryOTP generates fallback OTPs when Africa is talking sms fails
func (e *External) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	return e.otpExtension.GenerateRetryOTP(ctx, &payload.Phone, 2, nil)
}

// SendSMSViaTwilio makes a request to Twilio to send an SMS to a non-kenyan number
func (e *External) SendSMSViaTwilio(ctx context.Context, phonenumber, message string) error {
	return e.twilioExtension.SendSMS(ctx, phonenumber, message)
}

// SendInviteSMS is used to send an Invite SMS to a client
func (e *External) SendInviteSMS(ctx context.Context, phoneNumber, message string) error {
	if interserviceclient.IsKenyanNumber(phoneNumber) {
		_, err := e.SendSMS(ctx, phoneNumber, message, enumutils.SenderIDBewell)
		if err != nil {
			return fmt.Errorf("failed to send invite sms to recipient")
		}
	} else {
		// Make the request to twilio
		err := e.SendSMSViaTwilio(ctx, phoneNumber, message)
		if err != nil {
			return fmt.Errorf("sms not sent via twilio: %v", err)
		}
	}
	return nil
}

// SendFeedback sends the clients feed email
func (e *External) SendFeedback(ctx context.Context, subject, feedbackMessage string) (bool, error) {
	_, err := e.emailExtension.SimpleEmail(ctx, subject, feedbackMessage, nil, serverutils.MustGetEnvVar("SAVANNAH_ADMIN_EMAIL"))
	if err != nil {
		return false, fmt.Errorf("an erro occurred while sending the feedback: %v", err)
	}

	return true, nil
}

// GetLoggedInUserUID get the logged in user uid
func (e *External) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return firebasetools.GetLoggedInUserUID(ctx)
}
