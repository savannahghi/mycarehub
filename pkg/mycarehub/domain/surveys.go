package domain

import (
	"time"
)

// SurveyForm is contains the information about a survey form
type SurveyForm struct {
	ProjectID int    `json:"projectId"`
	XMLFormID string `json:"xmlFormId"`
	Name      string `json:"name"`
	EnketoID  string `json:"enketoId"`
}

// UserSurvey represents a user's surveys domain model
type UserSurvey struct {
	ID           string    `json:"id"`
	Active       bool      `json:"active"`
	Created      time.Time `json:"created"`
	Link         string    `json:"link"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	HasSubmitted bool      `json:"hasSubmitted"`
	UserID       string    `json:"userID"`
	Token        string    `json:"token"`
	ProjectID    int       `json:"projectID"`
	FormID       string    `json:"formID"`
	LinkID       int       `json:"linkID"`
}

// Submission represents a survey's submission domain model
type Submission struct {
	InstanceID  string    `json:"instanceId"`
	SubmitterID int       `json:"submitterId"`
	DeviceID    string    `json:"deviceId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ReviewState string    `json:"reviewState"`
	Submitter   Submitter `json:"submitter"`
}

// Submitter represents a survey's submitter domain model
type Submitter struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	DisplayName string    `json:"displayName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
}
