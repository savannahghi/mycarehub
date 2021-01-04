package exceptions

import (
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	return &domain.CustomError{
		Err:     err,
		Message: UserNotFoundErrMsg,
		Code:    int(base.UserNotFound),
	}
}

// ProfileNotFoundError returns an error message when a profile is not found
func ProfileNotFoundError(err error) error {
	return &domain.CustomError{
		Err:     err,
		Message: ProfileNotFoundErrMsg,
		Code:    int(base.ProfileNotFound),
	}
}

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &domain.CustomError{
		Err:     err,
		Message: NormalizeMSISDNErrMsg,
		Code:    int(base.Internal),
	}
}

// CheckPhoneNumberExistError check if phone number is registered to another user
func CheckPhoneNumberExistError(err error) error {
	return &domain.CustomError{
		Err:     err,
		Message: PhoneNUmberInUseErrMsg,
		Code:    int(base.PhoneNumberInUse),
	}
}

// PinNotFoundError displays error message when a pin is not found
func PinNotFoundError(err error) error {
	return &domain.CustomError{
		Err:     err,
		Message: PINNotFoundErrMsg,
		Code:    int(base.PINNotFound),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError(err error) error {
	return &domain.CustomError{
		Err:     err,
		Message: PINMismatchErrMsg,
		Code:    int(base.PINMismatch),
	}
}
