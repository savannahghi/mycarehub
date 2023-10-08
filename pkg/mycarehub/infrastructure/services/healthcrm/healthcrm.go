package healthcrm

import (
	"context"
	"strconv"

	"github.com/savannahghi/healthcrm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// IHealthCRMService holds the methods required to interact with healthcrm beckend service through healthcrm library
type IHealthCRMService interface {
	CreateFacility(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error)
	GetServicesOfferedInAFacility(ctx context.Context, facilityID string) (*domain.FacilityServicePage, error)
	GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error)
}

// IHealthCRMClient defines the signature of the methods in the healthcrm library that perform specifies actions
type IHealthCRMClient interface {
	CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	GetFacilityServices(ctx context.Context, facilityID string) (*healthcrm.FacilityServicePage, error)
	GetFacilityByID(ctx context.Context, id string) (*healthcrm.FacilityOutput, error)
}

// HealthCRMImpl is the implementation of health crm's service client
type HealthCRMImpl struct {
	clientSDK IHealthCRMClient
}

// NewHealthCRMService instantiates the healthCRM's service
func NewHealthCRMService(client IHealthCRMClient) *HealthCRMImpl {
	return &HealthCRMImpl{
		clientSDK: client,
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
			Contacts:    contacts,
			Identifiers: identifiers,
		})
	}

	var facilityOutput []*domain.Facility

	for _, facilityInput := range facilities {
		output, err := h.clientSDK.CreateFacility(ctx, facilityInput)
		if err != nil {
			return nil, err
		}

		var identifiers []*domain.FacilityIdentifier
		for _, id := range output.Identifiers {
			identifiers = append(identifiers, &domain.FacilityIdentifier{
				Type:   enums.FacilityIdentifierType(id.IdentifierType),
				Value:  id.IdentifierValue,
				Active: true,
			})
		}

		// Health CRM ID is also ann identifier, hence the mapping below
		identifiers = append(identifiers, &domain.FacilityIdentifier{
			Type:   enums.FacilityIdentifierTypeHealthCRM,
			Value:  output.ID,
			Active: true,
		})

		facilityOutput = append(facilityOutput, &domain.Facility{
			ID:                 &output.ID,
			Name:               output.Name,
			Phone:              output.Contacts[0].ContactValue,
			Active:             true,
			Country:            output.Country,
			County:             output.County,
			Address:            output.Address,
			Description:        output.Description,
			Identifiers:        identifiers,
			WorkStationDetails: domain.WorkStationDetails{},
			Coordinates: &domain.Coordinates{
				Lat: output.Coordinates.Latitude,
				Lng: output.Coordinates.Longitude,
			},
		})

	}

	return facilityOutput, nil
}

// GetServicesOfferedInAFacility retrieves the services offered in a facility
func (h *HealthCRMImpl) GetServicesOfferedInAFacility(ctx context.Context, facilityID string) (*domain.FacilityServicePage, error) {
	output, err := h.clientSDK.GetFacilityServices(ctx, facilityID)
	if err != nil {
		return nil, err
	}

	var facilityPage domain.FacilityServicePage
	var facilityServices []domain.FacilityService

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

		facilityServices = append(facilityServices, *facilityService)
	}

	facilityPage.Results = facilityServices
	facilityPage.Count = output.Count
	facilityPage.CurrentPage = output.CurrentPage
	facilityPage.EndIndex = output.EndIndex
	facilityPage.StartIndex = output.StartIndex
	facilityPage.Next = output.Next
	facilityPage.Previous = output.Previous
	facilityPage.TotalPages = output.TotalPages

	return &facilityPage, nil
}

// GetCRMFacilityByID is used to retrieve facility from health crm
func (h *HealthCRMImpl) GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error) {
	results, err := h.clientSDK.GetFacilityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := h.mapHealthCRMFacilityToMCHDomainFacility(results)

	return mapped[0], nil
}

// mapHealthCRMFacilityToMCHDomainFacility maps health crm facility to mch domain facility
func (h *HealthCRMImpl) mapHealthCRMFacilityToMCHDomainFacility(output *healthcrm.FacilityOutput) []*domain.Facility {
	var facilityOutput []*domain.Facility

	var operatingHours []domain.BusinessHours

	var facilityIdentifiers []*domain.FacilityIdentifier

	for _, identifier := range output.Identifiers {
		facilityIdentifiers = append(facilityIdentifiers, &domain.FacilityIdentifier{
			ID:     identifier.ID,
			Active: true,
			Type:   enums.FacilityIdentifierType(identifier.IdentifierType),
			Value:  identifier.IdentifierValue,
		})
	}

	for _, result := range output.BusinessHours {
		operatingHours = append(operatingHours, domain.BusinessHours{
			ID:          result.ID,
			Day:         result.Day,
			OpeningTime: result.OpeningTime,
			ClosingTime: result.ClosingTime,
			FacilityID:  result.FacilityID,
		})
	}

	facilityOutput = append(facilityOutput, &domain.Facility{
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
	})

	return facilityOutput
}
