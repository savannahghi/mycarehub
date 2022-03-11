package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// ScreeningToolQuestion defines a question for a screening tool
type ScreeningToolQuestion struct {
	ID               string                              `json:"questionID"`
	Question         string                              `json:"question"`
	ToolType         enums.ScreeningToolType             `json:"toolType"`
	ResponseChoices  map[string]interface{}              `json:"responseChoices"`
	ResponseType     enums.ScreeningToolResponseType     `json:"responseType"`
	ResponseCategory enums.ScreeningToolResponseCategory `json:"responseCategory"`
	Sequence         int                                 `json:"sequence"`
	Meta             map[string]interface{}              `json:"meta"`
	Active           bool                                `json:"active"`
}

// ValidateResponseQUestionCategory validates the response by category
func (q *ScreeningToolQuestion) ValidateResponseQUestionCategory(response string, category enums.ScreeningToolResponseCategory) error {
	switch category {
	case enums.ScreeningToolResponseCategorySingleChoice:
		if q.ResponseChoices == nil {
			return fmt.Errorf("response choices is nil")
		}
		_, ok := q.ResponseChoices[response]
		if !ok {
			return fmt.Errorf("invalid response: %s", response)
		}
	case enums.ScreeningToolResponseCategoryMultiChoice:
		if q.ResponseChoices == nil {
			return fmt.Errorf("response choices is nil")
		}
		responses := strings.Split(response, ",")

		for _, responseChoice := range responses {
			_, ok := q.ResponseChoices[string(responseChoice)]
			if !ok {
				return fmt.Errorf("invalid response: %s", response)
			}
		}
	case enums.ScreeningToolResponseCategoryOpenEnded:
		// no validation
	default:
		return fmt.Errorf("invalid response category: %s", category)
	}
	return nil
}

// ValidateResponseQUestionType validates the response by type
func (q *ScreeningToolQuestion) ValidateResponseQUestionType(response string, responseType enums.ScreeningToolResponseType) error {
	switch responseType {
	case enums.ScreeningToolResponseTypeInteger:
		_, err := strconv.ParseInt(response, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid response: %s", response)
		}
	case enums.ScreeningToolResponseTypeDate:
		_, err := time.Parse("02-01-2006", response)
		if err != nil {
			return fmt.Errorf("invalid response: %s", response)
		}
	case enums.ScreeningToolResponseTypeText:
		// no validation
	default:
		return fmt.Errorf("invalid response type: %s", responseType)
	}
	return nil
}

// ScreeningToolQuestionResponse defines a response for a screening tool
type ScreeningToolQuestionResponse struct {
	ID         string `json:"id"`
	QuestionID string `json:"questionID"`
	ClientID   string `json:"clientID"`
	Answer     string `json:"answer"`
	Active     bool   `json:"active"`
}
