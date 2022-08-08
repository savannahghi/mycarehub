package domain

import (
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Questionnaire defines the structure of a questionnaire
type Questionnaire struct {
	ID          string     `json:"id"`
	Active      bool       `json:"active"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

// ScreeningTool defines the structure of a screening tool that belongs to the questionnaire
type ScreeningTool struct {
	ID              string             `json:"id"`
	Active          bool               `json:"active"`
	QuestionnaireID string             `json:"questionnaireID"`
	Threshold       int                `json:"threshold"`
	ClientTypes     []enums.ClientType `json:"clientTypes"`
	Genders         []enumutils.Gender `json:"genders"`
	AgeRange        AgeRange           `json:"ageRange"`
	Questionnaire   Questionnaire      `json:"questionnaire"`
}

// Question represents a question within a questionnaire.
type Question struct {
	ID                string                          `json:"id"`
	Active            bool                            `json:"active"`
	QuestionnaireID   string                          `json:"questionnaireID"`
	Text              string                          `json:"text"`
	QuestionType      enums.QuestionType              `json:"questionType"`
	ResponseValueType enums.QuestionResponseValueType `json:"responseValue"`
	Required          bool                            `json:"required"`
	SelectMultiple    bool                            `json:"selectMultiple"`
	Sequence          int                             `json:"sequence"`
	Choices           []QuestionInputChoice           `json:"choices"`
}

// QuestionInputChoice defines the structure of choices for the Question
type QuestionInputChoice struct {
	ID         string `json:"id"`
	Active     bool   `json:"active"`
	QuestionID string `json:"questionID"`
	Choice     string `json:"choice"`
	Value      string `json:"value"`
	Score      int    `json:"score"`
}

// QuestionnaireScreeningToolResponse defines the response to the ScreeningTool question
// TODO: Rename to ScreeningToolResponse after removing old screening tool implementation
type QuestionnaireScreeningToolResponse struct {
	ID                string                                       `json:"id"`
	Active            bool                                         `json:"active"`
	ScreeningToolID   string                                       `json:"screeningToolID"`
	FacilityID        string                                       `json:"facilityID"`
	ClientID          string                                       `json:"clientID"`
	AggregateScore    int                                          `json:"aggregateScore"`
	QuestionResponses []QuestionnaireScreeningToolQuestionResponse `json:"questionResponses"`
}

// QuestionnaireScreeningToolQuestionResponse defines the structure of a screening tool question response
// TODO: Rename to ScreeningToolQuestionResponse after removing old screening tool implementation
type QuestionnaireScreeningToolQuestionResponse struct {
	ID                      string `json:"id"`
	Active                  bool   `json:"active"`
	ScreeningToolResponseID string `json:"screeningToolResponseID"`
	QuestionID              string `json:"questionID"`
	Response                string `json:"response"`
	Score                   int    `json:"score"`
}
