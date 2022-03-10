package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ScreeningToolType represents the type of screening tool question
type ScreeningToolType string

const (

	// ScreeningToolTypeTB denotes the value of the 'TB_ASSESSMENT' enum as a ScreeningToolType
	ScreeningToolTypeTB ScreeningToolType = "TB_ASSESSMENT"

	// ScreeningToolTypeGBV denotes the value of the 'VIOLENCE_ASSESSMENT' enum as a ScreeningToolType
	ScreeningToolTypeGBV ScreeningToolType = "VIOLENCE_ASSESSMENT"

	// ScreeningToolTypeCUI denotes the value of the 'CONTRACEPTIVE_ASSESSMENT' enum as a ScreeningToolType
	ScreeningToolTypeCUI ScreeningToolType = "CONTRACEPTIVE_ASSESSMENT"

	// ScreeningToolTypeAlcoholSubstanceAssessment denotes the value of the 'ALCOHOL_SUBSTANCE_ASSESSMENT' enum as a ScreeningToolType
	ScreeningToolTypeAlcoholSubstanceAssessment ScreeningToolType = "ALCOHOL_SUBSTANCE_ASSESSMENT"
)

//ScreeningToolQuestions is a set of a  valid and known screening tools question types.
var ScreeningToolQuestions = []ScreeningToolType{
	ScreeningToolTypeTB,
	ScreeningToolTypeGBV,
	ScreeningToolTypeCUI,
	ScreeningToolTypeAlcoholSubstanceAssessment,
}

// IsValid returns true if a screening tool question type is valid
func (m ScreeningToolType) IsValid() bool {
	switch m {
	case ScreeningToolTypeTB,
		ScreeningToolTypeGBV,
		ScreeningToolTypeCUI,

		ScreeningToolTypeAlcoholSubstanceAssessment:
		return true
	}
	return false
}

// String converts screening tool question type type to string
func (m ScreeningToolType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *ScreeningToolType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ScreeningToolType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid ScreeningToolType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m ScreeningToolType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
