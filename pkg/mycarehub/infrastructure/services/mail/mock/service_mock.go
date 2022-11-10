package mock

import (
	"context"
)

// MailServiceMock mocks the mailgun's client service library implementations
type MailServiceMock struct {
	MockSendFeedbackFn func(ctx context.Context, subject, feedbackMessage string) (bool, error)
}

// NewMailServiceMock initializes the mock client
func NewMailServiceMock() *MailServiceMock {
	return &MailServiceMock{
		MockSendFeedbackFn: func(ctx context.Context, subject, feedbackMessage string) (bool, error) {
			return true, nil
		},
	}
}

// SendFeedback mocks the mailgun's client service library implementations
func (m *MailServiceMock) SendFeedback(ctx context.Context, subject, feedbackMessage string) (bool, error) {
	return m.MockSendFeedbackFn(ctx, subject, feedbackMessage)
}
