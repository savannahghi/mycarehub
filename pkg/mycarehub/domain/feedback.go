package domain

import (
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// FeedbackResponse defines the field passed when sending feedback
type FeedbackResponse struct {
	UserID            string
	FeedbackType      enums.FeedbackType
	SatisfactionLevel int
	ServiceName       string
	Feedback          string
	RequiresFollowUp  bool
	PhoneNumber       string
	Facility          string
	Gender            enumutils.Gender
}
