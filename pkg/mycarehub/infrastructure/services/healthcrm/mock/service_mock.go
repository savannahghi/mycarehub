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
	MockCreateFacilityFn func(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error)
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
					Identifier:         []*domain.FacilityIdentifier{},
					WorkStationDetails: domain.WorkStationDetails{},
					Coordinates:        &domain.Coordinates{},
				},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of health crm create facility method
func (m *HealthCRMServiceMock) CreateFacility(ctx context.Context, facility []*domain.Facility) ([]*domain.Facility, error) {
	return m.MockCreateFacilityFn(ctx, facility)
}
