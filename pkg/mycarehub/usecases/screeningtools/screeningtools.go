package screeningtools

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetScreeningToolsQuestion represents the interface to get screening tools questions
type IGetScreeningToolsQuestion interface {
	GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error)
}

// UseCasesScreeningTools represents the usecases for screening tools
type UseCasesScreeningTools interface {
	IGetScreeningToolsQuestion
}

// ServiceScreeningToolsImpl represents screening tools implementation object
type ServiceScreeningToolsImpl struct {
	Query infrastructure.Query
}

// NewUseCasesScreeningTools is the controller for the screening tools usecases
func NewUseCasesScreeningTools(
	query infrastructure.Query,
) *ServiceScreeningToolsImpl {
	return &ServiceScreeningToolsImpl{
		Query: query,
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
