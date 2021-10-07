package facility

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

type FacilityUseCases interface {
	IFacilityList
	IFacilityRetrieve
	IFacilityCreate
	IFacilityUpdate
	IFacilityDelete
	IFacilityInactivate
	IFacilityReactivate
}

// IFacilityCreate ...
type IFacilityCreate interface {
	// TODO Ensure blank ID when creating
	// TODO Since `id` is optional, ensure pre-condition check
	Create(facility *domain.Facility) (*domain.Facility, error)
}

// IFacilityUpdate ...
type IFacilityUpdate interface {
	// TODO Pre-condition: ensure `id` is set and valid
	Update(facility *domain.Facility) (*domain.Facility, error)
}

// IFacilityDelete ...
type IFacilityDelete interface {
	// TODO Ensure delete is idempotent
	Delete(id string) (bool, error)
}

// IFacilityInactivate ...
type IFacilityInactivate interface {
	// TODO Toggle active boolean
	Inactivate(id string) (*domain.Facility, error)
}

// IFacilityReactivate ...
type IFacilityReactivate interface {
	Reactivate(id string) (*domain.Facility, error)
}

// IFacilityList ...
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

// IFacilityRetrieve ...
type IFacilityRetrieve interface {
	Retrieve(id string) (*domain.Facility, error)
}

// FacilityUseCaseImpl represents facility implementation object
type FacilityUseCaseImpl struct {
	Infrastructure infrastructure.Interactor
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(infra infrastructure.Interactor) FacilityUseCases {
	return &FacilityUseCaseImpl{
		Infrastructure: infra,
	}
}

// // Create creates a new facility
func (f *FacilityUseCaseImpl) Create(facility *domain.Facility) (*domain.Facility, error) {
	return nil, nil
}

// Update creates a new facility
func (f *FacilityUseCaseImpl) Update(facility *domain.Facility) (*domain.Facility, error) {
	return nil, nil
}

// Delete creates a new facility
func (f *FacilityUseCaseImpl) Delete(id string) (bool, error) {
	return false, nil
}

// Inactivate ...
// TODO Toggle active boolean
func (f *FacilityUseCaseImpl) Inactivate(id string) (*domain.Facility, error) {
	return nil, nil
}

// Reactivate ...
func (f *FacilityUseCaseImpl) Reactivate(id string) (*domain.Facility, error) {
	return nil, nil
}

// List ...
// TODO Document: callers should specify active
func (f *FacilityUseCaseImpl) List(
	// search
	searchTerm *string,
	// filter
	filter []*domain.FilterParam,
	// paginate
	page int,
) (*domain.FacilityPage, error) {
	return nil, nil
}

// Retrieve ...
func (f *FacilityUseCaseImpl) Retrieve(id string) (*domain.Facility, error) {
	return nil, nil
}
