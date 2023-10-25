package mail

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/savannahghi/serverutils"
)

var (
	mailGunFrom = serverutils.MustGetEnvVar("MAILGUN_FROM")
	mailGunTo   = serverutils.MustGetEnvVar("MYCAREHUB_ADMIN_EMAIL")
)

// IServiceMail holds the methods to interact with the MailGuns service
type IServiceMail interface {
	SendFeedback(ctx context.Context, subject, feedbackMessage string) (bool, error)
}

// IMailgunClient defines the methods used to communicate with the Mailgun service
type IMailgunClient interface {
	NewMessage(from string, subject string, text string, to ...string) *mailgun.Message
	Send(ctx context.Context, m *mailgun.Message) (string, string, error)
	SetHtml(html string)
}

// MailgunServiceImpl is a client for the Mailgun service
type MailgunServiceImpl struct {
	client IMailgunClient
}

// NewServiceMail initializes Mailgun client
func NewServiceMail(client IMailgunClient) *MailgunServiceImpl {
	return &MailgunServiceImpl{
		client: client,
	}
}

// SendFeedback sends an email to the feedback email address
func (mg *MailgunServiceImpl) SendFeedback(ctx context.Context, subject, feedbackMessage string) (bool, error) {
	m := mg.client.NewMessage(mailGunFrom, subject, "", mailGunTo)
	mg.client.SetHtml(feedbackMessage)

	_, _, err := mg.client.Send(ctx, m)
	if err != nil {
		return false, err
	}

	return true, nil
}
