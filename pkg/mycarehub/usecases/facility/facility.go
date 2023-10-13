package facility

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/converterandformatter"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/healthcrm"
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
	AddFacilityToProgram(ctx context.Context, facilityIDs []string, programID string) (bool, error)
	CreateFacilities(ctx context.Context, facilitiesInput []*dto.FacilityInput) ([]*domain.Facility, error)
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
	ListProgramFacilities(ctx context.Context, programID *string, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	SyncFacilities(ctx context.Context) error
}

// IFacilityRetrieve contains the method to retrieve a facility
type IFacilityRetrieve interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error)
}

// IUpdateFacility contains the methods for updating a facility
type IUpdateFacility interface {
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
	HealthCRM   healthcrm.IHealthCRMService
}

// NewFacilityUsecase returns a new facility service
func NewFacilityUsecase(
	create infrastructure.Create,
	query infrastructure.Query,
	delete infrastructure.Delete,
	update infrastructure.Update,
	pubsub pubsubmessaging.ServicePubsub,
	ext extension.ExternalMethodsExtension,
	healthcrmSvc healthcrm.IHealthCRMService,
) UseCasesFacility {
	return &UseCaseFacilityImpl{
		Create:      create,
		Query:       query,
		Delete:      delete,
		Update:      update,
		Pubsub:      pubsub,
		ExternalExt: ext,
		HealthCRM:   healthcrmSvc,
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
	output, err := f.Query.RetrieveFacility(ctx, id, isActive)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, identifier := range output.Identifiers {
		if identifier.Type == enums.FacilityIdentifierTypeHealthCRM {
			// Get facility services
			result, err := f.HealthCRM.GetServicesOfferedInAFacility(ctx, identifier.Value)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, err
			}

			output.Services = result.Results

			// Get facility Business hours
			facilityObj, err := f.HealthCRM.GetCRMFacilityByID(ctx, identifier.Value)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return nil, err
			}

			output.BusinessHours = facilityObj.BusinessHours
		}
	}

	return output, nil
}

// ListFacilities retrieves one or more facilities from the database based on a search parameter that can be either the
// facility name or the facility identifier
func (f *UseCaseFacilityImpl) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	pagination := &domain.Pagination{
		Limit:       paginationsInput.Limit,
		CurrentPage: paginationsInput.CurrentPage,
	}

	facilities, page, err := f.Query.ListFacilities(ctx, searchTerm, filterInput, pagination)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("failed to list facilities: %w", err))
		return nil, fmt.Errorf("failed to list facilities: %w", err)
	}

	return &domain.FacilityPage{
		Pagination: *page,
		Facilities: facilities,
	}, nil
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

// ListProgramFacilities is responsible for returning a list of paginated facilities
func (f *UseCaseFacilityImpl) ListProgramFacilities(ctx context.Context, programID *string, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	if programID == nil {
		loggedInUserID, err := f.ExternalExt.GetLoggedInUserUID(ctx)
		if err != nil {
			helpers.ReportErrorToSentry(fmt.Errorf("failed to get logged in user: %w", err))
			return nil, fmt.Errorf("failed to get logged in user: %w", err)
		}

		userProfile, err := f.Query.GetUserProfileByUserID(ctx, loggedInUserID)
		if err != nil {
			helpers.ReportErrorToSentry(fmt.Errorf("failed to get user profile: %w", err))
			return nil, fmt.Errorf("failed to get user profile: %w", err)
		}
		programID = &userProfile.CurrentProgramID
	}

	pagination := &domain.Pagination{
		Limit:       paginationsInput.Limit,
		CurrentPage: paginationsInput.CurrentPage,
	}

	facilities, page, err := f.Query.ListProgramFacilities(ctx, programID, searchTerm, filterInput, pagination)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("failed to list facilities: %w", err))
		return nil, fmt.Errorf("failed to list facilities: %w", err)
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

// CreateFacilities inserts multiple facility records together with the identifiers
func (f *UseCaseFacilityImpl) CreateFacilities(ctx context.Context, facilitiesInput []*dto.FacilityInput) ([]*domain.Facility, error) {
	if len(facilitiesInput) < 1 {
		helpers.ReportErrorToSentry(fmt.Errorf("empty facility details in input"))
		return nil, fmt.Errorf("empty facility details in input")
	}

	facilities := []*domain.Facility{}

	for _, facility := range facilitiesInput {
		lat, err := strconv.ParseFloat(facility.Coordinates.Lat, 64)
		if err != nil {
			return nil, err
		}

		lng, err := strconv.ParseFloat(facility.Coordinates.Lng, 64)
		if err != nil {
			return nil, err
		}

		var services []domain.FacilityService
		for _, service := range facility.Services {
			var serviceIdentifiers []domain.ServiceIdentifier
			for _, serviceIdentifier := range service.Identifiers {
				serviceIdentifiers = append(serviceIdentifiers, domain.ServiceIdentifier{
					IdentifierType:  serviceIdentifier.IdentifierType.String(),
					IdentifierValue: serviceIdentifier.IdentifierValue,
				})
			}

			services = append(services, domain.FacilityService{
				Name:        service.Name,
				Description: service.Description,
				Identifiers: serviceIdentifiers,
			})
		}

		var businessHours []domain.BusinessHours
		for _, businessHour := range facility.BusinessHours {
			businessHours = append(businessHours, domain.BusinessHours{
				Day:         businessHour.Day.String(),
				OpeningTime: businessHour.OpeningTime,
				ClosingTime: businessHour.ClosingTime,
			})
		}

		facilities = append(facilities, &domain.Facility{
			Name:               facility.Name,
			Phone:              facility.Phone,
			Active:             facility.Active,
			Country:            facility.Country,
			County:             facility.County,
			Address:            facility.Address,
			Description:        facility.Description,
			FHIROrganisationID: facility.FHIROrganisationID,
			Identifiers: []*domain.FacilityIdentifier{
				{
					Active: true,
					Type:   facility.Identifier.Type,
					Value:  facility.Identifier.Value,
				},
			},
			Coordinates: &domain.Coordinates{
				Lat: lat,
				Lng: lng,
			},
			Services:      services,
			BusinessHours: businessHours,
		})
	}

	// Create facility in the CRM first
	results, err := f.HealthCRM.CreateFacility(ctx, facilities)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	facilitiesObj, err := f.Create.CreateFacilities(ctx, results)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("failed to create facilities: %w", err))
		return nil, fmt.Errorf("failed to create facilities: %w", err)
	}

	for _, facility := range facilitiesObj {
		err = f.Pubsub.NotifyCreateOrganization(ctx, facility)
		if err != nil {
			helpers.ReportErrorToSentry(fmt.Errorf("failed to create publish organisation to clinical service: %w", err))
			return nil, fmt.Errorf("failed to create publish organisation to clinical service: %w", err)
		}
		err = f.Pubsub.NotifyCreateCMSFacility(ctx, &dto.CreateCMSFacilityPayload{FacilityID: *facility.ID, Name: facility.Name})
		if err != nil {
			helpers.ReportErrorToSentry(fmt.Errorf("failed to create facility in cms: %w", err))
			return nil, fmt.Errorf("failed to create facility in cms: %w", err)
		}
	}

	return facilitiesObj, nil
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

// AddFacilityToProgram is used to add a facility to a program
func (f *UseCaseFacilityImpl) AddFacilityToProgram(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
	facilities, err := f.Create.AddFacilityToProgram(ctx, programID, facilityIDs)
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
		ProgramID:  programID,
	}

	err = f.Pubsub.NotifyCMSAddFacilityToProgram(ctx, programFacilityPayload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}
