package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// PostgresMock struct implements mocks of `postgres's` internal methods.
type PostgresMock struct {
	//Get
	MockGetOrCreateFacilityFn         func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
	MockGetFacilitiesFn               func(ctx context.Context) ([]*domain.Facility, error)
	MockRetrieveFacilityFn            func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error)
	ListFacilitiesFn                  func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error)
	MockDeleteFacilityFn              func(ctx context.Context, id string) (bool, error)
	MockRetrieveFacilityByMFLCodeFn   func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
	MockGetUserProfileByPhoneNumberFn func(ctx context.Context, phoneNumber string) (*domain.User, error)
	MockGetUserPINByUserIDFn          func(ctx context.Context, userID string) (*domain.UserPIN, error)
	MockInactivateFacilityFn          func(ctx context.Context, mflCode *string) (bool, error)
}

// NewPostgresMock initializes a new instance of `GormMock` then mocking the case of success.
func NewPostgresMock() *PostgresMock {
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

	return &PostgresMock{
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
			return facilityInput, nil
		},
		MockGetFacilitiesFn: func(ctx context.Context) ([]*domain.Facility, error) {
			return facilitiesList, nil
		},
		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},
		ListFacilitiesFn: func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},
		MockDeleteFacilityFn: func(ctx context.Context, id string) (bool, error) {
			return true, nil
		},
		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
			return facilityInput, nil
		},
		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string) (*domain.UserPIN, error) {
			return &domain.UserPIN{
				ValidTo: time.Now().Add(time.Hour * 10),
			}, nil
		},
		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*domain.User, error) {
			id := uuid.New().String()
			return &domain.User{
				ID: &id,
			}, nil
		},
		MockInactivateFacilityFn: func(ctx context.Context, mflCode *string) (bool, error) {
			return true, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *PostgresMock) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return gm.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *PostgresMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (gm *PostgresMock) ListFacilities(
	ctx context.Context,
	searchTerm *string,
	filterInput []*dto.FiltersInput,
	paginationsInput *dto.PaginationsInput,
) (*domain.FacilityPage, error) {
	return gm.ListFacilitiesFn(ctx, searchTerm, filterInput, paginationsInput)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method
func (gm *PostgresMock) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return gm.MockGetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of deleting a facility by ID
func (gm *PostgresMock) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, id)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacilityByMFLCode method.
func (gm *PostgresMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return gm.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetUserProfileByPhoneNumber mocks the implementation of fetching a user profile by phonenumber
func (gm *PostgresMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// GetUserPINByUserID mocks the get user pin by ID implementation
func (gm *PostgresMock) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (gm *PostgresMock) InactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	return gm.MockInactivateFacilityFn(ctx, mflCode)
}
