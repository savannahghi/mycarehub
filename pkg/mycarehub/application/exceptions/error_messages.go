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
)
