package enums

import (
	"fmt"
	"io"
	"strconv"
)

//QuestionnaireQuestionTypeChoices is a list of QuestionnaireQuestionType choices
type QuestionnaireQuestionTypeChoices string

const (
	//OpenEnded is are questionnaire questions that are open ended
	OpenEnded QuestionnaireQuestionTypeChoices = "OPEN_ENDED"

	//CloseEnded is are questionnaire questions that are closed ended
	CloseEnded QuestionnaireQuestionTypeChoices = "CLOSE_ENDED"
)

// IsValid returns true if a QuestionnaireQuestionType is valid
func (q QuestionnaireQuestionTypeChoices) IsValid() bool {
	switch q {
	case OpenEnded, CloseEnded:
		return true
	}
	return false
}

// String converts the QuestionnaireQuestionType to a string
func (q QuestionnaireQuestionTypeChoices) String() string {
	return string(q)
}

// UnmarshalGQL converts the supplied value to a QuestionnaireQuestionType
func (q *QuestionnaireQuestionTypeChoices) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*q = QuestionnaireQuestionTypeChoices(str)
	if !q.IsValid() {
		return fmt.Errorf("%s is not a valid QuestionnaireQuestionType", str)
	}
	return nil
}

// MarshalGQL writes the QuestionnaireQuestionType to the supplied writer
func (q QuestionnaireQuestionTypeChoices) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}

// QuestionnaireResponseValueChoices is a custom type that defines the type of response values a questionnaire question can have
type QuestionnaireResponseValueChoices string

const (
	// StringResponseValue is a string response value
	StringResponseValue QuestionnaireResponseValueChoices = "STRING"
	// NumberResponseValue is a number response value
	NumberResponseValue QuestionnaireResponseValueChoices = "NUMBER"
	// BooleanResponseValue is a boolean response value
	BooleanResponseValue QuestionnaireResponseValueChoices = "BOOLEAN"
	// TimeResponseValue is a time response value
	TimeResponseValue QuestionnaireResponseValueChoices = "TIME"
	// DateResponseValue is a date response value
	DateResponseValue QuestionnaireResponseValueChoices = "DATE"
	// DateTimeResponseValue is a datetime response value
	DateTimeResponseValue QuestionnaireResponseValueChoices = "DATE_TIME"
)

// IsValid returns true if a QuestionnaireResponseValue is valid
func (q QuestionnaireResponseValueChoices) IsValid() bool {
	switch q {
	case StringResponseValue, NumberResponseValue, BooleanResponseValue, TimeResponseValue, DateResponseValue, DateTimeResponseValue:
		return true
	}
	return false
}

// String converts the QuestionnaireResponseValue to a string
func (q QuestionnaireResponseValueChoices) String() string {
	return string(q)
}

// UnmarshalGQL converts the supplied value to a QuestionnaireResponseValue
func (q *QuestionnaireResponseValueChoices) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*q = QuestionnaireResponseValueChoices(str)
	if !q.IsValid() {
		return fmt.Errorf("%s is not a valid QuestionnaireResponseValue", str)
	}
	return nil
}

// MarshalGQL writes the QuestionnaireResponseValue to the supplied writer
func (q QuestionnaireResponseValueChoices) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(q.String()))
}
