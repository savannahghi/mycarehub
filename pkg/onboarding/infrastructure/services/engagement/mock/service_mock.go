package mock

import (
	"net/http"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

// FakeServiceEngagement is an `engagement` service mock .
type FakeServiceEngagement struct {
	PublishKYCNudgeFn            func(uid string, payload base.Nudge) (*http.Response, error)
	PublishKYCFeedItemFn         func(uid string, payload base.Item) (*http.Response, error)
	ResolveDefaultNudgeByTitleFn func(UID string, flavour base.Flavour, nudgeTitle string) error
	SendMailFn                   func(email string, message string, subject string) error
	SendAlertToSupplierFn        func(input resources.EmailNotificationPayload) error
	NotifyAdminsFn               func(input resources.EmailNotificationPayload) error
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
func (f *FakeServiceEngagement) SendAlertToSupplier(input resources.EmailNotificationPayload) error {
	return f.SendAlertToSupplierFn(input)
}

// NotifyAdmins ...
func (f *FakeServiceEngagement) NotifyAdmins(input resources.EmailNotificationPayload) error {
	return f.NotifyAdminsFn(input)
}
