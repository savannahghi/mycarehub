package enums

import (
	"fmt"
	"io"
	"strconv"
)

// FilterDataType defines the various Filter data types
type FilterDataType string

const (
	// FilterDataTypeName represents a Name Filter data type
	FilterDataTypeName FilterDataType = "Name"

	// FilterDataTypeMFLCode represents an MFL Code Filter data type
	FilterDataTypeMFLCode FilterDataType = "Code"

	// FilterDataTypeActive represents the Active Filter data type
	FilterDataTypeActive FilterDataType = "Active"

	// FilterDataTypeCounty represents the County Filter data type
	FilterDataTypeCounty FilterDataType = "County"

	// Other Filter data Types
)

// FacilityFilterDataTypes represents a slice of all possible `FilterDataTypes` values
var FacilityFilterDataTypes = []FilterDataType{
	FilterDataTypeName,
	FilterDataTypeMFLCode,
	FilterDataTypeActive,
	FilterDataTypeCounty,
}

// IsValid returns true if an Filter data type is valid
func (e FilterDataType) IsValid() bool {
	switch e {
	case FilterDataTypeName,
		FilterDataTypeMFLCode,
		FilterDataTypeActive,
		FilterDataTypeCounty:
		return true
	}
	return false
}

// String ...
func (e FilterDataType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a filter data type.
func (e *FilterDataType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterDataType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterDataType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e FilterDataType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
