package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
)

// Appointment represents a single appointment
type Appointment struct {
	ID              string                  `json:"ID"`
	AppointmentUUID string                  `json:"appointmentUUID"`
	Type            string                  `json:"type"`
	Status          enums.AppointmentStatus `json:"status"`
	Reason          string                  `json:"reason"`
	Date            scalarutils.Date        `json:"date"`
	Start           time.Time               `json:"start"`
	End             time.Time               `json:"end"`

	ClientID   string
	FacilityID string
	Provider   string
}

//AppointmentsPage is a list of paginated appointments
type AppointmentsPage struct {
	Appointments []*Appointment `json:"appointments"`
	Pagination   Pagination     `json:"pagination"`
}

// AppointmentServiceRequests is a list of appointment service requests
type AppointmentServiceRequests struct {
	ID              string                  `json:"id"`
	AppointmentUUID string                  `json:"appointmentUUID"`
	Type            string                  `json:"AppointmentType"`
	Status          enums.AppointmentStatus `json:"Status"`
	Reason          string                  `json:"AppointmentReason"`
	Date            scalarutils.Date        `json:"SuggestedDate"`
	SuggestedTime   string                  `json:"SuggestedTime"`
	Provider        string                  `json:"Provider"`

	InProgressAt  *time.Time `json:"InProgressAt"`
	InProgressBy  *string    `json:"InProgressBy"`
	ResolvedAt    *time.Time `json:"ResolvedAt"`
	ResolvedBy    *string    `json:"ResolvedBy"`
	ClientName    *string    `json:"ClientName"`
	ClientContact *string    `json:"ClientContact"`
	CCCNumber     string     `json:"cccNumber"`
	MFLCODE       string     `json:"MFLCODE"`
}
