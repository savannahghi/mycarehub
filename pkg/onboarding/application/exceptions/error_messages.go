package exceptions

const (

	//UsernameInUseErrMsg  is the error message displayed when the provided userName is associated with a profile already
	UsernameInUseErrMsg = "provided username is already in use"

	//PhoneNUmberInUseErrMsg is the error message displayed when a phone number provided is associated with a profile already
	PhoneNUmberInUseErrMsg = "provided phone number is already in use"

	// UserNotFoundErrMsg is the error message displayed when a user is not found
	UserNotFoundErrMsg = "failed to get a user"

	// ProfileNotFoundErrMsg is the error message displayed when a user is not found
	ProfileNotFoundErrMsg = "failed to get a user profile"

	// PINNotFoundErrMsg is the error message displayed when a pin is not found
	PINNotFoundErrMsg = "failed to get a user pin"

	// CustomTokenErrMsg is the error message displayed when a
	// custom token is not created
	CustomTokenErrMsg = "failed to create custom token"

	// AuthenticateTokenErrMsg is the error message displayed when a
	// custom token is not authenticated
	AuthenticateTokenErrMsg = "failed to authenticate custom token"

	// UpdateProfileErrMsg is the error message displayed when a
	// user profile cannot be updated
	UpdateProfileErrMsg = "failed to update a user profile"

	// AddRecordErrMsg is the error message displayed when a
	// record fails to be added to the database
	AddRecordErrMsg = "failed to add the record to the database"

	// LikelyToRecommendErrMsg is the error message displayed that
	// occurs when the recommendation threshold is crossed
	LikelyToRecommendErrMsg = "the likelihood of recommending should be an int between 0 and 10"

	// ValidatePINLengthErrMsg  is the error message displayed when
	// an invalid Pin length is given
	ValidatePINLengthErrMsg = "pin should be of 4,5, or six digits"

	// ValidatePINDigitsErrMsg  is the error message displayed when
	// an invalid  pin digits are given
	ValidatePINDigitsErrMsg = "pin should be a valid number"

	// UsePinExistErrMsg  is the error message displayed when
	// user has pin already during set pin
	UsePinExistErrMsg = "the user has PIN already"

	// EncryptPINErrMsg  is the error message displayed when
	// pin encryption fails
	EncryptPINErrMsg = "unable to encrypt PIN"

	// RetrieveRecordErrMsg is the error message displayed when a
	// record fails to be retrieved from database
	RetrieveRecordErrMsg = "unable to retrieve newly created record"

	// ExistingPINErrMsg is the error message displayed when a
	// pin record fails to be retrieved from database
	ExistingPINErrMsg = "request for a PIN reset failed. User does not have an existing PIN"

	// CheckUserPINErrMsg is the error message displayed when
	// server unable to check if the user has a PIN
	CheckUserPINErrMsg = "unable to check if the user has a PIN"

	// GenerateAndSendOTPErrMsg is the error message displayed when a
	// generate and send otp fails
	GenerateAndSendOTPErrMsg = "failed to generate and send an otp"

	// NormalizeMSISDNErrMsg is the error message displayed when
	// normalize the msisdn(phone number) fails
	NormalizeMSISDNErrMsg = "unable to normalize the msisdn"

	// PINMismatchErrMsg is the error message displayed when
	// the user supplied PIN number does not match the PIN
	// record we have stored
	PINMismatchErrMsg = "wrong PIN credentials supplied"

	// InternalServerErrorMsg is an error message for database CRUD operations that
	// don't succeed e.g network latency
	InternalServerErrorMsg = "server error! unable to perform operation"

	// ValidatePuskTokenLengthErrMsg ...
	ValidatePuskTokenLengthErrMsg = "invalid push token detected"

	// WrongEnumErrMsg is an error message returned when a wrong enum
	// type is supplied
	WrongEnumErrMsg = "a wrong enum %s has been provided"

	// OTPVerificationErrMsg is an error message that is returned when
	// a given OTP code and Phone number fails verifciation
	OTPVerificationErrMsg = "failed to verify OTP"
)
