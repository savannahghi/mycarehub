package screeningtools

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetScreeningToolsQuestion represents the interface to get screening tools questions
type IGetScreeningToolsQuestion interface {
	GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error)
}

// IAnswerScreeningToolQuestion represents the interface to answer screening tools questions
type IAnswerScreeningToolQuestion interface {
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error)
}

// UseCasesScreeningTools represents the usecases for screening tools
type UseCasesScreeningTools interface {
	IGetScreeningToolsQuestion
	IAnswerScreeningToolQuestion
}

// ServiceScreeningToolsImpl represents screening tools implementation object
type ServiceScreeningToolsImpl struct {
	Query  infrastructure.Query
	Create infrastructure.Create
	Update infrastructure.Update
}

// NewUseCasesScreeningTools is the controller for the screening tools usecases
func NewUseCasesScreeningTools(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
) *ServiceScreeningToolsImpl {
	return &ServiceScreeningToolsImpl{
		Query:  query,
		Create: create,
		Update: update,
	}
}

// GetScreeningToolQuestions get all the screening tools questions
func (t *ServiceScreeningToolsImpl) GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error) {
	emptyString := ""
	if questionType == nil {
		questionType = &emptyString
	} else {
		ok := enums.ScreeningToolType(*questionType).IsValid()
		if !ok {
			return nil, fmt.Errorf("invalid question type: %s in input", *questionType)
		}
	}

	screeningToolsQuestions, err := t.Query.GetScreeningToolQuestions(ctx, *questionType)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ItemNotFoundErr(fmt.Errorf("failed to get screening tools questions: %v", err))
	}

	return screeningToolsQuestions, nil
}

// AnswerScreeningToolQuestions answer screening tools questions
func (t *ServiceScreeningToolsImpl) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error) {
	if len(screeningToolResponses) == 0 {
		return false, fmt.Errorf("no screening tool responses provided")
	}
	for _, screeningToolResponse := range screeningToolResponses {
		err := screeningToolResponse.Validate()
		if err != nil {
			return false, fmt.Errorf("invalid screening tool response: %v", err)
		}

		screeningToolQuestion, err := t.Query.GetScreeningToolQuestionByQuestionID(ctx, screeningToolResponse.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to get screening tool question: %v", err)
		}

		err = screeningToolQuestion.ValidateResponseQUestionCategory(screeningToolResponse.Response, screeningToolQuestion.ResponseCategory)
		if err != nil {
			return false, fmt.Errorf("invalid screening tool response: %v", err)
		}

		err = screeningToolQuestion.ValidateResponseQUestionType(screeningToolResponse.Response, screeningToolQuestion.ResponseType)
		if err != nil {
			return false, fmt.Errorf("invalid screening tool response: %v", err)
		}
		err = t.Update.InvalidateScreeningToolResponse(ctx, screeningToolResponse.ClientID, screeningToolResponse.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to invalidate previous screening tool response: %v", err)
		}
	}
	err := t.Create.AnswerScreeningToolQuestions(ctx, screeningToolResponses)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to answer screening tools questions: %v", err)
	}
	return true, nil
}
