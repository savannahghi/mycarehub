package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// MetricsUseCaseMock mocks the implementation of Metrics usecase methods.
type MetricsUseCaseMock struct {
	MockCollectMetricFn func(ctx context.Context, input domain.Metric) (bool, error)
}

// NewMetricsUseCaseMock creates in itializes create type mocks
func NewMetricsUseCaseMock() *MetricsUseCaseMock {
	return &MetricsUseCaseMock{
		MockCollectMetricFn: func(ctx context.Context, input domain.Metric) (bool, error) {
			return true, nil
		},
	}
}

// CollectMetric mocks the implementation of collecting metrics.
func (gm *MetricsUseCaseMock) CollectMetric(ctx context.Context, input domain.Metric) (bool, error) {
	return gm.MockCollectMetricFn(ctx, input)

}
