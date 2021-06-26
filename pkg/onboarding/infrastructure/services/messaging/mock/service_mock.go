package mock

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceMessaging is an `Messaging` service mock .
type FakeServiceMessaging struct {
	FetchSMSClientFn    func() *base.InterServiceClient
	FetchTwilioClientFn func() *base.InterServiceClient
	SendSMSFn           func(ctx context.Context, phoneNumbers []string, message string) error
}

// FetchSMSClient ...
func (f *FakeServiceMessaging) FetchSMSClient() *base.InterServiceClient {
	return f.FetchSMSClientFn()
}

// FetchTwilioClient ...
func (f *FakeServiceMessaging) FetchTwilioClient() *base.InterServiceClient {
	return f.FetchTwilioClientFn()
}

// SendSMS ...
func (f *FakeServiceMessaging) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.SendSMSFn(ctx, phoneNumbers, message)
}
