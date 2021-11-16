package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/profileutils"
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

	GenerateAndSendOTP(
		ctx context.Context,
		userID string,
		phoneNumber string,
		flavour feedlib.Flavour,
	) (string, error)
}

// UsecaseOTP defines otp service usecases interface
type UsecaseOTP interface {
	IGenerateOTP
	ISendOTP
	IverifyPhone
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
	userID string,
	phoneNumber string,
	flavour feedlib.Flavour,
) (string, error) {
	if !flavour.IsValid() {
		return "", exceptions.InvalidFlavourDefinedError()
	}

	otp, err := o.ExternalExt.GenerateAndSendOTP(ctx, phoneNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate and send OTP")
	}

	otpDataPayload := &domain.OTP{
		UserID:      userID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Minute * 10),
		Channel:     "SMS",
		Flavour:     flavour,
		PhoneNumber: phoneNumber,
	}

	err = o.Create.SaveOTP(ctx, otpDataPayload)
	if err != nil {
		return "", fmt.Errorf("failed to save otp")
	}

	return otp, nil
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

	otp, err := o.ExternalExt.GenerateAndSendOTP(ctx, *phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to generate and send OTP")
	}

	otpDataPayload := &domain.OTP{
		UserID:      *userProfile.ID,
		Valid:       true,
		GeneratedAt: time.Now(),
		ValidUntil:  time.Now().Add(time.Minute * 10),
		Channel:     "SMS",
		Flavour:     userProfile.Flavour,
		PhoneNumber: phone,
		OTP:         otp,
	}

	err = o.Create.SaveOTP(ctx, otpDataPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to save otp: %v", err)
	}

	//Return the OTP response
	return &profileutils.OtpResponse{
		OTP: otp,
	}, nil
}
