package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/savannahghi/profileutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
)

// OTP service endpoints
const (
	SendRetryOtp   = "internal/send_retry_otp/"
	SendOtp        = "internal/send_otp/"
	VerifyEmailOtp = "internal/verify_email_otp/"
	// VerifyOTPEndPoint ISC endpoint to verify OTP
	VerifyOTPEndPoint = "internal/verify_otp/"
)

// ServiceOTP represent the business logic required for management of OTP
type ServiceOTP interface {
	GenerateAndSendOTP(
		ctx context.Context,
		phone string,
	) (*profileutils.OtpResponse, error)
	SendRetryOTP(
		ctx context.Context,
		msisdn string,
		retryStep int,
	) (*profileutils.OtpResponse, error)
	VerifyOTP(ctx context.Context, phone, OTP string) (bool, error)
	VerifyEmailOTP(ctx context.Context, email, OTP string) (bool, error)
}

// ServiceOTPImpl represents OTP usecases
type ServiceOTPImpl struct {
	OtpExt  extension.ISCClientExtension
	baseExt extension.BaseExtension
}

// NewOTPService returns new instance of ServiceOTPImpl
func NewOTPService(otp extension.ISCClientExtension, ext extension.BaseExtension) *ServiceOTPImpl {
	return &ServiceOTPImpl{OtpExt: otp, baseExt: ext}
}

// GenerateAndSendOTP creates a new otp and sends it to the provided phone number.
func (o *ServiceOTPImpl) GenerateAndSendOTP(
	ctx context.Context,
	phone string,
) (*profileutils.OtpResponse, error) {
	body := map[string]interface{}{
		"msisdn": phone,
	}
	resp, err := o.OtpExt.MakeRequest(ctx, http.MethodPost, SendOtp, body)
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
func (o *ServiceOTPImpl) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
) (*profileutils.OtpResponse, error) {
	phoneNumber, err := o.baseExt.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	body := map[string]interface{}{
		"msisdn":    phoneNumber,
		"retryStep": retryStep,
	}
	resp, err := o.OtpExt.MakeRequest(ctx, http.MethodPost, SendRetryOtp, body)
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
func (o *ServiceOTPImpl) VerifyOTP(ctx context.Context, phone, otp string) (bool, error) {
	normalized, err := o.baseExt.NormalizeMSISDN(phone)
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

	resp, err := o.OtpExt.MakeRequest(ctx, http.MethodPost, VerifyOTPEndPoint, verifyPayload)
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
func (o *ServiceOTPImpl) VerifyEmailOTP(ctx context.Context, email, otp string) (bool, error) {

	type VerifyOTP struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verificationCode"`
	}

	verifyPayload := VerifyOTP{
		Email:            email,
		VerificationCode: otp,
	}

	resp, err := o.OtpExt.MakeRequest(ctx, http.MethodPost, VerifyEmailOtp, verifyPayload)
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
