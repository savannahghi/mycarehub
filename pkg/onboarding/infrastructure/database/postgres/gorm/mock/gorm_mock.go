package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

// GormMock struct implements mocks of `gorm's`internal methods.
//
// This mock struct should be separate from our own internal methods.
type GormMock struct {
	CreateFacilityFn   func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	RetrieveFacilityFn func(ctx context.Context, id *uuid.UUID) (*gorm.Facility, error)
	GetFacilitiesFn    func(ctx context.Context) ([]gorm.Facility, error)
	DeleteFacilityFn   func(ctx context.Context, mfl_code string) (bool, error)
	CollectMetricsFn   func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error)
	FindFacilityFn     func(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
func NewGormMock() *GormMock {
	return &GormMock{
		CreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			id := uuid.New()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &id,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},

		RetrieveFacilityFn: func(ctx context.Context, id *uuid.UUID) (*gorm.Facility, error) {
			facilityID := uuid.New()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			}, nil
		},
		GetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			var facilities []gorm.Facility
			facilityID := uuid.New()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			facilities = append(facilities, gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      true,
				County:      county,
				Description: description,
			})
			return facilities, nil
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
		DeleteFacilityFn: func(ctx context.Context, mfl_code string) (bool, error) {
			return true, nil
		},

		CollectMetricsFn: func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
			now := time.Now()
			metricID := uuid.New()
			return &gorm.Metric{
				MetricID:  &metricID,
				Type:      domain.EngagementMetrics,
				Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
				Timestamp: now,
				UID:       ksuid.New().String(),
			}, nil
		},
	}
}

// CreateFacility mocks the implementation of `gorm's` CreateFacility method.
func (gm *GormMock) CreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.CreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *uuid.UUID) (*gorm.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.GetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode string) (bool, error) {
	return gm.DeleteFacilityFn(ctx, mflcode)
}

// CollectMetrics mocks the implementation of  CollectMetrics method.
func (gm *GormMock) CollectMetrics(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
	return gm.CollectMetricsFn(ctx, metrics)
}

// FindFacility mocks the implementation of  FindFacility method.
func (gm *GormMock) FindFacility(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error) {
	return gm.FindFacilityFn(ctx, pagination, filter, sort)
}
