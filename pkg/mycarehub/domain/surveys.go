package domain

import (
	"net/http"
	"time"
)

// SurveysClient defines the fields required to access the surveys client
type SurveysClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// RequestHelperPayload is the payload that is sent to the surveys client
type RequestHelperPayload struct {
	Method string
	Path   string
	Body   interface{}
}

// SurveyForm is contains the information about a survey form
type SurveyForm struct {
	ProjectID int    `json:"projectId"`
	Name      string `json:"name"`
}

// UserSurveys represents a user's surveys domain model
type UserSurveys struct {
	ID           string    `json:"id"`
	Active       bool      `json:"active"`
	Created      time.Time `json:"created"`
	Link         string    `json:"link"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	HasSubmitted bool      `json:"hasSubmitted"`
	UserID       string    `json:"userID"`
}
