package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/healthcrm"
)

// HealthClientMock mocks the health CRM client library implementations
type HealthClientMock struct {
	MockCreateFacilityFn  func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	MockSearchFacilityFn  func(ctx context.Context, searchTerm string) (*healthcrm.FacilityPage, error)
	MockUpdateFacilityFn  func(ctx context.Context, id string, updatePayload *healthcrm.Facility) (*healthcrm.FacilityOutput, error)
	MockGetFacilityByIDFn func(ctx context.Context, id string) (*healthcrm.FacilityOutput, error)
}

// NewHealthCRMClientMock initializes our client mocks
func NewHealthCRMClientMock() *HealthClientMock {
	return &HealthClientMock{
		MockCreateFacilityFn: func(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
			return &healthcrm.FacilityOutput{
				ID:            gofakeit.UUID(),
				Created:       time.Now(),
				Name:          gofakeit.BeerName(),
				Description:   gofakeit.BeerName(),
				FacilityType:  "HOSPITAL",
				County:        gofakeit.CountryAbr(),
				Country:       gofakeit.CountryAbr(),
				Coordinates:   healthcrm.CoordinatesOutput{},
				Status:        "DRAFT",
				Address:       "12-MERU",
				Contacts:      []healthcrm.ContactsOutput{},
				Identifiers:   []healthcrm.IdentifiersOutput{},
				BusinessHours: []any{},
			}, nil
		},
		MockSearchFacilityFn: func(ctx context.Context, searchTerm string) (*healthcrm.FacilityPage, error) {
			return &healthcrm.FacilityPage{
				Count:       10,
				Next:        "/next",
				Previous:    nil,
				PageSize:    2,
				CurrentPage: 2,
				TotalPages:  2,
				StartIndex:  1,
				EndIndex:    1,
				Results: []healthcrm.FacilityOutput{
					{
						ID:            gofakeit.UUID(),
						Created:       time.Now(),
						Name:          gofakeit.BeerName(),
						Description:   gofakeit.BeerName(),
						FacilityType:  "HOSPITAL",
						County:        gofakeit.CountryAbr(),
						Country:       gofakeit.CountryAbr(),
						Coordinates:   healthcrm.CoordinatesOutput{},
						Status:        "DRAFT",
						Address:       "12-Meru",
						Contacts:      []healthcrm.ContactsOutput{},
						Identifiers:   []healthcrm.IdentifiersOutput{},
						BusinessHours: []any{},
					},
				},
			}, nil
		},
		MockUpdateFacilityFn: func(ctx context.Context, id string, updatePayload *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
			return &healthcrm.FacilityOutput{
				ID:            gofakeit.UUID(),
				Created:       time.Now(),
				Name:          gofakeit.BeerName(),
				Description:   gofakeit.BeerName(),
				FacilityType:  "HOSPITAL",
				County:        gofakeit.CountryAbr(),
				Country:       gofakeit.CountryAbr(),
				Coordinates:   healthcrm.CoordinatesOutput{},
				Status:        "DRAFT",
				Address:       "12-Meru",
				Contacts:      []healthcrm.ContactsOutput{},
				Identifiers:   []healthcrm.IdentifiersOutput{},
				BusinessHours: []any{},
			}, nil
		},
		MockGetFacilityByIDFn: func(ctx context.Context, id string) (*healthcrm.FacilityOutput, error) {
			return &healthcrm.FacilityOutput{
				ID:            gofakeit.UUID(),
				Created:       time.Now(),
				Name:          gofakeit.BeerName(),
				Description:   gofakeit.BeerName(),
				FacilityType:  "HOSPITAL",
				County:        gofakeit.CountryAbr(),
				Country:       gofakeit.CountryAbr(),
				Coordinates:   healthcrm.CoordinatesOutput{},
				Status:        "DRAFT",
				Address:       "12-Meru",
				Contacts:      []healthcrm.ContactsOutput{},
				Identifiers:   []healthcrm.IdentifiersOutput{},
				BusinessHours: []any{},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of creating a facility
func (sc HealthClientMock) CreateFacility(ctx context.Context, facility *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
	return sc.MockCreateFacilityFn(ctx, facility)
}

// SearchFacility mocks the implementation of searching of facility in health crm
func (sc *HealthClientMock) SearchFacilityByParam(ctx context.Context, searchTerm string) (*healthcrm.FacilityPage, error) {
	return sc.MockSearchFacilityFn(ctx, searchTerm)
}

// UpdateFacility mocks the implementation of updating a facility
func (sc HealthClientMock) UpdateFacility(ctx context.Context, id string, updatePayload *healthcrm.Facility) (*healthcrm.FacilityOutput, error) {
	return sc.MockUpdateFacilityFn(ctx, id, updatePayload)
}

// GetFacilityByID mocks the implementation of retrieving a facility with its id
func (sc HealthClientMock) GetFacilityByID(ctx context.Context, id string) (*healthcrm.FacilityOutput, error) {
	return sc.MockGetFacilityByIDFn(ctx, id)
}
