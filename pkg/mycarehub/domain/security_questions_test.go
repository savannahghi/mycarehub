package domain

import (
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestSecurityQuestion_Validate(t *testing.T) {
	type fields struct {
		SecurityQuestionID string
		QuestionStem       string
		Description        string
		Flavour            feedlib.Flavour
		Active             bool
		ResponseType       enums.SecurityQuestionResponseType
	}
	type args struct {
		response string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid: security question response type string",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "What is your mother's maiden name?",
				Description:        "Mother's maiden name",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeText,
			},

			args: args{
				response: "test",
			},
			wantErr: false,
		},
		{
			name: "valid: security question response type number",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "What is your age?",
				Description:        "Your age",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeNumber,
			},

			args: args{
				response: "24",
			},
			wantErr: false,
		},
		{
			name: "valid: security question response type boolean",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "Do you have children?",
				Description:        "general",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeBoolean,
			},

			args: args{
				response: "true",
			},
			wantErr: false,
		},
		{
			name: "invalid: security question response type boolean is invalid",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "Do you have children?",
				Description:        "general",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeBoolean,
			},

			args: args{
				response: "nah",
			},
			wantErr: true,
		},
		{
			name: "invalid: security question response type number",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "What is your age?",
				Description:        "Your age",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeNumber,
			},

			args: args{
				response: "twenty four",
			},
			wantErr: true,
		},
		{
			name: "valid: security question response type Date",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "When is your birthday?",
				Description:        "Your birthday",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeDate,
			},

			args: args{
				response: "24-02-1999",
			},
			wantErr: false,
		},
		{
			name: "invalid: security question response type Date, invalid date format",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "When is your birthday?",
				Description:        "Your birthday",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeDate,
			},

			args: args{
				response: "1999-02-24",
			},
			wantErr: true,
		},
		{
			name: "invalid: security question response type Date, invalid date format",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "When is your birthday?",
				Description:        "Your birthday",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeDate,
			},

			args: args{
				response: "24/02/1999",
			},
			wantErr: true,
		},
		{
			name: "invalid: security question response type Date, invalid date",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "When is your birthday?",
				Description:        "Your birthday",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeDate,
			},

			args: args{
				response: "34-02-1999",
			},
			wantErr: true,
		},
		{
			name: "invalid: security question response type Date, invalid date",
			fields: fields{
				SecurityQuestionID: "123",
				QuestionStem:       "When is your birthday?",
				Description:        "Your birthday",
				Flavour:            feedlib.FlavourConsumer,
				Active:             true,
				ResponseType:       enums.SecurityQuestionResponseTypeDate,
			},

			args: args{
				response: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SecurityQuestion{
				SecurityQuestionID: tt.fields.SecurityQuestionID,
				QuestionStem:       tt.fields.QuestionStem,
				Description:        tt.fields.Description,
				Flavour:            tt.fields.Flavour,
				Active:             tt.fields.Active,
				ResponseType:       tt.fields.ResponseType,
			}
			if err := s.Validate(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("SecurityQuestion.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
