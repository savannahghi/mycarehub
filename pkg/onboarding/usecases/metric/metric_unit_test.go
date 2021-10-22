package metric_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

func TestUsecaseMetricsImpl_CollectMetrics_Unittest(t *testing.T) {
	ctx := context.Background()

	f := testFakeInfrastructureInteractor

	metric := &dto.MetricInput{
		Type:      enums.EngagementMetrics,
		Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
		Timestamp: time.Now(),
		UID:       ksuid.New().String(),
	}

	invalidMetric := &dto.MetricInput{
		Type:      "",
		Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "bored"}`)),
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
			if tt.name == "Sad case" {
				fakeCreate.CollectMetricsFn = func(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
					return nil, fmt.Errorf("an error occurred while collecting metrics.")
				}
			}
			if tt.name == "Happy case" {
				fakeCreate.CollectMetricsFn = func(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
					metricID := uuid.New().String()
					return &domain.Metric{
						MetricID:  &metricID,
						Type:      enums.EngagementMetrics,
						Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "bored"}`)),
						Timestamp: time.Now(),
						UID:       ksuid.New().String(),
					}, nil
				}
			}

			_, err := f.CollectMetrics(tt.args.ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseMetricsImpl.CollectMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
