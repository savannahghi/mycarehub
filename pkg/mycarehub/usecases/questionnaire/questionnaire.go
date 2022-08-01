package questionnaire

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IQuestionnaire interface hold the methods that are used to interact with the questionnaire
type IQuestionnaire interface {
	ListQuestionnaires(ctx context.Context) ([]*domain.Questionnaire, error)
}

//UseCaseQuestionnaireImpl represents the questionnaire usecase object
type UseCaseQuestionnaireImpl struct {
	Create infrastructure.Create
	Update infrastructure.Update
	Query  infrastructure.Query
	Delete infrastructure.Delete
}

// NewUseCaseQuestionnaire returns a new instance of the UseCaseQuestionnaireImpl object
func NewUseCaseQuestionnaire(
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
	delete infrastructure.Delete,
) *UseCaseQuestionnaireImpl {
	return &UseCaseQuestionnaireImpl{
		Create: create,
		Update: update,
		Query:  query,
		Delete: delete,
	}
}

// ListQuestionnaires returns a list of all the questionnaires
func (q *UseCaseQuestionnaireImpl) ListQuestionnaires(ctx context.Context) ([]*domain.Questionnaire, error) {
	return q.Query.ListQuestionnaires(ctx)
}
