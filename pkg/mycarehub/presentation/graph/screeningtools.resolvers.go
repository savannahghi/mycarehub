package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) AnswerScreeningToolQuestion(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error) {
	return r.mycarehub.ScreeningTools.AnswerScreeningToolQuestions(ctx, screeningToolResponses)
}

func (r *queryResolver) GetScreeningToolQuestions(ctx context.Context, toolType *string) ([]*domain.ScreeningToolQuestion, error) {
	return r.mycarehub.ScreeningTools.GetScreeningToolQuestions(ctx, toolType)
}

func (r *queryResolver) GetAvailableScreeningToolQuestions(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error) {
	return r.mycarehub.ScreeningTools.GetAvailableScreeningToolQuestions(ctx, clientID)
}

func (r *queryResolver) GetAvailableFacilityScreeningTools(ctx context.Context, facilityID string) ([]*domain.AvailableScreeningTools, error) {
	return r.mycarehub.ScreeningTools.GetAvailableFacilityScreeningTools(ctx, facilityID)
}

func (r *queryResolver) GetAssessmentResponsesByToolType(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
	return r.mycarehub.ScreeningTools.GetAssessmentResponses(ctx, facilityID, toolType)
}

func (r *queryResolver) GetScreeningToolServiceRequestResponses(ctx context.Context, clientID *string, toolType *enums.ScreeningToolType) (*domain.ScreeningToolResponsePayload, error) {
	return r.mycarehub.ScreeningTools.GetScreeningToolServiceRequestResponses(ctx, *clientID, *toolType)
}
