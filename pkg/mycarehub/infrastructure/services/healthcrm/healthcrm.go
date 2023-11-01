package healthcrm

import (
	"context"
	"strconv"

	"github.com/savannahghi/healthcrm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

var facilityIdentifiersMap = map[string]string{
	"MFL Code":   enums.FacilityIdentifierTypeMFLCode.String(),
	"Health CRM": enums.FacilityIdentifierTypeHealthCRM.String(),
	"Slade Code": enums.FacilityIdentifierTypeSladeCode.String(),
}

// IHealthCRMService holds the methods required to interact with healthcrm beckend service through healthcrm library
type IHealthCRMService interface {
	CreateFacility(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error)
	GetServices(ctx context.Context, pagination *domain.Pagination) (*domain.FacilityServicePage, error)
	GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error)
	GetFacilities(ctx context.Context, location *dto.LocationInput, serviceIDs []string, searchParameter string, pagination *domain.Pagination) ([]*domain.Facility, error)
}

// IHealthCRMClient defines the signature of the methods in the healthcrm library that perform specifies actions
type IHealthCRMClient interface {
	CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	GetServices(ctx context.Context, pagination *healthcrm.Pagination) (*healthcrm.FacilityServicePage, error)
	GetFacilityByID(ctx context.Context, id string) (*healthcrm.FacilityOutput, error)
	GetFacilities(ctx context.Context, location *healthcrm.Coordinates, serviceIDs []string, searchParameter string, pagination *healthcrm.Pagination) (*healthcrm.FacilityPage, error)
}

// HealthCRMImpl is the implementation of health crm's service client
type HealthCRMImpl struct {
	client IHealthCRMClient
}

// NewHealthCRMService instantiates the healthCRM's service
func NewHealthCRMService(client IHealthCRMClient) *HealthCRMImpl {
	return &HealthCRMImpl{
		client: client,
	}
}

// CreateFacility creates facility in service registry
func (h *HealthCRMImpl) CreateFacility(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error) {
	var facilities []*healthcrm.Facility

	for _, facilityObj := range facility {
		var identifiers []healthcrm.Identifiers

		identifiers = append(identifiers, healthcrm.Identifiers{
			IdentifierType:  facilityObj.Identifiers[0].Type.String(),
			IdentifierValue: facilityObj.Identifiers[0].Value,
		})

		var contacts []healthcrm.Contacts

		contacts = append(contacts, healthcrm.Contacts{
			ContactType:  "PHONE_NUMBER",
			ContactValue: facilityObj.Phone,
			Role:         "PRIMARY_CONTACT",
		})

		var businessHours []healthcrm.BusinessHours
		for _, businessHour := range facilityObj.BusinessHours {
			businessHours = append(businessHours, healthcrm.BusinessHours{
				Day:         businessHour.Day,
				OpeningTime: businessHour.OpeningTime,
				ClosingTime: businessHour.ClosingTime,
			})
		}

		facilities = append(facilities, &healthcrm.Facility{
			Name:         facilityObj.Name,
			Description:  facilityObj.Description,
			FacilityType: "HOSPITAL",
			County:       facilityObj.County,
			Country:      facilityObj.Country,
			Address:      facilityObj.Address,
			Coordinates: &healthcrm.Coordinates{
				Latitude:  strconv.FormatFloat(facilityObj.Coordinates.Lat, 'f', -1, 64),
				Longitude: strconv.FormatFloat(facilityObj.Coordinates.Lng, 'f', -1, 64),
			},
			Contacts:      contacts,
			Identifiers:   identifiers,
			BusinessHours: businessHours,
		})
	}

	var facilityOutput []*domain.Facility

	for _, facilityInput := range facilities {
		output, err := h.client.CreateFacility(ctx, facilityInput)
		if err != nil {
			return nil, err
		}

		facilityOutput = append(facilityOutput, mapHealthCRMFacilityToMCHDomainFacility(output))
	}

	return facilityOutput, nil
}

// GetServices is used to fetch all the services available in health crm
// This function is also used to list all the services available in a facility if the ID of that facility is provided
func (h *HealthCRMImpl) GetServices(ctx context.Context, pagination *domain.Pagination) (*domain.FacilityServicePage, error) {
	paginationInput := &healthcrm.Pagination{
		Page:     strconv.Itoa(pagination.CurrentPage),
		PageSize: strconv.Itoa(pagination.Limit),
	}

	output, err := h.client.GetServices(ctx, paginationInput)
	if err != nil {
		return nil, err
	}

	var service []domain.FacilityService
	for _, result := range output.Results {
		var serviceIdentifiers []domain.ServiceIdentifier
		for _, serviceIdentifier := range result.Identifiers {
			serviceIdentifiers = append(serviceIdentifiers, domain.ServiceIdentifier{
				ID:              serviceIdentifier.ID,
				IdentifierType:  serviceIdentifier.IdentifierType,
				IdentifierValue: serviceIdentifier.IdentifierValue,
				ServiceID:       serviceIdentifier.ServiceID,
			})
		}

		facilityService := &domain.FacilityService{
			ID:          result.ID,
			Name:        result.Name,
			Description: result.Description,
			Identifiers: serviceIdentifiers,
		}

		service = append(service, *facilityService)
	}

	return &domain.FacilityServicePage{
		Results:     service,
		Count:       output.Count,
		Next:        output.Next,
		Previous:    output.Previous,
		PageSize:    output.PageSize,
		CurrentPage: output.CurrentPage,
		TotalPages:  output.TotalPages,
		StartIndex:  output.StartIndex,
		EndIndex:    output.EndIndex,
	}, nil
}

// GetCRMFacilityByID is used to retrieve facility from health crm
func (h *HealthCRMImpl) GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error) {
	results, err := h.client.GetFacilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapHealthCRMFacilityToMCHDomainFacility(results)

	return mapped, nil
}

// GetFacilities retrieves a list of facilities associated with MyCareHub
// stored in HealthCRM. The method allows for filtering facilities by location proximity and services offered.
// When invoking this method, we need to ensure that either the serviceIDs or the SearchParameter is passed, but not both or
// else an error will be returned
func (h *HealthCRMImpl) GetFacilities(ctx context.Context, location *dto.LocationInput, serviceIDs []string, searchParameter string, pagination *domain.Pagination) ([]*domain.Facility, error) {
	paginationInput := &healthcrm.Pagination{
		Page:     strconv.Itoa(pagination.CurrentPage),
		PageSize: strconv.Itoa(pagination.Limit),
	}

	var coordinates *healthcrm.Coordinates

	if location != nil {
		coordinates = &healthcrm.Coordinates{
			Latitude:  strconv.FormatFloat(location.Lat, 'f', -1, 64),
			Longitude: strconv.FormatFloat(location.Lng, 'f', -1, 64),
		}
		if location.Radius != nil {
			coordinates.Radius = strconv.FormatFloat(*location.Radius, 'f', -1, 64)
		}
	}

	facilities, err := h.client.GetFacilities(ctx, coordinates, serviceIDs, searchParameter, paginationInput)
	if err != nil {
		return nil, err
	}

	var facilityOutput []*domain.Facility

	for _, facility := range facilities.Results {
		facilityOutput = append(facilityOutput, mapHealthCRMFacilityToMCHDomainFacility(&facility))
	}

	return facilityOutput, nil
}

// mapHealthCRMFacilityToMCHDomainFacility is used to transform the output received from healthcrm library after retrieving a facility to the domain model of a facility in mycarehub
func mapHealthCRMFacilityToMCHDomainFacility(output *healthcrm.FacilityOutput) *domain.Facility {
	var facilityIdentifiers []*domain.FacilityIdentifier

	for _, identifier := range output.Identifiers {
		facilityIdentifier := &domain.FacilityIdentifier{
			ID:     identifier.ID,
			Active: true,
			Type:   enums.FacilityIdentifierType(facilityIdentifiersMap[identifier.IdentifierType]),
			Value:  identifier.IdentifierValue,
		}

		facilityIdentifiers = append(facilityIdentifiers, facilityIdentifier)
	}

	var operatingHours []domain.BusinessHours

	for _, result := range output.BusinessHours {
		operatingHours = append(operatingHours, domain.BusinessHours{
			ID:          result.ID,
			Day:         result.Day,
			OpeningTime: result.OpeningTime,
			ClosingTime: result.ClosingTime,
			FacilityID:  result.FacilityID,
		})
	}

	var services []domain.FacilityService

	for _, service := range output.Services {
		var serviceIdentifiers []domain.ServiceIdentifier

		for _, identifier := range service.Identifiers {
			serviceIdentifiers = append(serviceIdentifiers, domain.ServiceIdentifier{
				ID:              identifier.ID,
				IdentifierType:  identifier.IdentifierType,
				IdentifierValue: identifier.IdentifierValue,
				ServiceID:       identifier.ServiceID,
			})
		}
		services = append(services, domain.FacilityService{
			ID:          service.ID,
			Name:        service.Name,
			Description: service.Description,
			Identifiers: serviceIdentifiers,
		})
	}

	return &domain.Facility{
		ID:                 &output.ID,
		Name:               output.Name,
		Phone:              output.Contacts[0].ContactValue,
		Active:             true,
		Country:            output.Country,
		County:             output.County,
		Address:            output.Address,
		Description:        output.Description,
		Identifiers:        facilityIdentifiers,
		WorkStationDetails: domain.WorkStationDetails{},
		Coordinates: &domain.Coordinates{
			Lat: output.Coordinates.Latitude,
			Lng: output.Coordinates.Longitude,
		},
		BusinessHours: operatingHours,
		Services:      services,
	}
}
