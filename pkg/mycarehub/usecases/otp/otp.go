package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"

	// "github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/profileutils"
)

const (
	otpMessage = "%s is your MyCareHub verification code"
)

// IGenerateOTP specifies the method signature for generating an OTP
type IGenerateOTP interface {
	// TODO: ensure generated OTP is valid e.g valid until > generated at
	// metrics
	GenerateOTP(ctx context.Context) (string, error)
}

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
		phoneNumber string,
		flavour feedlib.Flavour,
	) (string, error)

	GenerateRetryOTP(
		ctx context.Context,
		payload *dto.SendRetryOTPPayload,
	) (string, error)
}

// UsecaseOTP defines otp service usecases interface
type UsecaseOTP interface {
	IGenerateOTP
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
}

// NewOTPUseCase initializes a new OTP service
func NewOTPUseCase(
	create infrastructure.Create,
	query infrastructure.Query,
	externalExt extension.ExternalMethodsExtension,
) *UseCaseOTPImpl {
	return &UseCaseOTPImpl{
		Create:      create,
		Query:       query,
		ExternalExt: externalExt,
	}
}

// GenerateAndSendOTP generates and send an otp to the intended user
func (o *UseCaseOTPImpl) GenerateAndSendOTP(
	ctx context.Context,
	phoneNumber string,
	flavour feedlib.Flavour,
) (string, error) {
	phone, err := converterandformatter.NormalizeMSISDN(phoneNumber)
	if err != nil {
		return "", exceptions.NormalizeMSISDNError(err)
	}

	if !flavour.IsValid() {
		return "", exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}

	userProfile, err := o.Query.GetUserProfileByPhoneNumber(ctx, *phone)
	if err != nil {
		return "", exceptions.UserNotFoundError(err)
	}

	otp, err := o.GenerateOTP(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate an OTP")
	}

	message := fmt.Sprintf(otpMessage, otp)
	otp, err = o.SendOTP(ctx, *phone, otp, message)
	if err != nil {
		return "", err
	}

	otpDataPayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Minute * 10),
		Channel:     "SMS",
		Flavour:     flavour,
		PhoneNumber: phoneNumber,
		OTP:         otp,
	}

	err = o.Create.SaveOTP(ctx, otpDataPayload)
	if err != nil {
		return "", fmt.Errorf("failed to save otp")
	}

	return otp, nil
}

// VerifyOTP verifies whether the supplied OTP is valid
func (o *UseCaseOTPImpl) VerifyOTP(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
	return o.Query.VerifyOTP(ctx, payload)
}

// GenerateOTP calls the engagement library to generate a random OTP
func (o *UseCaseOTPImpl) GenerateOTP(ctx context.Context) (string, error) {
	return o.ExternalExt.GenerateOTP(ctx)
}

// VerifyPhoneNumber checks validity of a phone number by sending an OTP to it
func (o *UseCaseOTPImpl) VerifyPhoneNumber(ctx context.Context, phone string, flavour feedlib.Flavour) (*profileutils.OtpResponse, error) {
	phoneNumber, err := converterandformatter.NormalizeMSISDN(phone)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	exists, err := o.Query.CheckIfPhoneNumberExists(ctx, *phoneNumber, true, flavour)
	if err != nil {
		return nil, fmt.Errorf("failed to check if phone exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("the provided phone number does not exist")
	}

	userProfile, err := o.Query.GetUserProfileByPhoneNumber(ctx, *phoneNumber)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	otp, err := o.GenerateOTP(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate an OTP")
	}

	message := fmt.Sprintf(otpMessage, otp)
	otp, err = o.SendOTP(ctx, phone, otp, message)
	if err != nil {
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
		return nil, fmt.Errorf("failed to save otp: %v", err)
	}

	return &profileutils.OtpResponse{
		OTP: otp,
	}, nil
}

// GenerateRetryOTP generates fallback OTPs when Africa is talking sms fails
func (o *UseCaseOTPImpl) GenerateRetryOTP(ctx context.Context, payload *dto.SendRetryOTPPayload) (string, error) {
	phoneNumber, err := converterandformatter.NormalizeMSISDN(payload.Phone)
	if err != nil {
		return "", err
	}

	validPayload := &dto.SendRetryOTPPayload{
		Phone:   *phoneNumber,
		Flavour: payload.Flavour,
	}

	retryResponseOTP, err := o.ExternalExt.GenerateRetryOTP(ctx, validPayload)
	if err != nil {
		return "", err
	}

	userProfile, err := o.Query.GetUserProfileByPhoneNumber(ctx, *phoneNumber)
	if err != nil {
		return "", exceptions.UserNotFoundError(err)
	}

	otpResponsePayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Hour * 1),
		Channel:     "SMS",
		Flavour:     payload.Flavour,
		PhoneNumber: *phoneNumber,
		OTP:         retryResponseOTP,
	}

	err = o.Create.SaveOTP(ctx, otpResponsePayload)
	if err != nil {
		return "", fmt.Errorf("failed to save otp: %v", err)
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
		_, err := o.ExternalExt.SendSMS(ctx, phoneNumber, message, enumutils.SenderIDBewell)
		if err != nil {
			return "", fmt.Errorf("failed to send OTP verification code to recipient")
		}
	} else {
		// Make the request to twilio
		err := o.ExternalExt.SendSMSViaTwilio(ctx, phoneNumber, message)
		if err != nil {
			return "", fmt.Errorf("sms not sent via twilio: %v", err)
		}
	}
	return code, nil
}
