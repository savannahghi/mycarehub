package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockGetOrCreateFacilityFn         func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	MockRetrieveFacilityFn            func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByMFLCodeFn   func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error)
	MockGetFacilitiesFn               func(ctx context.Context) ([]gorm.Facility, error)
	MockDeleteFacilityFn              func(ctx context.Context, mflCode int) (bool, error)
	MockListFacilitiesFn              func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
	MockGetUserProfileByPhoneNumberFn func(ctx context.Context, phoneNumber string) (*gorm.User, error)
	MockGetUserPINByUserIDFn          func(ctx context.Context, userID string) (*gorm.PINData, error)
	MockInactivateFacilityFn          func(ctx context.Context, mflCode *int) (bool, error)
	MockReactivateFacilityFn          func(ctx context.Context, mflCode *int) (bool, error)
	MockGetUserProfileByUserIDFn      func(ctx context.Context, UserID string) (*gorm.User, error)
	MockSaveTemporaryUserPinFn        func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetCurrentTermsFn             func(ctx context.Context) (*gorm.TermsOfService, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
//
// This initialization initializes all the good cases of your mock tests. i.e all success cases should be defined here.
func NewGormMock() *GormMock {

	/*
		In this section, you find commonly shared success case structs for mock tests
	*/

	ID := gofakeit.Number(300, 400)
	UUID := ksuid.New().String()
	name := gofakeit.Name()
	code := gofakeit.Number(0, 100)
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &UUID,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	var facilities []gorm.Facility
	facilities = append(facilities, *facility)

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
				ID:          &UUID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			},
		},
	}

	pinData := &gorm.PINData{
		PINDataID: &ID,
		UserID:    gofakeit.UUID(),
		HashedPIN: uuid.New().String(),
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   feedlib.FlavourConsumer,
	}

	return &GormMock{
		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			return facility, nil
		},

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {

			return facility, nil
		},
		MockGetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			return facilities, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, mflCode int) (bool, error) {
			return true, nil
		},

		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
			return facility, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},

		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},

		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string) (*gorm.PINData, error) {
			return pinData, nil
		},

		MockInactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockReactivateFacilityFn: func(ctx context.Context, mflCode *int) (bool, error) {
			return true, nil
		},
		MockGetCurrentTermsFn: func(ctx context.Context) (*gorm.TermsOfService, error) {
			termsID := ksuid.New().String()
			validFrom := time.Now()
			testText := "test"

			validTo := time.Now().AddDate(0, 0, 80)
			terms := &gorm.TermsOfService{
				TermsID:   &termsID,
				Text:      &testText,
				ValidFrom: &validFrom,
				ValidTo:   &validTo,
			}
			return terms, nil
		},
		MockGetUserProfileByUserIDFn: func(ctx context.Context, UserID string) (*gorm.User, error) {
			return &gorm.User{
				UserID: &UserID,
			}, nil
		},
		MockSaveTemporaryUserPinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *GormMock) GetOrCreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.MockGetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.MockGetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode int) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, mflcode)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (gm *GormMock) ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
	return gm.MockListFacilitiesFn(ctx, searchTerm, filter, pagination)
}

// GetUserProfileByPhoneNumber mocks the implementation of retrieving a user profile by phonenumber
func (gm *GormMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*gorm.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// GetUserPINByUserID mocks the implementation of retrieving a user pin by user ID
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string) (*gorm.PINData, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID)
}

// InactivateFacility mocks the implementation of inactivating the active status of a particular facility
func (gm *GormMock) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return gm.MockInactivateFacilityFn(ctx, mflCode)
}

// ReactivateFacility mocks the implementation of re-activating the active status of a particular facility
func (gm *GormMock) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	return gm.MockReactivateFacilityFn(ctx, mflCode)
}

//GetCurrentTerms mocks the implementation of getting all the current terms of service.
func (gm *GormMock) GetCurrentTerms(ctx context.Context) (*gorm.TermsOfService, error) {
	return gm.MockGetCurrentTermsFn(ctx)
}

// GetUserProfileByUserID mocks the implementation of retrieving a user profile by user ID
func (gm *GormMock) GetUserProfileByUserID(ctx context.Context, UserID string) (*gorm.User, error) {
	return gm.MockGetUserProfileByUserIDFn(ctx, UserID)
}

// SaveTemporaryUserPin mocks the implementation of saving a temporary user pin
func (gm *GormMock) SaveTemporaryUserPin(ctx context.Context, pinData *gorm.PINData) (bool, error) {
	return gm.MockSaveTemporaryUserPinFn(ctx, pinData)
}
