package exceptions

import (
	"fmt"

	"gitlab.slade360emr.com/go/base"
)

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: UserNotFoundErrMsg,
			Code:    int(base.UserNotFound),
		}
	}
	return nil
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
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: NormalizeMSISDNErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
}

// CheckPhoneNumberExistError check if phone number is registered to another user
func CheckPhoneNumberExistError() error {
	return &base.CustomError{
		Message: PhoneNumberInUseErrMsg,
		Code:    int(base.PhoneNumberInUse),
	}
}

// InternalServerError returns an error if something wrong happened in performing the operation
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
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: CustomTokenErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
}

// AuthenticateTokenError is the error message displayed when a
// custom token is not authenticated
func AuthenticateTokenError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: AuthenticateTokenErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
}

// UpdateProfileError is the error message displayed when a
// user profile cannot be updated
func UpdateProfileError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: UpdateProfileErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
}

// AddRecordError is the error message displayed when a
// record fails to be added to the database
func AddRecordError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: AddRecordErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
}

// RetrieveRecordError is the error message displayed when a
// failure occurs while retrieving records from the database
func RetrieveRecordError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: RetrieveRecordErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
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
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: GenerateAndSendOTPErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
}

// CheckUserPINError is the error message displayed when
// a server is unable to check if the user has a PIN
func CheckUserPINError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: CheckUserPINErrMsg,
			Code:    int(base.Internal),
		}
	}
	return nil
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
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: EncryptPINErrMsg,
			// TODO: add correct error code
			Code: int(base.Internal),
		}
	}
	return nil
}

// ValidatePINDigitsError  is the error message displayed when
// invalid  pin digits are given
func ValidatePINDigitsError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: ValidatePINDigitsErrMsg,
			// TODO: a give a correct code
			Code: int(base.UserNotFound),
		}
	}
	return nil
}

// ValidatePINLengthError  is the error message displayed when
// an invalid Pin length is given
func ValidatePINLengthError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: ValidatePINLengthErrMsg,
			// TODO: a give a correct code
			Code: int(base.UserNotFound),
		}
	}
	return nil
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
func WrongEnumTypeError(value string) error {
	return &base.CustomError{
		Err:     fmt.Errorf("%v", WrongEnumErrMsg),
		Message: fmt.Sprintf(WrongEnumErrMsg, value),
		// TODO: a give a correct code
		Code: int(base.Internal),
	}

}

// VerifyOTPError returns an error when OTP verification fails
func VerifyOTPError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: OTPVerificationErrMsg,
		// TODO: @salaton OTP verification error code
		Code: int(base.Internal),
	}
}

// MissingInputError returns an error when OTP verification fails
func MissingInputError(value string) error {
	return &base.CustomError{
		Err:     nil,
		Message: "expected `%s` to be defined",
		// TODO: @salaton error code
		Code: int(base.Internal),
	}
}

// InvalidFlavourDefinedError is the error message displayed when
// an invalid flavour is provided as input.
func InvalidFlavourDefinedError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid flavour defined"),
		Message: InvalidFlavourDefinedErrMsg,
		// TODO: a give a correct code
		Code: int(base.UndefinedArguments),
	}
}

// AddPartnerTypeError is an error message displayed when there is a
// failure to create a partner type
func AddPartnerTypeError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: AddPartnerTypeErrMsg,
			// TODO: provide a correct code
			Code: int(base.Internal),
		}
	}
	return nil
}

// InvalidPartnerTypeError is an error message displayed when an
// invalid partner type is provided
func InvalidPartnerTypeError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid `partnerType` provided"),
		Message: InvalidPartnerTypeErrMsg,
		// TODO: provide a correct code
		Code: int(base.Internal),
	}
}

// FetchDefaultCurrencyError is an error message displayed when
// the default currency is not found
func FetchDefaultCurrencyError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: FetchDefaultCurrencyErrMsg,
			// TODO: provide a correct code
			Code: int(base.Internal),
		}
	}
	return nil
}

// SupplierNotFoundError returns an error message when a supplier is not found
func SupplierNotFoundError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: SupplierNotFoundErrMsg,
			// TODO: provide a correct code
			Code: int(base.UserNotFound),
		}
	}
	return nil
}

// FindProviderError returns an error message when a provider is not found
func FindProviderError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: FindProviderErrMsg,
			// TODO: provide a correct code
			Code: int(base.Internal),
		}
	}
	return nil
}

// PublishKYCNudgeError returns an error message when there's a failure in
// creating a KYC nudge
func PublishKYCNudgeError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: PublishKYCNudgeErrMsg,
			// TODO: provide a correct code
			Code: int(base.Internal),
		}
	}
	return nil
}

// InvalidCredentialsError returns an error message when wrong credentials are provided
func InvalidCredentialsError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid credentials, expected a username AND password"),
		Message: InvalidCredentialsErrMsg,
		// TODO: provide a correct code
		Code: int(base.Internal),
	}
}

// SaveUserPinError returns an error message when we are unable to save a user pin
func SaveUserPinError(err error) error {
	if err != nil {
		return &base.CustomError{
			Err:     err,
			Message: SaveUserPinErrMsg,
			// TODO: provide a correct code
			Code: int(base.Internal),
		}
	}
	return nil
}

// CompleteSignUpError returns an error message when we are unable
// to CompleteSignup
func CompleteSignUpError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: BioDataErrMsg,
		// TODO: provide a correct code
		Code: int(base.Internal),
	}
}

// UsernameInUseError is the error message displayed when the provided username
// is associated with another profile
func UsernameInUseError() error {
	return &base.CustomError{
		Message: UsernameInUseErrMsg,
		// TODO: provide a correct code
		Code: int(base.UserNotFound),
	}
}
