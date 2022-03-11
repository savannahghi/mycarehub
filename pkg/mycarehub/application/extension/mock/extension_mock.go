package mock

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	openSourceDto "github.com/savannahghi/engagementcore/pkg/engagement/application/common/dto"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
)

// FakeExtensionImpl mocks the external calls logic
type FakeExtensionImpl struct {
	MockComparePINFn                      func(rawPwd string, salt string, encodedPwd string, options *extension.Options) bool
	MockCreateFirebaseCustomTokenFn       func(ctx context.Context, uid string) (string, error)
	MockAuthenticateCustomFirebaseTokenFn func(customAuthToken string) (*firebasetools.FirebaseUserTokens, error)
	MockGenerateTempPINFn                 func(ctx context.Context) (string, error)
	MockEncryptPINFn                      func(rawPwd string, options *extension.Options) (string, string)
	MockSendSMSFn                         func(ctx context.Context, phoneNumbers string, message string, from enumutils.SenderID) (*openSourceDto.SendMessageResponse, error)
	MockGenerateAndSendOTPFn              func(ctx context.Context, phoneNumber string) (string, error)
	MockGenerateOTPFn                     func(ctx context.Context) (string, error)
	MockGenerateRetryOTPFn                func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error)
	MockSendSMSViaTwilioFn                func(ctx context.Context, phonenumber, message string) error
	MockSendInviteSMSFn                   func(ctx context.Context, phoneNumber, message string) error
	MockSendFeedbackFn                    func(ctx context.Context, subject, feedbackMessage string) (bool, error)
	MockGetLoggedInUserUIDFn              func(ctx context.Context) (string, error)
	MockMakeRequestFn                     func(ctx context.Context, method string, path string, body interface{}) (*http.Response, error)
	MockLoginFn                           func(ctx context.Context) http.HandlerFunc
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
		MockSendSMSFn: func(ctx context.Context, phoneNumbers string, message string, from enumutils.SenderID) (*openSourceDto.SendMessageResponse, error) {
			return &openSourceDto.SendMessageResponse{
				SMSMessageData: &openSourceDto.SMS{
					Recipients: []openSourceDto.Recipient{
						{
							Number: interserviceclient.TestUserPhoneNumber,
						},
					},
				},
			}, nil
		},
		MockGenerateAndSendOTPFn: func(ctx context.Context, phoneNumber string) (string, error) {
			return "111222", nil
		},
		MockGenerateOTPFn: func(ctx context.Context) (string, error) {
			return "111222", nil
		},
		MockGenerateRetryOTPFn: func(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
			return "test-OTP", nil
		},
		MockSendSMSViaTwilioFn: func(ctx context.Context, phonenumber, message string) error {
			return nil
		},
		MockSendInviteSMSFn: func(ctx context.Context, phoneNumber, message string) error {
			return nil
		},
		MockSendFeedbackFn: func(ctx context.Context, subject, feedbackMessage string) (bool, error) {
			return true, nil
		},
		MockGetLoggedInUserUIDFn: func(ctx context.Context) (string, error) {
			return uuid.New().String(), nil
		},
		MockMakeRequestFn: func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
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

// GenerateTempPIN mocks the generate temp pin method
func (f *FakeExtensionImpl) GenerateTempPIN(ctx context.Context) (string, error) {
	return f.MockGenerateTempPINFn(ctx)
}

// EncryptPIN mocks the encrypt pin method
func (f *FakeExtensionImpl) EncryptPIN(rawPwd string, options *extension.Options) (string, string) {
	return f.MockEncryptPINFn(rawPwd, options)
}

// SendSMS mocks the send sms method
func (f *FakeExtensionImpl) SendSMS(ctx context.Context, phoneNumbers string, message string, from enumutils.SenderID) (*openSourceDto.SendMessageResponse, error) {
	return f.MockSendSMSFn(ctx, phoneNumbers, message, from)
}

// GenerateAndSendOTP mocks the generate and send OTP method
func (f *FakeExtensionImpl) GenerateAndSendOTP(ctx context.Context, phoneNumber string) (string, error) {
	return f.MockGenerateAndSendOTPFn(ctx, phoneNumber)
}

// GenerateOTP mocks the GenerateOTP implementation
func (f *FakeExtensionImpl) GenerateOTP(ctx context.Context) (string, error) {
	return f.MockGenerateOTPFn(ctx)
}

// GenerateRetryOTP mock the implementation of generating a retry OTP
func (f *FakeExtensionImpl) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	return f.MockGenerateRetryOTPFn(ctx, payload)
}

// SendSMSViaTwilio mocks the implementation of sending a SMS via twilio
func (f *FakeExtensionImpl) SendSMSViaTwilio(ctx context.Context, phonenumber, message string) error {
	return f.MockSendSMSViaTwilioFn(ctx, phonenumber, message)
}

// SendInviteSMS mocks the implementation of sending an invite sms
func (f *FakeExtensionImpl) SendInviteSMS(ctx context.Context, phoneNumber, message string) error {
	return f.MockSendInviteSMSFn(ctx, phoneNumber, message)
}

//SendFeedback mocks the implementation sending feedback
func (f *FakeExtensionImpl) SendFeedback(ctx context.Context, subject, feedbackMessage string) (bool, error) {
	return f.MockSendFeedbackFn(ctx, subject, feedbackMessage)
}

// GetLoggedInUserUID mocks the implementation of getting a logged in user
func (f *FakeExtensionImpl) GetLoggedInUserUID(ctx context.Context) (string, error) {
	return f.MockGetLoggedInUserUIDFn(ctx)
}

// MakeRequest mocks the implementation of making a http request
func (f *FakeExtensionImpl) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	return f.MockMakeRequestFn(ctx, method, path, body)
}

// Login mocks the login implementation to retrieve a token
func (f *FakeExtensionImpl) Login(ctx context.Context) http.HandlerFunc {
	return f.MockLoginFn(ctx)
}
