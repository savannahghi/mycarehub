package facility

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

// UseCasesFacility ...
type UseCasesFacility interface {
	IFacilityList
	IFacilityRetrieve
	IFacilityCreate
	IFacilityUpdate
	IFacilityDelete
	IFacilityInactivate
	IFacilityReactivate
}

// IFacilityCreate contains the method used to create a facility
type IFacilityCreate interface {
	// TODO Ensure blank ID when creating
	// TODO Since `id` is optional, ensure pre-condition check
	CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}

// IFacilityUpdate contains the method to update facility details
type IFacilityUpdate interface {
	// TODO Pre-condition: ensure `id` is set and valid
	Update(facility *domain.Facility) (*domain.Facility, error)
}

// IFacilityDelete contains the method to delete a facility
type IFacilityDelete interface {
	// TODO Ensure delete is idempotent
	Delete(id string) (bool, error)
}

// IFacilityInactivate contains the method to activate a facility
type IFacilityInactivate interface {
	// TODO Toggle active boolean
	Inactivate(id string) (*domain.Facility, error)
}

// IFacilityReactivate contains the method to re-activate a facility
type IFacilityReactivate interface {
	Reactivate(id string) (*domain.Facility, error)
}

// IFacilityList contains the method to list of facilities
type IFacilityList interface {
	// TODO Document: callers should specify active
	List(
		// search
		searchTerm *string,
		// filter
		filter []*domain.FilterParam,
		// paginate
		page int,
	) (*domain.FacilityPage, error)
}

// IFacilityRetrieve contains the method to retrieve a facility
type IFacilityRetrieve interface {
	Retrieve(id string) (*domain.Facility, error)
}

// UseCaseFacilityImpl represents facility implementation object
type UseCaseFacilityImpl struct {
	Infrastructure infrastructure.Infrastructure
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(infra infrastructure.Infrastructure) UseCasesFacility {
	return &UseCaseFacilityImpl{
		Infrastructure: infra,
	}
}

// CreateFacility creates a new facility
func (f *UseCaseFacilityImpl) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return f.Infrastructure.Database.CreateFacility(ctx, facility)
}

// Update creates a new facility
func (f *UseCaseFacilityImpl) Update(facility *domain.Facility) (*domain.Facility, error) {
	return nil, nil
}

// Delete creates a new facility
func (f *UseCaseFacilityImpl) Delete(id string) (bool, error) {
	return false, nil
}

// Inactivate inactivates the health facility
// TODO Toggle active boolean
func (f *UseCaseFacilityImpl) Inactivate(id string) (*domain.Facility, error) {
	return nil, nil
}

// Reactivate activates the inactivated health facility
func (f *UseCaseFacilityImpl) Reactivate(id string) (*domain.Facility, error) {
	return nil, nil
}

// List returns a list if health facility
// TODO Document: callers should specify active
func (f *UseCaseFacilityImpl) List(
	// search
	searchTerm *string,
	// filter
	filter []*domain.FilterParam,
	// paginate
	page int,
) (*domain.FacilityPage, error) {
	return nil, nil
}

// Retrieve find the health facility by ID
func (f *UseCaseFacilityImpl) Retrieve(id string) (*domain.Facility, error) {
	return nil, nil
}
