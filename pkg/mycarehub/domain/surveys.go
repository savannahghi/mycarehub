package domain

import (
	"time"
)

// SurveyForm is contains the information about a survey form
type SurveyForm struct {
	ProjectID int    `json:"projectID"`
	XMLFormID string `json:"xmlFormID"`
	Name      string `json:"name"`
	EnketoID  string `json:"enketoID"`
}

// UserSurvey represents a user's surveys domain model
type UserSurvey struct {
	ID             string    `json:"id"`
	Active         bool      `json:"active"`
	Created        time.Time `json:"created"`
	Link           string    `json:"link"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	HasSubmitted   bool      `json:"hasSubmitted"`
	UserID         string    `json:"userID"`
	Token          string    `json:"token"`
	ProjectID      int       `json:"projectID"`
	FormID         string    `json:"formID"`
	LinkID         int       `json:"linkID"`
	SubmittedAt    time.Time `json:"submittedAt"`
	ProgramID      string    `json:"programID"`
	OrganisationID string    `json:"organisationID"`
}

// Submission represents a survey's submission domain model
type Submission struct {
	InstanceID  string    `json:"instanceID"`
	SubmitterID int       `json:"submitterID"`
	DeviceID    string    `json:"deviceID"`
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

// SurveyRespondent represents a survey's respondent domain model
type SurveyRespondent struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	SubmittedAt time.Time `json:"submittedAt"`
	ProjectID   int       `json:"projectID"`
	SubmitterID int       `json:"submitterID"`
	FormID      string    `json:"formID"`
}

// SurveyRespondentPage represents a survey's respondent domain model
type SurveyRespondentPage struct {
	SurveyRespondents []*SurveyRespondent `json:"surveyRespondents"`
	Pagination        Pagination          `json:"pagination"`
}

// SurveyResponse represents a single survey submission
type SurveyResponse struct {
	Question     string   `json:"question"`
	QuestionType string   `json:"questionType"`
	Answer       []string `json:"answer"`
}

// SurveyServiceRequestUser is the models for a user(client) who has a survey service request
type SurveyServiceRequestUser struct {
	Name             string `json:"name"`
	FormID           string `json:"formID"`
	ProjectID        int    `json:"projectID"`
	SubmitterID      int    `json:"submitterID"`
	SurveyName       string `json:"surveyName"`
	ServiceRequestID string `json:"serviceRequestID"`
	PhoneNumber      string `json:"phoneNumber"`
}

// SurveyServiceRequestUserPage models the user's(client) survey service request page
type SurveyServiceRequestUserPage struct {
	Users      []*SurveyServiceRequestUser `json:"users"`
	Pagination Pagination                  `json:"pagination"`
}
