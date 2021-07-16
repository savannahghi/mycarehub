package mock

import (
	"context"

	"github.com/savannahghi/interserviceclient"
)

// FakeServiceMessaging is an `Messaging` service mock .
type FakeServiceMessaging struct {
	FetchSMSClientFn    func() *interserviceclient.InterServiceClient
	FetchTwilioClientFn func() *interserviceclient.InterServiceClient
	SendSMSFn           func(ctx context.Context, phoneNumbers []string, message string) error
}

// FetchSMSClient ...
func (f *FakeServiceMessaging) FetchSMSClient() *interserviceclient.InterServiceClient {
	return f.FetchSMSClientFn()
}

// FetchTwilioClient ...
func (f *FakeServiceMessaging) FetchTwilioClient() *interserviceclient.InterServiceClient {
	return f.FetchTwilioClientFn()
}

// SendSMS ...
func (f *FakeServiceMessaging) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.SendSMSFn(ctx, phoneNumbers, message)
}
