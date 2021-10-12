package mock

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	CreateFacilityFn   func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	GetFacilitiesFn    func(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityFn func(ctx context.Context, id *int64) (*domain.Facility, error)
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
	return &PostgresMock{
		CreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
			id := int64(1)
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &domain.Facility{
				ID:          id,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},
		GetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			id := int64(1)
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return []*domain.Facility{
				{
					ID:          id,
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			}, nil
		},
		RetrieveFacilityFn: func(ctx context.Context, id *int64) (*domain.Facility, error) {
			facilityID := int64(1)
			name := "test-facility"
			code := "t-100"
			county := "test-county"
			description := "test description"
			return &domain.Facility{
				ID:          facilityID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of `gorm's` CreateFacility method.
func (gm *PostgresMock) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.CreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` CreateFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *int64) (*domain.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id)
}
