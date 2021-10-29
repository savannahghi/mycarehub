package mock

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// GormMock struct implements mocks of `gorm's`internal methods.
//
// This mock struct should be separate from our own internal methods.
type GormMock struct {
	GetOrCreateFacilityFn       func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	RetrieveFacilityFn          func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	RetrieveFacilityByMFLCodeFn func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error)
	GetFacilitiesFn             func(ctx context.Context) ([]gorm.Facility, error)
	DeleteFacilityFn            func(ctx context.Context, mfl_code string) (bool, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
func NewGormMock() *GormMock {
	return &GormMock{

		GetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &id,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},

		RetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},
		GetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			var facilities []gorm.Facility
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			facilities = append(facilities, gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			})
			return facilities, nil
		},

		DeleteFacilityFn: func(ctx context.Context, mfl_code string) (bool, error) {
			return true, nil
		},

		RetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *GormMock) GetOrCreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.GetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
	return gm.RetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.GetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode string) (bool, error) {
	return gm.DeleteFacilityFn(ctx, mflcode)
}
