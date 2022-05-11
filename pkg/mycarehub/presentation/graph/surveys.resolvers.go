package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) SendClientSurveyLinks(ctx context.Context, facilityID string, formID string, projectID int, filterParams *dto.ClientFilterParamsInput) (bool, error) {
	return r.mycarehub.Surveys.SendClientSurveyLinks(ctx, &facilityID, &formID, &projectID, filterParams)
}

func (r *queryResolver) ListSurveys(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
	return r.mycarehub.Surveys.ListSurveys(ctx, &projectID)
}

func (r *queryResolver) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
	return r.mycarehub.Surveys.GetUserSurveyForms(ctx, userID)
}
