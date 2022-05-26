package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ScreeningToolResponseType represents the type of screening tool response
type ScreeningToolResponseType string

const (

	// ScreeningToolResponseTypeInteger denotes the value of the 'INTEGER' enum as a ScreeningToolResponseType
	ScreeningToolResponseTypeInteger ScreeningToolResponseType = "INTEGER"

	// ScreeningToolResponseTypeText denotes the value of the 'TEXT' enum as a ScreeningToolResponseType
	ScreeningToolResponseTypeText ScreeningToolResponseType = "TEXT"

	// ScreeningToolResponseTypeDate denotes the value of the 'DATE' enum as a ScreeningToolResponseType
	ScreeningToolResponseTypeDate ScreeningToolResponseType = "DATE"
)

// IsValid returns true if a screening tool response type is valid
func (m ScreeningToolResponseType) IsValid() bool {
	switch m {
	case ScreeningToolResponseTypeInteger,
		ScreeningToolResponseTypeText,
		ScreeningToolResponseTypeDate:
		return true
	}
	return false
}

// String converts screening tool response type type to string
func (m ScreeningToolResponseType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *ScreeningToolResponseType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ScreeningToolResponseType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid ScreeningToolResponseType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m ScreeningToolResponseType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
