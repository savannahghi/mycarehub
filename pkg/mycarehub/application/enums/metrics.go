package enums

import (
	"fmt"
	"io"
	"strconv"
)

// MetricType is a type of metric
type MetricType string

const (
	// MetricTypeContent is a metric for content
	MetricTypeContent MetricType = "CONTENT"

	// MetricTypeEngagement is a metric for content
	MetricTypeEngagement MetricType = "ENGAGEMENT"

	// MetricTypeSystem is a metric for content
	MetricTypeSystem MetricType = "SYSTEM"
)

// IsValid returns true if a metric type is valid
func (m MetricType) IsValid() bool {
	switch m {
	case MetricTypeContent, MetricTypeEngagement, MetricTypeSystem:
		return true
	}
	return false
}

// String converts the metric type enum to a string
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
		return fmt.Errorf("%s is not a valid metric type", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (m MetricType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
