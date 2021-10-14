package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	CreateFacilityFn   func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	GetFacilitiesFn    func(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityFn func(ctx context.Context, id *uuid.UUID) (*domain.Facility, error)
	FindFacilityFn     func(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error)
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
	return &PostgresMock{
		CreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
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
		GetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			id := uuid.New()
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
		FindFacilityFn: func(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error) {
			id := uuid.New()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"

			cursor := "1"
			startCursor := "1"
			endCursor := "1"

			return &dto.FacilityConnection{
				Edges: []*dto.FacilityEdge{
					{
						Cursor: &cursor,
						Node: &domain.Facility{
							ID:          id,
							Name:        name,
							Code:        code,
							Active:      true,
							County:      county,
							Description: description,
						},
					},
				},
				PageInfo: &firebasetools.PageInfo{
					HasNextPage:     false,
					HasPreviousPage: false,
					StartCursor:     &startCursor,
					EndCursor:       &endCursor,
				},
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of `gorm's` CreateFacility method.
func (gm *PostgresMock) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.CreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *PostgresMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return gm.GetFacilitiesFn(ctx)
}

// FindFacility mocks the implementation of  FindFacility method.
func (gm *PostgresMock) FindFacility(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error) {
	return gm.FindFacilityFn(ctx, pagination, filter, sort)
}
