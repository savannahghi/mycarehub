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
	SearchFacilityByParam(ctx context.Context, searchTerm string) (*domain.FacilityPage, error)
	UpdateCRMFacility(ctx context.Context, id string, updatePayload *domain.Facility) (*domain.Facility, error)
	GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error)
}

// IHealthCRMClient defines the signature of the methods in the healthcrm library that perform specifies actions
type IHealthCRMClient interface {
	CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	SearchFacility(ctx context.Context, searchTerm string) (*healthcrm.FacilityPage, error)
	UpdateFacility(ctx context.Context, id string, updatePayload *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
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
	facilities := h.mapMCHDomainFacilityToHealthCRMFacility(facility)

	var facilityOutput []*domain.Facility

	for _, facilityInput := range facilities {
		output, err := h.clientSDK.CreateFacility(ctx, facilityInput)
		if err != nil {
			return nil, err
		}

		facilityOutput = h.mapHealthCRMFacilityToMCHDomainFacility(output)

	}

	return facilityOutput, nil
}

// SearchFacilityByParam searches or a facility ih health crm using the provided search term
func (h *HealthCRMImpl) SearchFacilityByParam(ctx context.Context, searchTerm string) (*domain.FacilityPage, error) {
	results, err := h.clientSDK.SearchFacility(ctx, searchTerm)
	if err != nil {
		return nil, err
	}

	var page *domain.FacilityPage

	next, err := strconv.Atoi(results.Next)
	if err != nil {
		return nil, err
	}

	page.Pagination.Count = int64(results.Count)
	page.Pagination.CurrentPage = results.CurrentPage
	page.Pagination.NextPage = &next
	page.Pagination.Limit = results.EndIndex

	var facilityOutput []*domain.Facility

	for _, result := range results.Results {
		facilityOutput = append(facilityOutput, &domain.Facility{
			ID:          &result.ID,
			Name:        result.Name,
			Phone:       result.Contacts[0].ContactValue,
			Active:      true,
			Country:     result.Country,
			County:      result.County,
			Address:     result.Address,
			Description: result.Description,
			Identifier: domain.FacilityIdentifier{
				ID:    result.Identifiers[0].ID,
				Type:  enums.FacilityIdentifierType(result.Identifiers[0].IdentifierType),
				Value: result.Identifiers[0].IdentifierValue,
			},
			Coordinates: domain.Coordinates{
				Lat: result.Coordinates.Latitude,
				Lng: result.Coordinates.Longitude,
			},
		})
	}

	page.Facilities = facilityOutput

	return page, nil
}

// UpdateCRMFacility is used to update facility in healthCRM
func (h *HealthCRMImpl) UpdateCRMFacility(ctx context.Context, id string, updatePayload *domain.Facility) (*domain.Facility, error) {
	facilities := []*domain.Facility{
		updatePayload,
	}

	mapped := h.mapMCHDomainFacilityToHealthCRMFacility(facilities)

	facility, err := h.clientSDK.UpdateFacility(ctx, id, mapped[0])
	if err != nil {
		return nil, err
	}

	output := h.mapHealthCRMFacilityToMCHDomainFacility(facility)

	return output[0], nil
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

	facilityOutput = append(facilityOutput, &domain.Facility{
		ID:          &output.ID,
		Name:        output.Name,
		Phone:       output.Contacts[0].ContactValue,
		Country:     output.Country,
		County:      output.County,
		Address:     output.Address,
		Description: output.Description,
		Identifier: domain.FacilityIdentifier{
			Active: true,
			Type:   enums.FacilityIdentifierType(output.Identifiers[0].IdentifierType),
			Value:  output.Identifiers[0].IdentifierValue,
		},
		WorkStationDetails: domain.WorkStationDetails{},
		Coordinates: domain.Coordinates{
			Lat: output.Coordinates.Latitude,
			Lng: output.Coordinates.Longitude,
		},
	})

	return facilityOutput
}

// mapMCHDomainFacilityToHealthCRMFacility maps mch facility to health CRM's facility
func (h *HealthCRMImpl) mapMCHDomainFacilityToHealthCRMFacility(facility []*domain.Facility) []*healthcrm.Facility {
	var facilities []*healthcrm.Facility

	for _, facilityObj := range facility {
		var identifiers []healthcrm.Identifiers

		identifiers = append(identifiers, healthcrm.Identifiers{
			IdentifierType:  facilityObj.Identifier.Type.String(),
			IdentifierValue: facilityObj.Identifier.Value,
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
			Contacts:      contacts,
			Identifiers:   identifiers,
			BusinessHours: []any{},
		})
	}

	return facilities
}
