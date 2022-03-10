package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ScreeningToolResponseCategory represents the type of screening tool response
type ScreeningToolResponseCategory string

const (

	// ScreeningToolResponseCategorySingleChoice denotes the value of the 'SINGLE_CHOICE' enum as a ScreeningToolResponseCategory
	ScreeningToolResponseCategorySingleChoice ScreeningToolResponseCategory = "SINGLE_CHOICE"
	// ScreeningToolResponseCategoryMultiChoice denotes the value of the 'MULTI_CHOICE' enum as a ScreeningToolResponseCategory
	ScreeningToolResponseCategoryMultiChoice ScreeningToolResponseCategory = "MULTI_CHOICE"
	// ScreeningToolResponseCategoryOpenEnded denotes the value of the 'OPEN_ENDED' enum as a ScreeningToolResponseCategory
	ScreeningToolResponseCategoryOpenEnded ScreeningToolResponseCategory = "OPEN_ENDED"
)

//ScreeningToolResponseCategories is a set of a  valid and known screening tools response types.
var ScreeningToolResponseCategories = []ScreeningToolResponseCategory{
	ScreeningToolResponseCategorySingleChoice,
	ScreeningToolResponseCategoryMultiChoice,
	ScreeningToolResponseCategoryOpenEnded,
}

// IsValid returns true if a screening tool response type is valid
func (m ScreeningToolResponseCategory) IsValid() bool {
	switch m {
	case ScreeningToolResponseCategorySingleChoice,
		ScreeningToolResponseCategoryMultiChoice,
		ScreeningToolResponseCategoryOpenEnded:
		return true
	}
	return false
}

// String converts screening tool response type type to string
func (m ScreeningToolResponseCategory) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *ScreeningToolResponseCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ScreeningToolResponseCategory(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid ScreeningToolResponseCategory", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m ScreeningToolResponseCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
