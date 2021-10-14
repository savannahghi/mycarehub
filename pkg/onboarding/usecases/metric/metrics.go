package metric

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

// UsecasesMetrics contains a list of all interfaces that collect metrics.
type UsecasesMetrics interface {
	IMetricsCollect
}

// IMetricsCollect contains metrics collection methods
type IMetricsCollect interface {
	// TODO Check identifier; ensure idempotency
	CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error)
}

// UsecaseMetricsImpl represents the metrics implementation
type UsecaseMetricsImpl struct {
	Infrastructure infrastructure.Interactor
	//Hashing        utils.IHashing
}

// NewMetricUsecase returns a new metric usecase
func NewMetricUsecase(infra infrastructure.Interactor) *UsecaseMetricsImpl {
	return &UsecaseMetricsImpl{
		Infrastructure: infra,
	}
}

// CollectMetrics collects the various metrics.
func (m *UsecaseMetricsImpl) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
	//If the UID of the user is provided, hash it
	if metric.UID != "" {
		// Encrypting UID
		salt, encryptedUID := utils.EncryptUID(metric.UID, nil)

		// Comparing Salted UID
		matched := utils.CompareUID(metric.UID, salt, encryptedUID, nil)
		if !matched {
			return nil, fmt.Errorf("the provided UID do not match")

		}

		metricData := &dto.MetricInput{
			Type:      metric.Type,
			Payload:   metric.Payload,
			Timestamp: metric.Timestamp,
			UID:       encryptedUID,
		}

		return m.Infrastructure.CollectMetrics(ctx, metricData)
	}

	return m.Infrastructure.CollectMetrics(ctx, metric)
}
