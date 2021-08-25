package exceptions

const (

	//UsernameInUseErrMsg  is the error message displayed when the provided userName is associated with a profile already
	UsernameInUseErrMsg = "provided username is already in use"

	//PhoneNumberInUseErrMsg is the error message displayed when a phone number provided is associated with a profile already
	PhoneNumberInUseErrMsg = "provided phone number is already in use"

	//EmailInUseErrMsg is the error message displayed when a email provided is associated with a profile already
	EmailInUseErrMsg = "provided email address is already in use"

	// UserNotFoundErrMsg is the error message displayed when a user is not found
	UserNotFoundErrMsg = "failed to get a user"

	// ProfileNotFoundErrMsg is the error message displayed when a user is not found
	ProfileNotFoundErrMsg = "failed to get a user profile"

	// ProfileSuspenedFoundErrMsg is the error message displayed when a user's profile has been suspended
	ProfileSuspenedFoundErrMsg = "the user profile has been suspended"

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
	ExistingPINErrMsg = "user does not have an existing PIN"

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

	// ValidatePushTokenLengthErrMsg ...
	ValidatePushTokenLengthErrMsg = "invalid push token detected"

	// WrongEnumErrMsg is an error message returned when a wrong enum
	// type is supplied
	WrongEnumErrMsg = "a wrong enum %s has been provided"

	// OTPVerificationErrMsg is an error message that is returned when
	// a given OTP code and Phone number fails verifciation
	OTPVerificationErrMsg = "failed to verify OTP"

	// InvalidFlavourDefinedErrMsg for invalid flavour definitions
	InvalidFlavourDefinedErrMsg = "invalid flavour defined"

	// AddPartnerTypeErrMsg is an error message displayed when there is a
	// failure to create a partner type
	AddPartnerTypeErrMsg = "error occurred while adding partner type"

	// InvalidPartnerTypeErrMsg is an error message displayed when an
	// invalid partner type is provided
	InvalidPartnerTypeErrMsg = "invalid `partnerType` provided"

	// FetchDefaultCurrencyErrMsg is an error message displayed when
	// the default currency is not found
	FetchDefaultCurrencyErrMsg = "unable to fetch orgs default currency"

	// SupplierNotFoundErrMsg is an error message displayed when a supplier
	// profile is not found
	SupplierNotFoundErrMsg = "unable to get the user supplier profile"

	// CustomerNotFoundErrMsg  is an error message displayed when a customer
	// profile is not found
	CustomerNotFoundErrMsg = "unable to get the user customer profile"

	// SupplierKYCAlreadySubmittedErrMsg thrown when another kyc is been submitted whne an old exists
	// A supplier should have only one kyc
	SupplierKYCAlreadySubmittedErrMsg = "kyc already exists for supplier"

	// PublishKYCNudgeErrMsg is displayed if we are unable to publish a kyc nudge
	PublishKYCNudgeErrMsg = "unable to publish kyc nudge"

	// InvalidCredentialsErrMsg is displayed when wrong credentials are provided
	InvalidCredentialsErrMsg = "invalid credentials, expected a username AND password"

	// SaveUserPinErrMsg is displayed when a user pin is not saved
	SaveUserPinErrMsg = "unable to save user PIN"

	// GeneratePinErrMsg is displayed when a user pin is not generated
	GeneratePinErrMsg = "unable to generate temporary PIN for new user"

	// BioDataErrMsg is displayed when a user BioData is incomplete
	BioDataErrMsg = "incomplete Bio Data, expected first and last name"

	// ResourceUpdateErrMsg is throw when trying to update a resource with empty arguments
	ResourceUpdateErrMsg = "can't update resource with empty arguments"

	// InvalidSladeCodeErrMsg is displayed when slade code the edi user profile doesn't match with selected provider
	InvalidSladeCodeErrMsg = "invalid slade code, doesn't match with selected provider"

	// ResolveNudgeErrMsg is displayed when a nudge fails to be resolved
	// It takes the `flavour` and a `name` of the nudge to be more specific
	ResolveNudgeErrMsg = "failed to resolve %v `%s` nudge"

	// ResolveNudgeBadStatusErrMsg is displayed when a nudge fails to be resolved
	// because of an unexpected `Status code` being returned
	// It takes the `flavour`, `name` and a `StatusCode` f orthe nudge to be more specific
	ResolveNudgeBadStatusErrMsg = "failed to resolve %v `%s` nudge with status code %v"

	// RecordExistsErrMsg is the error message displayed an existing record is found
	RecordExistsErrMsg = "a similar record has been found in the database"

	// RecordDoesNotExistErrMsg is the error message displayed when a record cannot be found in
	// the database
	RecordDoesNotExistErrMsg = "the record cannot be found in the database"

	//SessionIDErrMsg is an error message displayed when a sessionId is not provided in
	// in the USSD details
	SessionIDErrMsg = "sessionID not provided"

	//RoleNotValidMsg is an error message displayed when the provided role is not valid
	// and the user doesn't have permissions to perform an action
	RoleNotValidMsg = "user role is not valid"

	//NavActionsError is an error message displayed when the system cannot update navigation actions
	NavActionsError = "navigation actions not updated"

	// AgentNotFoundErrMsg is an error message displayed when a supplier
	// profile is not found
	AgentNotFoundErrMsg = "user with agent role not found"
)
