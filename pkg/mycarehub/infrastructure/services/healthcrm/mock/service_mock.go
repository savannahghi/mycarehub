package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SMSServiceMock mocks the health CRM service mock methods
type HealthCRMServiceMock struct {
	MockCreateFacilityFn                func(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error)
	MockGetServicesFn                   func(ctx context.Context, facilityID string, pagination *domain.Pagination) (*domain.FacilityServicePage, error)
	MockGetCRMFacilityByIDFn            func(ctx context.Context, id string) (*domain.Facility, error)
	MockGetFacilitiesOfferingAServiceFn func(ctx context.Context, serviceID string, pagination *domain.Pagination) (*domain.FacilityPage, error)
}

// NewHealthServiceMock initializes the mock service
func NewHealthServiceMock() *HealthCRMServiceMock {
	UUID := uuid.New().String()
	return &HealthCRMServiceMock{
		MockCreateFacilityFn: func(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error) {
			return []*domain.Facility{
				{
					ID:                 &UUID,
					Name:               gofakeit.BeerName(),
					Phone:              interserviceclient.TestUserPhoneNumber,
					Active:             true,
					Country:            gofakeit.CountryAbr(),
					County:             gofakeit.Country(),
					Address:            gofakeit.Address().Address,
					Description:        gofakeit.BeerMalt(),
					FHIROrganisationID: gofakeit.UUID(),
					Identifiers:        []*domain.FacilityIdentifier{},
					WorkStationDetails: domain.WorkStationDetails{},
					Coordinates:        &domain.Coordinates{},
					BusinessHours: []domain.BusinessHours{
						{
							ID:          UUID,
							Day:         "MONDAY",
							OpeningTime: "12:12",
							ClosingTime: "12:12",
							FacilityID:  gofakeit.UUID(),
						},
					},
					Services: []domain.FacilityService{
						{
							ID:          UUID,
							Name:        "WEIGHT",
							Description: "Weight in Kgs",
							Identifiers: []domain.ServiceIdentifier{
								{
									ID:              UUID,
									IdentifierType:  "CIEL",
									IdentifierValue: "1234",
									ServiceID:       gofakeit.UUID(),
								},
							},
						},
					},
				},
			}, nil
		},
		MockGetServicesFn: func(ctx context.Context, facilityID string, pagination *domain.Pagination) (*domain.FacilityServicePage, error) {
			return &domain.FacilityServicePage{
				Results: []domain.FacilityService{
					{
						ID:          gofakeit.UUID(),
						Name:        "Service Test",
						Description: "Service Description",
						Identifiers: []domain.ServiceIdentifier{
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
		MockGetCRMFacilityByIDFn: func(ctx context.Context, id string) (*domain.Facility, error) {
			return &domain.Facility{
				ID:                 &UUID,
				Name:               gofakeit.BeerName(),
				Phone:              interserviceclient.TestUserPhoneNumber,
				Active:             true,
				Country:            gofakeit.CountryAbr(),
				County:             gofakeit.Country(),
				Address:            gofakeit.Address().Address,
				Description:        gofakeit.BeerMalt(),
				FHIROrganisationID: gofakeit.UUID(),
				Identifiers:        []*domain.FacilityIdentifier{},
				WorkStationDetails: domain.WorkStationDetails{},
				Coordinates:        &domain.Coordinates{},
				BusinessHours:      []domain.BusinessHours{},
			}, nil
		},
		MockGetFacilitiesOfferingAServiceFn: func(ctx context.Context, serviceID string, pagination *domain.Pagination) (*domain.FacilityPage, error) {
			return &domain.FacilityPage{
				Pagination: domain.Pagination{},
				Facilities: []*domain.Facility{
					{
						ID:                 &UUID,
						Name:               gofakeit.BeerName(),
						Phone:              interserviceclient.TestUserPhoneNumber,
						Active:             true,
						Country:            gofakeit.CountryAbr(),
						County:             gofakeit.Country(),
						Address:            gofakeit.Address().Address,
						Description:        gofakeit.BeerMalt(),
						FHIROrganisationID: gofakeit.UUID(),
						Identifiers:        []*domain.FacilityIdentifier{},
						WorkStationDetails: domain.WorkStationDetails{},
						Coordinates:        &domain.Coordinates{},
						BusinessHours:      []domain.BusinessHours{},
					},
				},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of health crm create facility method
func (m *HealthCRMServiceMock) CreateFacility(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error) {
	return m.MockCreateFacilityFn(ctx, facility)
}

// GetServices mocks the implementation of getting services offered in a facility
func (m *HealthCRMServiceMock) GetServices(ctx context.Context, facilityID string, pagination *domain.Pagination) (*domain.FacilityServicePage, error) {
	return m.MockGetServicesFn(ctx, facilityID, pagination)
}

// GetCRMFacilityByID mocks the implementation of retrieving a facility in health crm using its id
func (m *HealthCRMServiceMock) GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error) {
	return m.MockGetCRMFacilityByIDFn(ctx, id)
}

// GetFacilitiesOfferingAService mocks the implementation of fetching facilities offering a certain service
func (m *HealthCRMServiceMock) GetFacilitiesOfferingAService(ctx context.Context, serviceID string, pagination *domain.Pagination) (*domain.FacilityPage, error) {
	return m.MockGetFacilitiesOfferingAServiceFn(ctx, serviceID, pagination)
}
