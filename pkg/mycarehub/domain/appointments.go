package domain

import (
	"time"

	"github.com/savannahghi/scalarutils"
)

// Appointment represents a single appointment
type Appointment struct {
	ID     string           `json:"ID"`
	Type   string           `json:"type"`
	Status string           `json:"status"`
	Reason string           `json:"reason"`
	Date   scalarutils.Date `json:"date"`
	Start  time.Time        `json:"start"`
	End    time.Time        `json:"end"`

	ClientID   string
	FacilityID string
	Provider   string
}

//AppointmentsPage is a list of paginated appointments
type AppointmentsPage struct {
	Appointments []*Appointment `json:"appointments"`
	Pagination   Pagination     `json:"pagination"`
}
