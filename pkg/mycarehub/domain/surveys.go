package domain

import (
	"net/http"
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
