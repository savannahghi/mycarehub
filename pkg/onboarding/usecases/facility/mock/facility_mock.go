package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// CreateMock is a mock of the create methods
type CreateMock struct {
	CreateFacilityFn func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
}

// NewCreateMock initializes a new instance of `GormMock` then mocking the case of success.
func NewCreateMock() *CreateMock {
	return &CreateMock{
		CreateFacilityFn: func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
			id := uuid.New()
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
	}
}

// CreateFacility mocks the implementation of `gorm's` CreateFacility method.
func (f *CreateMock) CreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.CreateFacilityFn(ctx, facility)
}

// QueryMock is a mock of the query methods
type QueryMock struct {
	RetrieveFacilityFn func(ctx context.Context, id *uuid.UUID) (*domain.Facility, error)
	GetFacilitiesFn    func(ctx context.Context) ([]*domain.Facility, error)
}

// NewQueryMock initializes a new instance of `GormMock` then mocking the case of success.
func NewQueryMock() *QueryMock {
	return &QueryMock{

		RetrieveFacilityFn: func(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
			facilityID := uuid.New()
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
		GetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			facilityID := uuid.New()
			name := "test-facility"
			code := "t-100"
			county := "test-county"
			description := "test description"
			return []*domain.Facility{
				{
					ID:          facilityID,
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			}, nil
		},
	}
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (f *QueryMock) RetrieveFacility(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
	return f.RetrieveFacilityFn(ctx, id)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method
func (f *QueryMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.GetFacilitiesFn(ctx)
}
