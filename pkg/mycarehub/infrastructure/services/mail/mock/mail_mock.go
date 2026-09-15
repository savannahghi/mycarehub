package mock

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
)

// MailGunClientMock mocks the mailgun's client service library implementations
type MailGunClientMock struct {
	MockNewMessageFn func(from, subject, text string, to ...string) *mailgun.Message
	MockSendFn       func(ctx context.Context, m *mailgun.Message) (string, string, error)
	MockSetHtmlFn    func(html string)
}

// NewMailGunClientMock initializes the mock client
func NewMailGunClientMock() *MailGunClientMock {
	return &MailGunClientMock{
		MockNewMessageFn: func(from string, subject string, text string, to ...string) *mailgun.Message {
			return &mailgun.Message{}
		},
		MockSendFn: func(ctx context.Context, m *mailgun.Message) (string, string, error) {
			return "", "", nil
		},
		MockSetHtmlFn: func(html string) {},
	}
}

// NewMessage mocks the mailgun's client service library implementations
func (m *MailGunClientMock) NewMessage(from string, subject string, text string, to ...string) *mailgun.Message {
	return m.MockNewMessageFn(from, subject, text, to...)
}

// Send mocks the mailgun's client service library implementations
func (m *MailGunClientMock) Send(ctx context.Context, ms *mailgun.Message) (string, string, error) {
	return m.MockSendFn(ctx, ms)
}

// SetHtml mocks the mailgun client SetHTML method
func (m *MailGunClientMock) SetHtml(html string) {}
