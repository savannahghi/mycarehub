package otp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"gopkg.in/yaml.v2"
)

const otpService = "otp"

// OTP service endpoints
const (
	SendRetryOtp   = "internal/send_retry_otp/"
	SendOtp        = "internal/send_otp/"
	VerifyEmailOtp = "internal/verify_email_otp/"
)

// ServiceOTP represent the business logic required for management of OTP
type ServiceOTP interface {
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

// ServiceOTPImpl represents OTP usecases
type ServiceOTPImpl struct {
	Otp     *base.InterServiceClient
	baseExt extension.BaseExtension
}

// NewOTPService returns new instance of ServiceOTPImpl
func NewOTPService(r repository.OnboardingRepository, ext extension.BaseExtension) ServiceOTP {

	var config base.DepsConfig
	//os file and parse it to go type
	file, err := ioutil.ReadFile(filepath.Clean(base.PathToDepsFile()))
	if err != nil {
		log.Errorf("error occured while opening deps file %v", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

	var otpClient *base.InterServiceClient
	otpClient, err = base.SetupISCclient(config, otpService)
	if err != nil {
		log.Panicf("unable to initialize otp inter service client: %s", err)

	}

	return &ServiceOTPImpl{Otp: otpClient, baseExt: ext}
}

// GenerateAndSendOTP creates a new otp and sends it to the provided phone number.
func (o *ServiceOTPImpl) GenerateAndSendOTP(
	ctx context.Context,
	phone string,
) (*base.OtpResponse, error) {
	body := map[string]interface{}{
		"msisdn": phone,
	}
	resp, err := o.Otp.MakeRequest(http.MethodPost, SendOtp, body)
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
func (o *ServiceOTPImpl) SendRetryOTP(
	ctx context.Context,
	msisdn string,
	retryStep int,
) (*base.OtpResponse, error) {
	phoneNumber, err := o.baseExt.NormalizeMSISDN(msisdn)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	body := map[string]interface{}{
		"msisdn":    phoneNumber,
		"retryStep": retryStep,
	}
	resp, err := o.Otp.MakeRequest(http.MethodPost, SendRetryOtp, body)
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
func (o *ServiceOTPImpl) VerifyOTP(ctx context.Context, phone, otp string) (bool, error) {
	return base.VerifyOTP(phone, otp, o.Otp)
}

// VerifyEmailOTP checks the otp provided mathes the one sent to the user via email address
func (o *ServiceOTPImpl) VerifyEmailOTP(ctx context.Context, email, otp string) (bool, error) {

	type VerifyOTP struct {
		Email            string `json:"email"`
		VerificationCode string `json:"verificationCode"`
	}

	verifyPayload := VerifyOTP{
		Email:            email,
		VerificationCode: otp,
	}

	resp, err := o.Otp.MakeRequest(http.MethodPost, VerifyEmailOtp, verifyPayload)
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
