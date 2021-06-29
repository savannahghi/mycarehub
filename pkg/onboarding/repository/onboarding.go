package repository

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// SupplierRepository  defines signatures that relate to suppliers
type SupplierRepository interface {
	// supplier methods
	GetSupplierProfileByID(ctx context.Context, id string) (*base.Supplier, error)

	GetSupplierProfileByProfileID(ctx context.Context, profileID string) (*base.Supplier, error)

	UpdateSupplierProfile(ctx context.Context, profileID string, data *base.Supplier) error

	AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *base.PartnerType) (bool, error)

	AddSupplierAccountType(ctx context.Context, profileID string, accountType base.AccountType) (*base.Supplier, error)

	StageProfileNudge(ctx context.Context, nudge *base.Nudge) error

	StageKYCProcessingRequest(ctx context.Context, data *domain.KYCRequest) error

	// RemoveKYCProcessingRequest called when the user seeks to retire a kyc processing request.
	RemoveKYCProcessingRequest(ctx context.Context, supplierProfileID string) error

	// sets the active attribute of supplier profile to true
	ActivateSupplierProfile(profileID string, supplier base.Supplier) (*base.Supplier, error)

	FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error)

	FetchKYCProcessingRequestByID(ctx context.Context, id string) (*domain.KYCRequest, error)

	UpdateKYCProcessingRequest(ctx context.Context, sup *domain.KYCRequest) error
	CheckIfAdmin(profile *base.UserProfile) bool
}

// CustomerRepository  defines signatures that relate to customers
type CustomerRepository interface {
	// customer methods
	GetCustomerProfileByID(ctx context.Context, id string) (*base.Customer, error)

	GetCustomerProfileByProfileID(ctx context.Context, profileID string) (*base.Customer, error)

	UpdateCustomerProfile(
		ctx context.Context,
		profileID string,
		cus base.Customer,
	) (*base.Customer, error)
}

// OnboardingRepository interface that provide access to all persistent storage operations
type OnboardingRepository interface {
	UserProfileRepository

	SupplierRepository

	CustomerRepository

	AgentRepository

	// creates a user profile of using the provided phone number and uid
	CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error)

	// creates a new user profile that is pre-filled using the provided phone number
	CreateDetailedUserProfile(ctx context.Context, phoneNumber string, profile base.UserProfile) (*base.UserProfile, error)

	// creates an empty supplier profile
	CreateEmptySupplierProfile(ctx context.Context, profileID string) (*base.Supplier, error)

	// create a new supplier profile that is pre-filled using the provided profile ID
	CreateDetailedSupplierProfile(ctx context.Context, profileID string, supplier base.Supplier) (*base.Supplier, error)

	// creates an empty customer profile
	CreateEmptyCustomerProfile(ctx context.Context, profileID string) (*base.Customer, error)

	// fetches a user profile by uid
	GetUserProfileByUID(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error)

	// fetches a user profile by id. returns the unsuspend profile
	GetUserProfileByID(ctx context.Context, id string, suspended bool) (*base.UserProfile, error)

	// fetches a user profile by phone number
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error)

	// fetches a user profile by primary phone number
	GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phoneNumber string, suspend bool) (*base.UserProfile, error)

	// checks if a specific phone number has already been registered to another user
	CheckIfPhoneNumberExists(ctx context.Context, phone string) (bool, error)

	// checks if a specific email has already been registered to another user
	CheckIfEmailExists(ctx context.Context, phone string) (bool, error)

	// checks if a specific username has already been registered to another user
	CheckIfUsernameExists(ctx context.Context, phone string) (bool, error)

	GenerateAuthCredentialsForAnonymousUser(ctx context.Context) (*base.AuthCredentialResponse, error)

	GenerateAuthCredentials(ctx context.Context, phone string, profile *base.UserProfile) (*base.AuthCredentialResponse, error)

	FetchAdminUsers(ctx context.Context) ([]*base.UserProfile, error)

	// removes user completely. This should be used only under testing environment
	PurgeUserByPhoneNumber(ctx context.Context, phone string) error

	HardResetSecondaryPhoneNumbers(ctx context.Context, profile *base.UserProfile, newSecondaryPhones []string) error

	HardResetSecondaryEmailAddress(ctx context.Context, profile *base.UserProfile, newSecondaryEmails []string) error

	// PINs
	GetPINByProfileID(
		ctx context.Context,
		ProfileID string,
	) (*domain.PIN, error)

	// Record post visit survey
	RecordPostVisitSurvey(
		ctx context.Context,
		input dto.PostVisitSurveyInput,
		UID string,
	) error

	// User Pin methods
	SavePIN(ctx context.Context, pin *domain.PIN) (bool, error)
	UpdatePIN(ctx context.Context, id string, pin *domain.PIN) (bool, error)

	ExchangeRefreshTokenForIDToken(
		token string,
	) (*base.AuthCredentialResponse, error)

	GetCustomerOrSupplierProfileByProfileID(
		ctx context.Context,
		flavour base.Flavour,
		profileID string,
	) (*base.Customer, *base.Supplier, error)

	GetOrCreatePhoneNumberUser(ctx context.Context, phone string) (*dto.CreatedUserResponse, error)

	AddUserAsExperimentParticipant(ctx context.Context, profile *base.UserProfile) (bool, error)

	RemoveUserAsExperimentParticipant(ctx context.Context, profile *base.UserProfile) (bool, error)

	CheckIfExperimentParticipant(ctx context.Context, profileID string) (bool, error)

	AddNHIFDetails(
		ctx context.Context,
		input dto.NHIFDetailsInput,
		profileID string,
	) (*domain.NHIFDetails, error)

	GetNHIFDetailsByProfileID(
		ctx context.Context,
		profileID string,
	) (*domain.NHIFDetails, error)

	GetUserCommunicationsSettings(ctx context.Context, profileID string) (*base.UserCommunicationsSetting, error)

	SetUserCommunicationsSettings(ctx context.Context, profileID string,
		allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error)

	PersistIncomingSMSData(ctx context.Context, input *dto.AfricasTalkingMessage) error

	AddAITSessionDetails(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error)
	GetAITSessionDetails(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error)
	UpdateSessionLevel(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error)
	UpdateSessionPIN(ctx context.Context, sessionID string, pin string) (*domain.USSDLeadDetails, error)

	StageCRMPayload(ctx context.Context, payload dto.ContactLeadInput) error
	UpdateStageCRMPayload(ctx context.Context, phoneNumber string, payload *dto.ContactLeadInput) error
}

// UserProfileRepository interface that provide access to all persistent storage operations for user profile
type UserProfileRepository interface {
	UpdateUserName(ctx context.Context, id string, userName string) error
	UpdatePrimaryPhoneNumber(ctx context.Context, id string, phoneNumber string) error
	UpdatePrimaryEmailAddress(ctx context.Context, id string, emailAddress string) error
	UpdateSecondaryPhoneNumbers(ctx context.Context, id string, phoneNumbers []string) error
	UpdateSecondaryEmailAddresses(ctx context.Context, id string, emailAddresses []string) error
	UpdateVerifiedIdentifiers(ctx context.Context, id string, identifiers []base.VerifiedIdentifier) error
	UpdateVerifiedUIDS(ctx context.Context, id string, uids []string) error
	UpdateSuspended(ctx context.Context, id string, status bool) error
	UpdatePhotoUploadID(ctx context.Context, id string, uploadID string) error
	UpdateCovers(ctx context.Context, id string, covers []base.Cover) error
	UpdatePushTokens(ctx context.Context, id string, pushToken []string) error
	UpdatePermissions(ctx context.Context, id string, perms []base.PermissionType) error
	UpdateRole(ctx context.Context, id string, role base.RoleType) error
	UpdateBioData(ctx context.Context, id string, data base.BioData) error
	UpdateAddresses(ctx context.Context, id string, address base.Address, addressType base.AddressType) error
}

// AgentRepository  defines signatures that relate to agents
type AgentRepository interface {
	CreateAgentUserProfile(ctx context.Context, phoneNumber string) (*base.UserProfile, error)
}
