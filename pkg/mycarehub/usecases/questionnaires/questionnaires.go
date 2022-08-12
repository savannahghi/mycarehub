package questionnaires

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ICreateScreeningTools contains methods related to the screening tools
type ICreateScreeningTools interface {
	CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error)
	RespondToScreeningTool(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error)
}

// IGetScreeningTools contains methods related to the screening tools
type IGetScreeningTools interface {
	GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error)
	GetScreeningToolByID(ctx context.Context, id string) (*domain.ScreeningTool, error)
}

// UseCaseQuestionnaire contains questionnaire interfaces
type UseCaseQuestionnaire interface {
	ICreateScreeningTools
	IGetScreeningTools
}

// UseCaseQuestionnaireImpl represents the questionnaire implementations
type UseCaseQuestionnaireImpl struct {
	Query  infrastructure.Query
	Create infrastructure.Create
	Update infrastructure.Update
	Delete infrastructure.Delete
}

// NewUseCaseQuestionnaire is the controller function for the questionnaire usecase
func NewUseCaseQuestionnaire(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	delete infrastructure.Delete,
) UseCaseQuestionnaire {
	return &UseCaseQuestionnaireImpl{
		Query:  query,
		Create: create,
		Update: update,
		Delete: delete,
	}
}

// CreateScreeningTool creates the screening tool questionnaire
func (q *UseCaseQuestionnaireImpl) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	err := input.Questionnaire.Validate()
	if err != nil {
		return false, err
	}

	questions := []domain.Question{}
	for _, q := range input.Questionnaire.Questions {
		choices := []domain.QuestionInputChoice{}
		for _, c := range q.Choices {
			choices = append(choices, domain.QuestionInputChoice{
				Active: true,
				Choice: *c.Choice,
				Value:  c.Value,
				Score:  c.Score,
			})
		}
		questions = append(questions, domain.Question{
			Active:            true,
			Text:              q.Text,
			QuestionType:      q.QuestionType,
			ResponseValueType: q.ResponseValueType,
			Required:          q.Required,
			SelectMultiple:    q.SelectMultiple,
			Sequence:          q.Sequence,
			Choices:           choices,
		})
	}

	payload := &domain.ScreeningTool{
		Active:      true,
		Threshold:   input.Threshold,
		ClientTypes: input.ClientTypes,
		Genders:     input.Genders,
		AgeRange: domain.AgeRange{
			LowerBound: input.AgeRange.LowerBound,
			UpperBound: input.AgeRange.UpperBound,
		},
		Questionnaire: domain.Questionnaire{
			Active:      true,
			Name:        input.Questionnaire.Name,
			Description: input.Questionnaire.Description,
			Questions:   questions,
		},
	}

	err = q.Create.CreateScreeningTool(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to create screening tool questionnaire: %w", err)
	}
	return true, nil
}

// RespondToScreeningTool responds to the screening tool questionnaire
func (q *UseCaseQuestionnaireImpl) RespondToScreeningTool(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error) {
	err := input.Validate()
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	clientProfile, err := q.Query.GetClientProfileByClientID(ctx, input.ClientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get client profile: %w", err)
	}

	screeningTool, err := q.Query.GetScreeningToolByID(ctx, input.ScreeningToolID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get screening tool: %w", err)
	}

	payload := &domain.QuestionnaireScreeningToolResponse{
		Active:          true,
		ScreeningToolID: input.ScreeningToolID,
		FacilityID:      clientProfile.FacilityID,
		ClientID:        input.ClientID,
	}

	var aggregateScore int

	responses := []domain.QuestionnaireScreeningToolQuestionResponse{}
	for _, qr := range input.QuestionResponses {
		question, err := screeningTool.Questionnaire.GetQuestionByID(qr.QuestionID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to find question with id: %s", qr.QuestionID)
		}

		err = question.ValidateResponse(qr.Response)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to validate response: %w", err)
		}

		score := question.GetScore(qr.Response)
		aggregateScore += score

		responses = append(responses, domain.QuestionnaireScreeningToolQuestionResponse{
			Active:                  true,
			ScreeningToolResponseID: screeningTool.ID,
			QuestionID:              qr.QuestionID,
			Response:                qr.Response,
			Score:                   score,
		})
	}

	payload.AggregateScore = aggregateScore
	payload.QuestionResponses = responses

	responseID, err := q.Create.CreateScreeningToolResponse(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to create screening tool response: %w", err)
	}

	if aggregateScore >= screeningTool.Threshold {
		serviceRequest := fmt.Sprintf("%s has a score of %d for %s. They require your attention", clientProfile.User.Name, aggregateScore, screeningTool.Questionnaire.Name)
		err = q.Create.CreateServiceRequest(ctx, &dto.ServiceRequestInput{
			Active:      true,
			RequestType: enums.ServiceRequestTypeScreeningToolsRedFlag.String(),
			Status:      enums.ServiceRequestStatusPending.String(),
			Request:     serviceRequest,
			ClientID:    input.ClientID,
			FacilityID:  clientProfile.FacilityID,
			ClientName:  &clientProfile.User.Name,
			Flavour:     feedlib.FlavourConsumer,
			Meta: map[string]interface{}{
				"response_id": *responseID,
			},
		})
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create service request: %w", err)
		}
	}
	return true, nil
}

// GetAvailableScreeningTools returns the available screening tools
func (q *UseCaseQuestionnaireImpl) GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
	screeningTools, err := q.Query.GetAvailableScreeningTools(ctx, clientID, facilityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tools: %w", err)
	}

	return screeningTools, nil
}

// GetScreeningToolByID looks for a screening tool using the provided ID and returns the screening tool
func (q *UseCaseQuestionnaireImpl) GetScreeningToolByID(ctx context.Context, id string) (*domain.ScreeningTool, error) {
	screeningTool, err := q.Query.GetScreeningToolByID(ctx, id)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tool: %w", err)
	}
	return screeningTool, nil
}
