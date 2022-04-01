package screeningtools

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

var (
	questionMeta = map[string]interface{}{
		"violence_code": "GBV-EV",
	}
	gbvMetaTemplate = "The GBV code is " + questionMeta["violence_code"].(string) + ". "

	callClientString                    = "Consider calling Client for further discussion."
	repeatScreeningString               = "Repeat screening for client on subsequent visits."
	wantPositiveTBassessment            = "TB assessment: greater than or equal to 3 yes responses indicates positive TB cases. " + callClientString
	wantNegativeTBassessment            = "TB assessment: all no responses indicates no TB cases. " + repeatScreeningString
	wantPositiveGBVassessment           = "Violence assessment: greater than or equal to 1 yes responses indicates positive Violence cases. " + gbvMetaTemplate + callClientString
	wantNegativeGBVassessment           = "Violence assessment: all no responses indicates no GBV cases. " + repeatScreeningString
	wantPositiveContraceptiveAssessment = "Contraceptive assessment: " + "yes response to question number 4. " + callClientString
	wantPositiveAlcoholAssessment       = "Alcohol/Substance Assessment: greater than or equal to 3 yes responses indicates positive alcohol/substance cases. " + callClientString
	wantNegativeAlcoholAssessment       = "Alcohol/Substance Assessment: all no responses indicates no alcohol/substance cases. " + repeatScreeningString
)

func Test_serviceRequestTemplate(t *testing.T) {
	sequence := 1

	gbvQuestionMeta := enums.ScreeningToolTypeGBV.String() + "_question_number_" + strconv.Itoa(sequence) + "_question_meta"
	type args struct {
		question  *domain.ScreeningToolQuestion
		response  string
		condition map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_serviceRequestTemplate:  yes count greater than or equal to 3 for TB_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeTB,
					Sequence: sequence,
				},
				response: "0",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeTB.String() + "_yes_count": 3,
					enums.ScreeningToolTypeTB.String() + "_no_count":  1,
				},
			},
			want: wantPositiveTBassessment,
		},
		{
			name: "Test_serviceRequestTemplate:  no count equal to len questions for TB_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeTB,
					Sequence: sequence,
				},
				response: "1",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeTB.String() + "_yes_count": 0,
					enums.ScreeningToolTypeTB.String() + "_no_count":  4,
				},
			},
			want: wantNegativeTBassessment,
		},
		{
			name: "Test_serviceRequestTemplate:  yes count greater than 1 or equal to len questions for GBV_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeGBV,
					Sequence: sequence,
				},
				response: "0",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeGBV.String() + "_yes_count": 1,
					enums.ScreeningToolTypeGBV.String() + "_no_count":  3,
					gbvQuestionMeta: map[string]interface{}{
						"violence_code": "GBV-EV",
					},
				},
			},
			want: wantPositiveGBVassessment,
		},
		{
			name: "Test_serviceRequestTemplate:  no count equal to len questions for for GBV_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeGBV,
					Sequence: sequence,
				},
				response: "1",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeGBV.String() + "_yes_count": 0,
					enums.ScreeningToolTypeGBV.String() + "_no_count":  4,
				},
			},
			want: wantNegativeGBVassessment,
		},
		{
			name: "Test_serviceRequestTemplate:  yes for question 4 in  CONTRACEPTIVE_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeCUI,
					Sequence: 3,
				},
				response: "0",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeCUI.String() + "_yes_count":              1,
					enums.ScreeningToolTypeCUI.String() + "_no_count":               3,
					enums.ScreeningToolTypeCUI.String() + "_question_number_" + "3": "yes",
				},
			},
			want: wantPositiveContraceptiveAssessment,
		},
		{
			name: "Test_serviceRequestTemplate:  yes count >=3  ALCOHOL_SUBSTANCE_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeAlcoholSubstanceAssessment,
					Sequence: sequence,
				},
				response: "1",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeAlcoholSubstanceAssessment.String() + "_yes_count": 3,
					enums.ScreeningToolTypeAlcoholSubstanceAssessment.String() + "_no_count":  1,
				},
			},
			want: wantPositiveAlcoholAssessment,
		},
		{
			name: "Test_serviceRequestTemplate:  no count ==4 ALCOHOL_SUBSTANCE_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeAlcoholSubstanceAssessment,
					Sequence: sequence,
				},
				response: "1",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeAlcoholSubstanceAssessment.String() + "_yes_count": 0,
					enums.ScreeningToolTypeAlcoholSubstanceAssessment.String() + "_no_count":  4,
				},
			},
			want: wantNegativeAlcoholAssessment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serviceRequestTemplate(tt.args.question, tt.args.response, tt.args.condition); got != tt.want {
				t.Errorf("serviceRequestTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createServiceRequest(t *testing.T) {
	sequence := 1
	gbvQuestionMeta := enums.ScreeningToolTypeGBV.String() + "_question_number_" + strconv.Itoa(sequence) + "_question_meta"

	wantRedFlagTBRequest := &domain.ServiceRequest{
		RequestType: enums.ServiceRequestTypeScreeningTools.String(),
		Request:     wantPositiveTBassessment,
	}

	wantRedFlagGBVAssessment := &domain.ServiceRequest{
		RequestType: enums.ServiceRequestTypeScreeningTools.String(),
		Request:     wantPositiveGBVassessment,
	}

	wantRedFlagAlcoholAssessment := &domain.ServiceRequest{
		RequestType: enums.ServiceRequestTypeScreeningTools.String(),
		Request:     wantPositiveAlcoholAssessment,
	}

	type args struct {
		question      *domain.ScreeningToolQuestion
		responseValue string
		condition     map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *domain.ServiceRequest
	}{
		{
			name: "Test_createServiceRequest:  yes count greater than or equal to 3 for TB_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeTB,
					Sequence: sequence,
				},
				responseValue: "yes",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeTB.String() + "_yes_count": 3,
					enums.ScreeningToolTypeTB.String() + "_no_count":  1,
				},
			},
			want: wantRedFlagTBRequest,
		},
		{
			name: "Test_createServiceRequest:  no count equal to len questions for TB_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeTB,
					Sequence: sequence,
				},
				responseValue: "no",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeTB.String() + "_yes_count": 0,
					enums.ScreeningToolTypeTB.String() + "_no_count":  4,
				},
			},
			want: nil,
		},
		{
			name: "Test_createServiceRequest:  yes count greater than or equal to 1 for GBV_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeGBV,
					Sequence: sequence,
				},
				responseValue: "yes",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeGBV.String() + "_yes_count": 1,
					enums.ScreeningToolTypeGBV.String() + "_no_count":  3,
					gbvQuestionMeta: map[string]interface{}{
						"violence_code": "GBV-EV",
					},
				},
			},
			want: wantRedFlagGBVAssessment,
		},
		{
			name: "Test_createServiceRequest:  yes for question 4 in  CONTRACEPTIVE_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeCUI,
					Sequence: sequence,
				},
				responseValue: "yes",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeCUI.String() + "_yes_count":              1,
					enums.ScreeningToolTypeCUI.String() + "_no_count":               3,
					enums.ScreeningToolTypeCUI.String() + "_question_number_" + "4": "yes",
				},
			},
			want: nil,
		},
		{
			name: "Test_createServiceRequest:  yes count >=3  ALCOHOL_SUBSTANCE_ASSESSMENT",
			args: args{
				question: &domain.ScreeningToolQuestion{
					ToolType: enums.ScreeningToolTypeAlcoholSubstanceAssessment,
					Sequence: sequence,
				},
				responseValue: "yes",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeAlcoholSubstanceAssessment.String() + "_yes_count": 3,
					enums.ScreeningToolTypeAlcoholSubstanceAssessment.String() + "_no_count":  1,
				},
			},
			want: wantRedFlagAlcoholAssessment,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createServiceRequest(tt.args.question, tt.args.responseValue, tt.args.condition)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addServiceRequestCreateConditions(t *testing.T) {
	question := &domain.ScreeningToolQuestion{
		ID:       uuid.New().String(),
		Question: gofakeit.Sentence(1),
		ToolType: enums.ScreeningToolTypeTB,
		ResponseChoices: map[string]interface{}{
			"0": "Yes",
			"1": "No",
		},
		ResponseType:     enums.ScreeningToolResponseTypeInteger,
		ResponseCategory: enums.ScreeningToolResponseCategorySingleChoice,
		Sequence:         3,
		Meta:             map[string]interface{}{},
		Active:           false,
	}

	singleChoiceQuestionKey := question.ToolType.String() + "_question_number_" + strconv.Itoa(question.Sequence)

	sequence := question.Sequence
	tbQuestionMeta := enums.ScreeningToolTypeTB.String() + "_question_number_" + strconv.Itoa(sequence) + "_question_meta"

	wantYesCount := map[string]interface{}{
		enums.ScreeningToolTypeTB.String() + "_yes_count": 1,
		enums.ScreeningToolTypeTB.String() + "_no_count":  nil,
		singleChoiceQuestionKey:                           "0",
		tbQuestionMeta:                                    map[string]interface{}{},
	}
	wantNoCount := map[string]interface{}{
		enums.ScreeningToolTypeTB.String() + "_yes_count": nil,
		enums.ScreeningToolTypeTB.String() + "_no_count":  1,
		singleChoiceQuestionKey:                           "1",
		tbQuestionMeta:                                    map[string]interface{}{},
	}

	type args struct {
		question  *domain.ScreeningToolQuestion
		response  string
		condition map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "Test_addServiceRequestCreateConditions:  add yes counts",
			args: args{
				question: question,
				response: "0",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeTB.String() + "_yes_count": nil,
					enums.ScreeningToolTypeTB.String() + "_no_count":  nil,
				},
			},
			want: wantYesCount,
		},
		{
			name: "Test_addServiceRequestCreateConditions:  add no counts",
			args: args{
				question: question,
				response: "1",
				condition: map[string]interface{}{
					enums.ScreeningToolTypeTB.String() + "_yes_count": nil,
					enums.ScreeningToolTypeTB.String() + "_no_count":  nil,
				},
			},
			want: wantNoCount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addCondition(tt.args.question, tt.args.response, tt.args.condition); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addServiceRequestCreateConditions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initializeInt(t *testing.T) {
	type args struct {
		n interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test_initializeInt:  initialize existing int",
			args: args{
				n: 130,
			},
			want: 130,
		},
		{
			name: "Test_initializeInt:  initialize int",
			args: args{
				n: nil,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := interfaceToInt(tt.args.n); got != tt.want {
				t.Errorf("interfaceToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interfaceToString(t *testing.T) {
	type args struct {
		n interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_interfaceToString:  initialize existing string",
			args: args{
				n: "130",
			},
			want: "130",
		},
		{
			name: "Test_interfaceToString:  initialize string",
			args: args{
				n: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := interfaceToString(tt.args.n); got != tt.want {
				t.Errorf("interfaceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
