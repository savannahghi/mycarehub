package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ServiceRequestStatus is a list of all the service request status.
type ServiceRequestStatus string

const (
	//ServiceRequestStatusPending is the string type service request entry
	ServiceRequestStatusPending ServiceRequestStatus = "PENDING"
	//ServiceRequestStatusInProgress is the string type service request
	ServiceRequestStatusInProgress ServiceRequestStatus = "IN PROGRESS"
	//ServiceRequestStatusResolved is the string type service request
	ServiceRequestStatusResolved ServiceRequestStatus = "RESOLVED"
)

// AllServiceRequestStatus is a set of a  valid and known service request status.
var AllServiceRequestStatus = []ServiceRequestStatus{
	ServiceRequestStatusPending,
	ServiceRequestStatusInProgress,
	ServiceRequestStatusResolved,
}

// IsValid returns true if a request type is valid
func (m ServiceRequestStatus) IsValid() bool {
	switch m {
	case ServiceRequestStatusPending,
		ServiceRequestStatusInProgress,
		ServiceRequestStatusResolved:
		return true
	}
	return false
}

func (m ServiceRequestStatus) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a request status.
func (m *ServiceRequestStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ServiceRequestStatus(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid SortDataType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m ServiceRequestStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
