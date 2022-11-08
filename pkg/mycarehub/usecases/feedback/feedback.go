package feedback

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/mail"
)

// SendFeedback groups the methods to send feedbacks
type SendFeedback interface {
	SendFeedback(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error)
}

// UsecaseFeedback groups al the interfaces for the feedback usecase
type UsecaseFeedback interface {
	SendFeedback
}

// UsecaseFeedbackImpl represents the feedback implementation
type UsecaseFeedbackImpl struct {
	Query  infrastructure.Query
	Create infrastructure.Create
	Mail   mail.IServiceMail
}

// NewUsecaseFeedback is the controller function for the feedback usecase
func NewUsecaseFeedback(
	query infrastructure.Query,
	create infrastructure.Create,
	mail mail.IServiceMail,
) *UsecaseFeedbackImpl {
	return &UsecaseFeedbackImpl{
		Query:  query,
		Create: create,
		Mail:   mail,
	}
}

// SendFeedback sends the users feedback tothe admin
func (f *UsecaseFeedbackImpl) SendFeedback(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error) {
	if payload.Feedback == "" {
		return false, fmt.Errorf("feedback input cannot be empty")
	}

	userProfile, err := f.Query.GetUserProfileByUserID(ctx, payload.UserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to get user profile: %v", err)
	}

	feedbackSubject := fmt.Sprintf("%s's feedback", userProfile.Name)

	feedbackInput := &dto.FeedbackEmail{
		User:              userProfile.Name,
		FeedbackType:      payload.FeedbackType,
		SatisfactionLevel: payload.SatisfactionLevel,
		Feedback:          payload.Feedback,
	}
	if payload.FeedbackType == enums.ServiceFeedbackType {
		feedbackInput.ServiceName = payload.ServiceName
	}
	if payload.RequiresFollowUp {
		phoneNumber := fmt.Sprintf("Phone Number: %s", userProfile.Contacts.ContactValue)
		feedbackInput.PhoneNumber = phoneNumber
	}

	// Save feedback into the database before sending the email
	feedbackData := &domain.FeedbackResponse{
		UserID:            payload.UserID,
		FeedbackType:      payload.FeedbackType,
		SatisfactionLevel: payload.SatisfactionLevel,
		ServiceName:       payload.ServiceName,
		Feedback:          payload.Feedback,
		RequiresFollowUp:  payload.RequiresFollowUp,
		PhoneNumber:       userProfile.Contacts.ContactValue,
	}

	err = f.Create.SaveFeedback(ctx, feedbackData)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to save feedback: %w", err)
	}

	var writer bytes.Buffer
	tmpl := template.Must(template.New("FeedbackNotificationEmail").Parse(utils.FeedbackNotificationEmail))
	err = tmpl.Execute(&writer, feedbackInput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to create feedback email: %v", err)
	}

	feedbackSent, err := f.Mail.SendFeedback(ctx, feedbackSubject, writer.String())
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to send feedback: %v", err)
	}

	return feedbackSent, nil
}
