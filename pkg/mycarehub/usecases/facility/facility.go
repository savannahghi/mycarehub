package facility

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/sirupsen/logrus"
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
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}

// IFacilityUpdate contains the method to update facility details
type IFacilityUpdate interface {
	// TODO Pre-condition: ensure `id` is set and valid
	//Update(facility *domain.Facility) (*domain.Facility, error)
}

// IFacilityDelete contains the method to delete a facility
type IFacilityDelete interface {
	// TODO Ensure delete is idempotent
	DeleteFacility(ctx context.Context, id string) (bool, error)
}

// IFacilityInactivate contains the method to activate a facility
type IFacilityInactivate interface {
	// TODO Toggle active boolean
	//Inactivate(id string) (*domain.Facility, error)
}

// IFacilityReactivate contains the method to re-activate a facility
type IFacilityReactivate interface {
	//Reactivate(id string) (*domain.Facility, error)
}

// IFacilityList contains the method to list of facilities
type IFacilityList interface {
	// TODO Document: callers should specify active
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	FetchFacilities(ctx context.Context) ([]*domain.Facility, error)
}

// IFacilityRetrieve contains the method to retrieve a facility
type IFacilityRetrieve interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
}

// UseCaseFacilityImpl represents facility implementation object
type UseCaseFacilityImpl struct {
	Create infrastructure.Create
	Query  infrastructure.Query
	Delete infrastructure.Delete
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(create infrastructure.Create, query infrastructure.Query, delete infrastructure.Delete) *UseCaseFacilityImpl {
	return &UseCaseFacilityImpl{
		Create: create,
		Query:  query,
		Delete: delete,
	}
}

// GetOrCreateFacility creates a new facility
func (f *UseCaseFacilityImpl) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	logrus.Print("Facility>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	if facility.Code == "" {
		return nil, fmt.Errorf("facililty code cannot be nil")
	}
	fetchedFacility, err := f.RetrieveFacilityByMFLCode(ctx, facility.Code, facility.Active)
	if err != nil {
		if strings.Contains(err.Error(), "failed query and retrieve facility by MFLCode") {
			return f.Create.GetOrCreateFacility(ctx, facility)
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
	return f.Delete.DeleteFacility(ctx, id)
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
// 	pagination *firebasetools.PaginationsInput,
// 	filter []*dto.FacilityFilterInput,
// 	sort []*dto.FacilitySortInput,
// ) (*dto.FacilityConnection, error) {
// 	return nil, nil
// }

// RetrieveFacility find the health facility by ID
func (f *UseCaseFacilityImpl) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility id cannot be nil")
	}
	return f.Query.RetrieveFacility(ctx, id, isActive)
}

// FetchFacilities fetches healthcare facilities in platform
func (f *UseCaseFacilityImpl) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.Query.GetFacilities(ctx)
}

// RetrieveFacilityByMFLCode find the health facility by MFL Code
func (f *UseCaseFacilityImpl) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	if MFLCode == "" {
		return nil, fmt.Errorf("facility MFL code cannot be empty")
	}
	return f.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
}

//ListFacilities is responsible for returning a list of paginated facilities
func (f *UseCaseFacilityImpl) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	if searchTerm == nil {
		return nil, fmt.Errorf("search term cannot be nil")
	}

	if filterInput == nil {
		return nil, fmt.Errorf("filter input cannot be nil")
	}

	if paginationsInput == nil {
		return nil, fmt.Errorf("filter input cannot be nil")
	}

	return f.Query.ListFacilities(ctx, searchTerm, filterInput, paginationsInput)
}
