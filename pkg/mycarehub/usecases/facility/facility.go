package facility

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/healthcrm"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
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
	IFacilityRegistry
}

// IFacilityCreate contains the method used to create a facility
type IFacilityCreate interface {
	AddFacilityToProgram(ctx context.Context, facilityIDs []string, programID string) (bool, error)
	CreateFacilities(ctx context.Context, facilitiesInput []*dto.FacilityInput) ([]*domain.Facility, error)
	PublishFacilitiesToCMS(ctx context.Context, facilities []*domain.Facility) error
}

// IFacilityDelete contains the method to delete a facility
type IFacilityDelete interface {
	DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
}

// IFacilityInactivate contains the method to activate a facility
type IFacilityInactivate interface {
	InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
}

// IFacilityReactivate contains the method to re-activate a facility
type IFacilityReactivate interface {
	ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
}

// IFacilityList contains the method to list of facilities
type IFacilityList interface {
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

// IFacilityRegistry contains the methods that perform action related to health crm
type IFacilityRegistry interface {
	GetServices(ctx context.Context, pagination *dto.PaginationsInput) (*dto.FacilityServiceOutputPage, error)
	GetNearbyFacilities(ctx context.Context, locationInput *dto.LocationInput, serviceIDs []string, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error)
	SearchFacilitiesByService(ctx context.Context, locationInput *dto.LocationInput, serviceName string, pagination *dto.PaginationsInput) (*domain.FacilityPage, error)
	BookService(ctx context.Context, facilityID string, serviceIDs []string, serviceBookingTime time.Time) (*dto.BookingOutput, error)
	ListBookings(ctx context.Context, clientID string, bookingState enums.BookingState, pagination dto.PaginationsInput) (*dto.BookingPage, error)
	VerifyBookingCode(ctx context.Context, booking string, code string, programID string) (bool, error)
}

// UseCaseFacilityImpl represents facility implementation object
type UseCaseFacilityImpl struct {
	Create         infrastructure.Create
	Query          infrastructure.Query
	Delete         infrastructure.Delete
	Update         infrastructure.Update
	Pubsub         pubsubmessaging.ServicePubsub
	ExternalExt    extension.ExternalMethodsExtension
	HealthCRM      healthcrm.IHealthCRMService
	ServiceRequest servicerequest.UseCaseServiceRequest
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
	servicerequest servicerequest.UseCaseServiceRequest,
) UseCasesFacility {
	return &UseCaseFacilityImpl{
		Create:         create,
		Query:          query,
		Delete:         delete,
		Update:         update,
		Pubsub:         pubsub,
		ExternalExt:    ext,
		HealthCRM:      healthcrmSvc,
		ServiceRequest: servicerequest,
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
	output, err := f.HealthCRM.GetCRMFacilityByID(ctx, *id)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
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

	for _, facility := range facilities {
		for _, identifier := range facility.Identifiers {
			if identifier.Type == enums.FacilityIdentifierTypeHealthCRM {
				facilityObj, err := f.HealthCRM.GetCRMFacilityByID(ctx, identifier.Value)
				if err != nil {
					helpers.ReportErrorToSentry(err)
					return nil, err
				}

				facility.Services = facilityObj.Services
				facility.BusinessHours = facilityObj.BusinessHours
			}
		}
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
			Country:            facility.Country.String(),
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

// GetNearbyFacilities is used to show facility(ies) near my current location
func (f *UseCaseFacilityImpl) GetNearbyFacilities(ctx context.Context, locationInput *dto.LocationInput, serviceIDs []string, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	pagination := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	facilities, err := f.HealthCRM.GetFacilities(ctx, locationInput, serviceIDs, "", pagination)
	if err != nil {
		return nil, err
	}

	return &domain.FacilityPage{
		Pagination: *pagination,
		Facilities: facilities,
	}, nil
}

// GetServices is used to fetch all the services from health crm
func (f *UseCaseFacilityImpl) GetServices(ctx context.Context, pagination *dto.PaginationsInput) (*dto.FacilityServiceOutputPage, error) {
	page := &domain.Pagination{
		Limit:       pagination.Limit,
		CurrentPage: pagination.CurrentPage,
	}

	output, err := f.HealthCRM.GetServices(ctx, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return &dto.FacilityServiceOutputPage{
		Results: output.Results,
		Pagination: domain.Pagination{
			Limit:       output.PageSize,
			CurrentPage: output.CurrentPage,
			Count:       int64(output.Count),
			TotalPages:  output.TotalPages,
		},
	}, nil
}

// SearchFacilitiesByService is used to search for facilities offering a specific service by using the service name
// as the search parameter. If the location is provided, the response returned will order the facilities by the proximity
// to the user
func (f *UseCaseFacilityImpl) SearchFacilitiesByService(ctx context.Context, locationInput *dto.LocationInput, serviceName string, pagination *dto.PaginationsInput) (*domain.FacilityPage, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("missing required parameter: 'service name' is not provided")
	}

	page := &domain.Pagination{
		Limit:       pagination.Limit,
		CurrentPage: pagination.CurrentPage,
	}

	facilities, err := f.HealthCRM.GetFacilities(ctx, locationInput, []string{}, serviceName, page)
	if err != nil {
		return nil, err
	}

	return &domain.FacilityPage{
		Pagination: *page,
		Facilities: facilities,
	}, nil
}

// BookService is used to book for a service(s) in a facility
func (f *UseCaseFacilityImpl) BookService(ctx context.Context, facilityID string, serviceIDs []string, serviceBookingTime time.Time) (*dto.BookingOutput, error) {
	if len(serviceIDs) < 1 {
		return nil, fmt.Errorf("a booking should contain at least one service")
	}

	loggedInUser, err := f.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	userProfile, err := f.Query.GetUserProfileByUserID(ctx, loggedInUser)
	if err != nil {
		return nil, err
	}

	clientProfile, err := f.Query.GetClientProfile(ctx, *userProfile.ID, userProfile.CurrentProgramID)
	if err != nil {
		return nil, err
	}

	verificationCode, err := utils.GenerateTempPIN(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	// TODO: Check that the service exists ih health crm

	booking := &domain.Booking{
		Services: serviceIDs,
		Date:     serviceBookingTime,
		Facility: domain.Facility{
			ID: &facilityID,
		},
		Client: domain.ClientProfile{
			ID: clientProfile.ID,
		},
		OrganisationID:         clientProfile.OrganisationID,
		ProgramID:              clientProfile.ProgramID,
		VerificationCode:       verificationCode,
		VerificationCodeStatus: enums.UnVerified,
		BookingStatus:          enums.Pending,
	}

	result, err := f.Create.CreateBooking(ctx, booking)
	if err != nil {
		return nil, err
	}

	serviceRequestInput := &dto.ServiceRequestInput{
		ClientID:    *result.Client.ID,
		Flavour:     feedlib.FlavourConsumer,
		RequestType: enums.ServiceRequestBooking.String(),
		Request:     fmt.Sprintf("A new booking has been made by %s.", clientProfile.User.Name),
		FacilityID:  *clientProfile.DefaultFacility.ID,
		ClientName:  &clientProfile.User.Name,
		Meta: map[string]interface{}{
			"serviceIDs": serviceIDs,
			"bookingID":  result.ID,
		},
		ProgramID:      clientProfile.User.CurrentProgramID,
		OrganisationID: clientProfile.User.CurrentOrganizationID,
	}

	_, err = f.ServiceRequest.CreateServiceRequest(
		ctx,
		serviceRequestInput,
	)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	facility, err := f.HealthCRM.GetCRMFacilityByID(ctx, *result.Facility.ID)
	if err != nil {
		return nil, err
	}

	result.Facility = *facility

	var services []domain.FacilityService

	for _, service := range result.Services {
		serviceObj, err := f.HealthCRM.GetServiceByID(ctx, service)
		if err != nil {
			return nil, err
		}

		services = append(services, *serviceObj)
	}

	output := &dto.BookingOutput{
		ID:                     result.ID,
		Services:               services,
		Date:                   result.Date,
		Facility:               result.Facility,
		Client:                 result.Client,
		OrganisationID:         result.OrganisationID,
		ProgramID:              result.ProgramID,
		VerificationCode:       verificationCode,
		VerificationCodeStatus: result.VerificationCodeStatus,
		BookingStatus:          result.BookingStatus,
	}

	return output, nil
}

// ListBookings is used to show a paginated list of client bookings
func (f *UseCaseFacilityImpl) ListBookings(ctx context.Context, clientID string, bookingState enums.BookingState, pagination dto.PaginationsInput) (*dto.BookingPage, error) {
	pageInput := &domain.Pagination{
		Limit:       pagination.Limit,
		CurrentPage: pagination.CurrentPage,
	}

	results, page, err := f.Query.ListBookings(ctx, clientID, bookingState, pageInput)
	if err != nil {
		return nil, err
	}

	var output []dto.BookingOutput

	for _, result := range results {
		var services []domain.FacilityService

		for _, service := range result.Services {
			serviceObj, err := f.HealthCRM.GetServiceByID(ctx, service)
			if err != nil {
				return nil, err
			}

			services = append(services, *serviceObj)
		}

		facility, err := f.HealthCRM.GetCRMFacilityByID(ctx, *result.Facility.ID)
		if err != nil {
			return nil, err
		}

		result.Facility = *facility

		output = append(output, dto.BookingOutput{
			ID:                     result.ID,
			Active:                 result.Active,
			Services:               services,
			Date:                   result.Date,
			Facility:               result.Facility,
			Client:                 result.Client,
			OrganisationID:         result.OrganisationID,
			ProgramID:              result.ProgramID,
			VerificationCode:       result.VerificationCode,
			VerificationCodeStatus: result.VerificationCodeStatus,
			BookingStatus:          result.BookingStatus,
		})
	}

	return &dto.BookingPage{
		Results:    output,
		Pagination: *page,
	}, nil
}

// VerifyBookingCode is used to verify clients booking code upon their arrival in a facility
func (f *UseCaseFacilityImpl) VerifyBookingCode(ctx context.Context, bookingID string, code string, programID string) (bool, error) {
	payload := &domain.Booking{
		ID:               bookingID,
		ProgramID:        programID,
		VerificationCode: code,
	}

	updateData := map[string]interface{}{
		"verification_code_status": enums.Verified,
	}

	err := f.Update.UpdateBooking(ctx, payload, updateData)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, nil
	}

	return true, nil
}
