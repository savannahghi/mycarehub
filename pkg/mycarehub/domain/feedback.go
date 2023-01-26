package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// FeedbackResponse defines the field passed when sending feedback
type FeedbackResponse struct {
	UserID            string             `json:"userID"`
	FeedbackType      enums.FeedbackType `json:"feedbackType"`
	SatisfactionLevel int                `json:"satisfactionLevel"`
	ServiceName       string             `json:"serviceName"`
	Feedback          string             `json:"feedback"`
	RequiresFollowUp  bool               `json:"requiresFollowUp"`
	PhoneNumber       string             `json:"phoneNumber"`
	ProgramID         string             `json:"programID"`
	OrganisationID    string             `json:"organisationID"`
}
