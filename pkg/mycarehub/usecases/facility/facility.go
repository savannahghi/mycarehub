package facility

import (
	"context"
	"fmt"

	"github.com/savannahghi/converterandformatter"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
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
	AddFacilityToProgram(ctx context.Context, facilityIDs []string) (bool, error)
	CreateFacilities(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error)
	PublishFacilitiesToCMS(ctx context.Context, facilities []*domain.Facility) error
}

// IFacilityDelete contains the method to delete a facility
type IFacilityDelete interface {
	// TODO Ensure delete is idempotent
	DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
}

// IFacilityInactivate contains the method to activate a facility
type IFacilityInactivate interface {
	// TODO Toggle active boolean
	InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
}

// IFacilityReactivate contains the method to re-activate a facility
type IFacilityReactivate interface {
	ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
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
	RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error)
}

// IUpdateFacility contains the methods for updating a facility
type IUpdateFacility interface {
	UpdateFacility(ctx context.Context, updatePayload *dto.UpdateFacilityPayload) error
	AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error)
}

// UseCaseFacilityImpl represents facility implementation object
type UseCaseFacilityImpl struct {
	Create      infrastructure.Create
	Query       infrastructure.Query
	Delete      infrastructure.Delete
	Update      infrastructure.Update
	Pubsub      pubsubmessaging.ServicePubsub
	ExternalExt extension.ExternalMethodsExtension
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	update infrastructure.Update,
	pubsub pubsubmessaging.ServicePubsub,
	ext extension.ExternalMethodsExtension,
) UseCasesFacility {
	return &UseCaseFacilityImpl{
		Create:      create,
		Query:       query,
		Delete:      delete,
		Update:      update,
		Pubsub:      pubsub,
		ExternalExt: ext,
	}
}

// DeleteFacility deletes a facility from the database usinng the MFL Code
func (f *UseCaseFacilityImpl) DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.Delete.DeleteFacility(ctx, identifier)
}

// InactivateFacility inactivates the health facility
// TODO Toggle active boolean
func (f *UseCaseFacilityImpl) InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.Update.InactivateFacility(ctx, identifier)
}

// ReactivateFacility activates the inactivated health facility
func (f *UseCaseFacilityImpl) ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.Update.ReactivateFacility(ctx, identifier)
}

// RetrieveFacility find the health facility by ID
func (f *UseCaseFacilityImpl) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility id cannot be nil")
	}
	return f.Query.RetrieveFacility(ctx, id, isActive)
}

// SearchFacility retrieves one or more facilities from the database based on a search parameter that can be either the
// facility name or the facility identifier
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
		ID: &facilityUpdatePayload.FacilityID,
	}

	err := f.Update.UpdateFacility(ctx, facility, updatePayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// RetrieveFacilityByIdentifier find the health facility by MFL Code
func (f *UseCaseFacilityImpl) RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	if err := identifier.Validate(); err != nil {
		return nil, err
	}
	facility, err := f.Query.RetrieveFacilityByIdentifier(ctx, identifier, isActive)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return nil, err
	}
	return facility, nil
}

// ListFacilities is responsible for returning a list of paginated facilities
func (f *UseCaseFacilityImpl) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	pagination := &domain.Pagination{
		Limit:       paginationsInput.Limit,
		CurrentPage: paginationsInput.CurrentPage,
	}
	facilities, page, err := f.Query.ListFacilities(ctx, searchTerm, filterInput, pagination)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("%w", err))
		return nil, err
	}

	return &domain.FacilityPage{
		Pagination: *page,
		Facilities: facilities,
	}, nil
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
		ID: &facilityID,
	}

	err = f.Update.UpdateFacility(ctx, facility, update)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// AddFacilityToProgram is used to add a facility to a program that the currently logged in user (who should be a staff) is.
func (f *UseCaseFacilityImpl) AddFacilityToProgram(ctx context.Context, facilityIDs []string) (bool, error) {
	uid, err := f.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	userProfile, err := f.Query.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staffProfile, err := f.Query.GetStaffProfile(ctx, uid, userProfile.CurrentProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	facilities, err := f.Create.AddFacilityToProgram(ctx, staffProfile.User.CurrentProgramID, facilityIDs)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	var facilityList []string
	for _, facility := range facilities {
		facilityList = append(facilityList, *facility.ID)
	}

	programFacilityPayload := &dto.CMSLinkFacilityToProgramPayload{
		FacilityID: facilityList,
		ProgramID:  staffProfile.User.CurrentProgramID,
	}

	err = f.Pubsub.NotifyCMSAddFacilityToProgram(ctx, programFacilityPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// CreateFacilities inserts multiple facility records together with the identifiers
func (f *UseCaseFacilityImpl) CreateFacilities(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
	if len(facilities) < 1 {
		return []*domain.Facility{}, nil
	}
	return f.Create.CreateFacilities(ctx, facilities)
}

// PublishFacilitiesToCMS creates facilities in the CMS database
func (f *UseCaseFacilityImpl) PublishFacilitiesToCMS(ctx context.Context, facilities []*domain.Facility) error {
	for _, facility := range facilities {
		err := f.Pubsub.NotifyCreateCMSFacility(ctx, &dto.CreateCMSFacilityPayload{
			FacilityID: *facility.ID,
			Name:       facility.Name,
		})
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return err
		}
	}
	return nil
}
