package repository

import (
	"context"

	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// SupplierRepository  defines signatures that relate to suppliers
type SupplierRepository interface {
	// supplier methods
	GetSupplierProfileByID(ctx context.Context, id string) (*base.Supplier, error)

	GetSupplierProfileByProfileID(ctx context.Context, profileID string) (*base.Supplier, error)

	AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *base.PartnerType) (bool, error)

	UpdateSupplierProfile(ctx context.Context, data *base.Supplier) (*base.Supplier, error)

	StageProfileNudge(ctx context.Context, nudge map[string]interface{}) error

	StageKYCProcessingRequest(ctx context.Context, data *domain.KYCRequest) error

	// sets the active attribute of supplier profile to true
	ActivateSupplierProfile(ctx context.Context, profileID string) (*base.Supplier, error)

	FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error)

	FetchKYCProcessingRequestByID(ctx context.Context, id string) (*domain.KYCRequest, error)

	UpdateKYCProcessingRequest(ctx context.Context, sup *domain.KYCRequest) error
}

// CustomerRepository  defines signatures that relate to customers
type CustomerRepository interface {
	// customer methods
	GetCustomerProfileByID(ctx context.Context, id string) (*base.Customer, error)

	GetCustomerProfileByProfileID(ctx context.Context, profileID string) (*base.Customer, error)
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
	GetUserProfileByUID(ctx context.Context, uid string) (*base.UserProfile, error)

	// fetches a user profile by id
	GetUserProfileByID(ctx context.Context, id string) (*base.UserProfile, error)

	// fetches a user profile by phone number
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*base.UserProfile, error)

	// fetches a user profile by primary phone number
	GetUserProfileByPrimaryPhoneNumber(ctx context.Context, phoneNumber string) (*base.UserProfile, error)

	// checks if a specific phone number has already been registered to another user
	CheckIfPhoneNumberExists(ctx context.Context, phone string) (bool, error)

	// checks if a specific username has already been registered to another user
	CheckIfUsernameExists(ctx context.Context, phone string) (bool, error)

	GenerateAuthCredentialsForAnonymousUser(ctx context.Context) (*base.AuthCredentialResponse, error)

	GenerateAuthCredentials(ctx context.Context, phone string) (*base.AuthCredentialResponse, error)

	FetchAdminUsers(ctx context.Context) ([]*base.UserProfile, error)

	// removes user completely. This should be used only under testing environment
	PurgeUserByPhoneNumber(ctx context.Context, phone string) error

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
	SavePIN(ctx context.Context, pin *domain.PIN) (*domain.PIN, error)
	UpdatePIN(ctx context.Context, id string, pin *domain.PIN) (*domain.PIN, error)

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
	) (*auth.UserRecord, error)
}
