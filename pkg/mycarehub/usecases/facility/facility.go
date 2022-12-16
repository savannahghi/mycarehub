package facility

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/converterandformatter"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
)

// UseCasesFacility ...
type UseCasesFacility interface {
	IFacilityList
	IFacilityRetrieve
	IFacilityCreate
	IFacilityDelete
	IFacilityInactivate
	IFacilityReactivate
	IUpdateFacility
}

// IFacilityCreate contains the method used to create a facility
type IFacilityCreate interface {
	// TODO Ensure blank ID when creating
	// TODO Since `id` is optional, ensure pre-condition check
	GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}

// IFacilityDelete contains the method to delete a facility
type IFacilityDelete interface {
	// TODO Ensure delete is idempotent
	DeleteFacility(ctx context.Context, id int) (bool, error)
}

// IFacilityInactivate contains the method to activate a facility
type IFacilityInactivate interface {
	// TODO Toggle active boolean
	InactivateFacility(ctx context.Context, mflCode *int) (bool, error)
}

// IFacilityReactivate contains the method to re-activate a facility
type IFacilityReactivate interface {
	ReactivateFacility(ctx context.Context, mflCode *int) (bool, error)
}

// IFacilityList contains the method to list of facilities
type IFacilityList interface {
	// TODO Document: callers should specify active
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error)
	SyncFacilities(ctx context.Context) error
}

// IFacilityRetrieve contains the method to retrieve a facility
type IFacilityRetrieve interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error)
}

// IUpdateFacility contains the methods for updating a facility
type IUpdateFacility interface {
	UpdateFacility(ctx context.Context, updatePayload *dto.UpdateFacilityPayload) error
	AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error)
}

// UseCaseFacilityImpl represents facility implementation object
type UseCaseFacilityImpl struct {
	Create infrastructure.Create
	Query  infrastructure.Query
	Delete infrastructure.Delete
	Update infrastructure.Update
	Pubsub pubsubmessaging.ServicePubsub
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	update infrastructure.Update,
	pubsub pubsubmessaging.ServicePubsub,
) *UseCaseFacilityImpl {
	return &UseCaseFacilityImpl{
		Create: create,
		Query:  query,
		Delete: delete,
		Update: update,
		Pubsub: pubsub,
	}
}

// GetOrCreateFacility creates a new facility
func (f *UseCaseFacilityImpl) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	fetchedFacility, err := f.RetrieveFacilityByMFLCode(ctx, facility.Code, facility.Active)
	if err != nil {
		if strings.Contains(err.Error(), "failed query and retrieve facility by MFLCode") {
			return f.Create.GetOrCreateFacility(ctx, facility)
		}
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to retrieve facility")
	}
	return fetchedFacility, nil
}

// DeleteFacility deletes a facility from the database usinng the MFL Code
func (f *UseCaseFacilityImpl) DeleteFacility(ctx context.Context, id int) (bool, error) {
	return f.Delete.DeleteFacility(ctx, id)
}

// InactivateFacility inactivates the health facility
// TODO Toggle active boolean
func (f *UseCaseFacilityImpl) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return f.Update.InactivateFacility(ctx, mflCode)
}

// ReactivateFacility activates the inactivated health facility
func (f *UseCaseFacilityImpl) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return f.Update.ReactivateFacility(ctx, mflCode)
}

// RetrieveFacility find the health facility by ID
func (f *UseCaseFacilityImpl) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility id cannot be nil")
	}
	return f.Query.RetrieveFacility(ctx, id, isActive)
}

// SearchFacility retrieves one or more facilities from the database based on a search parameter that can be either the
// facility name or the facility mflcode
func (f *UseCaseFacilityImpl) SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
	return f.Query.SearchFacility(ctx, searchParameter)
}

// SyncFacilities gets a list of facilities without a fhir organisation ID from the database
// and pusblishes them to create organisation pubsub topic
func (f *UseCaseFacilityImpl) SyncFacilities(ctx context.Context) error {

	response, err := f.Query.GetFacilitiesWithoutFHIRID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	for _, facility := range response {
		err = f.Pubsub.NotifyCreateOrganization(ctx, facility)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return err
		}
	}

	return nil
}

// UpdateFacility updates the details of a facility or set of facilities specified
func (f *UseCaseFacilityImpl) UpdateFacility(ctx context.Context, facilityUpdatePayload *dto.UpdateFacilityPayload) error {
	updatePayload := map[string]interface{}{
		"fhir_organization_id": facilityUpdatePayload.FHIROrganisationID,
	}

	facility := &domain.Facility{
		ID: facilityUpdatePayload.FacilityID,
	}

	err := f.Update.UpdateFacility(ctx, facility, updatePayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// RetrieveFacilityByMFLCode find the health facility by MFL Code
func (f *UseCaseFacilityImpl) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
	if MFLCode == 0 {
		return nil, fmt.Errorf("facility MFL code cannot be empty")
	}
	return f.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
}

// ListFacilities is responsible for returning a list of paginated facilities
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

// AddFacilityContact adds/updates a facilities contact
func (f *UseCaseFacilityImpl) AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error) {
	phoneNumber, err := converterandformatter.NormalizeMSISDN(contact)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.NormalizeMSISDNError(err)
	}

	update := map[string]interface{}{
		"phone": *phoneNumber,
	}

	facility := &domain.Facility{
		ID: facilityID,
	}

	err = f.Update.UpdateFacility(ctx, facility, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}
