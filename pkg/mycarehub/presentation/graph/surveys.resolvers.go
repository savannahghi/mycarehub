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

func (r *mutationResolver) VerifySurveySubmission(ctx context.Context, input dto.VerifySurveySubmissionInput) (bool, error) {
	return r.mycarehub.Surveys.VerifySurveySubmission(ctx, input)
}

func (r *queryResolver) ListSurveys(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
	return r.mycarehub.Surveys.ListSurveys(ctx, &projectID)
}

func (r *queryResolver) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
	return r.mycarehub.Surveys.GetUserSurveyForms(ctx, userID)
}

func (r *queryResolver) ListSurveyRespondents(ctx context.Context, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyRespondentPage, error) {
	return r.mycarehub.Surveys.ListSurveyRespondents(ctx, projectID, formID, paginationInput)
}

func (r *queryResolver) GetSurveyResponse(ctx context.Context, input dto.SurveyResponseInput) ([]*domain.SurveyResponse, error) {
	return r.mycarehub.Surveys.GetSurveyResponse(ctx, input)
}

func (r *queryResolver) GetSurveyWithServiceRequest(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
	return r.mycarehub.Surveys.GetSurveysWithServiceRequests(ctx, facilityID)
}
