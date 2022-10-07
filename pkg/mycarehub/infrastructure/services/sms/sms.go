package sms

import (
	"context"

	"github.com/savannahghi/silcomms"
)

// IServiceSMS holds the methods to interact with the SILComms SMS service
type IServiceSMS interface {
	SendSMS(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error)
}

// ISILCommsClient defines the methods used to communicate with the SILComms service
type ISILCommsClient interface {
	SendBulkSMS(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error)
}

// SILCommsServiceImpl is a client for the SILComms service
type SILCommsServiceImpl struct {
	client ISILCommsClient
}

// NewServiceSMS initializes SILComms client
func NewServiceSMS(client ISILCommsClient) *SILCommsServiceImpl {
	return &SILCommsServiceImpl{
		client: client,
	}
}

// SendSMS sends an SMS to many recipients
func (sc *SILCommsServiceImpl) SendSMS(ctx context.Context, message string, recipients []string) (*silcomms.BulkSMSResponse, error) {
	resp, err := sc.client.SendBulkSMS(ctx, message, recipients)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
