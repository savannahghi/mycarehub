package domain

import (
	"fmt"
	"strconv"
	"time"

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

// Validate validates the security question response type
func (s *SecurityQuestion) Validate(response string) error {
	if s.ResponseType == enums.SecurityQuestionResponseTypeString {
		return nil
	}
	if s.ResponseType == enums.SecurityQuestionResponseTypeNumber {
		_, err := strconv.ParseInt(response, 10, 64)
		if err != nil {
			return fmt.Errorf("response type number %v is invalid: %v", response, err)
		}
	}
	if s.ResponseType == enums.SecurityQuestionResponseTypeDate {
		// the date format is DD-MM-YYYY
		_, err := time.Parse("02-01-2006", response)
		if err != nil {
			return fmt.Errorf("response type date %v is invalid: %v", response, err)
		}
	}
	return nil
}

// RecordSecurityQuestionResponse models the response to a security question
type RecordSecurityQuestionResponse struct {
	SecurityQuestionID string
	IsCorrect          bool
}

// SecurityQuestionResponse models the data that is expected from the security question response table
type SecurityQuestionResponse struct {
	ResponseID string `json:"id"`
	QuestionID string `json:"questionID"`
	UserID     string `json:"userID"`
	Active     bool   `json:"active"`
	Response   string `json:"response"`
}
