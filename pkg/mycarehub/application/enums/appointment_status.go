package enums

import (
	"fmt"
	"io"
	"strconv"
)

// AppointmentStatus represents the different states an appointment can be in
type AppointmentStatus string

// Collection of the various states an appointment can be in
const (
	AppointmentStatusScheduled       AppointmentStatus = "SCHEDULED"
	AppointmentStatusRescheduled     AppointmentStatus = "RESCHEDULED"
	AppointmentStatusWaiting         AppointmentStatus = "WAITING"
	AppointmentStatusMissed          AppointmentStatus = "MISSED"
	AppointmentStatusCompleted       AppointmentStatus = "COMPLETED"
	AppointmentStatusInConsultation  AppointmentStatus = "IN_CONSULTATION"
	AppointmentStatusWalkIn          AppointmentStatus = "WALK_IN"
	AppointmentStatusCancelled       AppointmentStatus = "CANCELLED"
	AppointmentStatusNeedsReschedule AppointmentStatus = "NEEDS_RESCHEDULE"
)

// IsValid returns true if a status is among the valid defined states
func (a AppointmentStatus) IsValid() bool {
	switch a {
	case AppointmentStatusNeedsReschedule, AppointmentStatusCancelled, AppointmentStatusWalkIn, AppointmentStatusInConsultation, AppointmentStatusScheduled, AppointmentStatusRescheduled, AppointmentStatusWaiting, AppointmentStatusMissed, AppointmentStatusCompleted:
		return true
	}
	return false
}

// String returns a string value of the response
func (a AppointmentStatus) String() string {
	return string(a)
}

// UnmarshalGQL converts the supplied value to a status type.
func (a AppointmentStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	a = AppointmentStatus(str)
	if !a.IsValid() {
		return fmt.Errorf("%s is not a valid Appointment Status", str)
	}
	return nil
}

// MarshalGQL writes the status type to the supplied
func (a AppointmentStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(a.String()))
}
