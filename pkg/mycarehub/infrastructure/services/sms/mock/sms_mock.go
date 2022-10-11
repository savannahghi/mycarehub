package mock

import (
	"context"

	"github.com/savannahghi/silcomms"
)

// SMSClientMock mocks the SILComms client library implementations
type SMSClientMock struct {
	MockSendBulkSMSFn func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error)
}

// NewSMSClientMock initializes our client mocks
func NewSMSClientMock() *SMSClientMock {
	return &SMSClientMock{
		MockSendBulkSMSFn: func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
			return &silcomms.BulkSMSResponse{
				GUID:       "",
				Sender:     "",
				Message:    message,
				Recipients: recipients,
				State:      "",
				SMS:        []string{},
				Created:    "",
				Updated:    "",
			}, nil
		},
	}
}

// SendBulkSMS mocks the implementation of sending bulk SMS
func (sc SMSClientMock) SendBulkSMS(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
	return sc.MockSendBulkSMSFn(ctx, message, recipients)
}
