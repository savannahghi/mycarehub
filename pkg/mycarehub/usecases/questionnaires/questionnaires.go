package questionnaires

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
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
	GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolPage, error)
	GetScreeningToolRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm *string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolRespondentsPage, error)
	GetScreeningToolResponse(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error)
}

// UseCaseQuestionnaire contains questionnaire interfaces
type UseCaseQuestionnaire interface {
	ICreateScreeningTools
	IGetScreeningTools
}

// UseCaseQuestionnaireImpl represents the questionnaire implementations
type UseCaseQuestionnaireImpl struct {
	Query       infrastructure.Query
	Create      infrastructure.Create
	Update      infrastructure.Update
	Delete      infrastructure.Delete
	ExternalExt extension.ExternalMethodsExtension
}

// NewUseCaseQuestionnaire is the controller function for the questionnaire usecase
func NewUseCaseQuestionnaire(
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	delete infrastructure.Delete,
	externalExt extension.ExternalMethodsExtension,
) UseCaseQuestionnaire {
	return &UseCaseQuestionnaireImpl{
		Query:       query,
		Create:      create,
		Update:      update,
		Delete:      delete,
		ExternalExt: externalExt,
	}
}

// CreateScreeningTool creates the screening tool questionnaire
func (q *UseCaseQuestionnaireImpl) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	err := input.Questionnaire.Validate()
	if err != nil {
		return false, err
	}

	userID, err := q.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	userProfile, err := q.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	questions := []domain.Question{}
	sequenceMap := make(map[int]int)
	for _, q := range input.Questionnaire.Questions {
		if _, ok := sequenceMap[q.Sequence]; ok {
			return false, fmt.Errorf("duplicate sequence found: %w", err)
		}
		sequenceMap[q.Sequence] = q.Sequence

		choices := []domain.QuestionInputChoice{}
		choiceMap := make(map[string]string)
		for _, c := range q.Choices {
			if _, ok := choiceMap[*c.Choice]; ok {
				return false, fmt.Errorf("duplicate choice found: %w", err)
			}
			choiceMap[*c.Choice] = *c.Choice

			choices = append(choices, domain.QuestionInputChoice{
				Active:    true,
				Choice:    *c.Choice,
				Value:     c.Value,
				Score:     c.Score,
				ProgramID: userProfile.CurrentProgramID,
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
			ProgramID:         userProfile.CurrentProgramID,
		})
	}

	payload := &domain.ScreeningTool{
		Active:      true,
		Threshold:   input.Threshold,
		ClientTypes: input.ClientTypes,
		Genders:     input.Genders,
		ProgramID:   userProfile.CurrentProgramID,
		AgeRange: domain.AgeRange{
			LowerBound: input.AgeRange.LowerBound,
			UpperBound: input.AgeRange.UpperBound,
		},
		Questionnaire: domain.Questionnaire{
			Active:      true,
			Name:        input.Questionnaire.Name,
			Description: input.Questionnaire.Description,
			Questions:   questions,
			ProgramID:   userProfile.CurrentProgramID,
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
		FacilityID:      clientProfile.DefaultFacility.ID,
		ClientID:        input.ClientID,
		ProgramID:       clientProfile.User.CurrentProgramID,
	}

	var aggregateScore int

	responses := []*domain.QuestionnaireScreeningToolQuestionResponse{}
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

		responses = append(responses, &domain.QuestionnaireScreeningToolQuestionResponse{
			Active:                  true,
			ScreeningToolResponseID: screeningTool.ID,
			QuestionID:              qr.QuestionID,
			Response:                qr.Response,
			Score:                   score,
			ProgramID:               clientProfile.User.CurrentProgramID,
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
			FacilityID:  clientProfile.DefaultFacility.ID,
			ClientName:  &clientProfile.User.Name,
			Flavour:     feedlib.FlavourConsumer,
			Meta: map[string]interface{}{
				"response_id":         *responseID,
				"screening_tool_name": screeningTool.Questionnaire.Name,
				"score":               aggregateScore,
			},
			ProgramID: clientProfile.User.CurrentProgramID,
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

// GetFacilityRespondedScreeningTools gets a list of  screening tools that have been responded to for a given facility
// These screening tools have a service request that has not been resolved yet
func (q *UseCaseQuestionnaireImpl) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolPage, error) {
	if err := paginationInput.Validate(); err != nil {
		return nil, err
	}

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	screeningTools, pageInfo, err := q.Query.GetFacilityRespondedScreeningTools(ctx, facilityID, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tools: %w", err)
	}
	return &domain.ScreeningToolPage{
		ScreeningTools: screeningTools,
		Pagination:     *pageInfo,
	}, nil
}

// GetScreeningToolRespondents returns the respondents for the screening tool
// the respondents belong to a given facility and they must have answered
// a given screening tool which has an unresolved service request
func (q *UseCaseQuestionnaireImpl) GetScreeningToolRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm *string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolRespondentsPage, error) {
	emptyString := ""
	if searchTerm == nil {
		searchTerm = &emptyString
	}
	if err := paginationInput.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	respondents, pageInfo, err := q.Query.GetScreeningToolRespondents(ctx, facilityID, screeningToolID, *searchTerm, paginationInput)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tool respondents: %w", err)
	}
	return &domain.ScreeningToolRespondentsPage{
		ScreeningToolRespondents: respondents,
		Pagination:               *pageInfo,
	}, nil
}

// GetScreeningToolResponse returns the screening tool response for the provided screening tool and client
// the response is in a human-readable format
func (q *UseCaseQuestionnaireImpl) GetScreeningToolResponse(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
	response, err := q.Query.GetScreeningToolResponseByID(ctx, id)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get screening tool response: %w", err)
	}
	return response, nil
}
