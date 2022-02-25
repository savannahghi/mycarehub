package exceptions

// ErrorCode are  used to determine the nature of an error, and why it occurred
// both the frontend and backend should be aware of these codes
type ErrorCode int

const (
	// OK is returned on success.
	OK ErrorCode = iota + 1

	// Internal errors means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	// it's value is 2
	Internal

	// UndefinedArguments errors means either one or more arguments to
	// a method have not been specified
	// it's value is 3
	UndefinedArguments

	// PhoneNumberInUse indicates that a phone number has an associated user profile.
	// this error can occur when fetching a user profile using a phone number, to check
	// that the phone number has not already been registered. The check usually runs
	// on both PRIMARY PHONE and SECONDARY PHONE
	// it's value is 4
	PhoneNumberInUse

	// EmailAddressInUse indicates that an email address has an associated user profile.
	// this error can occur when fetching a user profile using an email address, to check
	// that the email address has not already been registered. The check usually runs
	// on both PRIMARY EMAIL ADDRESS and SECONDARY EMAIL ADDRESS.
	// it's value is 5
	EmailAddressInUse

	// UsernameInUse indicates that a username has an associated user profile.
	// this error can occur when trying a update a user's username with a username that already has been taken
	// it's value is 6
	UsernameInUse

	// ProfileNotFound errors means a user profile does not exist with the provided parameters
	// This occures when fetching a user profile either by UID, ID , PHONE NUMBER or EMAIL and no
	// matching record is found
	// it's value is 7
	ProfileNotFound

	// PINMismatch errors means that the provided PINS do not match (are not similar)
	// it's value is 8
	PINMismatch

	// PINNotFound errors means a user PIN does not exist with the provided parameters
	// This occurs when fetching a PIN by the user's user profile ID and no
	// matching record is found. This should never occur and if it does then it means
	// there is a serious issue with our data
	// it's value is 9
	PINNotFound

	// UserNotFound errors means that a user's firebase auth account does not exists. This occurs
	// when fetching a firebase user by either a phone number or an email and their record is not found
	// it's value is 10
	UserNotFound

	// ProfileSuspended error means that user's profile has been suspended.
	// This may occur due to violation of terms or detection of suspicious activity
	// It's value is 11
	ProfileSuspended

	// PINError error means that some actions could not be performed on the PIN.
	// This may occur when the provided PIN cannot be encrypted, cannot be validated and/or is of invalid length
	// It's value is 12
	PINError

	// InvalidPushTokenLength means that an invalid push token was given.
	// This may occur when the lenth of the issued token is of less then the minimum character(1250)
	// It's error code is 13
	InvalidPushTokenLength

	// InvalidEnum means that the provided enumerator was of invalid.
	// This may occur when an invalid enum value has been defined. For example, PartnerType, LoginProviderType e.t.c
	// It's error code is 14
	InvalidEnum

	// OTPVerificationFailed means that the provide OTP could not be verified
	// This may occur when an incorrect OTP is supplied
	// It's error code is 15
	OTPVerificationFailed

	// MissingInput means that no OTP was submiited
	// This may occur when a user fails to provide an OTP but makes a submission
	// It's error code id 16
	MissingInput

	// InvalidFlavour means that the provide falvour is invalid
	// This may happen when the provided flavour is not consumer or pro
	// It's error code is 17
	InvalidFlavour

	// RecordNotFound means that the provided record is not found.
	// This may happen when the provided data e.g currency, user etc is not accepted
	// It's error code is 18
	RecordNotFound

	// UnableToFindProvider means that the selected provider could not be found
	// This may happen if the provider is not specified in the charge master
	// It's error code is 19
	UnableToFindProvider

	// PublishNudgeFailure means that there was an error while publishing a nudge
	// It's error code is 20
	PublishNudgeFailure

	// InvalidCredentials means that the provided credentials are invalid
	// This may happen when any of the customers provides wrong credentials
	// It's error code is 21
	InvalidCredentials

	// AddNewRecordError means that the record could not be saved
	// This may happen may be as a result of wrong credentials or biodata
	// It's error code is 22
	AddNewRecordError

	// RoleNotValid means that the user role does not match the role required
	// to perform the current operation that the user is trying to perform.
	// Its error code is 23
	RoleNotValid

	//UserNotAuthorizedToAccessThisResource means that the subject's
	//email has been found to not have access to the specified resource
	//Its error code is 24
	UserNotAuthorizedToAccessThisResource

	//UnableToCheckIfUserIsAnAdmin means that
	//checking to see if a user is an admin has failed
	//Its error code is 25
	UnableToCheckIfUserIsAnAdmin

	//LoggedInUserIsNotAnAdmin means that
	//the user currently logged in has been found to not be an admin
	//Its error code is 26
	LoggedInUserIsNotAnAdmin

	// UnableToRetrieveNotification means that
	//retrieving a node from firestore fails with this ID
	//Its error code is 27
	UnableToRetrieveNotification

	//UnableToSaveNotification means that
	//saving a notification after updating it to read has failed
	//Its error code is 28
	UnableToSaveNotification

	//NoConfirmedPhoneNumbers means that
	//a user's primary phone number is  nil
	//Its error code is 29
	NoConfirmedPhoneNumbers

	//InvalidPhoneNumberFormat means that
	//the phone number format is invalid
	//Its error code is 30
	InvalidPhoneNumberFormat

	// UnableToSendText means that
	//sending a text to the phone number in question has failed
	//Its  error code is 31
	UnableToSendText

	//UnknownStateProvided means that
	//an unknown state has been entered
	//Its error code is 32
	UnknownStateProvided

	//NavigationActionsError means that
	//the system is not able to update or retrive a users navigation actions
	//Its error code is 33
	NavigationActionsError

	// GetInviteLinkError means that the system is unable to get a user's invite link'
	// the flavour passed when generating the invite link is invalid
	// Its error code is 34
	GetInviteLinkError

	// SendInviteSMSError means that the system is unable to send an invite SMS' to a user'
	// the system failed to make a successful request to the messaging service
	// Its error code is 35
	SendInviteSMSError

	// GenerateTempPINError means that the system is unable to generate a temporary PIN'
	// the random number generator has failed (which is highly unlikely)
	// Its error code is 36
	GenerateTempPINError

	// InvalidResetPinPayloadError means that the system is unable to validate the reset pin input'
	// the user ID or the flavor are not provided
	// Its error code is 37
	InvalidResetPinPayloadError

	// EmptyUserIDInputError means that the system is unable to userID input'
	// the user ID is empty
	// Its error code is 38
	EmptyUserIDInputError

	// InvalidatePinError means that the system is unable to invalidate a reset pin'
	// the invalidation action has failed
	// Its error code is 39
	InvalidatePinError

	// ResetPinError means that the system is unable to reset a pin'
	// the reset action has failed
	// Its error code is 40
	ResetPinError

	// PINExpiredError means that the pin provided is expired'
	// the pin's expiration time' has passed
	// Its error code is 41
	PINExpiredError

	// EmptyInputError means that the system is unable to validate the input'
	// the input is empty
	// Its error code is 42
	EmptyInputError

	// NotOptedInError means that the system user has not opted in the input for a contact'
	// the user has not opted in for the input'
	// Its error code is 43
	NotOptedInError

	// NotActiveError means that the a field is not active'
	// the field is not active'
	// Its error code is 44
	NotActiveError

	// InvalidContactTypeError means that the system is unable to validate the contact type'
	// the contact type is invalid'
	// Its error code is 45
	InvalidContactTypeError

	// NoContactsError means that the system could not find any contacts'
	// the user has no contacts'
	// Its error code is 46
	NoContactsError

	// ContactNotFoundError means that the system could not find the contact'
	// the contact was not found'
	// Its error code is 47
	ContactNotFoundError

	// ExpiredPinError means that the system is unable to validate the pin'
	// The pin provided has expired'
	// Its error code is 48
	ExpiredPinError

	// LoginCountUpdateError means that the system is unable to update the login count'
	// The login count update has failed'
	// Its error code is 49
	LoginCountUpdateError

	// LoginTimeUpdateError means that the system is unable to update the login time'
	// The login time update has failed'
	// Its error code is 50
	LoginTimeUpdateError

	// NexAllowedLOginTimeError means that the system is unable to validate the login time'
	// The login time is not allowed'
	// Its error code is 51
	NexAllowedLOginTimeError

	// SendSMSError means that the system is unable to send an SMS'
	// The SMS sending has failed'
	// Its error code is 52
	SendSMSError

	// FailedToUpdateItemError means that the system is unable to update an item'
	// The update has failed'
	// Its error code is 53
	FailedToUpdateItemError

	// ItemNotFoundError means that the system is unable to find an item'
	// The item was not found'
	// Its error code is 54
	ItemNotFoundError

	// InputValidationError means that the system is unable to validate the input'
	// The input is invalid'
	// Its error code is 55
	InputValidationError

	// EncryptionError means that the system is unable to encrypt the input'
	// The encryption has failed'
	// Its error code is 56
	EncryptionError

	// FailedToSaveItemError means that the system is unable to save an item'
	// The save has failed'
	// Its error code is 57
	FailedToSaveItemError

	// GeneratePinError means that the system is unable to generate a pin'
	// The pin generation has failed'
	// Its error code is 58
	GeneratePinError

	// ValidatePINDigitsError means that the system is unable to validate the pin'
	// The pin is invalid'
	// Its error code is 59
	ValidatePINDigitsError

	// GetFAQContentError means that the system is unable to get the FAQ content'
	// The FAQ content retrieval has failed'
	// Its error code is 60
	GetFAQContentError

	// UpdateClientCaregiverError means that the system is unable to update the client caregiver'
	// The update has failed'
	// Its error code is 61
	UpdateClientCaregiverError

	// CreateClientCaregiverError means that the system is unable to create the client caregiver'
	// The creation has failed'
	// Its error code is 62
	CreateClientCaregiverError

	// GetLoggedInUserUIDError means that the system is unable to get the logged in user UID'
	// The retrieval has failed'
	// Its error code is 63
	GetLoggedInUserUIDError

	// CheckUserRoleError means that the system is unable to check the user role'
	// The check has failed'
	// Its error code is 64
	CheckUserRoleError

	// UserNotAuthorizedError means that the system is unable to validate the user'
	// The user is not authorized'
	// Its error code is 65
	UserNotAuthorizedError

	// CheckUserPermissionError means that the system is unable to check the user permission'
	// The check has failed'
	// Its error code is 66
	CheckUserPermissionError
	// AssignRolesError means that the system is unable to assign roles'
	// The assignment has failed'
	// Its error code is 67
	AssignRolesError

	// GetUserRolesError means that the system is unable to get the user roles'
	// The retrieval has failed'
	// Its error code is 68
	GetUserRolesError

	// GetUserPermissionsError means that the system is unable to get the user permissions'
	// The retrieval has failed'
	// Its error code is 69
	GetUserPermissionsError
)
