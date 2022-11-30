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
}

// CaregiverProfile is the profile for a caregiver with user's name
type CaregiverProfile struct {
	ID              string        `json:"id"`
	UserID          string        `json:"userID"`
	User            User          `json:"user"`
	CaregiverNumber string        `json:"caregiver_number"`
	IsClient        bool          `json:"is_client"`
	Consent         ConsentStatus `json:"consent"`
	ProgramID       string        `json:"program_id"`
}

// ConsentStatus is used to indicate the consent status of a caregiver
type ConsentStatus struct {
	ConsentStatus enums.ConsentState `json:"consentStatus"`
}

// CaregiverClient models the clients
type CaregiverClient struct {
	CaregiverID        string              `json:"caregiver_id"`
	ClientID           string              `json:"client_id"`
	Active             bool                `json:"active"`
	RelationshipType   enums.CaregiverType `json:"relationship_type"`
	CaregiverConsent   enums.ConsentState  `json:"caregiver_consent"`
	CaregiverConsentAt *time.Time          `json:"caregiver_consent_at"`
	ClientConsent      enums.ConsentState  `json:"client_consent"`
	ClientConsentAt    *time.Time          `json:"client_consent_at"`
	OrganisationID     string              `json:"organisation_id"`
	AssignedBy         string              `json:"assigned_by"`
	ProgramID          string              `json:"program_id"`
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
