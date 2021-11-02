package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockGetOrCreateFacilityFn         func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	MockRetrieveFacilityFn            func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByMFLCodeFn   func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error)
	MockGetFacilitiesFn               func(ctx context.Context) ([]gorm.Facility, error)
	MockDeleteFacilityFn              func(ctx context.Context, mfl_code string) (bool, error)
	MockRegisterClientFn              func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error)
	MockGetUserProfileByPhoneNumberFn func(ctx context.Context, phoneNumber string) (*gorm.User, error)
	MockSavePinFn                     func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	MockGetUserPINByUserIDFn          func(ctx context.Context, userID string) (*gorm.PINData, error)
	MockListFacilitiesFn              func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
//
// This initialization initializes all the good cases of your mock tests. i.e all success cases should be defined here.
func NewGormMock() *GormMock {

	/*
		In this section, you find commonly shared success case structs for mock tests
	*/

	ID := uuid.New().String()
	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
		Active:      strconv.FormatBool(true),
		County:      county,
		Description: description,
	}

	var facilities []gorm.Facility
	facilities = append(facilities, *facility)

	clientProfile := &gorm.ClientUserProfile{
		User: &gorm.User{
			FirstName:   gofakeit.FirstName(),
			LastName:    gofakeit.LastName(),
			Username:    gofakeit.Username(),
			MiddleName:  gofakeit.Name(),
			DisplayName: gofakeit.BeerAlcohol(),
			Gender:      enumutils.GenderMale,
		},
		Client: &gorm.ClientProfile{
			ClientType: enums.ClientTypeOvc,
		},
	}

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
		MockSavePinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},
		MockGetUserProfileByPhoneNumberFn: func(ctx context.Context, phoneNumber string) (*gorm.User, error) {
			ID := uuid.New().String()
			return &gorm.User{
				UserID: &ID,
			}, nil
		},

		MockRegisterClientFn: func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
			return clientProfile, nil
		},

		MockGetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			return facility, nil
		},

		MockRetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {

			return facility, nil
		},
		MockGetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			return facilities, nil
		},

		MockDeleteFacilityFn: func(ctx context.Context, mfl_code string) (bool, error) {
			return true, nil
		},

		MockRetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
			return facility, nil
		},
		MockListFacilitiesFn: func(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
			return facilitiesPage, nil
		},
		MockGetUserPINByUserIDFn: func(ctx context.Context, userID string) (*gorm.PINData, error) {
			return pinData, nil
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
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
	return gm.MockRetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.MockGetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode string) (bool, error) {
	return gm.MockDeleteFacilityFn(ctx, mflcode)
}

// RegisterClient mocks the implementation of RegisterClient method
func (gm *GormMock) RegisterClient(
	ctx context.Context,
	userInput *gorm.User,
	clientInput *gorm.ClientProfile,
) (*gorm.ClientUserProfile, error) {
	return gm.MockRegisterClientFn(ctx, userInput, clientInput)
}

// GetUserProfileByPhoneNumber mocks the implementation of retrieving a user profile by phonenumber
func (gm *GormMock) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*gorm.User, error) {
	return gm.MockGetUserProfileByPhoneNumberFn(ctx, phoneNumber)
}

// SavePin mocks the implementation of saving the pin to the database
func (gm *GormMock) SavePin(ctx context.Context, pinData *gorm.PINData) (bool, error) {
	return gm.MockSavePinFn(ctx, pinData)
}

// ListFacilities mocks the implementation of  ListFacilities method.
func (gm *GormMock) ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination *domain.FacilityPage) (*domain.FacilityPage, error) {
	return gm.MockListFacilitiesFn(ctx, searchTerm, filter, pagination)
}

// GetUserPINByUserID mocks the implementation of retrieving a user pin by user ID
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string) (*gorm.PINData, error) {
	return gm.MockGetUserPINByUserIDFn(ctx, userID)
}
