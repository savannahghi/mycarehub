package mock

import (
	"context"
)

// TwilioServiceMock mocks the twilio's client service library implementations
type TwilioServiceMock struct {
	MockSendSMSViaTwilioFn func(ctx context.Context, phonenumber, message string) error
}

// NewTwilioServiceMock initializes the mock client
func NewTwilioServiceMock() *TwilioServiceMock {
	return &TwilioServiceMock{
		MockSendSMSViaTwilioFn: func(ctx context.Context, phonenumber, message string) error {
			return nil
		},
	}
}

// SendSMSViaTwilio mocks the twilio's client service library implementations
func (m *TwilioServiceMock) SendSMSViaTwilio(ctx context.Context, phonenumber, message string) error {
	return m.MockSendSMSViaTwilioFn(ctx, phonenumber, message)
}
