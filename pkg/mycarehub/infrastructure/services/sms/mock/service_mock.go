package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/silcomms"
)

// SMSServiceMock mocks the SILComms SMS service mock methods
type SMSServiceMock struct {
	MockSendSMSFn func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error)
}

// NewSMSServiceMock initializes the mock service
func NewSMSServiceMock() *SMSServiceMock {
	UUID := uuid.New().String()
	return &SMSServiceMock{
		MockSendSMSFn: func(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
			return &silcomms.BulkSMSResponse{
				GUID:       UUID,
				Sender:     gofakeit.BeerName(),
				Message:    gofakeit.BeerStyle(),
				Recipients: []string{gofakeit.Phone()},
				State:      gofakeit.BS(),
				SMS:        []string{"test"},
				Created:    "2022-10-07T14:19:43.244460+03:00",
				Updated:    "2022-10-07T14:19:43.244460+03:00",
			}, nil
		},
	}
}

// SendSMS mocks the SILComms sending SMS service method
func (m *SMSServiceMock) SendSMS(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
	return m.MockSendSMSFn(ctx, message, recipients)
}
