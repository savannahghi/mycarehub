package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) AnswerScreeningToolQuestion(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error) {
	return r.mycarehub.ScreeningTools.AnswerScreeningToolQuestions(ctx, screeningToolResponses)
}

func (r *queryResolver) GetScreeningToolQuestions(ctx context.Context, toolType *string) ([]*domain.ScreeningToolQuestion, error) {
	return r.mycarehub.ScreeningTools.GetScreeningToolQuestions(ctx, toolType)
}
