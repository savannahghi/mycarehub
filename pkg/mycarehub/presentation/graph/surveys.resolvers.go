package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *queryResolver) ListSurveys(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
	return r.mycarehub.Surveys.ListSurveys(ctx, &projectID)
}

func (r *queryResolver) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurveys, error) {
	return r.mycarehub.Surveys.GetUserSurveyForms(ctx, userID)
}
