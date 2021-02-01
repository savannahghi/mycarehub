package messaging

import (
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

const smsService = "sms"
const twilioService = "sms"

// twilio isc paths
const (
	sendTwilioSMS = "internal/send_sms"
)

// sms isc paths
const (
	sendSMS = "internal/send_sms"
)

// ServiceMessaging represents messaging via sms and/or twilio
type ServiceMessaging interface {
	FetchSMSClient() *base.InterServiceClient
	FetchTwilioClient() *base.InterServiceClient
	SendSMS(phoneNumbers []string, message string) error
}

// ServiceMessagingImpl ...
type ServiceMessagingImpl struct {
	SMS    *base.InterServiceClient
	Twilio *base.InterServiceClient
}

// NewServiceMessagingImpl ...
func NewServiceMessagingImpl(baseExt extension.BaseExtension) ServiceMessaging {
	sms := utils.NewInterServiceClient(smsService, baseExt)
	tw := utils.NewInterServiceClient(twilioService, baseExt)
	return &ServiceMessagingImpl{SMS: sms, Twilio: tw}
}

// FetchSMSClient ...
func (s *ServiceMessagingImpl) FetchSMSClient() *base.InterServiceClient {
	return s.SMS
}

// FetchTwilioClient ...
func (s *ServiceMessagingImpl) FetchTwilioClient() *base.InterServiceClient {
	return s.Twilio
}

// SendSMS ...
func (s *ServiceMessagingImpl) SendSMS(phoneNumbers []string, message string) error {
	smsISC := base.SmsISC{
		Isc:      s.FetchSMSClient(),
		EndPoint: sendSMS,
	}

	twilioISC := base.SmsISC{
		Isc:      s.FetchTwilioClient(),
		EndPoint: sendTwilioSMS,
	}
	return base.SendSMS(phoneNumbers, message, smsISC, twilioISC)
}
