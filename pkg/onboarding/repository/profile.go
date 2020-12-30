package repository

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// OnboardingRepository interface that provide access to all persistent storage operations
type OnboardingRepository interface {
	base.UserProfileRepository
	// creates a user profile of using the provided phone number and uid
	CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error)

	// creates an empty supplier profile
	CreateEmptySupplierProfile(ctx context.Context, profileID string) (*domain.Supplier, error)

	// creates an empty customer profile
	CreateEmptyCustomerProfile(ctx context.Context, profileID string) (*domain.Customer, error)

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

	GenerateAuthCredentials(ctx context.Context, phone string) (*domain.AuthCredentialResponse, error)

	// supplier methods
	GetSupplierProfileByProfileID(ctx context.Context, profileID string) (*domain.Supplier, error)

	AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *domain.PartnerType) (bool, error)

	// PINs
	GetPINByProfileID(
		ctx context.Context,
		ProfileID string,
	) (*domain.PIN, error)

	// Record post visit survey
	RecordPostVisitSurvey(
		ctx context.Context,
		input *domain.PostVisitSurveyInput,
		UID string,
	) (bool, error)

	// User Pin methods
	SavePIN(ctx context.Context, pin *domain.PIN) (*domain.PIN, error)
}
