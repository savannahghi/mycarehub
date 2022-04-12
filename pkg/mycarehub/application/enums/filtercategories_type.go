package enums

import (
	"fmt"
	"io"
	"strconv"
)

// FilterSortCategoryType defines the various Filter category types
type FilterSortCategoryType string

const (
	// FilterSortCategoryTypeFacility represents a Facility Filter category type
	FilterSortCategoryTypeFacility FilterSortCategoryType = "Facility"

	// FilterSortCategoryTypeSortFacility represents a Facility Sort category type
	FilterSortCategoryTypeSortFacility FilterSortCategoryType = "SortFacility"
	// Other Filter category Types
)

// FacilityFilterCategoryTypes represents a slice of all possible `FilterCategoryTypes` values
var FacilityFilterCategoryTypes = []FilterSortCategoryType{
	FilterSortCategoryTypeFacility,
	FilterSortCategoryTypeSortFacility,
}

// IsValid returns true if an Filter category type is valid
func (e FilterSortCategoryType) IsValid() bool {
	switch e {
	case FilterSortCategoryTypeFacility,
		FilterSortCategoryTypeSortFacility:
		return true
	}
	return false
}

// String ...
func (e FilterSortCategoryType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a filter category type.
func (e *FilterSortCategoryType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterSortCategoryType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterSortCategoryType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e FilterSortCategoryType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
