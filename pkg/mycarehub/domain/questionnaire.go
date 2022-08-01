package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Questionnaire is the structure of a questionnaire response
type Questionnaire struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	StartDate   time.Time               `json:"start_date"`
	EndDate     time.Time               `json:"end_date"`
	Active      bool                    `json:"is_open"`
	Questions   []QuestionnaireQuestion `json:"questions"`
}

// QuestionnaireQuestion is the  of a questionnaire question
type QuestionnaireQuestion struct {
	ID             string                                 `json:"id"`
	Text           string                                 `json:"text"`
	SelectMultiple bool                                   `json:"select_multiple"`  // nullable, used by close-ended questions
	ValueType      string                                 `json:"question_type"`    // enum => string, integer, date, time, datetime, BooleanField
	QuestionType   enums.QuestionnaireQuestionTypeChoices `json:"question_type_id"` // enum => open-ended, closed-ended
	Sequence       int64                                  `json:"sequence"`
	HasCondition   bool                                   `json:"has_condition"` // if has condition and question type == close-ended
	Choices        []QuestionnaireQuestionChoice          `json:"choices"`       //nullable for open ended questions, not null && len >= 2 for closed questions. used by closed-ended questions
}

// QuestionnaireQuestionChoice is the  of a questionnaire question choice
type QuestionnaireQuestionChoice struct {
	ID         string `json:"id"`
	Choice     string `json:"choice"` // the actual choice to be selected
	Value      string `json:"value"`  // the value of the choice
	Score      int64  `json:"score"`  // the score of the choice. this will come into play when comparing the responses with the threshold
	QuestionID int64  `json:"question_id"`
}
