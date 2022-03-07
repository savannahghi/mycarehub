package screeningtools

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
)

func TestServiceScreeningToolsImpl_GetScreeningToolsQuestions(t *testing.T) {
	tbQuestionType := enums.ScreeningToolTypeTB.String()
	invalidToolType := enums.ScreeningToolType("invalid").String()
	type args struct {
		ctx      context.Context
		toolType *string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningToolQuestion
		wantErr bool
	}{
		{
			name: "happy case: get all screening tools questions",
			args: args{
				ctx:      context.Background(),
				toolType: &tbQuestionType,
			},
			wantErr: false,
		},
		{
			name: "happy case: get all screening tools questions, no params",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "sad case: get screening tools questions, invalid question type",
			args: args{
				ctx:      context.Background(),
				toolType: &invalidToolType,
			},
			wantErr: true,
		},
		{
			name: "sad case: get screening tools questions, error",
			args: args{
				ctx:      context.Background(),
				toolType: &tbQuestionType,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			tr := NewUseCasesScreeningTools(fakeDB)

			if tt.name == "sad case: get screening tools questions, error" {
				fakeDB.MockGetScreeningToolsQuestionsFn = func(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("failed to get screening tools questions")
				}
			}

			got, err := tr.GetScreeningToolQuestions(tt.args.ctx, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.GetScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ServiceScreeningToolsImpl.GetScreeningToolQuestions() = %v, want %v", got, tt.want)
			}
		})
	}
}
