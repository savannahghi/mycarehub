package enums

import (
	"fmt"
	"io"
	"strconv"
)

//QuestionType is a list of QuestionType choices
type QuestionType string

const (
	//QuestionTypeOpenEnded is are questionnaire questions that are open ended
	QuestionTypeOpenEnded QuestionType = "OPEN_ENDED"

	//QuestionTypeCloseEnded is are questionnaire questions that are closed ended
	QuestionTypeCloseEnded QuestionType = "CLOSE_ENDED"
)

// IsValid returns true if a QuestionType is valid
func (q QuestionType) IsValid() bool {
	switch q {
	case QuestionTypeOpenEnded, QuestionTypeCloseEnded:
		return true
	}
	return false
}

// String converts the QuestionType to a string
func (q QuestionType) String() string {
	return string(q)
}

// UnmarshalGQL converts the supplied value to a QuestionType
func (q *QuestionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*q = QuestionType(str)
	if !q.IsValid() {
		return fmt.Errorf("%s is not a valid QuestionType", str)
	}
	return nil
}

// MarshalGQL writes the QuestionType to the supplied writer
func (q QuestionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}

// QuestionResponseValueType is a custom type that defines the type of response values a questionnaire question can have
type QuestionResponseValueType string

const (
	// QuestionResponseValueTypeString is a string response value
	QuestionResponseValueTypeString QuestionResponseValueType = "STRING"
	// QuestionResponseValueTypeNumber is a number response value
	QuestionResponseValueTypeNumber QuestionResponseValueType = "NUMBER"
	// QuestionResponseValueTypeBoolean is a boolean response value
	QuestionResponseValueTypeBoolean QuestionResponseValueType = "BOOLEAN"
	// QuestionResponseValueTypeTime is a time response value
	QuestionResponseValueTypeTime QuestionResponseValueType = "TIME"
	// QuestionResponseValueTypeDate is a date response value
	QuestionResponseValueTypeDate QuestionResponseValueType = "DATE"
	// QuestionResponseValueTypeDateTime is a datetime response value
	QuestionResponseValueTypeDateTime QuestionResponseValueType = "DATE_TIME"
)

// IsValid returns true if a QuestionnaireResponseValue is valid
func (q QuestionResponseValueType) IsValid() bool {
	switch q {
	case QuestionResponseValueTypeString, QuestionResponseValueTypeNumber, QuestionResponseValueTypeBoolean, QuestionResponseValueTypeTime, QuestionResponseValueTypeDate, QuestionResponseValueTypeDateTime:
		return true
	}
	return false
}

// String converts the QuestionnaireResponseValue to a string
func (q QuestionResponseValueType) String() string {
	return string(q)
}

// UnmarshalGQL converts the supplied value to a QuestionnaireResponseValue
func (q *QuestionResponseValueType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*q = QuestionResponseValueType(str)
	if !q.IsValid() {
		return fmt.Errorf("%s is not a valid QuestionnaireResponseValue", str)
	}
	return nil
}

// MarshalGQL writes the QuestionnaireResponseValue to the supplied writer
func (q QuestionResponseValueType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}
