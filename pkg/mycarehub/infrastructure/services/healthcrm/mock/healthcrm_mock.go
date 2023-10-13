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
	MockCreateFacilityFn      func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	MockGetFacilityServicesFn func(ctx context.Context, facilityID string) (*healthcrm.FacilityServicePage, error)
	MockGetFacilityByIDFn     func(ctx context.Context, id string) (*healthcrm.FacilityOutput, error)
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
				BusinessHours: []healthcrm.BusinessHoursOutput{
					{
						ID:          gofakeit.UUID(),
						Day:         "MONDAY",
						OpeningTime: "12:12",
						ClosingTime: "17:19",
						FacilityID:  gofakeit.UUID(),
					},
				},
			}, nil
		},
		MockGetFacilityServicesFn: func(ctx context.Context, facilityID string) (*healthcrm.FacilityServicePage, error) {
			return &healthcrm.FacilityServicePage{
				Results: []healthcrm.FacilityService{
					{
						ID:          gofakeit.UUID(),
						Name:        "Client Service Test",
						Description: "Client Description",
						Identifiers: []*healthcrm.ServiceIdentifier{
							{
								ID:              gofakeit.UUID(),
								IdentifierType:  "CIEL",
								IdentifierValue: "123456",
								ServiceID:       gofakeit.UUID(),
							},
						},
					},
				},
				Count:       120,
				Next:        "/next/",
				Previous:    "/previous/",
				PageSize:    30,
				CurrentPage: 10,
				TotalPages:  40,
				StartIndex:  1,
				EndIndex:    30,
			}, nil
		},
		MockGetFacilityByIDFn: func(ctx context.Context, id string) (*healthcrm.FacilityOutput, error) {
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
				Address:      "12-Meru",
				Contacts: []healthcrm.ContactsOutput{
					{
						ID:           id,
						ContactType:  "PHONE",
						ContactValue: interserviceclient.TestUserPhoneNumber,
						Active:       true,
						Role:         "PRIMARY_CONTACT",
						FacilityID:   id,
					},
				},
				Identifiers: []healthcrm.IdentifiersOutput{
					{
						ID:              id,
						IdentifierType:  "MFL Code",
						IdentifierValue: "11094",
						FacilityID:      id,
					},
				},
				BusinessHours: []healthcrm.BusinessHoursOutput{
					{
						ID:          id,
						Day:         "MONDAY",
						OpeningTime: "06:40:00",
						ClosingTime: "18:30:00",
						FacilityID:  id,
					},
				},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of creating a facility
func (sc HealthClientMock) CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
	return sc.MockCreateFacilityFn(ctx, facility)
}

// GetFacilityServices mocks the implementation of getting facility services
func (sc HealthClientMock) GetFacilityServices(ctx context.Context, facilityID string) (*healthcrm.FacilityServicePage, error) {
	return sc.MockGetFacilityServicesFn(ctx, facilityID)
}

// GetFacilityByID mocks the implementation of retrieving a facility with its id
func (sc HealthClientMock) GetFacilityByID(ctx context.Context, id string) (*healthcrm.FacilityOutput, error) {
	return sc.MockGetFacilityByIDFn(ctx, id)
}
