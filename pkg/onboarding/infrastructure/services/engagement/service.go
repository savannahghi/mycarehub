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
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
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
)

// ServiceEngagement represents engagement usecases
type ServiceEngagement interface {
	PublishKYCNudge(
		uid string,
		payload base.Nudge,
	) (*http.Response, error)

	PublishKYCFeedItem(
		uid string,
		payload base.Item,
	) (*http.Response, error)

	ResolveDefaultNudgeByTitle(
		UID string,
		flavour base.Flavour,
		nudgeTitle string,
	) error

	SendMail(
		email string,
		message string,
		subject string,
	) error

	SendAlertToSupplier(input resources.EmailNotificationPayload) error

	NotifyAdmins(input resources.EmailNotificationPayload) error

	GenerateAndSendOTP(
		ctx context.Context,
		phone string,
	) (*base.OtpResponse, error)

	SendRetryOTP(
		ctx context.Context,
		msisdn string,
		retryStep int,
	) (*base.OtpResponse, error)

	VerifyOTP(ctx context.Context, phone, OTP string) (bool, error)

	VerifyEmailOTP(ctx context.Context, email, OTP string) (bool, error)
}

// ServiceEngagementImpl represents engagement usecases
type ServiceEngagementImpl struct {
	Engage  extension.ISCClientExtension
	baseExt extension.BaseExtension
}

// NewServiceEngagementImpl returns new instance of ServiceEngagementImpl
func NewServiceEngagementImpl(eng extension.ISCClientExtension, ext extension.BaseExtension) ServiceEngagement {
	return &ServiceEngagementImpl{Engage: eng, baseExt: ext}
}

// PublishKYCNudge calls the `engagement service` to publish
// a KYC nudge
func (en *ServiceEngagementImpl) PublishKYCNudge(
	uid string,
	payload base.Nudge,
) (*http.Response, error) {
	return en.Engage.MakeRequest(
		http.MethodPost,
		fmt.Sprintf(publishNudge, uid),
		payload,
	)
}

// PublishKYCFeedItem calls the `engagement service` to publish
// a KYC feed item
func (en *ServiceEngagementImpl) PublishKYCFeedItem(
	uid string,
	payload base.Item,
) (*http.Response, error) {
	return en.Engage.MakeRequest(
		http.MethodPost,
		fmt.Sprintf(publishItem, uid),
		payload,
	)
}

// ResolveDefaultNudgeByTitle calls the `engagement service`
// to resolve any default nudge by its `Title`
func (en *ServiceEngagementImpl) ResolveDefaultNudgeByTitle(
	UID string,
	flavour base.Flavour,
	nudgeTitle string,
) error {
	resp, err := en.Engage.MakeRequest(
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

	resp, err := en.Engage.MakeRequest(
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
func (en *ServiceEngagementImpl) SendAlertToSupplier(input resources.EmailNotificationPayload) error {
	var writer bytes.Buffer
	t := template.Must(template.New("acknowledgementKYCEmail").Parse(utils.AcknowledgementKYCEmail))
	_ = t.Execute(&writer, resources.EmailNotificationPayload{
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

	resp, err := en.Engage.MakeRequest(http.MethodPost, sendEmail, body)

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
func (en *ServiceEngagementImpl) NotifyAdmins(input resources.EmailNotificationPayload) error {
	adminEmail, err := base.GetEnvVar("SAVANNAH_ADMIN_EMAIL")
	if err != nil {
		return err
	}

	var writer bytes.Buffer
	t := template.Must(template.New("adminKYCSubmittedEmail").Parse(utils.AdminKYCSubmittedEmail))
	_ = t.Execute(&writer, resources.EmailNotificationPayload{
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

	resp, err := en.Engage.MakeRequest(http.MethodPost, sendEmail, body)

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
) (*base.OtpResponse, error) {
	body := map[string]interface{}{
		"msisdn": phone,
	}
	resp, err := en.Engage.MakeRequest(http.MethodPost, SendOtp, body)
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
	return &base.OtpResponse{OTP: OTP}, nil
}

// SendRetryOTP generates fallback OTPs when Africa is talking sms fails
func (en *ServiceEngagementImpl) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
) (*base.OtpResponse, error) {
	phoneNumber, err := en.baseExt.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	body := map[string]interface{}{
		"msisdn":    phoneNumber,
		"retryStep": retryStep,
	}
	resp, err := en.Engage.MakeRequest(http.MethodPost, SendRetryOtp, body)
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

	return &base.OtpResponse{OTP: RetryOTP}, nil
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

	resp, err := en.Engage.MakeRequest(http.MethodPost, VerifyOTPEndPoint, verifyPayload)
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

	resp, err := en.Engage.MakeRequest(http.MethodPost, VerifyEmailOtp, verifyPayload)
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
