package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Caregiver maps the schema for the table that stores the caregiver
type Caregiver struct {
	ID              string `json:"id"`
	UserID          string `json:"userID"`
	CaregiverNumber string `json:"caregiverNumber"`
	Active          bool   `json:"active"`
	ProgramID       string `json:"programID"`
	OrganisationID  string `json:"organisationID"`
}

// CaregiverProfile is the profile for a caregiver with user's name
type CaregiverProfile struct {
	ID              string        `json:"id"`
	UserID          string        `json:"userID"`
	User            User          `json:"user"`
	CaregiverNumber string        `json:"caregiverNumber"`
	IsClient        bool          `json:"isClient"`
	Consent         ConsentStatus `json:"consent"`
	CurrentClient   *string       `json:"currentClient"`
	CurrentFacility *string       `json:"currentFacility"`
}

// ConsentStatus is used to indicate the consent status of a caregiver
type ConsentStatus struct {
	ConsentStatus enums.ConsentState `json:"consentStatus"`
}

// CaregiverClient models the clients
type CaregiverClient struct {
	CaregiverID        string              `json:"caregiverID"`
	ClientID           string              `json:"clientID"`
	Active             bool                `json:"active"`
	RelationshipType   enums.CaregiverType `json:"relationshipType"`
	CaregiverConsent   enums.ConsentState  `json:"caregiverConsent"`
	CaregiverConsentAt *time.Time          `json:"caregiverConsentAt"`
	ClientConsent      enums.ConsentState  `json:"clientConsent"`
	ClientConsentAt    *time.Time          `json:"clientConsentAt"`
	OrganisationID     string              `json:"organisationID"`
	AssignedBy         string              `json:"assignedBy"`
	ProgramID          string              `json:"programID"`
}

// ManagedClient represents a client who is managed by a caregiver
type ManagedClient struct {
	ClientProfile      *ClientProfile     `json:"clientProfile"`
	CaregiverConsent   enums.ConsentState `json:"caregiverConsent"`
	ClientConsent      enums.ConsentState `json:"clientConsent"`
	WorkStationDetails WorkStationDetails `json:"workStationDetails"`
}

// ClientCaregivers is the model that holds the client's caregivers
type ClientCaregivers struct {
	Caregivers []*CaregiverProfile `json:"caregivers"`
}
