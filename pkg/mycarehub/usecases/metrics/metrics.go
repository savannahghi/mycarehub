package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IMetricsCollect is a usecase for collecting metrics
type IMetricsCollect interface {
	CollectMetric(ctx context.Context, input *domain.Metric) (bool, error)
	MetricCollector(ctx context.Context, metricType enums.MetricType) (*domain.Metric, func())
}

// UsecaseMetrics defines the interface required to implement the metrics features
type UsecaseMetrics interface {
	IMetricsCollect
}

// UsecaseMetricsImpl represents metrics implementation
type UsecaseMetricsImpl struct {
	Create infrastructure.Create
}

// NewUsecaseMetricsImpl initializes a new appointments usecase
func NewUsecaseMetricsImpl(create infrastructure.Create) UsecaseMetrics {
	return &UsecaseMetricsImpl{
		Create: create,
	}
}

// CollectMetric processes and saves a metric
func (m UsecaseMetricsImpl) CollectMetric(ctx context.Context, input *domain.Metric) (bool, error) {
	if valid := input.Type.IsValid(); !valid {
		return false, fmt.Errorf("not a valid metric type: %s", input.Type)
	}

	input.Timestamp = time.Now()

	err := m.Create.CreateMetric(ctx, input)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to save metric: %w", err)
	}

	return true, nil
}

// MetricCollector creates a metric object using the passed in metric type
// and a function to save the object in an async non blocking go routine
// It's up to the caller to execute the saving function preferably in a defer statement
func (m UsecaseMetricsImpl) MetricCollector(ctx context.Context, metricType enums.MetricType) (*domain.Metric, func()) {
	metric := &domain.Metric{
		UserID: nil,
		Type:   metricType,
		Event:  map[string]interface{}{},
	}

	collectFunc := func() {
		go func() {
			_, err := m.CollectMetric(context.Background(), metric)
			if err != nil {
				helpers.ReportErrorToSentry(err)
			}
		}()
	}

	return metric, collectFunc
}
