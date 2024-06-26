package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Booking is the booking class data model
type Booking struct {
	ID                     string                  `json:"id"`
	Active                 bool                    `json:"active"`
	Services               []string                `json:"services"`
	Date                   time.Time               `json:"date"`
	Facility               Facility                `json:"facility"`
	Client                 ClientProfile           `json:"clientProfile"`
	OrganisationID         string                  `json:"organisationID"`
	ProgramID              string                  `json:"programID"`
	VerificationCode       string                  `json:"verificationCode"`
	VerificationCodeStatus enums.BookingCodeStatus `json:"verificationCodeStatus"`
	BookingStatus          enums.BookingStatus     `json:"bookingStatus"`
}
