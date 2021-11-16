package domain

import (
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// SecurityQuestion models the security questions for the users
type SecurityQuestion struct {
	SecurityQuestionID string
	QuestionStem       string
	Description        string
	Flavour            feedlib.Flavour
	Active             bool
	ResponseType       enums.SecurityQuestionResponseType
}
