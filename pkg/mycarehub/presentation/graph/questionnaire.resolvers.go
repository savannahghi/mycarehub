package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	return r.mycarehub.Questionnaires.CreateScreeningTool(ctx, input)
}

func (r *mutationResolver) RespondToScreeningTool(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error) {
	return r.mycarehub.Questionnaires.RespondToScreeningTool(ctx, input)
}

func (r *queryResolver) GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
	return r.mycarehub.Questionnaires.GetAvailableScreeningTools(ctx, clientID, facilityID)
}
