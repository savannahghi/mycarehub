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
}

// CaregiverProfile is the profile for a caregiver with user's name
type CaregiverProfile struct {
	ID              string `json:"id"`
	UserID          string `json:"userID"`
	User            User   `json:"user"`
	CaregiverNumber string `json:"caregiver_number"`
}

// CaregiverClient models the clients
type CaregiverClient struct {
	CaregiverID        string              `json:"caregiver_id"`
	ClientID           string              `json:"client_id"`
	Active             bool                `json:"active"`
	RelationshipType   enums.CaregiverType `json:"relationship_type"`
	CaregiverConsent   *bool               `json:"caregiver_consent"`
	CaregiverConsentAt *time.Time          `json:"caregiver_consent_at"`
	ClientConsent      *bool               `json:"client_consent"`
	ClientConsentAt    *time.Time          `json:"client_consent_at"`
	OrganisationID     string              `json:"organisation_id"`
	AssignedBy         string              `json:"assigned_by"`
}

// ManagedClient represents a client who is managed by a caregiver
type ManagedClient struct {
	ClientProfile    *ClientProfile `json:"clientProfile"`
	CaregiverConsent *bool          `json:"caregiverConsent"`
	ClientConsent    *bool          `json:"clientConsent"`
}

// ClientCaregivers is the model that holds the client's caregivers
type ClientCaregivers struct {
	Caregivers []*CaregiverProfile `json:"caregivers"`
	Consented  bool                `json:"consent"`
}
