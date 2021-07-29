package engagement

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"

	"github.com/asaskevich/govalidator"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
)

const (
	// Feed ISC paths
	publishNudge         = "feed/%s/PRO/false/nudges/"
	publishItem          = "feed/%s/PRO/false/items/"
	resolveDefaultNudges = "feed/%s/%s/false/defaultnudges/%s/resolve/"

	// Communication ISC paths
	sendEmail = "internal/send_email"

	// SendRetryOtp ISC endpoint to send a retry OTP
	SendRetryOtp = "internal/send_retry_otp/"
	// SendOtp ISC endpoint to send OTP
	SendOtp = "internal/send_otp/"
	// VerifyEmailOtp ISC endpoint to verify email OTP
	VerifyEmailOtp = "internal/verify_email_otp/"
	// VerifyOTPEndPoint ISC endpoint to verify OTP
	VerifyOTPEndPoint = "internal/verify_otp/"

	sendSMS = "internal/send_sms"
)

// ServiceEngagement represents engagement usecases
type ServiceEngagement interface {
	PublishKYCNudge(
		ctx context.Context,
		uid string,
		payload feedlib.Nudge,
	) (*http.Response, error)

	PublishKYCFeedItem(
		ctx context.Context,
		uid string,
		payload feedlib.Item,
	) (*http.Response, error)

	ResolveDefaultNudgeByTitle(
		ctx context.Context,
		UID string,
		flavour feedlib.Flavour,
		nudgeTitle string,
	) error

	SendMail(
		ctx context.Context,
		email string,
		message string,
		subject string,
	) error

	SendAlertToSupplier(ctx context.Context, input dto.EmailNotificationPayload) error

	NotifyAdmins(ctx context.Context, input dto.EmailNotificationPayload) error

	NotifySupplierOnSuspension(ctx context.Context, input dto.EmailNotificationPayload) error

	GenerateAndSendOTP(
		ctx context.Context,
		phone string,
		appID *string,
	) (*profileutils.OtpResponse, error)

	SendRetryOTP(
		ctx context.Context,
		msisdn string,
		retryStep int,
		appID *string,
	) (*profileutils.OtpResponse, error)

	VerifyOTP(ctx context.Context, phone, OTP string) (bool, error)

	VerifyEmailOTP(ctx context.Context, email, OTP string) (bool, error)

	SendSMS(ctx context.Context, phoneNumbers []string, message string) error
}

// ServiceEngagementImpl represents engagement usecases
type ServiceEngagementImpl struct {
	Engage  extension.ISCClientExtension
	baseExt extension.BaseExtension
}

// NewServiceEngagementImpl returns new instance of ServiceEngagementImpl
func NewServiceEngagementImpl(
	eng extension.ISCClientExtension,
	ext extension.BaseExtension,
) *ServiceEngagementImpl {
	return &ServiceEngagementImpl{
		Engage:  eng,
		baseExt: ext,
	}
}

// PublishKYCNudge calls the `engagement service` to publish
// a KYC nudge
func (en *ServiceEngagementImpl) PublishKYCNudge(
	ctx context.Context,
	uid string,
	payload feedlib.Nudge,
) (*http.Response, error) {
	return en.Engage.MakeRequest(ctx,
		http.MethodPost,
		fmt.Sprintf(publishNudge, uid),
		payload,
	)
}

// PublishKYCFeedItem calls the `engagement service` to publish
// a KYC feed item
func (en *ServiceEngagementImpl) PublishKYCFeedItem(
	ctx context.Context,
	uid string,
	payload feedlib.Item,
) (*http.Response, error) {
	return en.Engage.MakeRequest(ctx,
		http.MethodPost,
		fmt.Sprintf(publishItem, uid),
		payload,
	)
}

// ResolveDefaultNudgeByTitle calls the `engagement service`
// to resolve any default nudge by its `Title`
func (en *ServiceEngagementImpl) ResolveDefaultNudgeByTitle(
	ctx context.Context,
	UID string,
	flavour feedlib.Flavour,
	nudgeTitle string,
) error {
	resp, err := en.Engage.MakeRequest(ctx,
		http.MethodPatch,
		fmt.Sprintf(
			resolveDefaultNudges,
			UID,
			flavour,
			nudgeTitle,
		),
		nil,
	)

	if err != nil {
		return exceptions.ResolveNudgeErr(
			err,
			flavour,
			nudgeTitle,
			nil,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return exceptions.ResolveNudgeErr(
			fmt.Errorf("unexpected status code %v", resp.StatusCode),
			flavour,
			nudgeTitle,
			&resp.StatusCode,
		)
	}

	return nil
}

// SendMail sends emails to communicate to our users
func (en *ServiceEngagementImpl) SendMail(
	ctx context.Context,
	email string,
	message string,
	subject string,
) error {
	if !govalidator.IsEmail(email) {
		return fmt.Errorf("invalid email address: %v", email)
	}

	body := map[string]interface{}{
		"to":      []string{email},
		"text":    message,
		"subject": subject,
	}

	resp, err := en.Engage.MakeRequest(ctx,
		http.MethodPost,
		sendEmail,
		body,
	)
	if err != nil {
		return fmt.Errorf("unable to send email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send email : %w, with status code %v",
			err,
			resp.StatusCode,
		)
	}

	return nil
}

//SendAlertToSupplier send email to supplier to acknowledgement receipt of
// KYC request/documents.
func (en *ServiceEngagementImpl) SendAlertToSupplier(ctx context.Context, input dto.EmailNotificationPayload) error {
	var writer bytes.Buffer
	t := template.Must(template.New("acknowledgementKYCEmail").Parse(utils.AcknowledgementKYCEmail))
	_ = t.Execute(&writer, dto.EmailNotificationPayload{
		SupplierName: input.SupplierName,
		PartnerType:  input.PartnerType,
		AccountType:  input.AccountType,
		EmailBody:    input.EmailBody,
		EmailAddress: input.EmailAddress,
		PrimaryPhone: input.PrimaryPhone,
	})

	body := map[string]interface{}{
		"to":      []string{input.EmailAddress},
		"text":    writer.String(),
		"subject": input.SubjectTitle,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send alert to supplier email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send alert to supplier email: %w", err)
	}

	return nil
}

//NotifyAdmins send email to admin notifying them of new
// KYC Request.
func (en *ServiceEngagementImpl) NotifyAdmins(ctx context.Context, input dto.EmailNotificationPayload) error {

	adminEmail, err := serverutils.GetEnvVar("SAVANNAH_ADMIN_EMAIL")
	if err != nil {
		return err
	}

	var writer bytes.Buffer
	t := template.Must(template.New("adminKYCSubmittedEmail").Parse(utils.AdminKYCSubmittedEmail))
	_ = t.Execute(&writer, dto.EmailNotificationPayload{
		SupplierName: input.SupplierName,
		PartnerType:  input.PartnerType,
		AccountType:  input.AccountType,
		EmailBody:    input.EmailBody,
		EmailAddress: input.EmailAddress,
		PrimaryPhone: input.PrimaryPhone,
	})

	body := map[string]interface{}{
		"to":      []string{adminEmail},
		"text":    writer.String(),
		"subject": input.SubjectTitle,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, sendEmail, body)
	if err != nil {
		return fmt.Errorf("unable to send alert to admin email: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send alert to admin email: %w", err)
	}

	return nil
}

// GenerateAndSendOTP creates a new otp and sends it to the provided phone number.
func (en *ServiceEngagementImpl) GenerateAndSendOTP(
	ctx context.Context,
	phone string,
	appID *string,
) (*profileutils.OtpResponse, error) {
	body := map[string]interface{}{
		"msisdn": phone,
		"appId":  appID,
	}
	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, SendOtp, body)
	if err != nil {
		return nil, exceptions.GenerateAndSendOTPError(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unable to generate and send otp, with status code %v", resp.StatusCode,
		)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var OTP string
	err = json.Unmarshal(code, &OTP)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}
	return &profileutils.OtpResponse{OTP: OTP}, nil
}

// SendRetryOTP generates fallback OTPs when Africa is talking sms fails
func (en *ServiceEngagementImpl) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
	appID *string,
) (*profileutils.OtpResponse, error) {
	phoneNumber, err := en.baseExt.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	body := map[string]interface{}{
		"msisdn":    phoneNumber,
		"retryStep": retryStep,
		"appId":     appID,
	}
	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, SendRetryOtp, body)
	if err != nil {
		return nil, exceptions.GenerateAndSendOTPError(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unable to generate and send fallback otp, with status code %v",
			resp.StatusCode,
		)
	}

	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to convert response to string: %v", err)
	}

	var RetryOTP string
	err = json.Unmarshal(code, &RetryOTP)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OTP: %v", err)
	}

	return &profileutils.OtpResponse{OTP: RetryOTP}, nil
}

// VerifyOTP takes a phone number and an OTP and checks for the validity of the OTP code
func (en *ServiceEngagementImpl) VerifyOTP(ctx context.Context, phone, otp string) (bool, error) {
	normalized, err := en.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, fmt.Errorf("invalid phone format: %w", err)
	}

	type VerifyOTP struct {
		Msisdn           string `json:"msisdn"`
		VerificationCode string `json:"verificationCode"`
	}

	verifyPayload := VerifyOTP{
		Msisdn:           *normalized,
		VerificationCode: otp,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, VerifyOTPEndPoint, verifyPayload)
	if err != nil {
		return false, fmt.Errorf(
			"can't complete OTP verification request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unable to verify OTP : %w, with status code %v", err, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("can't read OTP response data: %w", err)
	}

	type otpResponse struct {
		IsVerified bool `json:"IsVerified"`
	}

	var r otpResponse
	err = json.Unmarshal(data, &r)
	if err != nil {
		return false, fmt.Errorf(
			"can't unmarshal OTP response data from JSON: %w", err)
	}
	return r.IsVerified, nil
}

// VerifyEmailOTP checks the otp provided matches the one sent to the user via email address
func (en *ServiceEngagementImpl) VerifyEmailOTP(ctx context.Context, email, otp string) (bool, error) {

	type VerifyOTP struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verificationCode"`
	}

	verifyPayload := VerifyOTP{
		Email:            email,
		VerificationCode: otp,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, VerifyEmailOtp, verifyPayload)
	if err != nil {
		return false, fmt.Errorf(
			"can't complete OTP verification request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unable to verify OTP : %w, with status code %v", err, resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("can't read OTP response data: %w", err)
	}

	type otpResponse struct {
		IsVerified bool `json:"IsVerified"`
	}

	var r otpResponse
	err = json.Unmarshal(data, &r)
	if err != nil {
		return false, fmt.Errorf(
			"can't unmarshal OTP response data from JSON: %w", err)
	}

	return r.IsVerified, nil

}

//NotifySupplierOnSuspension send email to supplier notifying him of the
// suspension.
func (en *ServiceEngagementImpl) NotifySupplierOnSuspension(ctx context.Context, input dto.EmailNotificationPayload) error {
	var writer bytes.Buffer
	t := template.Must(template.New("supplierSuspensionEmail").Parse(utils.SupplierSuspensionEmail))
	_ = t.Execute(&writer, dto.EmailNotificationPayload{
		SupplierName: input.SupplierName,
		PartnerType:  input.PartnerType,
		AccountType:  input.AccountType,
		EmailBody:    input.EmailBody,
		EmailAddress: input.EmailAddress,
		PrimaryPhone: input.PrimaryPhone,
	})

	body := map[string]interface{}{
		"to":      []string{input.EmailAddress},
		"text":    writer.String(),
		"subject": input.SubjectTitle,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, sendEmail, body)

	if err != nil {
		return fmt.Errorf("unable to send alert to supplier email: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send alert to supplier email: %w", err)
	}

	return nil
}

// SendSMS does the actual delivery of messages to the provided phone numbers
func (en *ServiceEngagementImpl) SendSMS(ctx context.Context, phoneNumbers []string, message string) error {
	type PayloadRequest struct {
		To      []string           `json:"to"`
		Message string             `json:"message"`
		Sender  enumutils.SenderID `json:"sender"`
	}

	requestPayload := PayloadRequest{
		To:      phoneNumbers,
		Message: message,
		Sender:  enumutils.SenderIDBewell,
	}

	resp, err := en.Engage.MakeRequest(ctx, http.MethodPost, sendSMS, requestPayload)
	if err != nil {
		return fmt.Errorf("unable to send sms: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to send sms, with status code %v", resp.StatusCode)
	}

	return nil
}
