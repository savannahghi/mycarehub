package domain

import (
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestQuestion_ValidateResponse(t *testing.T) {
	type fields struct {
		ID                string
		Active            bool
		QuestionnaireID   string
		Text              string
		QuestionType      enums.QuestionType
		ResponseValueType enums.QuestionResponseValueType
		Required          bool
		SelectMultiple    bool
		Sequence          int
		Choices           []QuestionInputChoice
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
			name: "Happy case: open ended question",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeOpenEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
			},
			args: args{
				response: gofakeit.BS(),
			},
			wantErr: false,
		},
		{
			name: "Happy case: multiple choice question",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeOpenEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      1,
					},
				},
			},
			args: args{
				response: "0,",
			},
			wantErr: false,
		},
		{
			name: "Happy case: single choice question",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeOpenEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      1,
					},
				},
			},
			args: args{
				response: "0,",
			},
			wantErr: false,
		},
		{
			name: "Sad case: open ended question",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeOpenEnded,
				ResponseValueType: enums.QuestionResponseValueTypeNumber,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
			},
			args: args{
				response: gofakeit.BS(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: open ended question",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeOpenEnded,
				ResponseValueType: enums.QuestionResponseValueTypeBoolean,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
			},
			args: args{
				response: gofakeit.BS(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: multiple choice question, no coma in response value",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      1,
					},
				},
			},
			args: args{
				response: "0,1 3",
			},
			wantErr: true,
		},
		{
			name: "Sad case: multiple choice question, response value not in choices",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      1,
					},
				},
			},
			args: args{
				response: "0,7",
			},
			wantErr: true,
		},
		{
			name: "Sad case: single choice question, response value not in choices",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      1,
					},
				},
			},
			args: args{
				response: "7",
			},
			wantErr: true,
		},
		{
			name: "Sad case: single choice question, required response is empty",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      1,
					},
				},
			},
			args: args{
				response: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Question{
				ID:                tt.fields.ID,
				Active:            tt.fields.Active,
				QuestionnaireID:   tt.fields.QuestionnaireID,
				Text:              tt.fields.Text,
				QuestionType:      tt.fields.QuestionType,
				ResponseValueType: tt.fields.ResponseValueType,
				Required:          tt.fields.Required,
				SelectMultiple:    tt.fields.SelectMultiple,
				Sequence:          tt.fields.Sequence,
				Choices:           tt.fields.Choices,
			}
			if err := s.ValidateResponse(tt.args.response); (err != nil) != tt.wantErr {
				t.Errorf("Question.ValidateResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestion_GetScore(t *testing.T) {
	type fields struct {
		ID                string
		Active            bool
		QuestionnaireID   string
		Text              string
		QuestionType      enums.QuestionType
		ResponseValueType enums.QuestionResponseValueType
		Required          bool
		SelectMultiple    bool
		Sequence          int
		Choices           []QuestionInputChoice
	}
	type args struct {
		response string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			name: "Happy case: single choice question, response value in choices",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      10,
					},
				},
			},
			args: args{
				response: "1",
			},
			want: 10,
		},
		{
			name: "Happy case: multiple choice question, response value in choices",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      10,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      10,
					},
				},
			},
			args: args{
				response: "0,1",
			},
			want: 20,
		},
		{
			name: "Happy case: single choice question, response value in choices",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      0,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      10,
					},
				},
			},
			args: args{
				response: "",
			},
			want: 0,
		},
		{
			name: "Happy case: multiple choice question, response value in choices",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          0,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "0",
						Value:      gofakeit.BS(),
						Score:      10,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: uuid.NewString(),
						Choice:     "1",
						Value:      gofakeit.BS(),
						Score:      10,
					},
				},
			},
			args: args{
				response: "",
			},
			want: 0,
		},
		{
			name: "Happy case: skip scoring open ended question",
			fields: fields{
				ID:                uuid.NewString(),
				Active:            true,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeOpenEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				Sequence:          0,
			},
			args: args{
				response: gofakeit.BS(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Question{
				ID:                tt.fields.ID,
				Active:            tt.fields.Active,
				QuestionnaireID:   tt.fields.QuestionnaireID,
				Text:              tt.fields.Text,
				QuestionType:      tt.fields.QuestionType,
				ResponseValueType: tt.fields.ResponseValueType,
				Required:          tt.fields.Required,
				SelectMultiple:    tt.fields.SelectMultiple,
				Sequence:          tt.fields.Sequence,
				Choices:           tt.fields.Choices,
			}
			if got := s.GetScore(tt.args.response); got != tt.want {
				t.Errorf("Question.GetScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaire_GetQuestionByID(t *testing.T) {
	questionID1 := uuid.NewString()
	questionID2 := uuid.NewString()
	questions := []Question{
		{
			ID:                questionID1,
			Active:            true,
			QuestionnaireID:   uuid.NewString(),
			Text:              gofakeit.BS(),
			QuestionType:      enums.QuestionTypeOpenEnded,
			ResponseValueType: enums.QuestionResponseValueTypeString,
			Required:          true,
			Sequence:          0,
		},
		{
			ID:                questionID2,
			Active:            true,
			QuestionnaireID:   uuid.NewString(),
			Text:              gofakeit.BS(),
			QuestionType:      enums.QuestionTypeOpenEnded,
			ResponseValueType: enums.QuestionResponseValueTypeString,
			Required:          true,
			Sequence:          0,
		},
	}

	type fields struct {
		ID          string
		Active      bool
		Name        string
		Description string
		Questions   []Question
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Question
		wantErr bool
	}{
		{
			name: "Happy case: get question by id",
			fields: fields{
				ID:          uuid.NewString(),
				Active:      true,
				Name:        gofakeit.BS(),
				Description: gofakeit.BS(),
				Questions:   questions,
			},
			args: args{
				id: questionID1,
			},
			want: questions[0],
		},
		{
			name: "Sad case: question dies not exist",
			fields: fields{
				ID:          uuid.NewString(),
				Active:      true,
				Name:        gofakeit.BS(),
				Description: gofakeit.BS(),
				Questions:   questions,
			},
			args: args{
				id: uuid.NewString(),
			},
			want:    Question{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := Questionnaire{
				ID:          tt.fields.ID,
				Active:      tt.fields.Active,
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Questions:   tt.fields.Questions,
			}
			got, err := q.GetQuestionByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Questionnaire.GetQuestionByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Questionnaire.GetQuestionByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningTool_GetQuestion(t *testing.T) {
	question := "What is your name?"
	type fields struct {
		ID              string
		Active          bool
		QuestionnaireID string
		Threshold       int
		ClientTypes     []enums.ClientType
		Genders         []enumutils.Gender
		AgeRange        AgeRange
		Questionnaire   Questionnaire
	}
	type args struct {
		questionID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantNil bool
	}{
		{
			name: "Happy case: get question text",
			fields: fields{
				ID:              uuid.NewString(),
				Active:          true,
				QuestionnaireID: uuid.NewString(),
				Threshold:       3,
				ClientTypes:     []enums.ClientType{enums.ClientTypePmtct},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: AgeRange{
					LowerBound: 18,
					UpperBound: 25,
				},
				Questionnaire: Questionnaire{
					ID:          uuid.NewString(),
					Active:      false,
					Name:        gofakeit.BS(),
					Description: gofakeit.BS(),
					Questions: []Question{
						{
							ID:                "1",
							Active:            false,
							QuestionnaireID:   uuid.NewString(),
							Text:              question,
							QuestionType:      enums.QuestionTypeCloseEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    false,
							Sequence:          1,
							Choices: []QuestionInputChoice{
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "0",
									Value:      gofakeit.BS(),
									Score:      1,
								},
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "1",
									Value:      gofakeit.BS(),
									Score:      2,
								},
							},
						},
					},
				},
			},
			args: args{
				questionID: "1",
			},
			wantNil: false,
		},
		{
			name: "Sad case: Invalid question id",
			fields: fields{
				ID:              uuid.NewString(),
				Active:          true,
				QuestionnaireID: uuid.NewString(),
				Threshold:       3,
				ClientTypes:     []enums.ClientType{enums.ClientTypePmtct},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: AgeRange{
					LowerBound: 18,
					UpperBound: 25,
				},
				Questionnaire: Questionnaire{
					ID:          uuid.NewString(),
					Active:      false,
					Name:        gofakeit.BS(),
					Description: gofakeit.BS(),
					Questions: []Question{
						{
							ID:                "1",
							Active:            false,
							QuestionnaireID:   uuid.NewString(),
							Text:              question,
							QuestionType:      enums.QuestionTypeCloseEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    false,
							Sequence:          1,
							Choices: []QuestionInputChoice{
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "0",
									Value:      gofakeit.BS(),
									Score:      1,
								},
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "1",
									Value:      gofakeit.BS(),
									Score:      2,
								},
							},
						},
					},
				},
			},
			args: args{
				questionID: "invalid",
			},
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ScreeningTool{
				ID:              tt.fields.ID,
				Active:          tt.fields.Active,
				QuestionnaireID: tt.fields.QuestionnaireID,
				Threshold:       tt.fields.Threshold,
				ClientTypes:     tt.fields.ClientTypes,
				Genders:         tt.fields.Genders,
				AgeRange:        tt.fields.AgeRange,
				Questionnaire:   tt.fields.Questionnaire,
			}
			got := s.GetQuestion(tt.args.questionID)
			if (got == nil) != tt.wantNil {
				t.Errorf("ScreeningTool.GetQuestion(), wantNil %v", tt.wantNil)
				return
			}
		})
	}
}

func TestScreeningTool_GetNormalizedResponse(t *testing.T) {
	type fields struct {
		ID              string
		Active          bool
		QuestionnaireID string
		Threshold       int
		ClientTypes     []enums.ClientType
		Genders         []enumutils.Gender
		AgeRange        AgeRange
		Questionnaire   Questionnaire
	}
	type args struct {
		questionID string
		response   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantNil bool
	}{
		{
			name: "Happy case: Close ended question",
			fields: fields{
				ID:              uuid.NewString(),
				Active:          true,
				QuestionnaireID: uuid.NewString(),
				Threshold:       3,
				ClientTypes:     []enums.ClientType{enums.ClientTypePmtct},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: AgeRange{
					LowerBound: 18,
					UpperBound: 25,
				},
				Questionnaire: Questionnaire{
					ID:          uuid.NewString(),
					Active:      false,
					Name:        gofakeit.BS(),
					Description: gofakeit.BS(),
					Questions: []Question{
						{
							ID:                "1",
							Active:            false,
							QuestionnaireID:   uuid.NewString(),
							Text:              gofakeit.BS(),
							QuestionType:      enums.QuestionTypeCloseEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    false,
							Sequence:          1,
							Choices: []QuestionInputChoice{
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "0",
									Value:      "yes",
									Score:      1,
								},
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "1",
									Value:      "no",
									Score:      2,
								},
							},
						},
					},
				},
			},
			args: args{
				questionID: "1",
				response:   "1",
			},
			want: map[string]interface{}{
				"1": "no",
			},
		},
		{
			name: "Happy case: Close ended question",
			fields: fields{
				ID:              uuid.NewString(),
				Active:          true,
				QuestionnaireID: uuid.NewString(),
				Threshold:       3,
				ClientTypes:     []enums.ClientType{enums.ClientTypePmtct},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: AgeRange{
					LowerBound: 18,
					UpperBound: 25,
				},
				Questionnaire: Questionnaire{
					ID:          uuid.NewString(),
					Active:      false,
					Name:        gofakeit.BS(),
					Description: gofakeit.BS(),
					Questions: []Question{
						{
							ID:                "1",
							Active:            false,
							QuestionnaireID:   uuid.NewString(),
							Text:              gofakeit.BS(),
							QuestionType:      enums.QuestionTypeCloseEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    true,
							Sequence:          1,
							Choices: []QuestionInputChoice{
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "0",
									Value:      "yes",
									Score:      1,
								},
								{
									ID:         uuid.NewString(),
									Active:     true,
									QuestionID: "1",
									Choice:     "1",
									Value:      "no",
									Score:      2,
								},
							},
						},
					},
				},
			},
			args: args{
				questionID: "1",
				response:   "0,1",
			},
			want: map[string]interface{}{
				"0": "yes",
				"1": "no",
			},
		},
		{
			name: "Happy case: OPEN_ENDED ended question",
			fields: fields{
				ID:              uuid.NewString(),
				Active:          true,
				QuestionnaireID: uuid.NewString(),
				Threshold:       3,
				ClientTypes:     []enums.ClientType{enums.ClientTypePmtct},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: AgeRange{
					LowerBound: 18,
					UpperBound: 25,
				},
				Questionnaire: Questionnaire{
					ID:          uuid.NewString(),
					Active:      false,
					Name:        gofakeit.BS(),
					Description: gofakeit.BS(),
					Questions: []Question{
						{
							ID:                "1",
							Active:            false,
							QuestionnaireID:   uuid.NewString(),
							Text:              gofakeit.BS(),
							QuestionType:      enums.QuestionTypeOpenEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    false,
							Sequence:          1,
						},
					},
				},
			},
			args: args{
				questionID: "1",
				response:   "@here",
			},
			want: map[string]interface{}{
				"0": "@here",
			},
		},
		{
			name: "Sad case: invalid question ID",
			fields: fields{
				ID:              uuid.NewString(),
				Active:          true,
				QuestionnaireID: uuid.NewString(),
				Threshold:       3,
				ClientTypes:     []enums.ClientType{enums.ClientTypePmtct},
				Genders:         []enumutils.Gender{enumutils.GenderMale},
				AgeRange: AgeRange{
					LowerBound: 18,
					UpperBound: 25,
				},
				Questionnaire: Questionnaire{
					ID:          uuid.NewString(),
					Active:      false,
					Name:        gofakeit.BS(),
					Description: gofakeit.BS(),
					Questions: []Question{
						{
							ID:                "1",
							Active:            false,
							QuestionnaireID:   uuid.NewString(),
							Text:              gofakeit.BS(),
							QuestionType:      enums.QuestionTypeOpenEnded,
							ResponseValueType: enums.QuestionResponseValueTypeString,
							Required:          true,
							SelectMultiple:    false,
							Sequence:          1,
						},
					},
				},
			},
			args: args{
				questionID: "invlaid",
				response:   "@here",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ScreeningTool{
				ID:              tt.fields.ID,
				Active:          tt.fields.Active,
				QuestionnaireID: tt.fields.QuestionnaireID,
				Threshold:       tt.fields.Threshold,
				ClientTypes:     tt.fields.ClientTypes,
				Genders:         tt.fields.Genders,
				AgeRange:        tt.fields.AgeRange,
				Questionnaire:   tt.fields.Questionnaire,
			}
			if got := s.GetNormalizedResponse(tt.args.questionID, tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScreeningTool.GetNormalizedResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_GetNormalizedResponseForSingleChoice(t *testing.T) {
	type fields struct {
		ID                string
		Active            bool
		QuestionnaireID   string
		Text              string
		QuestionType      enums.QuestionType
		ResponseValueType enums.QuestionResponseValueType
		Required          bool
		SelectMultiple    bool
		Sequence          int
		Choices           []QuestionInputChoice
	}
	type args struct {
		response string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "Happy case: single choice question",
			fields: fields{
				ID:                "1",
				Active:            false,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          1,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "0",
						Value:      "yes",
						Score:      1,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "1",
						Value:      "no",
						Score:      2,
					},
				},
			},
			args: args{
				response: "1",
			},
			want: map[string]interface{}{
				"1": "no",
			},
		},
		{
			name: "Sad case: single choice question",
			fields: fields{
				ID:                "1",
				Active:            false,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          1,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "0",
						Value:      "yes",
						Score:      1,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "1",
						Value:      "no",
						Score:      2,
					},
				},
			},
			args: args{
				response: "invalid",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Question{
				ID:                tt.fields.ID,
				Active:            tt.fields.Active,
				QuestionnaireID:   tt.fields.QuestionnaireID,
				Text:              tt.fields.Text,
				QuestionType:      tt.fields.QuestionType,
				ResponseValueType: tt.fields.ResponseValueType,
				Required:          tt.fields.Required,
				SelectMultiple:    tt.fields.SelectMultiple,
				Sequence:          tt.fields.Sequence,
				Choices:           tt.fields.Choices,
			}
			if got := s.GetNormalizedResponseForSingleChoice(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Question.GetNormalizedResponseForSingleChoice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_GetNormalizedResponseForMultipleChoice(t *testing.T) {
	type fields struct {
		ID                string
		Active            bool
		QuestionnaireID   string
		Text              string
		QuestionType      enums.QuestionType
		ResponseValueType enums.QuestionResponseValueType
		Required          bool
		SelectMultiple    bool
		Sequence          int
		Choices           []QuestionInputChoice
	}
	type args struct {
		response string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name: "Happy case: multiple choice question",
			fields: fields{
				ID:                "1",
				Active:            false,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          1,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "0",
						Value:      "yes",
						Score:      1,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "1",
						Value:      "no",
						Score:      2,
					},
				},
			},
			args: args{
				response: "0,1",
			},
			want: map[string]interface{}{
				"0": "yes",
				"1": "no",
			},
		},
		{
			name: "Sad case: ,multiple choice question",
			fields: fields{
				ID:                "1",
				Active:            false,
				QuestionnaireID:   uuid.NewString(),
				Text:              gofakeit.BS(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    true,
				Sequence:          1,
				Choices: []QuestionInputChoice{
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "0",
						Value:      "yes",
						Score:      1,
					},
					{
						ID:         uuid.NewString(),
						Active:     true,
						QuestionID: "1",
						Choice:     "1",
						Value:      "no",
						Score:      2,
					},
				},
			},
			args: args{
				response: "invalid",
			},
			want: map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Question{
				ID:                tt.fields.ID,
				Active:            tt.fields.Active,
				QuestionnaireID:   tt.fields.QuestionnaireID,
				Text:              tt.fields.Text,
				QuestionType:      tt.fields.QuestionType,
				ResponseValueType: tt.fields.ResponseValueType,
				Required:          tt.fields.Required,
				SelectMultiple:    tt.fields.SelectMultiple,
				Sequence:          tt.fields.Sequence,
				Choices:           tt.fields.Choices,
			}
			if got := s.GetNormalizedResponseForMultipleChoice(tt.args.response); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Question.GetNormalizedResponseForMultipleChoice() = %v, want %v", got, tt.want)
			}
		})
	}
}
