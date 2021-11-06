package enums

import (
	"fmt"
	"io"
	"strconv"
)

// SortDataType is a list of all the sort types.
type SortDataType string

const (
	// SortDataTypeAsc is used to request ascending sort
	SortDataTypeAsc SortDataType = "asc"

	// SortDataTypeDesc is used to request descending sort
	SortDataTypeDesc SortDataType = "desc"
)

// AllSorts is a set of a  valid and known sort types.
var AllSorts = []SortDataType{
	SortDataTypeAsc,
	SortDataTypeDesc,
}

// IsValid returns true if a sort is valid
func (m SortDataType) IsValid() bool {
	switch m {
	case SortDataTypeAsc, SortDataTypeDesc:
		return true
	}
	return false
}

func (m SortDataType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *SortDataType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = SortDataType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid SortDataType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m SortDataType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
