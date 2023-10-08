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
}

// IHealthCRMClient defines the signature of the methods in the healthcrm library that perform specifies actions
type IHealthCRMClient interface {
	CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
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

	var facilityOutput []*domain.Facility

	for _, facilityInput := range facilities {
		output, err := h.clientSDK.CreateFacility(ctx, facilityInput)
		if err != nil {
			return nil, err
		}

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

	}

	return facilityOutput, nil
}
