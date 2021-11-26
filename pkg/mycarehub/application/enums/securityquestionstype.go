package enums

import (
	"fmt"
	"io"
	"strconv"
)

// SecurityQuestionResponseType is a list of all the security question response types.
type SecurityQuestionResponseType string

const (
	//SecurityQuestionResponseTypeText is the string type security question response
	SecurityQuestionResponseTypeText SecurityQuestionResponseType = "TEXT"
	//SecurityQuestionResponseTypeNumber is the number type security question response
	SecurityQuestionResponseTypeNumber SecurityQuestionResponseType = "NUMBER"
	//SecurityQuestionResponseTypeDate is the date type security question response
	SecurityQuestionResponseTypeDate SecurityQuestionResponseType = "DATE"
	//SecurityQuestionResponseTypeBoolean is the boolean type security question response
	SecurityQuestionResponseTypeBoolean SecurityQuestionResponseType = "BOOLEAN"
)

// AllSecurityQuestionResponseType is a set of a  valid and known security question types.
var AllSecurityQuestionResponseType = []SecurityQuestionResponseType{
	SecurityQuestionResponseTypeText,
	SecurityQuestionResponseTypeNumber,
	SecurityQuestionResponseTypeDate,
	SecurityQuestionResponseTypeBoolean,
}

// IsValid returns true if a sort is valid
func (m SecurityQuestionResponseType) IsValid() bool {
	switch m {
	case SecurityQuestionResponseTypeText,
		SecurityQuestionResponseTypeNumber,
		SecurityQuestionResponseTypeDate,
		SecurityQuestionResponseTypeBoolean:
		return true
	}
	return false
}

func (m SecurityQuestionResponseType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *SecurityQuestionResponseType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = SecurityQuestionResponseType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid SortDataType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m SecurityQuestionResponseType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
