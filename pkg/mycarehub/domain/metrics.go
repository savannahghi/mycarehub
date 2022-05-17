package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Metric is a recording of an event that occurs within the platform
type Metric struct {
	ID        int                    `json:"id"`
	UserID    *string                `json:"userID"`
	Type      enums.MetricType       `json:"type"`
	Event     map[string]interface{} `json:"event"`
	Timestamp time.Time              `json:"timestamp"`
}
