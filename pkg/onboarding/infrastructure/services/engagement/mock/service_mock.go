package mock

import (
	"context"
	"net/http"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/profileutils"
	"gitlab.slade360emr.com/go/apiclient"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

// FakeServiceEngagement is an `engagement` service mock .
type FakeServiceEngagement struct {
	PublishKYCNudgeFn            func(ctx context.Context, uid string, payload feedlib.Nudge) (*http.Response, error)
	PublishKYCFeedItemFn         func(ctx context.Context, uid string, payload feedlib.Item) (*http.Response, error)
	ResolveDefaultNudgeByTitleFn func(ctx context.Context, UID string, flavour feedlib.Flavour, nudgeTitle string) error
	SendMailFn                   func(ctx context.Context, email string, message string, subject string) error
	SendAlertToSupplierFn        func(ctx context.Context, input dto.EmailNotificationPayload) error
	NotifySupplierOnSuspensionFn func(ctx context.Context, input dto.EmailNotificationPayload) error
	NotifyAdminsFn               func(ctx context.Context, input dto.EmailNotificationPayload) error
	GenerateAndSendOTPFn         func(
		ctx context.Context,
		phone string,
	) (*profileutils.OtpResponse, error)

	SendRetryOTPFn func(
		ctx context.Context,
		msisdn string,
		retryStep int,
	) (*profileutils.OtpResponse, error)

	VerifyOTPFn func(ctx context.Context, phone, OTP string) (bool, error)

	VerifyEmailOTPFn func(ctx context.Context, email, OTP string) (bool, error)

	SendSMSFn func(ctx context.Context, phoneNumbers []string, message string) error

	GetSladerDataFn func(ctx context.Context, phoneNumber string) (*apiclient.Segment, error)
}

// PublishKYCNudge ...
func (f *FakeServiceEngagement) PublishKYCNudge(
	ctx context.Context,
	uid string,
	payload feedlib.Nudge,
) (*http.Response, error) {
	return f.PublishKYCNudgeFn(ctx, uid, payload)
}

// PublishKYCFeedItem ...
func (f *FakeServiceEngagement) PublishKYCFeedItem(
	ctx context.Context,
	uid string,
	payload feedlib.Item,
) (*http.Response, error) {
	return f.PublishKYCFeedItemFn(ctx, uid, payload)
}

// ResolveDefaultNudgeByTitle ...
func (f *FakeServiceEngagement) ResolveDefaultNudgeByTitle(
	ctx context.Context,
	UID string,
	flavour feedlib.Flavour,
	nudgeTitle string,
) error {
	return f.ResolveDefaultNudgeByTitleFn(
		ctx,
		UID,
		flavour,
		nudgeTitle,
	)
}

// SendMail ...
func (f *FakeServiceEngagement) SendMail(
	ctx context.Context,
	email string,
	message string,
	subject string,
) error {
	return f.SendMailFn(ctx, email, message, subject)
}

// SendAlertToSupplier ...
func (f *FakeServiceEngagement) SendAlertToSupplier(ctx context.Context, input dto.EmailNotificationPayload) error {
	return f.SendAlertToSupplierFn(ctx, input)
}

// NotifyAdmins ...
func (f *FakeServiceEngagement) NotifyAdmins(ctx context.Context, input dto.EmailNotificationPayload) error {
	return f.NotifyAdminsFn(ctx, input)
}

// GenerateAndSendOTP ...
func (f *FakeServiceEngagement) GenerateAndSendOTP(
	ctx context.Context,
	phone string,
) (*profileutils.OtpResponse, error) {
	return f.GenerateAndSendOTPFn(ctx, phone)
}

// SendRetryOTP ...
func (f *FakeServiceEngagement) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
) (*profileutils.OtpResponse, error) {
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
func (f *FakeServiceEngagement) NotifySupplierOnSuspension(ctx context.Context, input dto.EmailNotificationPayload) error {
	return f.NotifySupplierOnSuspensionFn(ctx, input)
}

// SendSMS ...
func (f *FakeServiceEngagement) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	return f.SendSMSFn(ctx, phoneNumbers, message)
}

// GetSladerData ...
func (f *FakeServiceEngagement) GetSladerData(ctx context.Context, phoneNumber string) (*apiclient.Segment, error) {
	return f.GetSladerDataFn(ctx, phoneNumber)
}
