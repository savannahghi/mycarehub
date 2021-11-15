package exceptions

import (
	"fmt"

	"github.com/savannahghi/errorcodeutil"
)

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NormalizeMSISDNErrMsg,
		Code:    int(Internal),
	}
}

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: UserNotFoundErrMsg,
		Code:    int(UserNotFound),
	}
}

// PinNotFoundError displays error message when a pin is not found
func PinNotFoundError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINNotFoundErrMsg,
		Code:    int(PINNotFound),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINMismatchErrMsg,
		Code:    int(PINMismatch),
	}
}

// InvalidFlavourDefinedError is the error message displayed when
// an invalid flavour is provided as input.
func InvalidFlavourDefinedError() error {
	return &errorcodeutil.CustomError{
		Err:     fmt.Errorf("invalid flavour defined"),
		Message: InvalidFlavourDefinedErrMsg,
		Code:    int(InvalidFlavour),
	}
}

// SaveUserPinError returns an error message when we are unable to save a user pin
func SaveUserPinError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: SaveUserPinErrMsg,
		Code:    int(PINError),
	}
}

// InvalidResetPinPayloadErr returns an error message when the provided reset pin payload is invalid
func InvalidResetPinPayloadErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: InvalidResetPinPayloadErrorMsg,
		Code:    int(InvalidResetPinPayloadError),
	}
}

// EmptyUserIDErr returns an error message when the user id is empty
func EmptyUserIDErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: EmptyUserIDInputErrorMsg,
		Code:    int(EmptyUserIDInputError),
	}
}

// ProfileNotFoundErr returns an error message when the profile is not found
func ProfileNotFoundErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ProfileNotFoundErrorMsg,
		Code:    int(ProfileNotFound),
	}
}

// InvalidatePinErr returns an error message when the reset pin is invalid
func InvalidatePinErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: InvalidatePinErrMsg,
		Code:    int(InvalidatePinError),
	}
}

// ResetPinErr returns an error message when the reset pin is invalid
func ResetPinErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ResetPinErrorMsg,
		Code:    int(ResetPinError),
	}
}

// PINExpiredErr returns an error message when the reset pin is invalid
func PINExpiredErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINExpiredErrorMsg,
		Code:    int(PINExpiredError),
	}
}

// EmptyInputErr returns an error message when an input is empty
func EmptyInputErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: EmptyInputErrorMsg,
		Code:    int(EmptyInputError),
	}
}

// PINErr returns an error message when the PIN is invalid
func PINErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINErrorMsg,
		Code:    int(PINError),
	}
}

// NotOptedInErr returns an error message when the user is not opted in
func NotOptedInErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NotOptedInErrorMsg,
		Code:    int(NotOptedInError),
	}
}

// NotActiveErr returns an error message when a field is not active
func NotActiveErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NotActiveErrorMsg,
		Code:    int(NotActiveError),
	}
}

// InvalidContactTypeErr returns an error message when the contact type is invalid
func InvalidContactTypeErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: InvalidContactTypeErrorMsg,
		Code:    int(InvalidContactTypeError),
	}
}

// NoContactsErr returns an error message when there are no contacts
func NoContactsErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NoContactsErrorMsg,
		Code:    int(NoContactsError),
	}
}

// ContactNotFoundErr returns an error message when the contact is not found
func ContactNotFoundErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ContactNotFoundErrorMsg,
		Code:    int(ContactNotFoundError),
	}
}

// GenerateTempPINErr returns an error message when the temp pin generation fails
func GenerateTempPINErr(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: GenerateTempPINErrMsg,
		Code:    int(GenerateTempPINError),
	}
}
