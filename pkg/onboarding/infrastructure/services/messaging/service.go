package messaging

import (
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

// TODO: Move `messaging` service to engagement after complete Twilio's migration

const (
	smsService    = "engagement"
	twilioService = "twilio"

	// twilio isc paths
	sendTwilioSMS = "internal/send_sms"

	// engagement isc paths
	sendSMS = "internal/send_sms"
)

// ServiceMessaging represents messaging via sms and/or twilio
type ServiceMessaging interface {
	FetchSMSClient() *base.InterServiceClient
	FetchTwilioClient() *base.InterServiceClient
	SendSMS(phoneNumbers []string, message string) error
}

// ServiceMessagingImpl represents our messaging struct
type ServiceMessagingImpl struct {
	SMS    *base.InterServiceClient
	Twilio *base.InterServiceClient
}

// NewServiceMessagingImpl returns new initialized instance of ServiceOnboardingImpl
func NewServiceMessagingImpl(baseExt extension.BaseExtension) ServiceMessaging {
	sms := utils.NewInterServiceClient(smsService, baseExt)
	tw := utils.NewInterServiceClient(twilioService, baseExt)
	return &ServiceMessagingImpl{SMS: sms, Twilio: tw}
}

// FetchSMSClient returns engagement's service SMS ISC
func (s *ServiceMessagingImpl) FetchSMSClient() *base.InterServiceClient {
	return s.SMS
}

// FetchTwilioClient returns twilio's service SMS ISC
func (s *ServiceMessagingImpl) FetchTwilioClient() *base.InterServiceClient {
	return s.Twilio
}

// SendSMS does the actual delvery of messages to the provided phone numbers
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
