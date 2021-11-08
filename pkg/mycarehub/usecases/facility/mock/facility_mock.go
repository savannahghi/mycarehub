package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// FacilityUsecaseMock mocks the implementation of facility usecase methods
type FacilityUsecaseMock struct {
	MockGetOrCreateFacilityFn       func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	MockRetrieveFacilityFn          func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	MockRetrieveFacilityByMFLCodeFn func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	MockGetFacilitiesFn             func(ctx context.Context) ([]*domain.Facility, error)
	MockListFacilitiesFn            func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	DeleteFacilityFn                func(ctx context.Context, id string) (bool, error)
	FetchFacilitiesFn               func(ctx context.Context) ([]*domain.Facility, error)
	MockInactivateFacilityFn        func(ctx context.Context, mflCode *string) (bool, error)
	MockReactivateFacilityFn        func(ctx context.Context, mflCode *string) (bool, error)
}

// NewFacilityUsecaseMock initializes a new instance of `GormMock` then mocking the case of success.
func NewFacilityUsecaseMock() *FacilityUsecaseMock {
	ID := uuid.New().String()
	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facilityInput := &domain.Facility{
		ID:          &ID,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	var facilitiesList []*domain.Facility
	facilitiesList = append(facilitiesList, facilityInput)

	nextPage := 3
	previousPage := 1
	facilitiesPage := &domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:        1,
			CurrentPage:  2,
			Count:        3,
			TotalPages:   3,
			NextPage:     &nextPage,
			PreviousPage: &previousPage,
		},
		Facilities: []domain.Facility{
			{
				ID:          &ID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			},
		},
	}

	return &FacilityUsecaseMock{
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockGetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},

		DeleteFacilityFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},

		FetchFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, mflCode *string) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, mflCode *string) (bool, error) {
			return true, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (f *FacilityUsecaseMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return f.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (f *FacilityUsecaseMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacilityByMFLCode method.
func (f *FacilityUsecaseMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method
func (f *FacilityUsecaseMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.MockGetFacilitiesFn(ctx)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (f *FacilityUsecaseMock) ListFacilities(
	ctx context.Context,
	searchTerm *string,
	filterInput []*dto.FiltersInput,
	paginationsInput *dto.PaginationsInput,
) (*domain.FacilityPage, error) {
	return f.MockListFacilitiesFn(ctx, searchTerm, filterInput, paginationsInput)
}

// DeleteFacility mocks the implementation of deleting a facility by ID
func (f *FacilityUsecaseMock) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return f.DeleteFacilityFn(ctx, id)
}

// FetchFacilities mocks the implementation of fetch all facilities
func (f *FacilityUsecaseMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.FetchFacilitiesFn(ctx)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (f *FacilityUsecaseMock) InactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	return f.MockInactivateFacilityFn(ctx, mflCode)
}

// ReactivateFacility mocks the implementation of re-activating the active status of a particular facility
func (f *FacilityUsecaseMock) ReactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	return f.MockReactivateFacilityFn(ctx, mflCode)
}
