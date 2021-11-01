package enums

import (
	"fmt"
	"io"
	"strconv"
)

// FilterCategoryType defines the various Filter category types
type FilterCategoryType string

const (
	// FilterCategoryTypeFacility represents a Facility Filter category type
	FilterCategoryTypeFacility FilterCategoryType = "Facility"
	// Other Filter category Types
)

// FacilityFilterCategoryTypes represents a slice of all possible `FilterCategoryTypes` values
var FacilityFilterCategoryTypes = []FilterCategoryType{
	FilterCategoryTypeFacility,
}

// IsValid returns true if an Filter category type is valid
func (e FilterCategoryType) IsValid() bool {
	switch e {
	case FilterCategoryTypeFacility:
		return true
	}
	return false
}

// String ...
func (e FilterCategoryType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a filter category type.
func (e *FilterCategoryType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterCategoryType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterCategoryType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e FilterCategoryType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
