package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"gopkg.in/go-playground/validator.v9"
)

// Questionnaire defines the structure of a questionnaire
type Questionnaire struct {
	ID             string     `json:"id"`
	Active         bool       `json:"active"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Questions      []Question `json:"questions"`
	ProgramID      string     `json:"programID"`
	OrganisationID string     `json:"organisationID"`
}

// GetQuestionByID returns a question by ID
func (q Questionnaire) GetQuestionByID(id string) (Question, error) {
	for _, q := range q.Questions {
		if q.ID == id {
			return q, nil
		}
	}
	return Question{}, fmt.Errorf("question not found")
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
	ProgramID       string             `json:"programID"`
	OrganisationID  string             `json:"organisationID"`
}

// GetQuestion returns the question details for a given screening tool question
func (s ScreeningTool) GetQuestion(questionID string) *Question {
	question, err := s.Questionnaire.GetQuestionByID(questionID)
	if err != nil {
		return nil
	}
	return &question
}

// GetNormalizedResponse returns the human readable response for a given question response
func (s ScreeningTool) GetNormalizedResponse(questionID, response string) map[string]interface{} {
	question, err := s.Questionnaire.GetQuestionByID(questionID)
	if err != nil {
		return nil
	}
	switch question.QuestionType {
	case enums.QuestionTypeCloseEnded:
		if question.SelectMultiple {
			return question.GetNormalizedResponseForMultipleChoice(response)
		}
		return question.GetNormalizedResponseForSingleChoice(response)
	case enums.QuestionTypeOpenEnded:
		return map[string]interface{}{
			"0": response,
		}
	}
	return nil
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
	ProgramID         string                          `json:"programID"`
	OrganisationID    string                          `json:"organisationID"`
}

// ValidateResponse helps with validation of a question response input
func (s Question) ValidateResponse(response string) error {
	v := validator.New()
	err := v.Struct(s)
	switch s.ResponseValueType {
	case enums.QuestionResponseValueTypeNumber:
		_, err := strconv.Atoi(response)
		if err != nil {
			return fmt.Errorf("response value must be a number")
		}
	case enums.QuestionResponseValueTypeBoolean:
		if _, err := strconv.ParseBool(response); err != nil {
			return fmt.Errorf("response value must be a boolean")
		}
	}

	choicesMap := make(map[string]bool)
	for _, c := range s.Choices {
		choicesMap[c.Choice] = true
	}

	multiChoiceMap := make(map[string]string)
	for _, c := range strings.Split(response, ",") {
		if c == "" {
			continue
		}
		multiChoiceMap[c] = c
	}

	switch s.QuestionType {
	case enums.QuestionTypeCloseEnded:
		if s.SelectMultiple {
			for _, c := range multiChoiceMap {
				if !choicesMap[c] {
					return fmt.Errorf("response value must be one of the choices")
				}
			}
		} else {
			if !choicesMap[response] {
				return fmt.Errorf("response value must be one of the choices")
			}
		}

	}

	if s.Required && response == "" {
		return fmt.Errorf("response is required")
	}

	return err
}

// GetScore returns the score for a given question response
func (s Question) GetScore(response string) int {
	switch s.QuestionType {
	case enums.QuestionTypeCloseEnded:
		if s.SelectMultiple {
			return s.GetScoreForMultipleChoice(response)
		}
		return s.GetScoreForSingleChoice(response)
	}
	return 0
}

// GetScoreForSingleChoice returns the score for a single choice question response
func (s Question) GetScoreForSingleChoice(response string) int {
	for _, c := range s.Choices {
		if c.Choice == response {
			return c.Score
		}
	}
	return 0
}

// GetScoreForMultipleChoice returns the score for a multiple choice question response
func (s Question) GetScoreForMultipleChoice(response string) int {
	var score int
	for _, c := range s.Choices {
		if strings.Contains(response, c.Choice) {
			score += c.Score
		}
	}
	return score
}

// GetNormalizedResponseForSingleChoice returns the human readable response for a single choice question response
func (s Question) GetNormalizedResponseForSingleChoice(response string) map[string]interface{} {
	for _, c := range s.Choices {
		if c.Choice == response {
			return map[string]interface{}{
				c.Choice: c.Value,
			}
		}
	}
	return nil
}

// GetNormalizedResponseForMultipleChoice returns the human readable response for a multiple choice question response
func (s Question) GetNormalizedResponseForMultipleChoice(response string) map[string]interface{} {
	choices := map[string]interface{}{}
	for _, c := range s.Choices {
		if strings.Contains(response, c.Choice) {
			choices[c.Choice] = c.Value
		}
	}
	return choices
}

// QuestionInputChoice defines the structure of choices for the Question
type QuestionInputChoice struct {
	ID             string `json:"id"`
	Active         bool   `json:"active"`
	QuestionID     string `json:"questionID"`
	Choice         string `json:"choice"`
	Value          string `json:"value"`
	Score          int    `json:"score"`
	ProgramID      string `json:"programID"`
	OrganisationID string `json:"organisationID"`
}

// QuestionnaireScreeningToolResponse defines the response to the ScreeningTool question
// TODO: Rename to ScreeningToolResponse after removing old screening tool implementation
type QuestionnaireScreeningToolResponse struct {
	ID                string                                        `json:"id"`
	Active            bool                                          `json:"active"`
	ScreeningToolID   string                                        `json:"screeningToolID"`
	FacilityID        string                                        `json:"facilityID"`
	ClientID          string                                        `json:"clientID"`
	DateOfResponse    time.Time                                     `json:"dateOfResponse"`
	AggregateScore    int                                           `json:"aggregateScore"`
	QuestionResponses []*QuestionnaireScreeningToolQuestionResponse `json:"questionResponses"`
	ProgramID         string                                        `json:"programID"`
	OrganisationID    string                                        `json:"organisationID"`
	CaregiverID       *string                                       `json:"caregiverID"`
}

// QuestionnaireScreeningToolQuestionResponse defines the structure of a screening tool question response
// TODO: Rename to ScreeningToolQuestionResponse after removing old screening tool implementation
type QuestionnaireScreeningToolQuestionResponse struct {
	ID                      string                          `json:"id"`
	Active                  bool                            `json:"active"`
	ScreeningToolResponseID string                          `json:"screeningToolResponseID"`
	QuestionID              string                          `json:"questionID"`
	QuestionType            enums.QuestionType              `json:"questionType"`
	SelectMultiple          bool                            `json:"selectMultiple"`
	ResponseValueType       enums.QuestionResponseValueType `json:"responseValueType"`
	Sequence                int                             `json:"sequence"`
	QuestionText            string                          `json:"questionText"`
	Response                string                          `json:"response"`
	NormalizedResponse      map[string]interface{}          `json:"normalizedResponse"`
	Score                   int                             `json:"score"`
	ProgramID               string                          `json:"programID"`
	OrganisationID          string                          `json:"organisationID"`
	FacilityID              string                          `json:"facilityID"`
}

// ScreeningToolRespondent defines the structure of a screening tool respondent
// the respondent is the client who has a service request for the screening tool
type ScreeningToolRespondent struct {
	ClientID                string `json:"clientID"`
	ScreeningToolResponseID string `json:"screeningToolResponseID"`
	ServiceRequestID        string `json:"serviceRequestID"`
	Name                    string `json:"name"`
	PhoneNumber             string `json:"phoneNumber"`
	ServiceRequest          string `json:"serviceRequest"`
}

// ScreeningToolPage represents the screening tool pagination mode
type ScreeningToolPage struct {
	ScreeningTools []*ScreeningTool `json:"screeningTools"`
	Pagination     Pagination       `json:"pagination"`
}

// ScreeningToolRespondentsPage represents a list of respondents
// Pagination is included
type ScreeningToolRespondentsPage struct {
	ScreeningToolRespondents []*ScreeningToolRespondent `json:"screeningToolRespondents"`
	Pagination               Pagination                 `json:"pagination"`
}
