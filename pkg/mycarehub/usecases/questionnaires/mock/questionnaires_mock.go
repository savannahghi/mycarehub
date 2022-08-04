package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// QuestionnaireUseCaseMock mocks the questionnaire instance
type QuestionnaireUseCaseMock struct {
	MockCreateScreeningToolFn func(ctx context.Context, input dto.ScreeningToolInput) (bool, error)
}

// NewServiceRequestUseCaseMock initializes a new questionnaire instance mock
func NewServiceRequestUseCaseMock() *QuestionnaireUseCaseMock {
	return &QuestionnaireUseCaseMock{
		MockCreateScreeningToolFn: func(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
			return true, nil
		},
	}
}

// CreateScreeningTool mocks the CreateScreeningTool method
func (q *QuestionnaireUseCaseMock) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	return q.MockCreateScreeningToolFn(ctx, input)
}
