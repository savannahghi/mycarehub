package mock

import (
	"context"
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
	NotifySupplierOnSuspensionFn func(input resources.EmailNotificationPayload) error
	NotifyAdminsFn               func(input resources.EmailNotificationPayload) error
	GenerateAndSendOTPFn         func(
		ctx context.Context,
		phone string,
	) (*base.OtpResponse, error)

	SendRetryOTPFn func(
		ctx context.Context,
		msisdn string,
		retryStep int,
	) (*base.OtpResponse, error)

	VerifyOTPFn func(ctx context.Context, phone, OTP string) (bool, error)

	VerifyEmailOTPFn func(ctx context.Context, email, OTP string) (bool, error)

	SendSMSFn func(phoneNumbers []string, message string) error
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

// GenerateAndSendOTP ...
func (f *FakeServiceEngagement) GenerateAndSendOTP(
	ctx context.Context,
	phone string,
) (*base.OtpResponse, error) {
	return f.GenerateAndSendOTPFn(ctx, phone)
}

// SendRetryOTP ...
func (f *FakeServiceEngagement) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
) (*base.OtpResponse, error) {
	return f.SendRetryOTPFn(ctx, msisdn, retryStep)
}

// VerifyOTP ...
func (f *FakeServiceEngagement) VerifyOTP(ctx context.Context, phone, OTP string) (bool, error) {
	return f.VerifyOTPFn(ctx, phone, OTP)
}

// VerifyEmailOTP ...
func (f *FakeServiceEngagement) VerifyEmailOTP(ctx context.Context, email, OTP string) (bool, error) {
	return f.VerifyEmailOTPFn(ctx, email, OTP)
}

// NotifySupplierOnSuspension ...
func (f *FakeServiceEngagement) NotifySupplierOnSuspension(input resources.EmailNotificationPayload) error {
	return f.NotifySupplierOnSuspensionFn(input)
}

// SendSMS ...
func (f *FakeServiceEngagement) SendSMS(phoneNumbers []string, message string) error {
	return f.SendSMSFn(phoneNumbers, message)
}
