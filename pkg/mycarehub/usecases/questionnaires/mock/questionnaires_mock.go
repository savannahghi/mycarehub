package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// QuestionnaireUseCaseMock mocks the questionnaire instance
type QuestionnaireUseCaseMock struct {
	MockCreateScreeningToolFn                func(ctx context.Context, input dto.ScreeningToolInput) (bool, error)
	MockRespondToScreeningToolFn             func(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error)
	MockGetAvailableScreeningToolsFn         func(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error)
	MockGetScreeningToolByIDFn               func(ctx context.Context, id string) (*domain.ScreeningTool, error)
	MockGetFacilityRespondedScreeningToolsFn func(ctx context.Context, facilityID string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolPage, error)
	MockGetScreeningToolRespondentsFn        func(ctx context.Context, facilityID string, screeningToolID string, searchTerm *string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolRespondentsPage, error)
	MockGetScreeningToolResponseFn           func(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error)
}

// NewServiceRequestUseCaseMock initializes a new questionnaire instance mock
func NewServiceRequestUseCaseMock() *QuestionnaireUseCaseMock {
	UUID := gofakeit.UUID()
	bs := gofakeit.BS()
	now := time.Now()
	questionInputChoice := domain.QuestionInputChoice{
		ID:         UUID,
		Active:     true,
		QuestionID: UUID,
		Choice:     "0",
		Value:      bs,
		Score:      1,
		ProgramID:  UUID,
	}
	question := domain.Question{
		ID:                UUID,
		Active:            true,
		QuestionnaireID:   UUID,
		Text:              bs,
		QuestionType:      enums.QuestionTypeCloseEnded,
		ResponseValueType: enums.QuestionResponseValueTypeBoolean,
		Required:          true,
		SelectMultiple:    false,
		Sequence:          1,
		Choices:           []domain.QuestionInputChoice{questionInputChoice},
		ProgramID:         "",
	}
	questionnaire := domain.Questionnaire{
		ID:          UUID,
		Active:      true,
		Name:        bs,
		Description: bs,
		Questions:   []domain.Question{question},
		ProgramID:   UUID,
	}
	screeningTool := domain.ScreeningTool{
		ID:              UUID,
		Active:          true,
		QuestionnaireID: UUID,
		Threshold:       1,
		ClientTypes:     []enums.ClientType{enums.ClientTypeDreams},
		Genders:         []enumutils.Gender{enumutils.GenderMale},
		AgeRange: domain.AgeRange{
			LowerBound: 20,
			UpperBound: 25,
		},
		Questionnaire: questionnaire,
		ProgramID:     UUID,
	}
	pagination := domain.Pagination{
		Limit:       1,
		CurrentPage: 1,
	}
	screeningToolRespondent := domain.ScreeningToolRespondent{
		ClientID:                UUID,
		ScreeningToolResponseID: UUID,
		ServiceRequestID:        UUID,
		Name:                    bs,
		PhoneNumber:             "0999999999",
		ServiceRequest:          bs,
	}

	questionnaireScreeningToolQuestionResponse := domain.QuestionnaireScreeningToolQuestionResponse{
		ID:                      UUID,
		Active:                  true,
		ScreeningToolResponseID: UUID,
		QuestionID:              UUID,
		QuestionType:            enums.QuestionTypeCloseEnded,
		SelectMultiple:          false,
		ResponseValueType:       "0",
		Sequence:                1,
		QuestionText:            bs,
		Response:                bs,
		NormalizedResponse:      map[string]interface{}{},
		Score:                   0,
		ProgramID:               UUID,
	}
	questionnaireScreeningToolResponse := domain.QuestionnaireScreeningToolResponse{
		ID:                UUID,
		Active:            true,
		ScreeningToolID:   UUID,
		FacilityID:        UUID,
		ClientID:          UUID,
		DateOfResponse:    now,
		AggregateScore:    1,
		QuestionResponses: []*domain.QuestionnaireScreeningToolQuestionResponse{&questionnaireScreeningToolQuestionResponse},
		ProgramID:         "",
	}

	return &QuestionnaireUseCaseMock{
		MockCreateScreeningToolFn: func(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
			return true, nil
		},
		MockRespondToScreeningToolFn: func(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error) {
			return true, nil
		},
		MockGetAvailableScreeningToolsFn: func(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
			return []*domain.ScreeningTool{&screeningTool}, nil
		},
		MockGetScreeningToolByIDFn: func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
			return &screeningTool, nil
		},
		MockGetFacilityRespondedScreeningToolsFn: func(ctx context.Context, facilityID string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolPage, error) {
			return &domain.ScreeningToolPage{
				ScreeningTools: []*domain.ScreeningTool{&screeningTool},
				Pagination:     pagination,
			}, nil
		},
		MockGetScreeningToolRespondentsFn: func(ctx context.Context, facilityID string, screeningToolID string, searchTerm *string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolRespondentsPage, error) {
			return &domain.ScreeningToolRespondentsPage{
				ScreeningToolRespondents: []*domain.ScreeningToolRespondent{&screeningToolRespondent},
				Pagination:               pagination,
			}, nil
		},
		MockGetScreeningToolResponseFn: func(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
			return &questionnaireScreeningToolResponse, nil
		},
	}
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	return q.MockCreateScreeningToolFn(ctx, input)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) RespondToScreeningTool(ctx context.Context, input dto.QuestionnaireScreeningToolResponseInput) (bool, error) {
	return q.MockRespondToScreeningToolFn(ctx, input)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
	return q.MockGetAvailableScreeningToolsFn(ctx, clientID, facilityID)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) GetScreeningToolByID(ctx context.Context, id string) (*domain.ScreeningTool, error) {
	return q.MockGetScreeningToolByIDFn(ctx, id)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) GetFacilityRespondedScreeningTools(ctx context.Context, facilityID string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolPage, error) {
	return q.MockGetFacilityRespondedScreeningToolsFn(ctx, facilityID, paginationInput)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) GetScreeningToolRespondents(ctx context.Context, facilityID string, screeningToolID string, searchTerm *string, paginationInput *dto.PaginationsInput) (*domain.ScreeningToolRespondentsPage, error) {
	return q.MockGetScreeningToolRespondentsFn(ctx, facilityID, screeningToolID, searchTerm, paginationInput)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (q *QuestionnaireUseCaseMock) GetScreeningToolResponse(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
	return q.MockGetScreeningToolResponseFn(ctx, id)
}
