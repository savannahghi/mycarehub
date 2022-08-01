package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain/model"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *questionnaireInputResolver) QuestionnaireScreeningTool(ctx context.Context, obj *dto.QuestionnaireInput, data *model.QuestionnaireScreeningToolsInput) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *questionnaireQuestionInputResolver) QuestionType(ctx context.Context, obj *dto.QuestionnaireQuestionInput, data model.QuestionnaireQuestionType) error {
	panic(fmt.Errorf("not implemented"))
}

func (r *questionnaireQuestionInputResolver) Choices(ctx context.Context, obj *dto.QuestionnaireQuestionInput, data []*model.QuestionnaireQuestionInputChoiceInput) error {
	panic(fmt.Errorf("not implemented"))
}

// QuestionnaireInput returns generated.QuestionnaireInputResolver implementation.
func (r *Resolver) QuestionnaireInput() generated.QuestionnaireInputResolver {
	return &questionnaireInputResolver{r}
}

// QuestionnaireQuestionInput returns generated.QuestionnaireQuestionInputResolver implementation.
func (r *Resolver) QuestionnaireQuestionInput() generated.QuestionnaireQuestionInputResolver {
	return &questionnaireQuestionInputResolver{r}
}

type questionnaireInputResolver struct{ *Resolver }
type questionnaireQuestionInputResolver struct{ *Resolver }
