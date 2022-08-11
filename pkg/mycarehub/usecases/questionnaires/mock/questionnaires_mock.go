package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// QuestionnaireUseCaseMock mocks the questionnaire instance
type QuestionnaireUseCaseMock struct {
	MockCreateScreeningToolFn        func(ctx context.Context, input dto.ScreeningToolInput) (bool, error)
	MockGetAvailableScreeningToolsFn func(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error)
}

// NewServiceRequestUseCaseMock initializes a new questionnaire instance mock
func NewServiceRequestUseCaseMock() *QuestionnaireUseCaseMock {
	return &QuestionnaireUseCaseMock{
		MockCreateScreeningToolFn: func(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
			return true, nil
		},
		MockGetAvailableScreeningToolsFn: func(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
			return []*domain.ScreeningTool{
				{
					ID:              uuid.New().String(),
					Active:          true,
					QuestionnaireID: uuid.New().String(),
					Threshold:       13,
					ClientTypes:     []enums.ClientType{"PMTCT", "PEP"},
					Genders:         []enumutils.Gender{"MALE"},
					AgeRange: domain.AgeRange{
						LowerBound: 14,
						UpperBound: 22,
					},
					Questionnaire: domain.Questionnaire{},
				},
			}, nil
		},
	}
}

// CreateScreeningTool mocks the CreateScreeningTool method
func (q *QuestionnaireUseCaseMock) CreateScreeningTool(ctx context.Context, input dto.ScreeningToolInput) (bool, error) {
	return q.MockCreateScreeningToolFn(ctx, input)
}

// GetAvailableScreeningTools mocks the GetAvailableScreeningTools method
func (q *QuestionnaireUseCaseMock) GetAvailableScreeningTools(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
	return q.MockGetAvailableScreeningToolsFn(ctx, clientID, facilityID)
}
