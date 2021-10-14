package metric_test

import (
	"context"
	"testing"
	"time"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
	"gorm.io/datatypes"
)

func TestUsecaseMetricsImpl_CollectMetrics_Integration_Test(t *testing.T) {
	ctx := context.Background()

	m := testInfrastructureInteractor

	metric := &dto.MetricInput{
		Type:      domain.EngagementMetrics,
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

	// Encrypt UID
	salt, encryptedData := utils.EncryptUID("1234", nil)
	assert.NotNil(t, encryptedData)
	assert.NotNil(t, salt)

	// Check whether the UID has been encrypted
	isEncypted := utils.CompareUID("1234", salt, encryptedData, nil)
	assert.True(t, isEncypted)
	assert.Equal(t, true, isEncypted)

	// Collect valid metrics
	collectedMetric, err := m.CollectMetrics(ctx, metric)
	assert.Nil(t, err)
	assert.NotNil(t, collectedMetric)

	// Collect invalid metrics
	invalidmetric, err := m.CollectMetrics(ctx, invalidMetric)
	assert.NotNil(t, err)
	assert.Nil(t, invalidmetric)

}
