package mock

import (
	"context"
	"strconv"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// GormMock struct implements mocks of `gorm's`internal methods.
type GormMock struct {
	MockGetOrCreateFacilityFn       func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	MockRetrieveFacilityFn          func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	MockRetrieveFacilityByMFLCodeFn func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error)
	MockGetFacilitiesFn             func(ctx context.Context) ([]gorm.Facility, error)
	MockDeleteFacilityFn            func(ctx context.Context, mfl_code string) (bool, error)
	MockRegisterClientFn            func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
//
// This initialization initializes all the good cases of your mock tests. i.e all success cases should be defined here.
func NewGormMock() *GormMock {

	/*
		In this section, you find commonly shared success case structs for mock tests
	*/

	facilityID := uuid.New().String()
	name := gofakeit.Name()
	code := "KN001"
	county := "Kanairo"
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &facilityID,
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

	return &GormMock{
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
