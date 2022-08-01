package questionnaire

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// IQuestionnaire interface hold the methods that are used to interact with the questionnaire
type IQuestionnaire interface {
	CreateQuestionnaire(input *dto.QuestionnaireInput) error
	GetQuestionnaireByID(id string) (*domain.Questionnaire, error)
	updatedQuestionnaire(id string, input *dto.QuestionnaireInput) error
	DeleteQuestionnaire(id string) error
}
