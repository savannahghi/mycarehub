package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// UserPINUseCases represents all the business logic that touch on user PIN Management
type UserPINUseCases interface {
	SetUserPIN(ctx context.Context, pin string, phone string) (bool, error)
	ChangeUserPIN(ctx context.Context, phone string, pin string) (*domain.PIN, error)
	RequestPINReset(ctx context.Context, phone string) (string, error)
}

// UserPinUseCaseImpl represents usecase implementation object
type UserPinUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
	otpUseCases          OTPUseCases
	profileUseCases      ProfileUseCase
}

// NewUserPinUseCase returns a new UserPin usecase
func NewUserPinUseCase(r repository.OnboardingRepository, otp OTPUseCases, p ProfileUseCase) UserPINUseCases {
	return &UserPinUseCaseImpl{
		onboardingRepository: r,
		otpUseCases:          otp,
		profileUseCases:      p,
	}
}

// SetUserPIN receives phone number and pin from phonenumber sign up
func (u *UserPinUseCaseImpl) SetUserPIN(ctx context.Context, pin string, phone string) (bool, error) {

	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return false, &domain.CustomError{
			Err:     err,
			Message: exceptions.NormalizeMSISDNErrMsg,
			Code:    int(base.Internal),
		}
	}

	pr, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return false, &domain.CustomError{
			Err:     err,
			Message: exceptions.ProfileNotFoundErrMsg,
			Code:    int(base.ProfileNotFound),
		}
	}

	if err := utils.ValidatePINLength(pin); err != nil {
		return false, err
	}

	if err = utils.ValidatePINDigits(pin); err != nil {
		return false, err
	}

	// EncryptPIN the PIN
	salt, encryptedPin := utils.EncryptPIN(pin, nil)

	pinPayload := &domain.PIN{
		ProfileID: pr.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	if _, err := u.onboardingRepository.SavePIN(ctx, pinPayload); err != nil {
		return false, fmt.Errorf("unable to save user PIN: %v", err)
	}

	return true, nil
}

// RequestPINReset sends a request given an existing user's phone number,
// sends an otp to the phone number that is then used in the process of
// updating their old PIN to a new one
func (u *UserPinUseCaseImpl) RequestPINReset(ctx context.Context, phone string) (string, error) {
	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return "", &domain.CustomError{
			Err:     err,
			Message: exceptions.NormalizeMSISDNErrMsg,
			Code:    int(base.Internal),
		}
	}

	pr, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return "", &domain.CustomError{
			Err:     err,
			Message: exceptions.ProfileNotFoundErrMsg,
			Code:    int(base.ProfileNotFound),
		}
	}

	exists, err := u.CheckHasPIN(ctx, pr.ID)
	if err != nil {
		return "", &domain.CustomError{
			Err:     err,
			Message: exceptions.CheckUserPINErrMsg,
			Code:    int(base.Internal),
		}
	}
	if !exists {
		return "", &domain.CustomError{
			Err:     err,
			Message: exceptions.ExistingPINErrMsg,
			Code:    int(base.PINNotFound),
		}
	}

	// generate and send otp to the phone number
	code, err := u.otpUseCases.GenerateAndSendOTP(ctx, phone)
	if err != nil {
		return "", &domain.CustomError{
			Err:     err,
			Message: exceptions.GenerateAndSendOTPErrMsg,
			Code:    int(base.Internal),
		}
	}

	return code, nil
}

// ChangeUserPIN resets a user's pin with the newly supplied pin
func (u *UserPinUseCaseImpl) ChangeUserPIN(ctx context.Context, phone string, pin string) (*domain.PIN, error) {
	phoneNumber, err := base.NormalizeMSISDN(phone)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.NormalizeMSISDNErrMsg,
			Code:    int(base.Internal),
		}
	}

	profile, err := u.onboardingRepository.GetUserProfileByPrimaryPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.ProfileNotFoundErrMsg,
			Code:    int(base.ProfileNotFound),
		}
	}

	exists, err := u.CheckHasPIN(ctx, profile.ID)
	if !exists {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.ExistingPINErrMsg,
			Code:    int(base.PINNotFound),
		}
	}

	_, err = u.onboardingRepository.
		GetPINByProfileID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to read PIN: %w", err)
	}
	// EncryptPIN the PIN
	salt, encryptedPin := utils.EncryptPIN(pin, nil)
	if err != nil {
		return nil, &domain.CustomError{
			Err:     err,
			Message: exceptions.EncryptPINErrMsg,
			// TODO: correct error code
			Code: int(base.UserNotFound),
		}
	}

	pinPayload := &domain.PIN{
		ProfileID: profile.ID,
		PINNumber: encryptedPin,
		Salt:      salt,
	}
	return u.onboardingRepository.UpdatePIN(ctx, pinPayload)
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
