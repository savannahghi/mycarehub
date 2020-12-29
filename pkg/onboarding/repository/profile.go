package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// OnboardingRepository interface that provide access to all persistent storage operations
type OnboardingRepository interface {
	base.UserProfileRepository
	// creates a user profile of using the provided phone number and uid
	CreateUserProfile(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error)

	GetUserProfile(ctx context.Context, uid string) (*base.UserProfile, error)
	GetUserProfileByID(ctx context.Context, id string) (*base.UserProfile, error)
	CheckIfPhoneNumberExists(ctx context.Context, phone string) (bool, error)

	// supplier methods
	GetSupplierProfileByProfileID(ctx context.Context, profileID string) (*domain.Supplier, *firestore.DocumentSnapshot, error)
	AddPartnerType(ctx context.Context, profileID string, name *string, partnerType *domain.PartnerType) (bool, error)
}
