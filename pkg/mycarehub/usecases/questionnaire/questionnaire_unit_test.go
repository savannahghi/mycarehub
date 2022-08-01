package questionnaire_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaire"
)

func TestUseCaseQuestionnaireImpl_ListQuestionnaires(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()

	q := questionnaire.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Questionnaire
		wantErr bool
	}{
		{
			name: "Happy case: ListQuestionnaires",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to ListQuestionnaires",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to ListQuestionnaires" {
				fakeDB.MockListQuestionnairesFn = func(ctx context.Context) ([]*domain.Questionnaire, error) {
					return nil, fmt.Errorf("unable to list questionnaires")
				}
			}
			got, err := q.ListQuestionnaires(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.ListQuestionnaires() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got %v", got)
				return
			}
		})
	}
}
