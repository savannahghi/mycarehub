package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ServiceRequestType is a list of all the service request types.
type ServiceRequestType string

const (
	//ServiceRequestTypeHealthDiaryEntry represents a health diary entry
	ServiceRequestTypeHealthDiaryEntry ServiceRequestType = "HEALTH_DIARY_ENTRY"
	//ServiceRequestTypeRedFlag represents a health diary entry
	ServiceRequestTypeRedFlag ServiceRequestType = "RED_FLAG"
)

// AllServiceRequestType is a set of a  valid and known service request types.
var AllServiceRequestType = []ServiceRequestType{
	ServiceRequestTypeHealthDiaryEntry,
	ServiceRequestTypeRedFlag,
}

// IsValid returns true if a request type is valid
func (m ServiceRequestType) IsValid() bool {
	switch m {
	case ServiceRequestTypeHealthDiaryEntry,
		ServiceRequestTypeRedFlag:
		return true
	}
	return false
}

func (m ServiceRequestType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a request type.
func (m *ServiceRequestType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ServiceRequestType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid SortDataType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m ServiceRequestType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
