package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// Caregiver maps the schema for the table that stores the caregiver
type Caregiver struct {
	ID              string `json:"id"`
	UserID          string
	CaregiverNumber string

	// Deprecated
	FirstName     string              `json:"firstName"`
	LastName      string              `json:"lastName"`
	PhoneNumber   string              `json:"phoneNumber"`
	CaregiverType enums.CaregiverType `json:"caregiverType"`
	Active        bool                `json:"active"`
}

// CaregiverProfile is the profile for a caregiver with user's name
type CaregiverProfile struct {
	ID              string `json:"id"`
	User            User   `json:"user"`
	CaregiverNumber string `json:"caregiver_number"`
}
