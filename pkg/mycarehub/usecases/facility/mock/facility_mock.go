package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// FacilityUsecaseMock mocks the implementation of facility usecase methods
type FacilityUsecaseMock struct {
	MockRetrieveFacilityFn             func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	MockRetrieveFacilityByIdentifierFn func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error)
	MockGetFacilitiesFn                func(ctx context.Context) ([]*domain.Facility, error)
	MockListFacilitiesFn               func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	DeleteFacilityFn                   func(ctx context.Context, id int) (bool, error)
	FetchFacilitiesFn                  func(ctx context.Context) ([]*domain.Facility, error)
	MockInactivateFacilityFn           func(ctx context.Context, mflCode *int) (bool, error)
	MockUpdateFacilityFn               func(ctx context.Context, updateFacilityData *domain.UpdateFacilityPayload) error
	MockAddFacilityToProgramFn         func(ctx context.Context, facilityID []string) (bool, error)
}

// NewFacilityUsecaseMock initializes a new instance of `GormMock` then mocking the case of success.
func NewFacilityUsecaseMock() *FacilityUsecaseMock {
	ID := uuid.New().String()
	name := gofakeit.Name()
	country := "Kenya"
	phone := interserviceclient.TestUserPhoneNumber
	description := gofakeit.HipsterSentence(15)
	FHIROrganisationID := uuid.New().String()

	facilityInput := &domain.Facility{
		ID:                 &ID,
		Name:               name,
		Phone:              phone,
		Active:             true,
		Country:            country,
		Description:        description,
		FHIROrganisationID: FHIROrganisationID,
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
		Facilities: []*domain.Facility{
			{
				ID:          &ID,
				Name:        name,
				Phone:       phone,
				Active:      true,
				Country:     country,
				Description: description,
			},
		},
	}

	return &FacilityUsecaseMock{

		MockUpdateFacilityFn: func(ctx context.Context, updateFacilityData *domain.UpdateFacilityPayload) error {
			return nil
		},

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockRetrieveFacilityByIdentifierFn: func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockGetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},

		DeleteFacilityFn: func(ctx context.Context, id int) (bool, error) {
			return true, nil
		},

		FetchFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockAddFacilityToProgramFn: func(ctx context.Context, facilityID []string) (bool, error) {
			return true, nil
		},
	}
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (f *FacilityUsecaseMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByIdentifier mocks the implementation of `gorm's` RetrieveFacilityByIdentifier method.
func (f *FacilityUsecaseMock) RetrieveFacilityByIdentifier(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityByIdentifierFn(ctx, MFLCode, isActive)
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
func (f *FacilityUsecaseMock) DeleteFacility(ctx context.Context, id int) (bool, error) {
	return f.DeleteFacilityFn(ctx, id)
}

// FetchFacilities mocks the implementation of fetch all facilities
func (f *FacilityUsecaseMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.FetchFacilitiesFn(ctx)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (f *FacilityUsecaseMock) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return f.MockInactivateFacilityFn(ctx, mflCode)
}

// UpdateFacility mocks the implementation of updating a facility
func (f *FacilityUsecaseMock) UpdateFacility(ctx context.Context, updateFacilityData *domain.UpdateFacilityPayload) error {
	return f.MockUpdateFacilityFn(ctx, updateFacilityData)
}

// AddFacilityToProgram mocks the implementation of adding a facility to a program
func (f *FacilityUsecaseMock) AddFacilityToProgram(ctx context.Context, facilityID []string) (bool, error) {
	return f.MockAddFacilityToProgramFn(ctx, facilityID)
}
