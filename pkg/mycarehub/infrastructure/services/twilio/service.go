package twilio

import (
	"context"
	"net/url"

	twilioClient "github.com/kevinburke/twilio-go"
	"github.com/savannahghi/serverutils"
)

var (
	smsNumber = serverutils.MustGetEnvVar("TWILIO_SMS_NUMBER")
)

// ITwilioService is aN interface that holds methods for sending SMS messages via Twilio
type ITwilioService interface {
	SendSMSViaTwilio(ctx context.Context, phonenumber, message string) error
}

// ITwilioClient is an interface that holds methods for sending SMS messages via Twilio client
type ITwilioClient interface {
	SendMessage(from string, to string, body string, mediaURLs []*url.URL) (*twilioClient.Message, error)
}

// ServiceImpl defines the implementation of the Twilio service
type ServiceImpl struct {
	client ITwilioClient
}

// NewServiceTwilio initializes Twilio Service
func NewServiceTwilio(client ITwilioClient) *ServiceImpl {
	return &ServiceImpl{
		client: client,
	}
}

// SendSMSViaTwilio sends an outbound SMS message via Twilio
func (t *ServiceImpl) SendSMSViaTwilio(ctx context.Context, phonenumber, message string) error {
	_, err := t.client.SendMessage(smsNumber, phonenumber, message, nil)
	if err != nil {
		return err
	}

	return nil
}
