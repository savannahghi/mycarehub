package repository

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"

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
	base.UserProfileRepository

	SupplierRepository

	CustomerRepository
	// creates a user profile of using the provided phone number and uid
	CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error)

	// creates an empty supplier profile
	CreateEmptySupplierProfile(ctx context.Context, profileID string) (*base.Supplier, error)

	// creates an empty customer profile
	CreateEmptyCustomerProfile(ctx context.Context, profileID string) (*base.Customer, error)

	// fetches a user profile by uid
	GetUserProfileByUID(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error)

	// fetches a user profile by id. returns the unsuspend profile
	GetUserProfileByID(ctx context.Context, id string, suspended bool) (*base.UserProfile, error)

	// fetches a user profile by phone number
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error)

	// fetches a user profile by primary phone number
	GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phoneNumber string, suspended bool) (*base.UserProfile, error)

	// checks if a specific phone number has already been registered to another user
	CheckIfPhoneNumberExists(ctx context.Context, phone string) (bool, error)

	// checks if a specific email has already been registered to another user
	CheckIfEmailExists(ctx context.Context, phone string) (bool, error)

	// checks if a specific username has already been registered to another user
	CheckIfUsernameExists(ctx context.Context, phone string) (bool, error)

	GenerateAuthCredentialsForAnonymousUser(ctx context.Context) (*base.AuthCredentialResponse, error)

	GenerateAuthCredentials(ctx context.Context, phone string) (*base.AuthCredentialResponse, error)

	FetchAdminUsers(ctx context.Context) ([]*base.UserProfile, error)

	// removes user completely. This should be used only under testing environment
	PurgeUserByPhoneNumber(ctx context.Context, phone string) error

	HardResetSecondaryPhoneNumbers(ctx context.Context, id string, newSecondaryPhones []string) error

	HardResetSecondaryEmailAddress(ctx context.Context, id string, newSecondaryEmails []string) error

	// PINs
	GetPINByProfileID(
		ctx context.Context,
		ProfileID string,
	) (*domain.PIN, error)

	// Record post visit survey
	RecordPostVisitSurvey(
		ctx context.Context,
		input resources.PostVisitSurveyInput,
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

	GetOrCreatePhoneNumberUser(
		ctx context.Context,
		phone string,
	) (*resources.CreatedUserResponse, error)

	AddUserAsExperimentParticipant(ctx context.Context, profile *base.UserProfile) (bool, error)

	RemoveUserAsExperimentParticipant(ctx context.Context, profile *base.UserProfile) (bool, error)

	CheckIfExperimentParticipant(ctx context.Context, profileID string) (bool, error)

	AddNHIFDetails(
		ctx context.Context,
		input resources.NHIFDetailsInput,
		profileID string,
	) (*domain.NHIFDetails, error)

	GetNHIFDetailsByProfileID(
		ctx context.Context,
		profileID string,
	) (*domain.NHIFDetails, error)

	GetUserCommunicationsSettings(ctx context.Context, profileID string) (*base.UserCommunicationsSetting, error)

	SetUserCommunicationsSettings(ctx context.Context, profileID string,
		allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error)
}
