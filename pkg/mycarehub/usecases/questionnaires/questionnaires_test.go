package questionnaires_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires"
)

func TestUseCaseQuestionnaireImpl_CreateScreeningTool(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)
	choice := "YES"
	questionnare := &dto.ScreeningToolInput{
		Questionnaire: dto.QuestionnaireInput{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.Sentence(20),
			Questions: []*dto.QuestionInput{
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &choice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &choice,
							Value:  "YES",
							Score:  1,
						},
					},
				},
			},
		},
		Threshold:   3,
		ClientTypes: []enums.ClientType{enums.ClientTypePmtct, enums.ClientTypePmtct},
		Genders:     []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: dto.AgeRangeInput{
			LowerBound: 10,
			UpperBound: 20,
		},
	}
	invalidQuestionnare := &dto.ScreeningToolInput{
		Questionnaire: dto.QuestionnaireInput{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.Sentence(20),
			Questions: []*dto.QuestionInput{
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeOpenEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &choice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &choice,
							Value:  "YES",
							Score:  1,
						},
					},
				},
			},
		},
		Threshold:   3,
		ClientTypes: []enums.ClientType{enums.ClientTypePmtct, enums.ClientTypePmtct},
		Genders:     []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: dto.AgeRangeInput{
			LowerBound: 10,
			UpperBound: 20,
		},
	}
	type args struct {
		ctx   context.Context
		input dto.ScreeningToolInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: Create screening tool",
			args: args{
				ctx:   context.Background(),
				input: *questionnare,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: unable to create screening tool",
			args: args{
				ctx:   context.Background(),
				input: *questionnare,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to validate question",
			args: args{
				ctx:   context.Background(),
				input: *invalidQuestionnare,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to create screening tool" {
				fakeDB.MockCreateScreeningToolFn = func(ctx context.Context, input *domain.ScreeningTool) error {
					return errors.New("unable to create screening tool")
				}
			}
			got, err := q.CreateScreeningTool(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.CreateScreeningTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseQuestionnaireImpl.CreateScreeningTool() = %v, want %v", got, tt.want)
			}
		})
	}
}
