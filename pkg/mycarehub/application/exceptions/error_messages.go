package exceptions

const (

	// UserNotFoundErrMsg is the error message displayed when a user is not found
	UserNotFoundErrMsg = "failed to get a user"

	// PINNotFoundErrMsg is the error message displayed when a pin is not found
	PINNotFoundErrMsg = "failed to get a user pin"

	// NormalizeMSISDNErrMsg is the error message displayed when
	// normalize the msisdn(phone number) fails
	NormalizeMSISDNErrMsg = "unable to normalize the msisdn"

	// PINMismatchErrMsg is the error message displayed when
	// the user supplied PIN number does not match the PIN
	// record we have stored
	PINMismatchErrMsg = "wrong PIN credentials supplied"

	// InvalidFlavourDefinedErrMsg for invalid flavour definitions
	InvalidFlavourDefinedErrMsg = "invalid flavour defined"

	// SaveUserPinErrMsg is displayed when a user pin is not saved
	SaveUserPinErrMsg = "unable to save user PIN"

	// InvalidResetPinPayloadErrorMsg is displayed when a pin payload is invalid
	InvalidResetPinPayloadErrorMsg = "failed to validate reset pin payload"

	// EmptyUserIDInputErrorMsg is displayed when a user id input is empty
	EmptyUserIDInputErrorMsg = "user id input is empty"

	// ProfileNotFoundErrorMsg is displayed when a profile is not found
	ProfileNotFoundErrorMsg = "user profile not found"

	// StaffProfileNotFoundErrorMsg is displayed when a staff profile is not found
	StaffProfileNotFoundErrorMsg = "staff profile not found"

	// InvalidatePinErrMsg is displayed when the invalidate action for reset pin fails
	InvalidatePinErrMsg = "unable to invalidate reset pin"

	// ResetPinErrorMsg is displayed when the reset pin action fails
	ResetPinErrorMsg = "unable to reset pin"

	// PINExpiredErrorMsg is the error message displayed when a PIN is expired
	PINExpiredErrorMsg = "pin expired"

	// EmptyInputErrorMsg is the error message displayed when an input is empty
	EmptyInputErrorMsg = "input is empty"

	// PINErrorMsg is the error message displayed when a PIN is invalid
	PINErrorMsg = "invalid pin"

	// NotOptedInErrorMsg  is the error message displayed when a user is not opted in
	NotOptedInErrorMsg = "user not opted in"

	// NotActiveErrorMsg  is the error message displayed when a field is not active
	NotActiveErrorMsg = "field active is false"

	// InvalidContactTypeErrorMsg is the error message displayed when a contact type is invalid
	InvalidContactTypeErrorMsg = "invalid contact type"

	// NoContactsErrorMsg is the error message displayed when a user has no contacts
	NoContactsErrorMsg = "user has no contacts"

	// ContactNotFoundErrorMsg is the error message displayed when a contact is not found
	ContactNotFoundErrorMsg = "contact not found"

	// GenerateTempPINErrMsg is the error message displayed when a temp pin is not generated
	GenerateTempPINErrMsg = "unable to generate temporary pin"

	// ExpiredPinErrorMsg is the error message displayed when a pin is expired
	ExpiredPinErrorMsg = "pin expired"

	// LoginCountUpdateErrorMsg is the error message displayed when a login count update fails
	LoginCountUpdateErrorMsg = "unable to update login count"

	// LoginTimeUpdateErrorMsg is the error message displayed when a login time update fails
	LoginTimeUpdateErrorMsg = "unable to update login time"

	// NexAllowedLOginTimeErrorMsg is the error message displayed when a login time is not allowed
	NexAllowedLOginTimeErrorMsg = "login time not allowed"

	// SendSMSErrorMsg is the error message displayed when a SMS is not sent
	SendSMSErrorMsg = "unable to send SMS"

	// FailedToUpdateItemErrorMsg is the error message displayed when an item is not updated
	FailedToUpdateItemErrorMsg = "failed to update item"

	// ItemNotFoundErrorMsg is the error message displayed when an item is not found
	ItemNotFoundErrorMsg = "item not found"

	// InputValidationErrorMsg is the error message displayed when an input is invalid
	InputValidationErrorMsg = "input validation failed"

	// EncryptionErrorMsg is the error message displayed when an encryption fails
	EncryptionErrorMsg = "encryption failed"

	// FailedToSaveItemErrorMsg is the error message displayed when an item is not saved
	FailedToSaveItemErrorMsg = "failed to save item"

	// GeneratePinErrorMsg is the error message displayed when a pin is not generated
	GeneratePinErrorMsg = "unable to generate pin"

	// GetInviteLinkErrorMsg is the error message displayed when an invite link is not generated
	GetInviteLinkErrorMsg = "unable to generate invite link"

	// ValidatePINDigitsErrorMsg is the error message displayed when a pin is invalid
	ValidatePINDigitsErrorMsg = "invalid pin digits"

	// ExistingPINErrMsg is the error message displayed when a
	// pin record fails to be retrieved from database
	ExistingPINErrMsg = "user does not have an existing PIN"

	// ClientProfileNotFoundErrorMsg is displayed when a profile is not found
	ClientProfileNotFoundErrorMsg = "client profile not found"

	// ClientCCCIdentifierNotFoundErrorMsg is displayed when a clients CCC number identifier is not found
	ClientCCCIdentifierNotFoundErrorMsg = "client ccc identifier not found"

	// InternalErrorMsg is the error message displayed when an internal server error occurs
	InternalErrorMsg = "internal error"

	// UpdateClientCaregiverErrorMsg is the error message displayed when a caregiver is not updated
	UpdateClientCaregiverErrorMsg = "unable to update caregiver"

	// CreateClientCaregiverErrorMsg is the error message displayed when a caregiver is not created
	CreateClientCaregiverErrorMsg = "unable to create caregiver"

	// GetLoggedInUserUIDErrorMsg is the error message displayed when a logged in user uid is not found
	GetLoggedInUserUIDErrorMsg = "unable to get logged in user uid"

	// CheckUserRoleErrorMsg is the error message displayed when a user role is not found
	CheckUserRoleErrorMsg = "unable to check user role"

	// UserNotAuthorizedErrorMsg is the error message displayed when a user is not authorized
	UserNotAuthorizedErrorMsg = "user not authorized"

	// CheckUserPermissionErrorMsg is the error message displayed when a user permission is not found
	CheckUserPermissionErrorMsg = "unable to check user permission"
	//AssignRolesErrorMsg is the error message displayed when a user role assignment fails
	AssignRolesErrorMsg = "unable to assign roles"
	// GetUserRolesErrorMsg is the error message displayed when a user role is not found
	GetUserRolesErrorMsg = "unable to get user roles"
	// GetUserPermissionsErrorMsg is the error message displayed when a user permission is not found
	GetUserPermissionsErrorMsg = "unable to get user permissions"

	// RevokeRolesErrorMsg is the error message displayed when a user roles' revocation fails
	RevokeRolesErrorMsg = "unable to revoke roles"

	// NickNameExistsErrorMsg is the error message displayed when a user attempts to set a nickname that
	// has already been taken
	NickNameExistsErrorMsg = "username has already been taken"

	// ClientHasUnresolvedPinResetRequestErrorMsg is the error message displayed when a client has an unresolved pin reset request
	ClientHasUnresolvedPinResetRequestErrorMsg = "client has unresolved pin reset request"

	// RetryLoginErrorMsg is the error message displayed when a user attempts to login when they have exponential backoff
	RetryLoginErrorMsg = "retry login failed due to exponential backoff"

	// GetAllRolesErrorMsg is the error message displayed when a user role is not found
	GetAllRolesErrorMsg = "unable to query all roles"

	// RecordNotFoundErrorMsg is the error message displayed when a record is not found
	RecordNotFoundErrorMsg = "record not found"

	// FailedSecurityCountExceededErrorMsg is the error message displayed when a user has failed security count
	FailedSecurityCountExceededErrorMsg = "you have reached the maximum number of attempts"

	// SecurityQuestionResponseMismatchErrorMsg is the error message displayed when a user's security question response does not match
	SecurityQuestionResponseMismatchErrorMsg = "security question response does not match"

	// SecurityQuestionNotFoundErrorMsg is the error message displayed when a user's security question is not found
	SecurityQuestionNotFoundErrorMsg = "security question not found"

	// UpdateProfileErrorMsg is the error message displayed when a user's profile is not updated
	UpdateProfileErrorMsg = "unable to update profile"

	// StaffHasUnresolvedPinResetRequestErrorMsg is the error message displayed when a staff has an unresolved pin reset request
	StaffHasUnresolvedPinResetRequestErrorMsg = "staff has unresolved pin reset request"

	// FailedToCreateAnOrganizationErrorMsg is the error message displayed when an organization is not created
	FailedToCreateAnOrganizationErrorMsg = "failed to create an organization"

	// DuplicateOrganisationCodeErrorMessage is the error message to be displayed when an existing organisation code is used to register a new organisation
	DuplicateOrganisationCodeErrorMessage = "the provided organisation code is associated with another organisation"

	// DuplicateOrganisationNameErrorMessage is the error message to be displayed when an existing organisation name is used to create a new organisation
	DuplicateOrganisationNameErrorMessage = "the provided name is associated with another organisation"

	// DuplicateOrganisationPhoneNumberErrorMessage is the error message to be displayed when an existing phone number is used to create a new organisation
	DuplicateOrganisationPhoneNumberErrorMessage = "the provided phone number is associated with another organisation"

	// DuplicateOrganisationEmailAddressErrorMessage is the error message to be displayed when an existing email is used to create a new organisation
	DuplicateOrganisationEmailAddressErrorMessage = "the provided email is associated with another organisation"
)
