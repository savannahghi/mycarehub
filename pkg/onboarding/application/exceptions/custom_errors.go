package exceptions

import (
	"fmt"

	"gitlab.slade360emr.com/go/base"
)

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: UserNotFoundErrMsg,
		Code:    int(base.UserNotFound),
	}
}

// ProfileNotFoundError returns an error message when a profile is not found
func ProfileNotFoundError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: ProfileNotFoundErrMsg,
		Code:    int(base.ProfileNotFound),
	}
}

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: NormalizeMSISDNErrMsg,
		Code:    int(base.Internal),
	}
}

// CheckPhoneNumberExistError check if phone number is registered to another user
func CheckPhoneNumberExistError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: PhoneNUmberInUseErrMsg,
		Code:    int(base.PhoneNumberInUse),
	}
}

// InternalServerError returns an error if something wrong happened in performing the operration
func InternalServerError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: InternalServerErrorMsg,
		Code:    int(base.Internal),
	}
}

// PinNotFoundError displays error message when a pin is not found
func PinNotFoundError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: PINNotFoundErrMsg,
		Code:    int(base.PINNotFound),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: PINMismatchErrMsg,
		Code:    int(base.PINMismatch),
	}
}

// CustomTokenError is the error message displayed when a
// custom token is not created
func CustomTokenError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: CustomTokenErrMsg,
		Code:    int(base.Internal),
	}
}

// AuthenticateTokenError is the error message displayed when a
// custom token is not authenticated
func AuthenticateTokenError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: AuthenticateTokenErrMsg,
		Code:    int(base.Internal),
	}
}

// UpdateProfileError is the error message displayed when a
// user profile cannot be updated
func UpdateProfileError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: UpdateProfileErrMsg,
		Code:    int(base.Internal),
	}
}

// AddRecordError is the error message displayed when a
// record fails to be added to the database
func AddRecordError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: AddRecordErrMsg,
		Code:    int(base.Internal),
	}
}

// RetrieveRecordError is the error message displayed when a
// failure occurs while retrieving records from the database
func RetrieveRecordError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: RetrieveRecordErrMsg,
		Code:    int(base.Internal),
	}
}

// LikelyToRecommendError is the error message displayed that
// occurs when the recommendation threshold is crossed
func LikelyToRecommendError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: LikelyToRecommendErrMsg,
		Code:    int(base.UndefinedArguments),
	}
}

// GenerateAndSendOTPError is the error message displayed when a
// generate and send otp fails
func GenerateAndSendOTPError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: GenerateAndSendOTPErrMsg,
		Code:    int(base.Internal),
	}
}

// CheckUserPINError is the error message displayed when
// a server is unable to check if the user has a PIN
func CheckUserPINError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: CheckUserPINErrMsg,
		Code:    int(base.Internal),
	}
}

// ExistingPINError is the error message displayed when a
// pin record fails to be retrieved from database
func ExistingPINError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: ExistingPINErrMsg,
		Code:    int(base.PINNotFound),
	}
}

// EncryptPINError  is the error message displayed when
// pin encryption fails
func EncryptPINError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: EncryptPINErrMsg,
		// TODO: add correct error code
		Code: int(base.Internal),
	}
}

// ValidatePINDigitsError  is the error message displayed when
// invalid  pin digits are given
func ValidatePINDigitsError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: ValidatePINDigitsErrMsg,
		// TODO: a give a correct code
		Code: int(base.UserNotFound),
	}
}

// ValidatePINLengthError  is the error message displayed when
// an invalid Pin length is given
func ValidatePINLengthError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: ValidatePINLengthErrMsg,
		// TODO: a give a correct code
		Code: int(base.UserNotFound),
	}
}

// InValidPushTokenLengthError  is the error message displayed when
// an invalid push token is given
func InValidPushTokenLengthError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid push token length"),
		Message: ValidatePuskTokenLengthErrMsg,
		// TODO: a give a correct code
		Code: int(base.UndefinedArguments),
	}
}

// WrongEnumTypeError  is the error message displayed when
// an invalid enum is given
func WrongEnumTypeError(value string, err error) error {
	return &base.CustomError{
		Err:     err,
		Message: fmt.Sprintf(WrongEnumErrMsg, value),
		// TODO: a give a correct code
		Code: int(base.Internal),
	}
}
