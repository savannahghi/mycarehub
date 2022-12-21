package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	serviceSMS "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms"
	serviceTwilio "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio"
	"github.com/savannahghi/serverutils"

	"github.com/savannahghi/profileutils"
)

const (
	otpMessage = "%s is your %v verification code %v"
)

var (
	consumerAppIdentifier = serverutils.MustGetEnvVar("CONSUMER_APP_IDENTIFIER")
	proAppIdentifier      = serverutils.MustGetEnvVar("PRO_APP_IDENTIFIER")
	consumerAppName       = serverutils.MustGetEnvVar("CONSUMER_APP_NAME")
	proAppName            = serverutils.MustGetEnvVar("PRO_APP_NAME")
)

// IverifyPhone specifies the method signature for verifying phone via OTP.
type IverifyPhone interface {
	VerifyPhoneNumber(ctx context.Context, phone string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error)
}

// ISendOTP is used to send an OTP
type ISendOTP interface {
	// delegate to GenerateOTP
	// clients should call: SendOTP
	// send on the primary channel
	// metrics
	// the middle parameter is an error code e.g if rate limited
	SendOTP(
		ctx context.Context,
		phoneNumber string,
		code string,
		message string,
	) (string, error)

	GenerateAndSendOTP(
		ctx context.Context,
		username string,
		flavour feedlib.Flavour,
	) (string, error)

	GenerateRetryOTP(
		ctx context.Context,
		payload *dto.SendRetryOTPPayload,
	) (string, error)
}

// UsecaseOTP defines otp service usecases interface
type UsecaseOTP interface {
	ISendOTP
	IVerifyOTP
	IverifyPhone
}

// IVerifyOTP specifies the method responsible for verifying the OTP
type IVerifyOTP interface {
	VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error)
}

// UseCaseOTPImpl is the OTP service implementation
type UseCaseOTPImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	ExternalExt extension.ExternalMethodsExtension
	SMS         serviceSMS.IServiceSMS
	Twilio      serviceTwilio.ITwilioService
}

// NewOTPUseCase initializes a new OTP service
func NewOTPUseCase(
	create infrastructure.Create,
	query infrastructure.Query,
	externalExt extension.ExternalMethodsExtension,
	sms serviceSMS.IServiceSMS,
	twilio serviceTwilio.ITwilioService,
) *UseCaseOTPImpl {
	return &UseCaseOTPImpl{
		Create:      create,
		Query:       query,
		ExternalExt: externalExt,
		SMS:         sms,
		Twilio:      twilio,
	}
}

// GenerateAndSendOTP generates and send an otp to the intended user
func (o *UseCaseOTPImpl) GenerateAndSendOTP(
	ctx context.Context,
	username string,
	flavour feedlib.Flavour,
) (string, error) {

	if !flavour.IsValid() {
		return "", exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err := o.Query.GetUserProfileByUsername(ctx, username)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.UserNotFoundError(err)
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", fmt.Errorf("failed to generate an OTP")
	}

	phone, err := o.Query.GetContactByUserID(ctx, userProfile.ID, "PHONE")
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.ContactNotFoundErr(err)
	}

	var message string
	switch flavour {
	case feedlib.FlavourConsumer:
		message = fmt.Sprintf(otpMessage, otp, consumerAppName, consumerAppIdentifier)
	case feedlib.FlavourPro:
		message = fmt.Sprintf(otpMessage, otp, proAppName, proAppIdentifier)
	}

	otp, err = o.SendOTP(ctx, phone.ContactValue, otp, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", err
	}

	otpDataPayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Minute * 10),
		Channel:     "SMS",
		Flavour:     flavour,
		PhoneNumber: phone.ContactValue,
		OTP:         otp,
	}

	err = o.Create.SaveOTP(ctx, otpDataPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", fmt.Errorf("failed to save otp: %w", err)
	}

	return otp, nil
}

// VerifyOTP verifies whether the supplied OTP is valid
func (o *UseCaseOTPImpl) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return o.Query.VerifyOTP(ctx, payload)
}

// VerifyPhoneNumber checks validity of a phone number by sending an OTP to it
func (o *UseCaseOTPImpl) VerifyPhoneNumber(ctx context.Context, phone string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
	phoneNumber, err := converterandformatter.NormalizeMSISDN(phone)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	exists, err := o.Query.CheckIfPhoneNumberExists(ctx, *phoneNumber, true, flavour)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to check if phone exists: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("the provided phone number does not exist")
	}

	userProfile, err := o.Query.GetUserProfileByPhoneNumber(ctx, *phoneNumber)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.UserNotFoundError(err)
	}

	otp, err := utils.GenerateOTP()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to generate an OTP")
	}

	var message string
	switch flavour {
	case feedlib.FlavourConsumer:
		message = fmt.Sprintf(otpMessage, otp, consumerAppName, consumerAppIdentifier)
	case feedlib.FlavourPro:
		message = fmt.Sprintf(otpMessage, otp, proAppName, proAppIdentifier)
	}

	otp, err = o.SendOTP(ctx, phone, otp, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	otpDataPayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Minute * 10),
		Channel:     "SMS",
		Flavour:     flavour,
		PhoneNumber: phone,
		OTP:         otp,
	}

	err = o.Create.SaveOTP(ctx, otpDataPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to save otp: %w", err)
	}

	return &profileutils.OtpResponse{
		OTP: otp,
	}, nil
}

// GenerateRetryOTP generates fallback OTPs when Africa is talking sms fails
func (o *UseCaseOTPImpl) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	retryResponseOTP, err := utils.GenerateOTP()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", err
	}

	userProfile, err := o.Query.GetUserProfileByUsername(ctx, payload.Username)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.UserNotFoundError(err)
	}

	phone, err := o.Query.GetContactByUserID(ctx, userProfile.ID, "PHONE")
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", exceptions.ContactNotFoundErr(err)
	}

	// send retry otp
	_, err = o.SMS.SendSMS(ctx, retryResponseOTP, []string{phone.ContactValue})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", fmt.Errorf("failed to send OTP verification code to recipient %w", err)
	}

	otpResponsePayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Hour * 1),
		Channel:     "SMS",
		Flavour:     payload.Flavour,
		PhoneNumber: phone.ContactValue,
		OTP:         retryResponseOTP,
	}

	err = o.Create.SaveOTP(ctx, otpResponsePayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return "", fmt.Errorf("failed to save otp: %w", err)
	}

	return retryResponseOTP, nil
}

// SendOTP sends an OTP message to the specified phonenumber. It checks to see whether the provided
// phone number is Kenyan and if true, it uses AIT else for foreign numbers, it uses twilio to send
// the otp
func (o *UseCaseOTPImpl) SendOTP(
	ctx context.Context,
	phoneNumber string,
	code string,
	message string,
) (string, error) {
	if interserviceclient.IsKenyanNumber(phoneNumber) {
		_, err := o.SMS.SendSMS(ctx, message, []string{phoneNumber})
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return "", fmt.Errorf("failed to send OTP verification code to recipient: %w", err)
		}

	} else {
		// Make the request to twilio
		err := o.Twilio.SendSMSViaTwilio(ctx, phoneNumber, message)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return "", fmt.Errorf("sms not sent via twilio: %w", err)
		}
	}
	return code, nil
}
