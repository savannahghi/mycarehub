package infrastructure

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pg "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	RegisterClient(
		ctx context.Context,
		userInput *dto.UserInput,
		clientInput *dto.ClientProfileInput,
	) (*domain.ClientUserProfile, error)
	SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error)
}

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id string) (bool, error)
}

// ServiceCreateImpl represents create contract implementation object
type ServiceCreateImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceCreateImpl returns new instance of ServiceCreateImpl
func NewServiceCreateImpl(on pg.OnboardingDb) Create {
	return &ServiceCreateImpl{
		onboarding: on,
	}
}

// GetOrCreateFacility is responsible for creating a representation of a facility
func (f ServiceCreateImpl) GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.onboarding.GetOrCreateFacility(ctx, &facility)
}

// RegisterClient creates a client user and saves the details in the database
func (f ServiceCreateImpl) RegisterClient(
	ctx context.Context,
	userInput *dto.UserInput,
	clientInput *dto.ClientProfileInput,
) (*domain.ClientUserProfile, error) {
	return f.onboarding.RegisterClient(ctx, userInput, clientInput)
}

// SavePin saves a user's pin in the database
func (f ServiceCreateImpl) SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error) {
	return f.onboarding.SavePin(ctx, pinInput)
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput) (*domain.FacilityPage, error)
}

// ServiceQueryImpl contains implementation for the Query interface
type ServiceQueryImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceQueryImpl is the initializer for Service query
func NewServiceQueryImpl(on pg.OnboardingDb) *ServiceQueryImpl {
	return &ServiceQueryImpl{
		onboarding: on,
	}
}

// RetrieveFacility  is a repository implementation method for RetrieveFacility
func (q ServiceQueryImpl) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return q.onboarding.RetrieveFacility(ctx, id, isActive)
}

// GetUserProfileByPhoneNumber returns a user profile based on the provided phonenumber
func (q ServiceQueryImpl) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return q.onboarding.GetUserProfileByPhoneNumber(ctx, phoneNumber)
}

// RetrieveFacilityByMFLCode  is a repository implementation method for RetrieveFacilityByMFLCode
func (q ServiceQueryImpl) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return q.onboarding.RetrieveByFacilityMFLCode(ctx, MFLCode, isActive)
}

//GetFacilities is responsible for returning a slice of healthcare facilities in the platform.
func (q ServiceQueryImpl) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return q.onboarding.GetFacilities(ctx)
}

//ListFacilities is responsible for returning a list of paginated facilities
func (q ServiceQueryImpl) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, PaginationsInput dto.PaginationsInput,
) (*domain.FacilityPage, error) {
	return q.onboarding.ListFacilities(ctx, searchTerm, filterInput, PaginationsInput)
}

// DeleteFacility is responsible for deletion of a facility from the database using the facility's id
func (f ServiceDeleteImpl) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return f.onboarding.DeleteFacility(ctx, id)
}

// ServiceDeleteImpl represents delete facility implementation object
type ServiceDeleteImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceDeleteImpl returns new instance of NewServiceDeleteImpl
func NewServiceDeleteImpl(on pg.OnboardingDb) Delete {
	return &ServiceDeleteImpl{
		onboarding: on,
	}
}
