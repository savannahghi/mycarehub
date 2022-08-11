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

var surveyBaseURL = serverutils.MustGetEnvVar("SURVEYS_BASE_URL")

// IListSurveys lists the surveys available for a given project
type IListSurveys interface {
	ListSurveys(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error)
	GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error)
	SendClientSurveyLinks(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error)
}

// IVerifySurveySubmission contains all the methods that can be used to update a survey
type IVerifySurveySubmission interface {
	VerifySurveySubmission(ctx context.Context, input dto.VerifySurveySubmissionInput) (bool, error)
}

// UsecaseSurveys groups al the interfaces for the Surveys usecase
type UsecaseSurveys interface {
	IListSurveys
	IVerifySurveySubmission
}

// UsecaseSurveysImpl represents the Surveys implementation
type UsecaseSurveysImpl struct {
	Surveys      surveys.Surveys
	Query        infrastructure.Query
	Create       infrastructure.Create
	Update       infrastructure.Update
	Notification notification.UseCaseNotification
}

// NewUsecaseSurveys is the controller function for the Surveys usecase
func NewUsecaseSurveys(
	surveys surveys.Surveys,
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	notification notification.UseCaseNotification,
) *UsecaseSurveysImpl {
	return &UsecaseSurveysImpl{
		Surveys:      surveys,
		Query:        query,
		Create:       create,
		Update:       update,
		Notification: notification,
	}
}

// GetUserSurveyForms lists the surveys available for a given project
func (u *UsecaseSurveysImpl) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
	userSurveys, err := u.Query.GetUserSurveyForms(ctx, userID, nil, nil, nil)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return userSurveys, nil
}

// VerifySurveySubmission method is used to verify whether a user has filled a survey.
// If the user has filled the survey and submitted their data, the method marks (in the database), that the survey has been  submitted.
// This method is called when the user goes back from the page that used to fill surveys.
func (u *UsecaseSurveysImpl) VerifySurveySubmission(ctx context.Context, input dto.VerifySurveySubmissionInput) (bool, error) {
	submitters, err := u.Surveys.ListSubmitters(ctx, input.ProjectID, input.FormID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	for _, submitter := range submitters {
		if submitter.ID == input.SubmitterID {
			survey := &domain.UserSurvey{
				LinkID:    input.SubmitterID,
				ProjectID: input.ProjectID,
				FormID:    input.FormID,
			}

			updateData := map[string]interface{}{
				"has_submitted": true,
			}
			err := u.Update.UpdateUserSurveys(ctx, survey, updateData)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}
			break
		}
	}

	return true, nil
}

// ListSurveys lists the surveys available for a given project
func (u *UsecaseSurveysImpl) ListSurveys(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error) {
	allSurveys, err := u.Surveys.ListSurveyForms(ctx, *projectID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	return allSurveys, nil
}

// SendClientSurveyLinks sends survey links to clients
func (u *UsecaseSurveysImpl) SendClientSurveyLinks(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error) {
	clients, err := u.Query.GetClientsByFilterParams(ctx, facilityID, filterParams)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error getting clients: %w", err)
	}

	if len(clients) == 0 {
		return true, nil
	}

	surveyForm, err := u.Surveys.GetSurveyForm(ctx, *projectID, *formID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error getting survey form: %w", err)
	}

	// alert is a temporary type that holds survey notification input
	// the notifications should be sent out after successful survey creation
	type alert struct {
		client domain.ClientProfile
		survey dto.UserSurveyInput
	}
	clientsNotifications := []alert{}

	userSurveyInputs := []*dto.UserSurveyInput{}

	for _, client := range clients {
		// validate if they have an existing survey that has been sent
		// If a survey exists for a client, continue to the next client
		hasSubmitted := false
		userSurveys, err := u.Query.GetUserSurveyForms(ctx, client.UserID, projectID, formID, &hasSubmitted)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

		if len(userSurveys) > 0 {
			continue
		}

		accessTokenInput := dto.SurveyLinkInput{
			ProjectID:   *projectID,
			FormID:      *formID,
			DisplayName: client.UserID,
			OnceOnly:    true,
		}

		publicAccessToken, err := u.Surveys.GeneratePublicAccessLink(ctx, accessTokenInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("error generating public access link for user: %w", err)
		}

		link := fmt.Sprintf("%s/-/single/%s?st=%s", surveyBaseURL, surveyForm.EnketoID, publicAccessToken.Token)

		userSurveyInput := &dto.UserSurveyInput{
			UserID:    client.UserID,
			ProjectID: *projectID,
			FormID:    *formID,
			Title:     surveyForm.Name,
			Link:      link,
			LinkID:    publicAccessToken.ID,
			Token:     publicAccessToken.Token,
		}
		userSurveyInputs = append(userSurveyInputs, userSurveyInput)

		clientsNotifications = append(clientsNotifications, alert{client: *client, survey: *userSurveyInput})

	}

	err = u.Create.CreateUserSurveys(ctx, userSurveyInputs)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error creating user survey: %w", err)
	}

	for _, alert := range clientsNotifications {
		notificationInput := notification.ClientNotificationInput{
			Survey: &domain.UserSurvey{
				Link:   alert.survey.Link,
				Title:  surveyForm.Name,
				UserID: alert.client.UserID,
			},
		}

		composedNotification := notification.ComposeClientNotification(enums.NotificationTypeSurveys, notificationInput)

		err := u.Notification.NotifyUser(ctx, alert.client.User, composedNotification)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}

	return true, nil
}
