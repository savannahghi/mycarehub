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
	MockCreateFacilityFn        func(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error)
	MockSearchFacilityByParamFn func(ctx context.Context, searchTerm string) (*domain.FacilityPage, error)
	MockUpdateCRMFacilityFn     func(ctx context.Context, id string, updatePayload *domain.Facility) (*domain.Facility, error)
	MockGetCRMFacilityByIDFn    func(ctx context.Context, id string) (*domain.Facility, error)
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
					Identifier:         domain.FacilityIdentifier{},
					WorkStationDetails: domain.WorkStationDetails{},
					Coordinates:        domain.Coordinates{},
				},
			}, nil
		},
		MockSearchFacilityByParamFn: func(ctx context.Context, searchTerm string) (*domain.FacilityPage, error) {
			return &domain.FacilityPage{
				Pagination: domain.Pagination{
					Limit:        5,
					CurrentPage:  2,
					Count:        20,
					TotalPages:   20,
					NextPage:     new(int),
					PreviousPage: new(int),
					Sort:         &domain.SortParam{},
				},
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
						Identifier:         domain.FacilityIdentifier{},
						WorkStationDetails: domain.WorkStationDetails{},
						Coordinates:        domain.Coordinates{},
					},
				},
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
				Identifier:         domain.FacilityIdentifier{},
				WorkStationDetails: domain.WorkStationDetails{},
				Coordinates:        domain.Coordinates{},
			}, nil
		},
		MockUpdateCRMFacilityFn: func(ctx context.Context, id string, updatePayload *domain.Facility) (*domain.Facility, error) {
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
				Identifier:         domain.FacilityIdentifier{},
				WorkStationDetails: domain.WorkStationDetails{},
				Coordinates:        domain.Coordinates{},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of health crm create facility method
func (m *HealthCRMServiceMock) CreateFacility(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error) {
	return m.MockCreateFacilityFn(ctx, facility)
}

// SearchFacilityByParam mocks the implementation of searching a facility in health crm
func (m *HealthCRMServiceMock) SearchFacilityByParam(ctx context.Context, searchTerm string) (*domain.FacilityPage, error) {
	return m.MockSearchFacilityByParamFn(ctx, searchTerm)
}

// GetCRMFacilityByID mocks the implementation of retrieving a facility in health crm using its id
func (m *HealthCRMServiceMock) GetCRMFacilityByID(ctx context.Context, id string) (*domain.Facility, error) {
	return m.MockGetCRMFacilityByIDFn(ctx, id)
}

// UpdateCRMFacility mocks the implementation of updating facility in health crm
func (m *HealthCRMServiceMock) UpdateCRMFacility(ctx context.Context, id string, updatePayload *domain.Facility) (*domain.Facility, error) {
	return m.MockUpdateCRMFacilityFn(ctx, id, updatePayload)
}
