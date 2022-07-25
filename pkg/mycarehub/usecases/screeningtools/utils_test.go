package screeningtools

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func Test_validateResponses(t *testing.T) {
	type args struct {
		ctx          context.Context
		questions    []*domain.ScreeningToolQuestion
		answersInput []*dto.ScreeningToolQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: valid input",
			args: args{
				ctx: context.Background(),
				questions: []*domain.ScreeningToolQuestion{
					{
						ID:       "1",
						Question: "question",
						ToolType: enums.ScreeningToolTypeTB,
						ResponseChoices: map[string]interface{}{
							"1": "response1",
							"2": "response2",
						},
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						Sequence:         0,
						Meta:             map[string]interface{}{},
						Active:           true,
					},
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         uuid.NewString(),
						QuestionID:       "1",
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeTB,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: response out of range",
			args: args{
				ctx: context.Background(),
				questions: []*domain.ScreeningToolQuestion{
					{
						ID:       "1",
						Question: "question",
						ToolType: enums.ScreeningToolTypeTB,
						ResponseChoices: map[string]interface{}{
							"1": "response1",
							"2": "response2",
						},
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						Sequence:         0,
						Meta:             map[string]interface{}{},
						Active:           true,
					},
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         uuid.NewString(),
						QuestionID:       "1",
						Response:         "4",
						ToolType:         enums.ScreeningToolTypeTB,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: response type is invalid",
			args: args{
				ctx: context.Background(),
				questions: []*domain.ScreeningToolQuestion{
					{
						ID:       "1",
						Question: "question",
						ToolType: enums.ScreeningToolTypeTB,
						ResponseChoices: map[string]interface{}{
							"1": "response1",
							"2": "response2",
						},
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						Sequence:         0,
						Meta:             map[string]interface{}{},
						Active:           true,
					},
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         uuid.NewString(),
						QuestionID:       "1",
						Response:         "invalid",
						ToolType:         enums.ScreeningToolTypeTB,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateResponses(tt.args.ctx, tt.args.questions, tt.args.answersInput); (err != nil) != tt.wantErr {
				t.Errorf("validateResponses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_generateServiceRequest(t *testing.T) {
	id := uuid.NewString()
	type args struct {
		ctx           context.Context
		clientProfile *domain.ClientProfile
		answersInput  []*dto.ScreeningToolQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}{
		{
			name: "happy case: generate tb service request",
			args: args{
				ctx: context.Background(),
				clientProfile: &domain.ClientProfile{
					ID: &id,
					User: &domain.User{
						ID: &id, Name: gofakeit.Name(),
					},
					FacilityID: uuid.NewString(),
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeTB,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeTB,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 1,
					},
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeTB,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 2,
					},
				},
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "happy case: generate gbv service request",
			args: args{
				ctx: context.Background(),
				clientProfile: &domain.ClientProfile{
					ID: &id,
					User: &domain.User{
						ID: &id, Name: gofakeit.Name(),
					},
					FacilityID: uuid.NewString(),
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeGBV,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
				},
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "happy case: generate alcohol and substance service request",
			args: args{
				ctx: context.Background(),
				clientProfile: &domain.ClientProfile{
					ID: &id,
					User: &domain.User{
						ID: &id, Name: gofakeit.Name(),
					},
					FacilityID: uuid.NewString(),
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeAlcoholSubstanceAssessment,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeAlcoholSubstanceAssessment,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 1,
					},
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeAlcoholSubstanceAssessment,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 2,
					},
				},
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "happy case: generate contraceptive use service request",
			args: args{
				ctx: context.Background(),
				clientProfile: &domain.ClientProfile{
					ID: &id,
					User: &domain.User{
						ID: &id, Name: gofakeit.Name(),
					},
					FacilityID: uuid.NewString(),
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "yes",
						ToolType:         enums.ScreeningToolTypeCUI,
						ResponseType:     enums.ScreeningToolResponseTypeText,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 3,
					},
				},
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "sad case: invalid tool type",
			args: args{
				ctx: context.Background(),
				clientProfile: &domain.ClientProfile{
					ID: &id,
					User: &domain.User{
						ID: &id, Name: gofakeit.Name(),
					},
					FacilityID: uuid.NewString(),
				},
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "yes",
						ToolType:         enums.ScreeningToolType(enums.ScreeningToolResponseType.String("invalid")),
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 3,
					},
				},
			},
			wantErr: true,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateServiceRequest(tt.args.ctx, tt.args.clientProfile, tt.args.answersInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("generateServiceRequest() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_calculateToolScore(t *testing.T) {
	id := uuid.NewString()
	type args struct {
		ctx          context.Context
		answersInput []*dto.ScreeningToolQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: context.Background(),
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeAlcoholSubstanceAssessment,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 0,
					},
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeAlcoholSubstanceAssessment,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 1,
					},
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "1",
						ToolType:         enums.ScreeningToolTypeAlcoholSubstanceAssessment,
						ResponseType:     enums.ScreeningToolResponseTypeInteger,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 2,
					},
				},
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "sad case: non response type integer passed",
			args: args{
				answersInput: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:         id,
						QuestionID:       uuid.NewString(),
						Response:         "yes",
						ToolType:         enums.ScreeningToolTypeCUI,
						ResponseType:     enums.ScreeningToolResponseTypeText,
						ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
						QuestionSequence: 3,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateToolScore(tt.args.ctx, tt.args.answersInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateToolScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateToolScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
