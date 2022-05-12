package surveys

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
)

// IListSurveys lists the surveys available for a given project
type IListSurveys interface {
	ListSurveys(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error)
	GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurveys, error)
}

// UsecaseSurveys groups al the interfaces for the Surveys usecase
type UsecaseSurveys interface {
	IListSurveys
}

// UsecaseSurveysImpl represents the Surveys implementation
type UsecaseSurveysImpl struct {
	Surveys surveys.Surveys
	Query   infrastructure.Query
}

// NewUsecaseSurveys is the controller function for the Surveys usecase
func NewUsecaseSurveys(
	surveys surveys.Surveys,
	query infrastructure.Query,
) *UsecaseSurveysImpl {
	return &UsecaseSurveysImpl{
		Surveys: surveys,
		Query:   query,
	}
}

// GetUserSurveyForms lists the surveys available for a given project
func (u *UsecaseSurveysImpl) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurveys, error) {
	surveys, err := u.Query.GetUserSurveyForms(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return surveys, nil
}

// ListSurveys lists the surveys available for a given project
func (u *UsecaseSurveysImpl) ListSurveys(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error) {
	surveys, err := u.Surveys.ListSurveyForms(ctx, *projectID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	return surveys, nil
}
