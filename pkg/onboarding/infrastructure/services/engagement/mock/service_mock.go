package mock

import (
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceEngagement is an `engagement` service mock .
type FakeServiceEngagement struct {
	PublishKYCNudgeFn            func(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItemFn         func(uid string, payload base.Item) (*http.Response, error)
	ResolveDefaultNudgeByTitleFn func(UID string, flavour base.Flavour, nudgeTitle string) error
	SendMailFn                   func(email string, message string, subject string) error
	SendAlertToSupplierFn        func(supplierName string, partnerType string, accountType string, subjectTitle string, emailBody string, emailAddress string) error
}

// PublishKYCNudge ...
func (f *FakeServiceEngagement) PublishKYCNudge(
	uid string,
	payload base.Nudge,
) (*http.Response, error) {
	return f.PublishKYCNudgeFn(uid, payload)
}

// PublishKYCFeedItem ...
func (f *FakeServiceEngagement) PublishKYCFeedItem(
	uid string,
	payload base.Item,
) (*http.Response, error) {
	return f.PublishKYCFeedItemFn(uid, payload)
}

// ResolveDefaultNudgeByTitle ...
func (f *FakeServiceEngagement) ResolveDefaultNudgeByTitle(
	UID string,
	flavour base.Flavour,
	nudgeTitle string,
) error {
	return f.ResolveDefaultNudgeByTitleFn(
		UID,
		flavour,
		nudgeTitle,
	)
}

// SendMail ...
func (f *FakeServiceEngagement) SendMail(
	email string,
	message string,
	subject string,
) error {
	return f.SendMailFn(email, message, subject)
}

// SendAlertToSupplier ...
func (f *FakeServiceEngagement) SendAlertToSupplier(
	supplierName string,
	partnerType string,
	accountType string,
	subjectTitle string,
	emailBody string,
	emailAddress string,
) error {
	return f.SendAlertToSupplierFn(supplierName, partnerType, accountType, subjectTitle, emailBody, emailAddress)
}
