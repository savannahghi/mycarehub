package facility

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
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
	GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
}

// IFacilityUpdate contains the method to update facility details
type IFacilityUpdate interface {
	// TODO Pre-condition: ensure `id` is set and valid
	Update(facility *domain.Facility) (*domain.Facility, error)
}

// IFacilityDelete contains the method to delete a facility
type IFacilityDelete interface {
	// TODO Ensure delete is idempotent
	DeleteFacility(ctx context.Context, id string) (bool, error)
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
	// // TODO Document: callers should specify active
	// List(
	// 	// search
	// 	searchTerm *string,
	// 	// filter
	// 	filter []*domain.FilterParam,
	// 	// paginate
	// 	page int,
	// ) (*domain.FacilityPage, error)
	FetchFacilities(ctx context.Context) ([]*domain.Facility, error)
}

// IFacilityRetrieve contains the method to retrieve a facility
type IFacilityRetrieve interface {
	RetrieveFacility(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
}

// UseCaseFacilityImpl represents facility implementation object
type UseCaseFacilityImpl struct {
	Infrastructure infrastructure.Interactor
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(infra infrastructure.Interactor) *UseCaseFacilityImpl {
	return &UseCaseFacilityImpl{
		Infrastructure: infra,
	}
}

// GetOrCreateFacility creates a new facility
func (f *UseCaseFacilityImpl) GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	fetchedFacility, err := f.Infrastructure.RetrieveFacilityByMFLCode(ctx, facility.Code, facility.Active)
	if err != nil {
		if strings.Contains(err.Error(), "failed query and retrieve facility by MFLCode") {
			return f.Infrastructure.GetOrCreateFacility(ctx, facility)
		}
		return nil, fmt.Errorf("failed to retrieve facility")
	}
	return fetchedFacility, nil
}

// Update creates a new facility
func (f *UseCaseFacilityImpl) Update(facility *domain.Facility) (*domain.Facility, error) {
	return nil, nil
}

// DeleteFacility deletes a facility from the database usinng the MFL Code
func (f *UseCaseFacilityImpl) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return f.Infrastructure.DeleteFacility(ctx, id)
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

// // List returns a list if health facility
// // TODO Document: callers should specify active
// func (f *UseCaseFacilityImpl) List(
// 	pagination *firebasetools.PaginationInput,
// 	filter []*dto.FacilityFilterInput,
// 	sort []*dto.FacilitySortInput,
// ) (*dto.FacilityConnection, error) {
// 	return nil, nil
// }

// RetrieveFacility find the health facility by ID
func (f *UseCaseFacilityImpl) RetrieveFacility(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error) {
	return f.Infrastructure.RetrieveFacility(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode find the health facility by MFL Code
func (f *UseCaseFacilityImpl) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return f.Infrastructure.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
}

// FetchFacilities fetches healthcare facilities in platform
func (f *UseCaseFacilityImpl) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.Infrastructure.GetFacilities(ctx)
}
