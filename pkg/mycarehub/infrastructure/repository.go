package infrastructure

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id string) (bool, error)
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error)
	GetCurrentTerms(ctx context.Context, flavour enums.Flavour) (string, error)
}

// Update represents all the update action interfaces
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *string) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *string) (bool, error)
}
