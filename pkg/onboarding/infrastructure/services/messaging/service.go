package messaging

import (
	"context"

	"github.com/savannahghi/interserviceclient"
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
	FetchSMSClient() *interserviceclient.InterServiceClient
	FetchTwilioClient() *interserviceclient.InterServiceClient
	SendSMS(ctx context.Context, phoneNumbers []string, message string) error
}

// ServiceMessagingImpl represents our messaging struct
type ServiceMessagingImpl struct {
	SMS    *interserviceclient.InterServiceClient
	Twilio *interserviceclient.InterServiceClient
}

// NewServiceMessagingImpl returns new initialized instance of ServiceOnboardingImpl
func NewServiceMessagingImpl(baseExt extension.BaseExtension) ServiceMessaging {
	sms := utils.NewInterServiceClient(smsService, baseExt)
	tw := utils.NewInterServiceClient(twilioService, baseExt)
	return &ServiceMessagingImpl{SMS: sms, Twilio: tw}
}

// FetchSMSClient returns engagement's service SMS ISC
func (s *ServiceMessagingImpl) FetchSMSClient() *interserviceclient.InterServiceClient {
	return s.SMS
}

// FetchTwilioClient returns twilio's service SMS ISC
func (s *ServiceMessagingImpl) FetchTwilioClient() *interserviceclient.InterServiceClient {
	return s.Twilio
}

// SendSMS does the actual delvery of messages to the provided phone numbers
func (s *ServiceMessagingImpl) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	smsISC := interserviceclient.SmsISC{
		Isc:      s.FetchSMSClient(),
		EndPoint: sendSMS,
	}

	twilioISC := interserviceclient.SmsISC{
		Isc:      s.FetchTwilioClient(),
		EndPoint: sendTwilioSMS,
	}
	return interserviceclient.SendSMS(ctx, phoneNumbers, message, smsISC, twilioISC)
}
