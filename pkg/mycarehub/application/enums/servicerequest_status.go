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
	//ServiceRequestStatusRejected is the string type service request
	ServiceRequestStatusRejected ServiceRequestStatus = "REJECTED"
)

// IsValid returns true if a request type is valid
func (m ServiceRequestStatus) IsValid() bool {
	switch m {
	case ServiceRequestStatusPending,
		ServiceRequestStatusInProgress,
		ServiceRequestStatusResolved,
		ServiceRequestStatusRejected:
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

// PINResetVerificationStatus  status for verifying pin reset service request
type PINResetVerificationStatus string

// Valid PINResetVerificationStatus
const (
	PINResetVerificationStatusApproved PINResetVerificationStatus = "APPROVED"
	PINResetVerificationStatusRejected PINResetVerificationStatus = "REJECTED"
)

// IsValid checks if the PINResetVerificationStatus is valid
func (e PINResetVerificationStatus) IsValid() bool {
	switch e {
	case PINResetVerificationStatusApproved, PINResetVerificationStatusRejected:
		return true
	}
	return false
}

// String returns the string
func (e PINResetVerificationStatus) String() string {
	return string(e)
}

// UnmarshalGQL converts the input, if valid, into an PINResetVerificationStatus value
func (e *PINResetVerificationStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PINResetVerificationStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PINResetVerificationStatus", str)
	}
	return nil
}

// MarshalGQL converts PINResetVerificationStatus into a valid JSON string
func (e PINResetVerificationStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
