package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ScreeningToolsUseCaseMock mocks the implementation of ScreeningTools usecase methods.
type ScreeningToolsUseCaseMock struct {
	MockGetScreeningToolQuestionsFn               func(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error)
	MockGetAvailableScreeningToolQuestionsFn      func(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error)
	MockGetAvailableFacilityScreeningToolsFn      func(ctx context.Context, facilityID string) ([]*domain.AvailableScreeningTools, error)
	MockAnswerScreeningToolQuestionsFn            func(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error)
	MockGetAssessmentResponsesFn                  func(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error)
	MockGetScreeningToolServiceRequestResponsesFn func(ctx context.Context, clientID string, toolType enums.ScreeningToolType) (*domain.ScreeningToolResponsePayload, error)
}

// NewScreeningToolsUseCaseMock creates in itializes create type mocks
func NewScreeningToolsUseCaseMock() *ScreeningToolsUseCaseMock {
	UUID := gofakeit.UUID()
	now := time.Now()
	screeningToolQuestion := domain.ScreeningToolQuestion{
		ID:               UUID,
		Question:         gofakeit.BS(),
		ToolType:         enums.ScreeningToolTypeGBV,
		ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
		ResponseType:     enums.ScreeningToolResponseTypeText,
		ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
		Sequence:         0,
		Meta:             map[string]interface{}{},
		Active:           true,
	}

	availableScreenigtools := domain.AvailableScreeningTools{
		ToolType: enums.ScreeningToolTypeGBV,
	}

	screeningtoolAssesmentResponse := domain.ScreeningToolAssessmentResponse{
		ClientName:   gofakeit.BS(),
		DateAnswered: now,
		ClientID:     UUID,
	}

	screeningToolResponse := domain.ScreeningToolResponse{
		ToolIndex: 0,
		Tool:      UUID,
		Response:  "0",
	}
	screeningToolResponsePayload := domain.ScreeningToolResponsePayload{
		ServiceRequestID:       UUID,
		ClientContact:          "0999999999",
		ScreeningToolResponses: []*domain.ScreeningToolResponse{&screeningToolResponse},
	}
	return &ScreeningToolsUseCaseMock{
		MockGetScreeningToolQuestionsFn: func(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error) {
			return []*domain.ScreeningToolQuestion{&screeningToolQuestion}, nil
		},
		MockGetAvailableScreeningToolQuestionsFn: func(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error) {
			return []*domain.AvailableScreeningTools{&availableScreenigtools}, nil
		},
		MockGetAvailableFacilityScreeningToolsFn: func(ctx context.Context, facilityID string) ([]*domain.AvailableScreeningTools, error) {
			return []*domain.AvailableScreeningTools{&availableScreenigtools}, nil
		},
		MockAnswerScreeningToolQuestionsFn: func(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error) {
			return true, nil
		},
		MockGetAssessmentResponsesFn: func(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
			return []*domain.ScreeningToolAssessmentResponse{&screeningtoolAssesmentResponse}, nil
		},
		MockGetScreeningToolServiceRequestResponsesFn: func(ctx context.Context, clientID string, toolType enums.ScreeningToolType) (*domain.ScreeningToolResponsePayload, error) {
			return &screeningToolResponsePayload, nil
		},
	}
}

// GetScreeningToolQuestions mock the implementation of the GetScreeningToolQuestions method
func (gm *ScreeningToolsUseCaseMock) GetScreeningToolQuestions(ctx context.Context, questionType *string) ([]*domain.ScreeningToolQuestion, error) {
	return gm.MockGetScreeningToolQuestionsFn(ctx, questionType)
}

// GetAvailableScreeningToolQuestions mock the implementation of the GetAvailableScreeningToolQuestions method
func (gm *ScreeningToolsUseCaseMock) GetAvailableScreeningToolQuestions(ctx context.Context, clientID string) ([]*domain.AvailableScreeningTools, error) {
	return gm.MockGetAvailableFacilityScreeningToolsFn(ctx, clientID)
}

// GetAvailableFacilityScreeningTools mock the implementation of the GetAvailableFacilityScreeningTools method
func (gm *ScreeningToolsUseCaseMock) GetAvailableFacilityScreeningTools(ctx context.Context, facilityID string) ([]*domain.AvailableScreeningTools, error) {
	return gm.MockGetAvailableFacilityScreeningToolsFn(ctx, facilityID)
}

// AnswerScreeningToolQuestions mock the implementation of the AnswerScreeningToolQuestions method
func (gm *ScreeningToolsUseCaseMock) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) (bool, error) {
	return gm.MockAnswerScreeningToolQuestionsFn(ctx, screeningToolResponses)
}

// GetAssessmentResponses mock the implementation of the GetAssessmentResponses method
func (gm *ScreeningToolsUseCaseMock) GetAssessmentResponses(ctx context.Context, facilityID string, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
	return gm.MockGetAssessmentResponsesFn(ctx, facilityID, toolType)
}

// GetScreeningToolServiceRequestResponses mock the implementation of the GetScreeningToolServiceRequestResponses method
func (gm *ScreeningToolsUseCaseMock) GetScreeningToolServiceRequestResponses(ctx context.Context, clientID string, toolType enums.ScreeningToolType) (*domain.ScreeningToolResponsePayload, error) {
	return gm.MockGetScreeningToolServiceRequestResponsesFn(ctx, clientID, toolType)
}
