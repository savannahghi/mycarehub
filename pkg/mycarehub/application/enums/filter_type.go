package enums

import (
	"fmt"
	"io"
	"strconv"
)

// FilterSortDataType defines the various Filter data types
type FilterSortDataType string

// Note: the constant values should match the table field name
const (
	// FilterSortDataTypeCreatedAt represents created at Filter data type
	FilterSortDataTypeCreatedAt FilterSortDataType = "created"

	// FilterSortDataTypeUpdatedAt represents  updated at Filter data type
	FilterSortDataTypeUpdatedAt FilterSortDataType = "updated"

	// FilterSortDataTypeName represents a Name Filter data type
	FilterSortDataTypeName FilterSortDataType = "name"

	// FilterSortDataTypeMFLCode represents an MFL Code Filter data type
	FilterSortDataTypeMFLCode FilterSortDataType = "mfl_code"

	// FilterSortDataTypeActive represents the Active Filter data type
	FilterSortDataTypeActive FilterSortDataType = "active"

	// FilterSortDataTypeCounty represents the County Filter data type
	FilterSortDataTypeCounty FilterSortDataType = "county"

	// Other Filter data Types
)

// FacilityFilterDataTypes represents a slice of all possible `FilterDataTypes` values
var FacilityFilterDataTypes = []FilterSortDataType{
	FilterSortDataTypeName,
	FilterSortDataTypeMFLCode,
	FilterSortDataTypeActive,
	FilterSortDataTypeCounty,
}

// FacilitySortDataTypes represents a slice of all possible `SortDataTypes` values
var FacilitySortDataTypes = []FilterSortDataType{
	FilterSortDataTypeCreatedAt,
	FilterSortDataTypeUpdatedAt,
	FilterSortDataTypeName,
	FilterSortDataTypeMFLCode,
	FilterSortDataTypeActive,
	FilterSortDataTypeCounty,
}

// IsValid returns true if an Filter data type is valid
func (e FilterSortDataType) IsValid() bool {
	switch e {
	case FilterSortDataTypeCreatedAt,
		FilterSortDataTypeUpdatedAt,
		FilterSortDataTypeName,
		FilterSortDataTypeMFLCode,
		FilterSortDataTypeActive,
		FilterSortDataTypeCounty:
		return true
	}
	return false
}

// String ...
func (e FilterSortDataType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a filter data type.
func (e *FilterSortDataType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterSortDataType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterSortDataType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e FilterSortDataType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
