package mock

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// QuestionnaireUseCaseMock mocks the implementation of questionnaire usecase methods
type QuestionnaireUseCaseMock struct {
	MockListQuestionnairesFn func(ctx context.Context) ([]*domain.Questionnaire, error)
}

func NewQuestionnaireUsecaseMock() *QuestionnaireUseCaseMock {
	UUID := uuid.New().String()
	questionnaire := &domain.Questionnaire{
		ID:          UUID,
		Name:        "Test Questionnaire",
		Description: "Test Questionnaire Description",
		StartDate:   time.Now(),
		EndDate:     time.Now(),
		Active:      true,
		Questions:   []domain.QuestionnaireQuestion{},
	}
	return &QuestionnaireUseCaseMock{
		MockListQuestionnairesFn: func(ctx context.Context) ([]*domain.Questionnaire, error) {
			return []*domain.Questionnaire{questionnaire}, nil
		},
	}
}

// ListQuestionnaires mocks the implementation of questionnaire usecase methods
func (m *QuestionnaireUseCaseMock) ListQuestionnaires(ctx context.Context) ([]*domain.Questionnaire, error) {
	return m.MockListQuestionnairesFn(ctx)
}
