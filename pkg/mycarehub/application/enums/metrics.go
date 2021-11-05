package enums

import (
	"fmt"
	"io"
	"strconv"
)

// MetricType is a list of all the metrics type to be colected.
type MetricType string

// metrics type constants
const (
	EngagementMetrics      MetricType = "Engagement"
	SatisfactionMetrics    MetricType = "Satisfaction"
	UserInteractionMetrics MetricType = "UserInteraction"
	PerformanceMetrics     MetricType = "Performance"
)

// AllMetrics is a set of a  valid and known metric types.
var AllMetrics = []MetricType{
	EngagementMetrics,
	SatisfactionMetrics,
	UserInteractionMetrics,
	PerformanceMetrics,
}

// IsValid returns true if a metric is valid
func (m MetricType) IsValid() bool {
	switch m {
	case EngagementMetrics, SatisfactionMetrics, UserInteractionMetrics, PerformanceMetrics:
		return true
	}
	return false
}

func (m MetricType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a metric type.
func (m *MetricType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = MetricType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid MetricType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (m MetricType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
