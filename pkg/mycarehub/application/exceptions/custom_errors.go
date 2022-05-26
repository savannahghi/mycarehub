package exceptions

import (
	"fmt"

	"github.com/savannahghi/errorcodeutil"
)

// NormalizeMSISDNError returns an error when normalizing the msisdn fails
func NormalizeMSISDNError(err error) error {
	return &CustomError{
		Err:     err,
		Message: NormalizeMSISDNErrMsg,
		Code:    int(Internal),
	}
}

// UserNotFoundError returns an error message when a user is not found
func UserNotFoundError(err error) error {
	return &CustomError{
		Err:     err,
		Message: UserNotFoundErrMsg,
		Code:    int(UserNotFound),
	}
}

// PinNotFoundError displays error message when a pin is not found
func PinNotFoundError(err error) error {
	return &CustomError{
		Err:     err,
		Message: PINNotFoundErrMsg,
		Code:    int(PINNotFound),
	}
}

// PinMismatchError displays an error when the supplied PIN
// does not match the PIN stored
func PinMismatchError() error {
	return &CustomError{
		Err:     nil,
		Message: PINMismatchErrMsg,
		Code:    int(PINMismatch),
	}
}

// InvalidFlavourDefinedErr is the error message displayed when
// an invalid flavour is provided as input.
func InvalidFlavourDefinedErr(err error) error {
	return &CustomError{
		Err:     fmt.Errorf("invalid flavour defined"),
		Message: InvalidFlavourDefinedErrMsg,
		Code:    int(InvalidFlavour),
	}
}

// SaveUserPinError returns an error message when we are unable to save a user pin
func SaveUserPinError(err error) error {
	return &CustomError{
		Err:     err,
		Message: SaveUserPinErrMsg,
		Code:    int(PINError),
	}
}

// InvalidResetPinPayloadErr returns an error message when the provided reset pin payload is invalid
func InvalidResetPinPayloadErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: InvalidResetPinPayloadErrorMsg,
		Code:    int(InvalidResetPinPayloadError),
	}
}

// EmptyUserIDErr returns an error message when the user id is empty
func EmptyUserIDErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: EmptyUserIDInputErrorMsg,
		Code:    int(EmptyUserIDInputError),
	}
}

// ProfileNotFoundErr returns an error message when the profile is not found
func ProfileNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ProfileNotFoundErrorMsg,
		Code:    int(ProfileNotFound),
	}
}

// ClientProfileNotFoundErr returns an error message when the client profile is not found
func ClientProfileNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ClientProfileNotFoundErrorMsg,
		Code:    int(ProfileNotFound),
	}
}

// ClientCCCIdentifierNotFoundErr returns an error message when the client profile is not found
func ClientCCCIdentifierNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ClientCCCIdentifierNotFoundErrorMsg,
		Code:    int(CCCIdentifierNotFoundError),
	}
}

// StaffProfileNotFoundErr returns an error message when the client profile is not found
func StaffProfileNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: StaffProfileNotFoundErrorMsg,
		Code:    int(ProfileNotFound),
	}
}

// InvalidatePinErr returns an error message when the reset pin is invalid
func InvalidatePinErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: InvalidatePinErrMsg,
		Code:    int(InvalidatePinError),
	}
}

// ResetPinErr returns an error message when the reset pin is invalid
func ResetPinErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ResetPinErrorMsg,
		Code:    int(ResetPinError),
	}
}

// PINExpiredErr returns an error message when the reset pin is invalid
func PINExpiredErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: PINExpiredErrorMsg,
		Code:    int(PINExpiredError),
	}
}

// EmptyInputErr returns an error message when an input is empty
func EmptyInputErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: EmptyInputErrorMsg,
		Code:    int(EmptyInputError),
	}
}

// PINErr returns an error message when the PIN is invalid
func PINErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: PINErrorMsg,
		Code:    int(PINError),
	}
}

// NotOptedInErr returns an error message when the user is not opted in
func NotOptedInErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: NotOptedInErrorMsg,
		Code:    int(NotOptedInError),
	}
}

// NotActiveErr returns an error message when a field is not active
func NotActiveErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: NotActiveErrorMsg,
		Code:    int(NotActiveError),
	}
}

// InvalidContactTypeErr returns an error message when the contact type is invalid
func InvalidContactTypeErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: InvalidContactTypeErrorMsg,
		Code:    int(InvalidContactTypeError),
	}
}

// NoContactsErr returns an error message when there are no contacts
func NoContactsErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: NoContactsErrorMsg,
		Code:    int(NoContactsError),
	}
}

// ContactNotFoundErr returns an error message when the contact is not found
func ContactNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ContactNotFoundErrorMsg,
		Code:    int(ContactNotFoundError),
	}
}

// GenerateTempPINErr returns an error message when the temp pin generation fails
func GenerateTempPINErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GenerateTempPINErrMsg,
		Code:    int(GenerateTempPINError),
	}
}

// ExpiredPinErr returns an error message when the pin is expired
func ExpiredPinErr() error {
	return &CustomError{
		Err:     nil,
		Message: ExpiredPinErrorMsg,
		Code:    int(ExpiredPinError),
	}
}

// LoginCountUpdateErr returns an error message when the login count update fails
func LoginCountUpdateErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: LoginCountUpdateErrorMsg,
		Code:    int(LoginCountUpdateError),
	}
}

// LoginTimeUpdateErr returns an error message when the login time update fails
func LoginTimeUpdateErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: LoginTimeUpdateErrorMsg,
		Code:    int(LoginTimeUpdateError),
	}
}

// NexAllowedLoginTimeErr returns an error message when the login time update fails
func NexAllowedLoginTimeErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: NexAllowedLOginTimeErrorMsg,
		Code:    int(NexAllowedLoginTimeError),
	}
}

// SendSMSErr returns an error message when the SMS sending fails
func SendSMSErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: SendSMSErrorMsg,
		Code:    int(SendSMSError),
	}
}

// FailedToUpdateItemErr returns an error message when the item update fails
func FailedToUpdateItemErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: FailedToUpdateItemErrorMsg,
		Code:    int(FailedToUpdateItemError),
	}
}

// ItemNotFoundErr returns an error message when the item is not found
func ItemNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ItemNotFoundErrorMsg,
		Code:    int(ItemNotFoundError),
	}
}

// InputValidationErr returns an error message when the input is invalid
func InputValidationErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: InputValidationErrorMsg,
		Code:    int(InputValidationError),
	}
}

// EncryptionErr returns an error message when the encryption fails
func EncryptionErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: EncryptionErrorMsg,
		Code:    int(EncryptionError),
	}
}

// FailedToSaveItemErr returns an error message when the item save fails
func FailedToSaveItemErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: FailedToSaveItemErrorMsg,
		Code:    int(FailedToSaveItemError),
	}
}

// GeneratePinErr returns an error message when the pin generation fails
func GeneratePinErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GeneratePinErrorMsg,
		Code:    int(GeneratePinError),
	}
}

// GetInviteLinkErr returns an error message when the invite link generation fails
func GetInviteLinkErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GetInviteLinkErrorMsg,
		Code:    int(GetInviteLinkError),
	}
}

// ValidatePINDigitsErr returns an error message when the pin digits are invalid
func ValidatePINDigitsErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: ValidatePINDigitsErrorMsg,
		Code:    int(ValidatePINDigitsError),
	}
}

// ExistingPINError is the error message displayed when a
// pin record fails to be retrieved from dataerrorcodeutil
func ExistingPINError(err error) error {
	return &CustomError{
		Err:     err,
		Message: ExistingPINErrMsg,
		Code:    int(errorcodeutil.PINNotFound),
	}

}

// InternalErr returns an error message when the server fails
func InternalErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: InternalErrorMsg,
		Code:    int(Internal),
	}
}

// GetFAQContentErr returns an error message when the faq content fails
func GetFAQContentErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GetFAQContentErrorMsg,
		Code:    int(GetFAQContentError),
	}
}

// UpdateClientCaregiverErr returns an error message when client caregiver update fails
func UpdateClientCaregiverErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: UpdateClientCaregiverErrorMsg,
		Code:    int(UpdateClientCaregiverError),
	}
}

// CreateClientCaregiverErr returns an error message when the client caregiver creation fails
func CreateClientCaregiverErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: CreateClientCaregiverErrorMsg,
		Code:    int(CreateClientCaregiverError),
	}
}

// GetLoggedInUserUIDErr returns an error message when the logged in user uid check fails
func GetLoggedInUserUIDErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GetLoggedInUserUIDErrorMsg,
		Code:    int(GetLoggedInUserUIDError),
	}
}

// CheckUserRoleErr returns an error message when the user role check fails
func CheckUserRoleErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: CheckUserRoleErrorMsg,
		Code:    int(CheckUserRoleError),
	}
}

// UserNotAuthorizedErr returns an error message when the user is not authorized to perform the action
func UserNotAuthorizedErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: UserNotAuthorizedErrorMsg,
		Code:    int(UserNotAuthorizedError),
	}
}

// CheckUserPermissionErr returns an error message when the user permission check fails
func CheckUserPermissionErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: CheckUserPermissionErrorMsg,
		Code:    int(CheckUserPermissionError),
	}
}

// AssignRolesErr returns an error message when the user role assignment fails
func AssignRolesErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: AssignRolesErrorMsg,
		Code:    int(AssignRolesError),
	}
}

// GetUserRolesErr returns an error message when the user role retrieval fails
func GetUserRolesErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GetUserRolesErrorMsg,
		Code:    int(GetUserRolesError),
	}
}

// GetUserPermissionsErr returns an error message when the user permission retrieval fails
func GetUserPermissionsErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GetUserPermissionsErrorMsg,
		Code:    int(GetUserPermissionsError),
	}
}

// RevokeRolesErr returns an error message when the user roles' revocation fails
func RevokeRolesErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: RevokeRolesErrorMsg,
		Code:    int(RevokeRolesError),
	}
}

// UserNameExistsErr returns an error message when the item update fails
func UserNameExistsErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: NickNameExistsErrorMsg,
		Code:    int(NickNameExistsError),
	}
}

// ClientHasUnresolvedPinResetRequestErr returns an error message when the client has an unresolved pin reset request
func ClientHasUnresolvedPinResetRequestErr() error {
	return &CustomError{
		Err:     nil,
		Message: ClientHasUnresolvedPinResetRequestErrorMsg,
		Code:    int(ClientHasUnresolvedPinResetRequestError),
	}
}

// RetryLoginErr returns an error message when the user is not authorized to perform the action
func RetryLoginErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: RetryLoginErrorMsg,
		Code:    int(RetryLoginError),
	}
}

// GetAllRolesErr returns an error message when the user role retrieval fails
func GetAllRolesErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: GetAllRolesErrorMsg,
		Code:    int(GetAllRolesError),
	}
}

// FailedSecurityCountExceededErr returns an error message when the user is not authorized to verify the security question response
func FailedSecurityCountExceededErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: FailedSecurityCountExceededErrorMsg,
		Code:    int(FailedSecurityCountExceededError),
	}
}

// SecurityQuestionResponseMismatchErr returns an error message when the security question response does not match
func SecurityQuestionResponseMismatchErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: SecurityQuestionResponseMismatchErrorMsg,
		Code:    int(SecurityQuestionResponseMismatchError),
	}
}

// SecurityQuestionNotFoundErr returns an error message when the security question is not found
func SecurityQuestionNotFoundErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: SecurityQuestionNotFoundErrorMsg,
		Code:    int(SecurityQuestionNotFoundError),
	}
}

// UpdateProfileErr returns an error message when the user profile update fails
func UpdateProfileErr(err error) error {
	return &CustomError{
		Err:     err,
		Message: UpdateProfileErrorMsg,
		Code:    int(UpdateProfileError),
	}
}

// StaffHasUnresolvedPinResetRequestErr returns an error message when the staff has an unresolved pin reset request
func StaffHasUnresolvedPinResetRequestErr() error {
	return &CustomError{
		Err:     nil,
		Message: StaffHasUnresolvedPinResetRequestErrorMsg,
		Code:    int(StaffHasUnresolvedPinResetRequestError),
	}
}
