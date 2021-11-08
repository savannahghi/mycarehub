package domain

import (
	"fmt"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// FiltersParam contains the inputs for filter parameters
type FiltersParam struct {
	Name     string
	DataType enums.FilterSortDataType
	Value    string // TODO: Clear spec on validation e.g dates must be ISO 8601. This is the actual data being filtered
}

// Validate is a filter param method that performs validations
func (f FiltersParam) Validate() error {
	if f.DataType == enums.FilterSortDataTypeName {
		if f.Value == "" {
			return fmt.Errorf("name cannot be empty")
		}
	}
	if f.DataType == enums.FilterSortDataTypeMFLCode {
		if f.Value == "" {
			return fmt.Errorf("MFL code cannot be empty")
		}
	}
	if f.DataType == enums.FilterSortDataTypeActive {
		_, err := strconv.ParseBool(f.Value)
		if err != nil {
			return fmt.Errorf("failed to convert to bool %v: %v", f.Value, err)
		}
	}
	if f.DataType == enums.FilterSortDataTypeCounty {
		ok := enums.CountyType(f.Value).IsValid()
		if !ok {
			return fmt.Errorf("invalid county passed: %v", f.Value)
		}
	}
	// Validate enums
	// TODO: Very strict validation of data <-> data type
	// 	     this is a good candidate for TDD with unit tests
	// TODO: make sure this is always called before filter params are used
	return nil
}

// SortParam includes the fields required for sorting the different types of fields
type SortParam struct {
	Field     enums.FilterSortDataType
	Direction enums.SortDataType
}
