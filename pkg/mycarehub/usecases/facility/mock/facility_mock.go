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
	MockRetrieveFacilityByIdentifierFn func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error)
	MockGetFacilitiesFn                func(ctx context.Context) ([]*domain.Facility, error)
	MockListFacilitiesFn               func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockDeleteFacilityFn               func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	FetchFacilitiesFn                  func(ctx context.Context) ([]*domain.Facility, error)
	MockInactivateFacilityFn           func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockReactivateFacilityFn           func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error)
	MockUpdateFacilityFn               func(ctx context.Context, updatePayload *dto.UpdateFacilityPayload) error
	MockAddFacilityContactFn           func(ctx context.Context, facilityID string, contact string) (bool, error)
	MockSearchFacilityFn               func(ctx context.Context, searchParameter *string) ([]*domain.Facility, error)
	MockSyncFacilitiesFn               func(ctx context.Context) error
	MockCreateFacilitiesFn             func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error)
	MockPublishFacilitiesToCMSFn       func(ctx context.Context, facilities []*domain.Facility) error
	MockAddFacilityToProgramFn         func(ctx context.Context, facilityIDs []string, programID string) (bool, error)
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

		MockUpdateFacilityFn: func(ctx context.Context, updatePayload *dto.UpdateFacilityPayload) error {
			return nil
		},

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockRetrieveFacilityByIdentifierFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},

		MockGetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},

		FetchFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
			return true, nil
		},
		MockAddFacilityContactFn: func(ctx context.Context, facilityID string, contact string) (bool, error) {
			return true, nil
		},
		MockSearchFacilityFn: func(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockSyncFacilitiesFn: func(ctx context.Context) error {
			return nil
		},
		MockCreateFacilitiesFn: func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockPublishFacilitiesToCMSFn: func(ctx context.Context, facilities []*domain.Facility) error {
			return nil
		},
		MockAddFacilityToProgramFn: func(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
			return true, nil
		},
	}
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (f *FacilityUsecaseMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByIdentifier mocks the implementation of `gorm's` RetrieveFacilityByIdentifier method.
func (f *FacilityUsecaseMock) RetrieveFacilityByIdentifier(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	return f.MockRetrieveFacilityByIdentifierFn(ctx, identifier, isActive)
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
func (f *FacilityUsecaseMock) DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.MockDeleteFacilityFn(ctx, identifier)
}

// FetchFacilities mocks the implementation of fetch all facilities
func (f *FacilityUsecaseMock) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return f.FetchFacilitiesFn(ctx)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (f *FacilityUsecaseMock) InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.MockInactivateFacilityFn(ctx, identifier)
}

// ReactivateFacility mocks the implementation of reactivating the active status of a particular facility
func (f *FacilityUsecaseMock) ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	return f.MockReactivateFacilityFn(ctx, identifier)
}

// UpdateFacility mocks the implementation of updating a facility
func (f *FacilityUsecaseMock) UpdateFacility(ctx context.Context, updatePayload *dto.UpdateFacilityPayload) error {
	return f.MockUpdateFacilityFn(ctx, updatePayload)
}

// AddFacilityContact mock the implementation of the AddFacilityContact method
func (f *FacilityUsecaseMock) AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error) {
	return f.MockAddFacilityContactFn(ctx, facilityID, contact)
}

// SearchFacility mock the implementation of the SearchFacility method
func (f *FacilityUsecaseMock) SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
	return f.MockSearchFacilityFn(ctx, searchParameter)
}

// SyncFacilities mock the implementation of the SyncFacilities method
func (f *FacilityUsecaseMock) SyncFacilities(ctx context.Context) error {
	return f.MockSyncFacilitiesFn(ctx)
}

// CreateFacilities Mocks the implementation of CreateFacilities method
func (f *FacilityUsecaseMock) CreateFacilities(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
	return f.MockCreateFacilitiesFn(ctx, facilities)
}

// PublishFacilitiesToCMS mock the implementation of the PublishFacilitiesToCMS method
func (f *FacilityUsecaseMock) PublishFacilitiesToCMS(ctx context.Context, facilities []*domain.Facility) error {
	return f.MockPublishFacilitiesToCMSFn(ctx, facilities)
}

// CmdAddFacilityToProgram mock the implementation of the CmdAddFacilityToProgram method
func (f *FacilityUsecaseMock) AddFacilityToProgram(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
	return f.MockAddFacilityToProgramFn(ctx, facilityIDs, programID)
}
