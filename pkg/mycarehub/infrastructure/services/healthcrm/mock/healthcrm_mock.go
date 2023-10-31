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
	MockCreateFacilityFn  func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	MockGetServicesFn     func(ctx context.Context, pagination *healthcrm.Pagination) (*healthcrm.FacilityServicePage, error)
	MockGetFacilityByIDFn func(ctx context.Context, id string) (*healthcrm.FacilityOutput, error)
	MockGetFacilitiesFn   func(ctx context.Context, location *healthcrm.Coordinates, serviceIDs []string, searchParameter string, pagination *healthcrm.Pagination) (*healthcrm.FacilityPage, error)
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
		MockGetServicesFn: func(ctx context.Context, pagination *healthcrm.Pagination) (*healthcrm.FacilityServicePage, error) {
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
		MockGetFacilitiesFn: func(ctx context.Context, location *healthcrm.Coordinates, serviceIDs []string, searchParameter string, pagination *healthcrm.Pagination) (*healthcrm.FacilityPage, error) {
			latitude := -1.2979512335313856
			longitude := 36.78882506563385
			return &healthcrm.FacilityPage{
				Count:       0,
				Next:        "",
				Previous:    nil,
				PageSize:    0,
				CurrentPage: 0,
				TotalPages:  0,
				StartIndex:  0,
				EndIndex:    0,
				Results: []healthcrm.FacilityOutput{
					{
						ID:           gofakeit.UUID(),
						Created:      time.Time{},
						Name:         gofakeit.Name(),
						Description:  "",
						FacilityType: "",
						County:       "",
						Country:      "",
						Coordinates: healthcrm.CoordinatesOutput{
							Latitude:  latitude,
							Longitude: longitude,
						},
						Status:  "",
						Address: "",
						Contacts: []healthcrm.ContactsOutput{
							{
								ID:           "",
								ContactType:  "PHONE",
								ContactValue: gofakeit.PhoneFormatted(),
								Active:       false,
								Role:         "",
								FacilityID:   "",
							},
						},
						Identifiers: []healthcrm.IdentifiersOutput{
							{
								IdentifierType:  "MFL_CODE",
								IdentifierValue: "12345",
							},
						},
						BusinessHours: []healthcrm.BusinessHoursOutput{
							{
								Day:         "MONDAY",
								OpeningTime: "08:00",
								ClosingTime: "10:00",
							},
						},
						Services: []healthcrm.FacilityService{
							{
								ID:          gofakeit.UUID(),
								Name:        "PEP",
								Description: "",
								Identifiers: []*healthcrm.ServiceIdentifier{
									{
										ID:              "",
										IdentifierType:  "CIEL",
										IdentifierValue: "1234",
										ServiceID:       "",
									},
								},
							},
						},
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

// GetServices mocks the implementation of getting services
func (sc HealthClientMock) GetServices(ctx context.Context, pagination *healthcrm.Pagination) (*healthcrm.FacilityServicePage, error) {
	return sc.MockGetServicesFn(ctx, pagination)
}

// GetFacilityByID mocks the implementation of retrieving a facility with its id
func (sc HealthClientMock) GetFacilityByID(ctx context.Context, id string) (*healthcrm.FacilityOutput, error) {
	return sc.MockGetFacilityByIDFn(ctx, id)
}

// GetFacilities mocks the implementation of retrieving facilities
func (sc HealthClientMock) GetFacilities(ctx context.Context, location *healthcrm.Coordinates, serviceIDs []string, searchParameter string, pagination *healthcrm.Pagination) (*healthcrm.FacilityPage, error) {
	return sc.MockGetFacilitiesFn(ctx, location, serviceIDs, searchParameter, pagination)
}
