package mock

import (
	"gitlab.slade360emr.com/go/base"
)

// FakeServiceMailgun is an `Mailgun` service mock .
type FakeServiceMailgun struct {
	SendMailFn    func(email string, message string, subject string) error
	FetchClientFn func() *base.InterServiceClient
}

// SendMail ...
func (f *FakeServiceMailgun) SendMail(email string, message string, subject string) error {
	return f.SendMailFn(email, message, subject)
}

// FetchClient ...
func (f *FakeServiceMailgun) FetchClient() *base.InterServiceClient {
	return f.FetchClientFn()
}
