package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// UserPINUseCases represents all the business logic that touch on user PIN Management
type UserPINUseCases interface {
	SetUserPIN(ctx context.Context, pin string, phone string) (bool, error)
	ResetUserPIN(
		ctx context.Context,
		phone string,
		PIN string,
		OTP string,
	) (bool, error)
	ChangeUserPIN(ctx context.Context, phone string, pin string) (bool, error)
	RequestPINReset(ctx context.Context, phone string) (*base.OtpResponse, error)
	CheckHasPIN(ctx context.Context, profileID string) (bool, error)
}

// UserPinUseCaseImpl represents usecase implementation object
type UserPinUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
	otpUseCases          otp.ServiceOTP
	profileUseCases      ProfileUseCase
	baseExt              extension.BaseExtension
	pinExt               extension.PINExtension
}

// NewUserPinUseCase returns a new UserPin usecase
func NewUserPinUseCase(
	r repository.OnboardingRepository,
	otp otp.ServiceOTP, p ProfileUseCase,
	ext extension.BaseExtension, pin extension.PINExtension) UserPINUseCases {
	return &UserPinUseCaseImpl{
		onboardingRepository: r,
		otpUseCases:          otp,
		profileUseCases:      p,
		baseExt:              ext,
		pinExt:               pin,
	}
}

// SetUserPIN receives phone number and pin from phonenumber sign up
func (u *UserPinUseCaseImpl) SetUserPIN(
	ctx context.Context,
	pin string,
	phone string,
) (bool, error) {
	phoneNumber, err := u.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}

	if err := extension.ValidatePINLength(pin); err != nil {
		return false, exceptions.ValidatePINLengthError(err)
	}

	if err = extension.ValidatePINDigits(pin); err != nil {
		return false, exceptions.ValidatePINDigitsError(err)
	}

	pr, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, *phoneNumber, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return false, err
	}
	// EncryptPIN the PIN
	salt, encryptedPin := u.pinExt.EncryptPIN(pin, nil)

	pinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: pr.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	if _, err := u.onboardingRepository.SavePIN(ctx, pinPayload); err != nil {
		return false, exceptions.SaveUserPinError(err)
	}

	return true, nil
}

// RequestPINReset sends a request given an existing user's phone number,
// sends an otp to the phone number that is then used in the process of
// updating their old PIN to a new one
func (u *UserPinUseCaseImpl) RequestPINReset(
	ctx context.Context,
	phone string,
) (*base.OtpResponse, error) {
	phoneNumber, err := u.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	pr, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, *phoneNumber, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	exists, err := u.CheckHasPIN(ctx, pr.ID)
	if !exists {
		return nil, exceptions.ExistingPINError(err)
	}
	// generate and send otp to the phone number
	otpResp, err := u.otpUseCases.GenerateAndSendOTP(ctx, phone)
	if err != nil {
		return nil, exceptions.GenerateAndSendOTPError(err)
	}

	return otpResp, nil
}

// ResetUserPIN resets a user's PIN with the newly supplied PIN
func (u *UserPinUseCaseImpl) ResetUserPIN(
	ctx context.Context,
	phone string,
	PIN string,
	OTP string,
) (bool, error) {
	phoneNumber, err := u.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}

	verified, err := u.otpUseCases.VerifyOTP(ctx, phone, OTP)
	if err != nil {
		return false, exceptions.VerifyOTPError(err)
	}

	if !verified {
		return false, exceptions.VerifyOTPError(nil)
	}

	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return false, err
	}

	_, err = u.CheckHasPIN(ctx, profile.ID)
	if err != nil {
		return false, exceptions.EncryptPINError(err)
	}

	// EncryptPIN the PIN
	salt, encryptedPin := u.pinExt.EncryptPIN(PIN, nil)

	pinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: profile.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	_, err = u.onboardingRepository.UpdatePIN(ctx, profile.ID, pinPayload)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}

// ChangeUserPIN updates authenticated user's pin with the newly supplied pin
func (u *UserPinUseCaseImpl) ChangeUserPIN(
	ctx context.Context,
	phone string,
	pin string,
) (bool, error) {
	phoneNumber, err := u.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}

	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(
		ctx,
		*phoneNumber,
		false,
	)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return false, err
	}

	_, err = u.CheckHasPIN(ctx, profile.ID)
	if err != nil {
		return false, exceptions.EncryptPINError(err)
	}

	salt, encryptedPin := u.pinExt.EncryptPIN(pin, nil)

	pinPayload := &domain.PIN{
		ID:        uuid.New().String(),
		ProfileID: profile.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	_, err = u.onboardingRepository.UpdatePIN(ctx, profile.ID, pinPayload)
	if err != nil {
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}

// CheckHasPIN given a phone number checks if the phonenumber is present in our collections
// which essentially means that the number has an already existing PIN
func (u *UserPinUseCaseImpl) CheckHasPIN(ctx context.Context, profileID string) (bool, error) {

	pinData, err := u.onboardingRepository.GetPINByProfileID(ctx, profileID)
	if err != nil {
		return false, err
	}

	if pinData == nil {
		return false, fmt.Errorf("%v", base.PINNotFound)
	}

	return true, nil
}
