package mock

import (
	"github.com/mailgun/mailgun-go"
)

// MailGunClientMock mocks the mailgun's client service library implementations
type MailGunClientMock struct {
	MockNewMessageFn func(from, subject, text string, to ...string) *mailgun.Message
	MockSendFn       func(m *mailgun.Message) (string, string, error)
}

// NewMailGunClientMock initializes the mock client
func NewMailGunClientMock() *MailGunClientMock {
	return &MailGunClientMock{
		MockNewMessageFn: func(from, subject, text string, to ...string) *mailgun.Message {
			return &mailgun.Message{}
		},
		MockSendFn: func(m *mailgun.Message) (string, string, error) {
			return "", "", nil
		},
	}
}

// NewMessage mocks the mailgun's client service library implementations
func (m *MailGunClientMock) NewMessage(from, subject, text string, to ...string) *mailgun.Message {
	return m.MockNewMessageFn(from, subject, text, to...)
}

// Send mocks the mailgun's client service library implementations
func (m *MailGunClientMock) Send(ms *mailgun.Message) (string, string, error) {
	return m.MockSendFn(ms)
}
