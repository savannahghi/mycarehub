package domain

import (
	"fmt"
	"io"
	"log"
	"strconv"
)

// AllMetrics is a set of a  valid and known metric types.
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
	_, err := fmt.Fprint(w, strconv.Quote(m.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}

// FacilityIdentifiers a list of valid facility filters
type FacilityIdentifiers string

// facility filter constants
const (
	//MFLCode Identifies a facility by a unique code
	MFLCode FacilityIdentifiers = "MFL_CODE"
	Active  FacilityIdentifiers = "ACTIVE"
	County  FacilityIdentifiers = "COUNTY"
)

var ( // Identifies a facility by status active

)

// AllFacilityIdentifiers is the set of valid facility filters
var AllFacilityIdentifiers = []FacilityIdentifiers{
	MFLCode,
	Active,
	County,
}

// IsValid returns true if a facility field param is valid
func (e FacilityIdentifiers) IsValid() bool {
	switch e {
	case MFLCode, Active, County:
		return true
	}
	return false
}

func (e FacilityIdentifiers) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a facility field param
func (e *FacilityIdentifiers) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FacilityIdentifiers(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FacilityIdentifiers", str)
	}
	return nil
}

// MarshalGQL writes the facility field param
func (e FacilityIdentifiers) MarshalGQL(w io.Writer) {
	_, err := fmt.Fprint(w, strconv.Quote(e.String()))
	if err != nil {
		log.Printf("%v\n", err)
	}
}
