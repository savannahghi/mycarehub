package usecases

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"gopkg.in/yaml.v2"
)

const otpService = "otp"

// OTP service endpoints
const (
	SendEmail    = "internal/send_email"
	SendRetryOtp = "internal/send_retry_otp/"
	SendOtp      = "internal/send_otp/"
)

// OTPUseCases represent the business logic required for management of OTP
type OTPUseCases interface {
	SendRetryOTP(ctx context.Context, msisdn string, retryStep int) (string, error)
	GenerateAndSendOTP(phone string) (string, error)
	// TODO consider moving this to OTP service or making an isc for it
	VerifyEmailOTP(ctx context.Context, msisdn, otp, flavour string) (bool, error)
}

// OTPUseCasesImpl represents OTP usecases
type OTPUseCasesImpl struct {
	Otp *base.InterServiceClient
}

// NewOTPUseCasesImpl returns new instance of OTPUseCasesImpl
func NewOTPUseCasesImpl(r repository.OnboardingRepository) *OTPUseCasesImpl {

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

	return &OTPUseCasesImpl{Otp: otpClient}
}

// GenerateAndSendOTP creates a new otp and sends it to the provided phone number.
func (o *OTPUseCasesImpl) GenerateAndSendOTP(ctx context.Context, phone string) (string, error) {
	body := map[string]interface{}{
		"msisdn": phone,
	}
	defaultOTP := ""
	resp, err := o.Otp.MakeRequest(http.MethodPost, SendOtp, body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to generate and send otp: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return defaultOTP, fmt.Errorf("unable to generate and send otp, with status code %v", resp.StatusCode)
	}
	code, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return defaultOTP, fmt.Errorf("unable to convert response to string: %v", err)
	}

	return string(code), nil
}
