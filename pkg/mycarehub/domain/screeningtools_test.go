package domain

import (
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestScreeningToolQuestion_ValidateResponseQUestionCategory(t *testing.T) {
	type fields struct {
		ID               string
		Question         string
		ToolType         enums.ScreeningToolType
		ResponseChoices  map[string]interface{}
		ResponseType     enums.ScreeningToolResponseType
		ResponseCategory enums.ScreeningToolResponseCategory
		Sequence         int
		Meta             map[string]interface{}
		Active           bool
	}
	type args struct {
		response string
		category enums.ScreeningToolResponseCategory
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "validate single choice response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "0",
				category: enums.ScreeningToolResponseCategorySingleChoice,
			},
			wantErr: false,
		},
		{
			name: "validate multi choice response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategoryMultiChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "0,1",
				category: enums.ScreeningToolResponseCategoryMultiChoice,
			},
			wantErr: false,
		},
		{
			name: "validate open ended response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeText,
				ResponseCategory: enums.ScreeningToolResponseCategoryOpenEnded,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "open ended response",
				category: enums.ScreeningToolResponseCategoryOpenEnded,
			},
			wantErr: false,
		},
		{
			name: "invalid single choice response, out of range",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "2",
				category: enums.ScreeningToolResponseCategorySingleChoice,
			},
			wantErr: true,
		},
		{
			name: "invalid single choice response, invalid response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "invalid",
				category: enums.ScreeningToolResponseCategorySingleChoice,
			},
			wantErr: true,
		},
		{
			name: "invalid multi choice response, out of range",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategoryMultiChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "2,3",
				category: enums.ScreeningToolResponseCategoryMultiChoice,
			},
			wantErr: true,
		},
		{
			name: "invalid multi choice response, invalid response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategoryMultiChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response: "invalid",
				category: enums.ScreeningToolResponseCategoryMultiChoice,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &ScreeningToolQuestion{
				ID:               tt.fields.ID,
				Question:         tt.fields.Question,
				ToolType:         tt.fields.ToolType,
				ResponseChoices:  tt.fields.ResponseChoices,
				ResponseType:     tt.fields.ResponseType,
				ResponseCategory: tt.fields.ResponseCategory,
				Sequence:         tt.fields.Sequence,
				Meta:             tt.fields.Meta,
				Active:           tt.fields.Active,
			}
			if err := q.ValidateResponseQuestionCategory(tt.args.response, tt.args.category); (err != nil) != tt.wantErr {
				t.Errorf("ScreeningToolQuestion.ValidateResponseQuestionCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScreeningToolQuestion_ValidateResponseQUestionType(t *testing.T) {
	type fields struct {
		ID               string
		Question         string
		ToolType         enums.ScreeningToolType
		ResponseChoices  map[string]interface{}
		ResponseType     enums.ScreeningToolResponseType
		ResponseCategory enums.ScreeningToolResponseCategory
		Sequence         int
		Meta             map[string]interface{}
		Active           bool
	}
	type args struct {
		response     string
		responseType enums.ScreeningToolResponseType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "validate integer response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseChoices:  map[string]interface{}{"0": "yes", "1": "no"},
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response:     "0",
				responseType: enums.ScreeningToolResponseTypeInteger,
			},
			wantErr: false,
		},
		{
			name: "validate date response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeDate,
				ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response:     "02-01-2006",
				responseType: enums.ScreeningToolResponseTypeDate,
			},
			wantErr: false,
		},
		{
			name: "validate text response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeText,
				ResponseCategory: enums.ScreeningToolResponseCategoryOpenEnded,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response:     "text",
				responseType: enums.ScreeningToolResponseTypeText,
			},
			wantErr: false,
		},
		{
			name: "invalid integer response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeInteger,
				ResponseCategory: enums.ScreeningToolResponseCategoryOpenEnded,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response:     "invalid",
				responseType: enums.ScreeningToolResponseTypeInteger,
			},
			wantErr: true,
		},
		{
			name: "invalid date response",
			fields: fields{
				ID:               "questionID",
				Question:         "question",
				ToolType:         enums.ScreeningToolTypeTB,
				ResponseType:     enums.ScreeningToolResponseTypeDate,
				ResponseCategory: enums.ScreeningToolResponseCategoryOpenEnded,
				Sequence:         1,
				Meta:             map[string]interface{}{"meta": "meta"},
				Active:           true,
			},
			args: args{
				response:     "INVALID",
				responseType: enums.ScreeningToolResponseTypeDate,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &ScreeningToolQuestion{
				ID:               tt.fields.ID,
				Question:         tt.fields.Question,
				ToolType:         tt.fields.ToolType,
				ResponseChoices:  tt.fields.ResponseChoices,
				ResponseType:     tt.fields.ResponseType,
				ResponseCategory: tt.fields.ResponseCategory,
				Sequence:         tt.fields.Sequence,
				Meta:             tt.fields.Meta,
				Active:           tt.fields.Active,
			}
			if err := q.ValidateResponseQUestionType(tt.args.response, tt.args.responseType); (err != nil) != tt.wantErr {
				t.Errorf("ScreeningToolQuestion.ValidateResponseQUestionType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
