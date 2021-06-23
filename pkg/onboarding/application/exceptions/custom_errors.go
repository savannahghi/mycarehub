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

// ProfileSuspendFoundError is returned is the user profile has been suspended.
func ProfileSuspendFoundError() error {
	return &base.CustomError{
		Message: ProfileSuspenedFoundErrMsg,
		Code:    int(base.ProfileSuspended),
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
func CheckPhoneNumberExistError() error {
	return &base.CustomError{
		Message: PhoneNumberInUseErrMsg,
		Code:    int(base.PhoneNumberInUse),
	}
}

// CheckEmailExistError returned when the provided email already exists.
func CheckEmailExistError() error {
	return &base.CustomError{
		Message: EmailInUseErrMsg,
		Code:    int(base.EmailAddressInUse),
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
		Code:    int(base.PINError),
	}
}

// ValidatePINDigitsError  is the error message displayed when
// invalid  pin digits are given
func ValidatePINDigitsError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: ValidatePINDigitsErrMsg,
		Code:    int(base.PINError),
	}

}

// ValidatePINLengthError  is the error message displayed when
// an invalid Pin length is given
func ValidatePINLengthError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: ValidatePINLengthErrMsg,
		Code:    int(base.PINError),
	}

}

// InValidPushTokenLengthError  is the error message displayed when
// an invalid push token is given
func InValidPushTokenLengthError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid push token length"),
		Message: ValidatePushTokenLengthErrMsg,
		Code:    int(base.InvalidPushTokenLength),
	}
}

// WrongEnumTypeError  is the error message displayed when
// an invalid enum is given
func WrongEnumTypeError(value string) error {
	return &base.CustomError{
		Err:     fmt.Errorf("%v", WrongEnumErrMsg),
		Message: fmt.Sprintf(WrongEnumErrMsg, value),
		Code:    int(base.InvalidEnum),
	}

}

// VerifyOTPError returns an error when OTP verification fails
func VerifyOTPError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: OTPVerificationErrMsg,
		Code:    int(base.OTPVerificationFailed),
	}
}

// MissingInputError returns an error when OTP verification fails
func MissingInputError(value string) error {
	return &base.CustomError{
		Err:     nil,
		Message: "expected `%s` to be defined",
		Code:    int(base.OTPVerificationFailed),
	}
}

// InvalidFlavourDefinedError is the error message displayed when
// an invalid flavour is provided as input.
func InvalidFlavourDefinedError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid flavour defined"),
		Message: InvalidFlavourDefinedErrMsg,
		Code:    int(base.InvalidFlavour),
	}
}

// AddPartnerTypeError is an error message displayed when there is a
// failure to create a partner type
func AddPartnerTypeError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: AddPartnerTypeErrMsg,
		Code:    int(base.InvalidEnum),
	}

}

// InvalidPartnerTypeError is an error message displayed when an
// invalid partner type is provided
func InvalidPartnerTypeError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid `partnerType` provided"),
		Message: InvalidPartnerTypeErrMsg,
		Code:    int(base.InvalidEnum),
	}
}

// FetchDefaultCurrencyError is an error message displayed when
// the default currency is not found
func FetchDefaultCurrencyError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: FetchDefaultCurrencyErrMsg,
		Code:    int(base.RecordNotFound),
	}
}

// SupplierNotFoundError returns an error message when a supplier is not found
func SupplierNotFoundError() error {
	return &base.CustomError{
		Message: SupplierNotFoundErrMsg,
		Code:    int(base.ProfileNotFound),
	}

}

// CustomerNotFoundError returns an error message when a customer is not found
func CustomerNotFoundError() error {
	return &base.CustomError{
		Message: CustomerNotFoundErrMsg,
		Code:    int(base.ProfileNotFound),
	}
}

// SupplierKYCAlreadySubmittedNotFoundError is returned when the user trys to
// submit another KCY when then is one already submitted
func SupplierKYCAlreadySubmittedNotFoundError() error {
	return &base.CustomError{
		Message: SupplierKYCAlreadySubmittedErrMsg,
		Code:    int(base.KYCAlreadySubmitted),
	}
}

// FindProviderError returns an error message when a provider is not found
func FindProviderError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: FindProviderErrMsg,
		Code:    int(base.UnableToFindProvider),
	}
}

// PublishKYCNudgeError returns an error message when there's a failure in
// creating a KYC nudge
func PublishKYCNudgeError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: PublishKYCNudgeErrMsg,
		Code:    int(base.PublishNudgeFailure),
	}
}

// InvalidCredentialsError returns an error message when wrong credentials are provided
func InvalidCredentialsError() error {
	return &base.CustomError{
		Err:     fmt.Errorf("invalid credentials, expected a username AND password"),
		Message: InvalidCredentialsErrMsg,
		Code:    int(base.InvalidCredentials),
	}
}

// SaveUserPinError returns an error message when we are unable to save a user pin
func SaveUserPinError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: SaveUserPinErrMsg,
		Code:    int(base.PINError),
	}
}

// GeneratePinError returns an error message when we are unable to generate a temporary PIN
func GeneratePinError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: GeneratePinErrMsg,
		Code:    int(base.PINError),
	}
}

// CompleteSignUpError returns an error message when we are unable
// to CompleteSignup
func CompleteSignUpError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: BioDataErrMsg,
		Code:    int(base.AddNewRecordError),
	}
}

// UsernameInUseError is the error message displayed when the provided username
// is associated with another profile
func UsernameInUseError() error {
	return &base.CustomError{
		Message: UsernameInUseErrMsg,
		Code:    int(base.UsernameInUse),
	}
}

// SecondaryResourceHardResetError this error is returned when there argument to reset a resource has a length of 0
// resource here means secondary phone numbers and secondary emails
func SecondaryResourceHardResetError() error {
	return &base.CustomError{
		Message: ResourceUpdateErrMsg,
		Code:    int(base.UndefinedArguments),
	}
}

// InvalidSladeCodeError when the slade code the edi user profile doesn't match with selected provider
func InvalidSladeCodeError() error {
	return &base.CustomError{
		Message: InvalidSladeCodeErrMsg,
		Code:    int(base.InvalidSladeCode),
	}
}

// ResolveNudgeErr is the error that represents the failure of not
// being able to resolve a given nudge
func ResolveNudgeErr(
	err error,
	flavour base.Flavour,
	name string,
	statusCode *int,
) error {
	if statusCode != nil {
		return &base.CustomError{
			Err: err,
			Message: fmt.Sprintf(
				ResolveNudgeBadStatusErrMsg,
				flavour,
				name,
				statusCode,
			),
			Code: int(base.Internal),
		}
	}

	return &base.CustomError{
		Err: err,
		Message: fmt.Sprintf(
			ResolveNudgeErrMsg,
			flavour,
			name,
		),
		Code: int(base.Internal),
	}
}

// RecordExistsError is the error message displayed when a
// similar record is found in the DB
func RecordExistsError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: RecordExistsErrMsg,
		Code:    int(base.Internal),
	}
}

// RecordDoesNotExistError is the error message displayed when a
// record is not found in the DB
func RecordDoesNotExistError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: RecordDoesNotExistErrMsg,
		Code:    int(base.Internal),
	}
}

// SessionIDError return an error when a ussd sessionId is not provided
func SessionIDError(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: SessionIDErrMsg,
		Code:    int(base.Internal),
	}
}

// RoleNotValid return an error when a user does not have the required role
func RoleNotValid(err error) error {
	return &base.CustomError{
		Err:     err,
		Message: RoleNotValidMsg,
		Code:    int(base.RoleNotValid),
	}
}
