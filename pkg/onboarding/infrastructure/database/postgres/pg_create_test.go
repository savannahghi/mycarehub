package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

func TestOnboardingDb_CreateFacility(t *testing.T) {
	ctx := context.Background()
	name := "Kanairo One"
	code := "KN001"
	county := "Kanairo"
	description := "This is just for mocking"
	type args struct {
		ctx      context.Context
		facility *dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)
			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if tt.name == "sad case - facility code not defined" {
				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return facility, nil
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facility to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facility not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_CollectMetrics_Unittest(t *testing.T) {
	ctx := context.Background()

	metric := &dto.MetricInput{
		Type:      domain.EngagementMetrics,
		Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
		Timestamp: time.Now(),
		UID:       ksuid.New().String(),
	}

	invalidMetric := &dto.MetricInput{
		Type:      "",
		Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
		Timestamp: time.Now(),
		UID:       ksuid.New().String(),
	}

	type args struct {
		ctx    context.Context
		metric *dto.MetricInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				metric: metric,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				metric: invalidMetric,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Happy case" {
				fakeGorm.CollectMetricsFn = func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
					now := time.Now()
					metricID := uuid.New()
					return &gorm.Metric{
						MetricID:  &metricID,
						Type:      domain.EngagementMetrics,
						Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
						Timestamp: now,
						UID:       ksuid.New().String(),
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeGorm.CollectMetricsFn = func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
					return nil, fmt.Errorf("an error occurred while collecting metrics")
				}
			}

			_, err := d.CollectMetrics(tt.args.ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.CollectMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
