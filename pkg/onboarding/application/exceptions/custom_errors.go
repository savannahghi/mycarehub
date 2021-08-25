package exceptions

import (
	"fmt"

	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
)

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: UserNotFoundErrMsg,
		Code:    int(errorcodeutil.UserNotFound),
	}
}

// ProfileSuspendFoundError is returned is the user profile has been suspended.
func ProfileSuspendFoundError() error {
	return &errorcodeutil.CustomError{
		Message: ProfileSuspenedFoundErrMsg,
		Code:    int(errorcodeutil.ProfileSuspended),
	}
}

// ProfileNotFoundError returns an error message when a profile is not found
func ProfileNotFoundError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ProfileNotFoundErrMsg,
		Code:    int(errorcodeutil.ProfileNotFound),
	}
}

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NormalizeMSISDNErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// CheckPhoneNumberExistError check if phone number is registered to another user
func CheckPhoneNumberExistError() error {
	return &errorcodeutil.CustomError{
		Message: PhoneNumberInUseErrMsg,
		Code:    int(errorcodeutil.PhoneNumberInUse),
	}
}

// CheckEmailExistError returned when the provided email already exists.
func CheckEmailExistError() error {
	return &errorcodeutil.CustomError{
		Message: EmailInUseErrMsg,
		Code:    int(errorcodeutil.EmailAddressInUse),
	}
}

// InternalServerError returns an error if something wrong happened in performing the operation
func InternalServerError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: InternalServerErrorMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// PinNotFoundError displays error message when a pin is not found
func PinNotFoundError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINNotFoundErrMsg,
		Code:    int(errorcodeutil.PINNotFound),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PINMismatchErrMsg,
		Code:    int(errorcodeutil.PINMismatch),
	}
}

// CustomTokenError is the error message displayed when a
// custom token is not created
func CustomTokenError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: CustomTokenErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// AuthenticateTokenError is the error message displayed when a
// custom token is not authenticated
func AuthenticateTokenError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: AuthenticateTokenErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// UpdateProfileError is the error message displayed when a
// user profile cannot be updated
func UpdateProfileError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: UpdateProfileErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// AddRecordError is the error message displayed when a
// record fails to be added to the dataerrorcodeutil
func AddRecordError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: AddRecordErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// RetrieveRecordError is the error message displayed when a
// failure occurs while retrieving records from the dataerrorcodeutil
func RetrieveRecordError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: RetrieveRecordErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// LikelyToRecommendError is the error message displayed that
// occurs when the recommendation threshold is crossed
func LikelyToRecommendError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: LikelyToRecommendErrMsg,
		Code:    int(errorcodeutil.UndefinedArguments),
	}
}

// GenerateAndSendOTPError is the error message displayed when a
// generate and send otp fails
func GenerateAndSendOTPError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: GenerateAndSendOTPErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// CheckUserPINError is the error message displayed when
// a server is unable to check if the user has a PIN
func CheckUserPINError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: CheckUserPINErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// ExistingPINError is the error message displayed when a
// pin record fails to be retrieved from dataerrorcodeutil
func ExistingPINError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ExistingPINErrMsg,
		Code:    int(errorcodeutil.PINNotFound),
	}

}

// EncryptPINError  is the error message displayed when
// pin encryption fails
func EncryptPINError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: EncryptPINErrMsg,
		Code:    int(errorcodeutil.PINError),
	}
}

// ValidatePINDigitsError  is the error message displayed when
// invalid  pin digits are given
func ValidatePINDigitsError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ValidatePINDigitsErrMsg,
		Code:    int(errorcodeutil.PINError),
	}

}

// ValidatePINLengthError  is the error message displayed when
// an invalid Pin length is given
func ValidatePINLengthError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: ValidatePINLengthErrMsg,
		Code:    int(errorcodeutil.PINError),
	}

}

// InValidPushTokenLengthError  is the error message displayed when
// an invalid push token is given
func InValidPushTokenLengthError() error {
	return &errorcodeutil.CustomError{
		Err:     fmt.Errorf("invalid push token length"),
		Message: ValidatePushTokenLengthErrMsg,
		Code:    int(errorcodeutil.InvalidPushTokenLength),
	}
}

// WrongEnumTypeError  is the error message displayed when
// an invalid enum is given
func WrongEnumTypeError(value string) error {
	return &errorcodeutil.CustomError{
		Err:     fmt.Errorf("%v", WrongEnumErrMsg),
		Message: fmt.Sprintf(WrongEnumErrMsg, value),
		Code:    int(errorcodeutil.InvalidEnum),
	}

}

// VerifyOTPError returns an error when OTP verification fails
func VerifyOTPError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: OTPVerificationErrMsg,
		Code:    int(errorcodeutil.OTPVerificationFailed),
	}
}

// MissingInputError returns an error when OTP verification fails
func MissingInputError(value string) error {
	return &errorcodeutil.CustomError{
		Err:     nil,
		Message: "expected `%s` to be defined",
		Code:    int(errorcodeutil.OTPVerificationFailed),
	}
}

// InvalidFlavourDefinedError is the error message displayed when
// an invalid flavour is provided as input.
func InvalidFlavourDefinedError() error {
	return &errorcodeutil.CustomError{
		Err:     fmt.Errorf("invalid flavour defined"),
		Message: InvalidFlavourDefinedErrMsg,
		Code:    int(errorcodeutil.InvalidFlavour),
	}
}

// AddPartnerTypeError is an error message displayed when there is a
// failure to create a partner type
func AddPartnerTypeError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: AddPartnerTypeErrMsg,
		Code:    int(errorcodeutil.InvalidEnum),
	}

}

// InvalidPartnerTypeError is an error message displayed when an
// invalid partner type is provided
func InvalidPartnerTypeError() error {
	return &errorcodeutil.CustomError{
		Err:     fmt.Errorf("invalid `partnerType` provided"),
		Message: InvalidPartnerTypeErrMsg,
		Code:    int(errorcodeutil.InvalidEnum),
	}
}

// FetchDefaultCurrencyError is an error message displayed when
// the default currency is not found
func FetchDefaultCurrencyError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: FetchDefaultCurrencyErrMsg,
		Code:    int(errorcodeutil.RecordNotFound),
	}
}

// SupplierNotFoundError returns an error message when a supplier is not found
func SupplierNotFoundError() error {
	return &errorcodeutil.CustomError{
		Message: SupplierNotFoundErrMsg,
		Code:    int(errorcodeutil.ProfileNotFound),
	}

}

// CustomerNotFoundError returns an error message when a customer is not found
func CustomerNotFoundError() error {
	return &errorcodeutil.CustomError{
		Message: CustomerNotFoundErrMsg,
		Code:    int(errorcodeutil.ProfileNotFound),
	}
}

// SupplierKYCAlreadySubmittedNotFoundError is returned when the user trys to
// submit another KCY when then is one already submitted
func SupplierKYCAlreadySubmittedNotFoundError() error {
	return &errorcodeutil.CustomError{
		Message: SupplierKYCAlreadySubmittedErrMsg,
		Code:    int(errorcodeutil.KYCAlreadySubmitted),
	}
}

// PublishKYCNudgeError returns an error message when there's a failure in
// creating a KYC nudge
func PublishKYCNudgeError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: PublishKYCNudgeErrMsg,
		Code:    int(errorcodeutil.PublishNudgeFailure),
	}
}

// InvalidCredentialsError returns an error message when wrong credentials are provided
func InvalidCredentialsError() error {
	return &errorcodeutil.CustomError{
		Err:     fmt.Errorf("invalid credentials, expected a username AND password"),
		Message: InvalidCredentialsErrMsg,
		Code:    int(errorcodeutil.InvalidCredentials),
	}
}

// SaveUserPinError returns an error message when we are unable to save a user pin
func SaveUserPinError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: SaveUserPinErrMsg,
		Code:    int(errorcodeutil.PINError),
	}
}

// GeneratePinError returns an error message when we are unable to generate a temporary PIN
func GeneratePinError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: GeneratePinErrMsg,
		Code:    int(errorcodeutil.PINError),
	}
}

// CompleteSignUpError returns an error message when we are unable
// to CompleteSignup
func CompleteSignUpError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: BioDataErrMsg,
		Code:    int(errorcodeutil.AddNewRecordError),
	}
}

// UsernameInUseError is the error message displayed when the provided username
// is associated with another profile
func UsernameInUseError() error {
	return &errorcodeutil.CustomError{
		Message: UsernameInUseErrMsg,
		Code:    int(errorcodeutil.UsernameInUse),
	}
}

// SecondaryResourceHardResetError this error is returned when there argument to reset a resource has a length of 0
// resource here means secondary phone numbers and secondary emails
func SecondaryResourceHardResetError() error {
	return &errorcodeutil.CustomError{
		Message: ResourceUpdateErrMsg,
		Code:    int(errorcodeutil.UndefinedArguments),
	}
}

// InvalidSladeCodeError when the slade code the edi user profile doesn't match with selected provider
func InvalidSladeCodeError() error {
	return &errorcodeutil.CustomError{
		Message: InvalidSladeCodeErrMsg,
		Code:    int(errorcodeutil.InvalidSladeCode),
	}
}

// ResolveNudgeErr is the error that represents the failure of not
// being able to resolve a given nudge
func ResolveNudgeErr(
	err error,
	flavour feedlib.Flavour,
	name string,
	statusCode *int,
) error {
	if statusCode != nil {
		return &errorcodeutil.CustomError{
			Err: err,
			Message: fmt.Sprintf(
				ResolveNudgeBadStatusErrMsg,
				flavour,
				name,
				statusCode,
			),
			Code: int(errorcodeutil.Internal),
		}
	}

	return &errorcodeutil.CustomError{
		Err: err,
		Message: fmt.Sprintf(
			ResolveNudgeErrMsg,
			flavour,
			name,
		),
		Code: int(errorcodeutil.Internal),
	}
}

// RecordExistsError is the error message displayed when a
// similar record is found in the DB
func RecordExistsError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: RecordExistsErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// RecordDoesNotExistError is the error message displayed when a
// record is not found in the DB
func RecordDoesNotExistError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: RecordDoesNotExistErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// SessionIDError return an error when a ussd sessionId is not provided
func SessionIDError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: SessionIDErrMsg,
		Code:    int(errorcodeutil.Internal),
	}
}

// RoleNotValid return an error when a user does not have the required role
func RoleNotValid(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: RoleNotValidMsg,
		Code:    int(errorcodeutil.RoleNotValid),
	}
}

// NavigationActionsError return an error when user navigation actions can not be manipulated
func NavigationActionsError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: NavActionsError,
		Code:    int(errorcodeutil.NavigationActionsError),
	}
}

// AgentNotFoundError returns an error message when an agent is not found
func AgentNotFoundError(err error) error {
	return &errorcodeutil.CustomError{
		Err:     err,
		Message: AgentNotFoundErrMsg,
		Code:    int(errorcodeutil.ProfileNotFound),
	}

}
