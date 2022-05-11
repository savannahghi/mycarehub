package surveys

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
	"github.com/savannahghi/serverutils"
)

// IListSurveys lists the surveys available for a given project
type IListSurveys interface {
	ListSurveys(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error)
	GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error)
	SendClientSurveyLinks(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error)
}

// UsecaseSurveys groups al the interfaces for the Surveys usecase
type UsecaseSurveys interface {
	IListSurveys
}

// UsecaseSurveysImpl represents the Surveys implementation
type UsecaseSurveysImpl struct {
	Surveys      surveys.Surveys
	Query        infrastructure.Query
	Create       infrastructure.Create
	Notification notification.UseCaseNotification
}

// NewUsecaseSurveys is the controller function for the Surveys usecase
func NewUsecaseSurveys(
	surveys surveys.Surveys,
	query infrastructure.Query,
	create infrastructure.Create,
	notification notification.UseCaseNotification,
) *UsecaseSurveysImpl {
	return &UsecaseSurveysImpl{
		Surveys:      surveys,
		Query:        query,
		Create:       create,
		Notification: notification,
	}
}

// GetUserSurveyForms lists the surveys available for a given project
func (u *UsecaseSurveysImpl) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
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

// SendClientSurveyLinks sends survey links to clients
func (u *UsecaseSurveysImpl) SendClientSurveyLinks(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error) {

	var (
		surveyBaseURL    = serverutils.MustGetEnvVar("SURVEYS_BASE_URL")
		userSurveyInputs = []*dto.UserSurveyInput{}
	)

	clients, err := u.Query.GetClientsByFilterParams(ctx, facilityID, filterParams)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error getting clients: %w", err)
	}

	surveyForm, err := u.Surveys.GetSurveyForm(ctx, *projectID, *formID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error getting survey form: %w", err)
	}

	for _, client := range clients {
		odkUserAccessTokenInput := dto.SurveyLinkInput{
			ProjectID:   *projectID,
			FormID:      *formID,
			DisplayName: client.UserID,
			OnceOnly:    true,
		}

		odkUserPublicAccessToken, err := u.Surveys.GeneratePublickAccessLink(ctx, odkUserAccessTokenInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error generating public access link for user: %w", err)
		}

		link := fmt.Sprintf("%s/-/single/%s?st=%s", surveyBaseURL, surveyForm.EnketoID, odkUserPublicAccessToken.Token)

		userSurveyInput := &dto.UserSurveyInput{
			UserID:    client.UserID,
			ProjectID: *projectID,
			FormID:    *formID,
			Title:     surveyForm.Name,
			Link:      link,
			LinkID:    odkUserPublicAccessToken.ID,
			Token:     odkUserPublicAccessToken.Token,
		}
		userSurveyInputs = append(userSurveyInputs, userSurveyInput)

		notificationArgs := notification.ClientNotificationArgs{
			Survey: &domain.UserSurvey{
				Link:   link,
				Title:  surveyForm.Name,
				UserID: client.UserID,
			},
		}

		// TODO:  implement batch notifications after saving the surveys
		composedNotification := notification.ComposeClientNotification(enums.NotificationTypeSurveys, notificationArgs)

		notificationErr := u.Notification.NotifyUser(ctx, client.User, composedNotification)
		if notificationErr != nil {
			helpers.ReportErrorToSentry(notificationErr)
		}

	}

	err = u.Create.CreateUserSurveys(ctx, userSurveyInputs)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error creating user survey: %w", err)
	}

	return true, nil
}
