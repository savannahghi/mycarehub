package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	//Get
	GetOrCreateFacilityFn func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	GetFacilitiesFn       func(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityFn    func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	ListFacilitiesFn      func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination domain.FacilityPage) (*domain.FacilityPage, error)
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
	return &PostgresMock{
		GetOrCreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := enums.CountyTypeNairobi
			description := "This is just for mocking"
			return &domain.Facility{
				ID:          &id,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},
		GetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := enums.CountyTypeNairobi
			description := "This is just for mocking"
			return []*domain.Facility{
				{
					ID:          &id,
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			}, nil
		},
		RetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := enums.CountyTypeNairobi
			description := "test description"
			return &domain.Facility{
				ID:          &facilityID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},
		ListFacilitiesFn: func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination domain.FacilityPage) (*domain.FacilityPage, error) {
			facilityID := uuid.New().String()
			name := "test-facility"
			code := "t-100"
			county := enums.CountyTypeNairobi
			description := "test description"
			nextPage := 1
			previousPage := 1
			return &domain.FacilityPage{
				Pagination: domain.Pagination{
					Limit:        1,
					CurrentPage:  1,
					Count:        1,
					TotalPages:   1,
					NextPage:     &nextPage,
					PreviousPage: &previousPage,
				},
				Facilities: []domain.Facility{
					{
						ID:          &facilityID,
						Name:        name,
						Code:        code,
						Active:      true,
						County:      county,
						Description: description,
					},
				},
			}, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *PostgresMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.GetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id, isActive)
}

// ListFacilities mocks the implementation of  `gorm's` ListFacilities method.
func (gm *PostgresMock) ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination domain.FacilityPage) (*domain.FacilityPage, error) {
	return gm.ListFacilitiesFn(ctx, searchTerm, filter, pagination)
}
