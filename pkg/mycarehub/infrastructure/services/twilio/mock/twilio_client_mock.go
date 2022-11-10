package mock

import (
	"net/url"

	twilioClient "github.com/kevinburke/twilio-go"
)

// TwilioClientMock mocks the twilio's client service library implementations
type TwilioClientMock struct {
	MockSendMessageFn func(from string, to string, body string, mediaURLs []*url.URL) (*twilioClient.Message, error)
}

// NewTwilioClientMock initializes the mock client
func NewTwilioClientMock() *TwilioClientMock {
	return &TwilioClientMock{
		MockSendMessageFn: func(from string, to string, body string, mediaURLs []*url.URL) (*twilioClient.Message, error) {
			return &twilioClient.Message{
				From: "+254700000000",
			}, nil
		},
	}
}

// SendMessage mocks the twilio's client service library implementations
func (m *TwilioClientMock) SendMessage(from string, to string, body string, mediaURLs []*url.URL) (*twilioClient.Message, error) {
	return m.MockSendMessageFn(from, to, body, mediaURLs)
}
