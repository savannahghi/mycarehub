package customerrors

import (
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
)

// SaveUserPinError returns an error message when we are unable to save a user pin
func SaveUserPinError(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to save user PIN",
		Code:    int(exceptions.PINError),
	}
}

// InvalidResetPinPayloadErr returns an error message when the provided reset pin payload is invalid
func InvalidResetPinPayloadErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "failed to validate reset pin payload",
		Code:    int(exceptions.InvalidResetPinPayloadError),
	}
}

// PinNotFoundError displays error message when a pin is not found
func PinNotFoundError(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "failed to get a user pin",
		Code:    int(exceptions.PINNotFound),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError() error {
	return &exceptions.CustomError{
		Err:     nil,
		Message: "wrong PIN credentials supplied",
		Code:    int(exceptions.PINMismatch),
	}
}

// InvalidatePinErr returns an error message when the reset pin is invalid
func InvalidatePinErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to invalidate reset pin",
		Code:    int(exceptions.InvalidatePinError),
	}
}

// ResetPinErr returns an error message when the reset pin is invalid
func ResetPinErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to reset pin",
		Code:    int(exceptions.ResetPinError),
	}
}

// PINExpiredErr returns an error message when the reset pin is invalid
func PINExpiredErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "pin expired",
		Code:    int(exceptions.PINExpiredError),
	}
}

// PINErr returns an error message when the PIN is invalid
func PINErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "invalid pin",
		Code:    int(exceptions.PINError),
	}
}

// GenerateTempPINErr returns an error message when the temp pin generation fails
func GenerateTempPINErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to generate temporary pin",
		Code:    int(exceptions.GenerateTempPINError),
	}
}

// ExpiredPinErr returns an error message when the pin is expired
func ExpiredPinErr() error {
	return &exceptions.CustomError{
		Err:     nil,
		Message: "pin expired",
		Code:    int(exceptions.ExpiredPinError),
	}
}

// GeneratePinErr returns an error message when the pin generation fails
func GeneratePinErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "unable to generate pin",
		Code:    int(exceptions.GeneratePinError),
	}
}

// ValidatePINDigitsErr returns an error message when the pin digits are invalid
func ValidatePINDigitsErr(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "invalid pin digits",
		Code:    int(exceptions.ValidatePINDigitsError),
	}
}

// ExistingPINError is the error message displayed when a
// pin record fails to be retrieved from dataerrorcodeutil
func ExistingPINError(err error) error {
	return &exceptions.CustomError{
		Err:     err,
		Message: "user does not have an existing PIN",
		Code:    int(errorcodeutil.PINNotFound),
	}

}
