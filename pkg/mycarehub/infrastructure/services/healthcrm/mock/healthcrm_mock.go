package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/healthcrm"
	"github.com/savannahghi/interserviceclient"
)

// HealthClientMock mocks the health CRM client library implementations
type HealthClientMock struct {
	MockCreateFacilityFn func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
}

// NewHealthCRMClientMock initializes our client mocks
func NewHealthCRMClientMock() *HealthClientMock {
	return &HealthClientMock{
		MockCreateFacilityFn: func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
			return &healthcrm.FacilityOutput{
				ID:           gofakeit.UUID(),
				Created:      time.Now(),
				Name:         gofakeit.BeerName(),
				Description:  gofakeit.BeerName(),
				FacilityType: "HOSPITAL",
				County:       gofakeit.CountryAbr(),
				Country:      gofakeit.CountryAbr(),
				Coordinates:  healthcrm.CoordinatesOutput{},
				Status:       "DRAFT",
				Address:      "12-MERU",
				Contacts: []healthcrm.ContactsOutput{
					{
						ID:           gofakeit.UUID(),
						ContactType:  "PHONE_NUMBER",
						ContactValue: interserviceclient.TestUserPhoneNumber,
						Active:       true,
						Role:         "PRIMARY_CONTACT",
						FacilityID:   gofakeit.UUID(),
					},
				},
				Identifiers: []healthcrm.IdentifiersOutput{
					{
						ID:              gofakeit.UUID(),
						IdentifierType:  "MFL Code",
						IdentifierValue: "11094",
						FacilityID:      gofakeit.UUID(),
					},
				},
				BusinessHours: []any{},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of creating a facility
func (sc HealthClientMock) CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
	return sc.MockCreateFacilityFn(ctx, facility)
}
