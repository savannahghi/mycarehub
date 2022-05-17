package metrics

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IMetricsCollect is a usecase for collecting metrics
type IMetricsCollect interface {
	CollectMetric(ctx context.Context, input domain.Metric) (bool, error)
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
func (m UsecaseMetricsImpl) CollectMetric(ctx context.Context, input domain.Metric) (bool, error) {
	if valid := input.Type.IsValid(); !valid {
		return false, fmt.Errorf("not a valid metric type: %s", input.Type)
	}

	//TODO: metric event validation

	err := m.Create.CreateMetric(ctx, &input)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to save metric: %w", err)
	}

	return true, nil
}
