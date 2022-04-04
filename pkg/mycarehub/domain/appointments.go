package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// Appointment represents a single appointment
type Appointment struct {
	ID         string           `json:"ID"`
	ExternalID string           `json:"externalID"`
	Reason     string           `json:"reason"`
	Date       scalarutils.Date `json:"date"`

	ClientID                  string
	FacilityID                string
	Provider                  string
	HasRescheduledAppointment bool `json:"hasRescheduledApointment"`
}

//AppointmentsPage is a list of paginated appointments
type AppointmentsPage struct {
	Appointments []*Appointment `json:"appointments"`
	Pagination   Pagination     `json:"pagination"`
}

// AppointmentServiceRequests is a list of appointment service requests
type AppointmentServiceRequests struct {
	ID         string           `json:"id"`
	ExternalID string           `json:"AppointmentID"`
	Reason     string           `json:"AppointmentReason"`
	Date       scalarutils.Date `json:"SuggestedDate"`

	Status        string     `json:"Status"`
	InProgressAt  *time.Time `json:"InProgressAt"`
	InProgressBy  *string    `json:"InProgressBy"`
	ResolvedAt    *time.Time `json:"ResolvedAt"`
	ResolvedBy    *string    `json:"ResolvedBy"`
	ClientName    *string    `json:"ClientName"`
	ClientContact *string    `json:"ClientContact"`
	CCCNumber     string     `json:"CCCNumber"`
	MFLCODE       string     `json:"MFLCODE"`
}
