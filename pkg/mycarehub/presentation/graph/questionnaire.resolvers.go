package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateScreeningTool is the resolver for the createScreeningTool field.
func (r *mutationResolver) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	return r.mycarehub.Questionnaires.CreateScreeningTool(ctx, input)
}

// RespondToScreeningTool is the resolver for the respondToScreeningTool field.
func (r *mutationResolver) RespondToScreeningTool(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error) {
	return r.mycarehub.Questionnaires.RespondToScreeningTool(ctx, input)
}

// GetAvailableScreeningTools is the resolver for the getAvailableScreeningTools field.
func (r *queryResolver) GetAvailableScreeningTools(ctx context.Context, clientID *string) ([]*domain.ScreeningTool, error) {
	return r.mycarehub.Questionnaires.GetAvailableScreeningTools(ctx, clientID)
}

// GetScreeningToolByID is the resolver for the getScreeningToolByID field.
func (r *queryResolver) GetScreeningToolByID(ctx context.Context, id string) (*domain.ScreeningTool, error) {
	return r.mycarehub.Questionnaires.GetScreeningToolByID(ctx, id)
}

// GetFacilityRespondedScreeningTools is the resolver for the getFacilityRespondedScreeningTools field.
func (r *queryResolver) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, paginationInput dto.PaginationsInput) (*domain.ScreeningToolPage, error) {
	return r.mycarehub.Questionnaires.GetFacilityRespondedScreeningTools(ctx, facilityID, &paginationInput)
}

// GetScreeningToolRespondents is the resolver for the getScreeningToolRespondents field.
func (r *queryResolver) GetScreeningToolRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm *string, paginationInput dto.PaginationsInput) (*domain.ScreeningToolRespondentsPage, error) {
	return r.mycarehub.Questionnaires.GetScreeningToolRespondents(ctx, facilityID, screeningToolID, searchTerm, &paginationInput)
}

// GetScreeningToolResponse is the resolver for the getScreeningToolResponse field.
func (r *queryResolver) GetScreeningToolResponse(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
	return r.mycarehub.Questionnaires.GetScreeningToolResponse(ctx, id)
}
