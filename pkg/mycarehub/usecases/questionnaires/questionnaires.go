package questionnaires

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IScreeningTools contains methods related to the screening tools
type IScreeningTools interface {
	CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error)
}

// UseCaseQuestionnaire contains questionnaire interfaces
type UseCaseQuestionnaire interface {
	IScreeningTools
}

// UseCaseQuestionnaireImpl represents the questionnaire implementations
type UseCaseQuestionnaireImpl struct {
	Query  infrastructure.Query
	Create infrastructure.Create
	Update infrastructure.Update
	Delete infrastructure.Delete
}

// NewUseCaseQuestionnaire is the controller function for the questionnaire usecase
func NewUseCaseQuestionnaire(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	delete infrastructure.Delete,
) *UseCaseQuestionnaireImpl {
	return &UseCaseQuestionnaireImpl{
		Query:  query,
		Create: create,
		Update: update,
		Delete: delete,
	}
}

// CreateScreeningTool creates the screening tool questionnaire
func (q *UseCaseQuestionnaireImpl) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	err := input.Questionnaire.Validate()
	if err != nil {
		return false, err
	}

	questions := []domain.Question{}
	for _, q := range input.Questionnaire.Questions {
		choices := []domain.QuestionInputChoice{}
		for _, c := range q.Choices {
			choices = append(choices, domain.QuestionInputChoice{
				Active: true,
				Choice: *c.Choice,
				Value:  c.Value,
				Score:  c.Score,
			})
		}
		questions = append(questions, domain.Question{
			Active:            true,
			Text:              q.Text,
			QuestionType:      q.QuestionType,
			ResponseValueType: q.ResponseValueType,
			Required:          q.Required,
			SelectMultiple:    q.SelectMultiple,
			Sequence:          q.Sequence,
			Choices:           choices,
		})
	}

	payload := &domain.ScreeningTool{
		Active:      true,
		Threshold:   input.Threshold,
		ClientTypes: input.ClientTypes,
		Genders:     input.Genders,
		AgeRange: domain.AgeRange{
			LowerBound: input.AgeRange.LowerBound,
			UpperBound: input.AgeRange.UpperBound,
		},
		Questionnaire: domain.Questionnaire{
			Active:      true,
			Name:        input.Questionnaire.Name,
			Description: input.Questionnaire.Description,
			Questions:   questions,
		},
	}

	err = q.Create.CreateScreeningTool(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to create screening tool questionnaire: %w", err)
	}
	return true, nil
}
