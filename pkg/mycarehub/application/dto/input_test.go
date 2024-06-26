package dto

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
	"github.com/segmentio/ksuid"
)

func TestFacilityInput_Validate(t *testing.T) {
	longWord := gofakeit.Sentence(100)
	veryLongWord := gofakeit.Sentence(500)

	type fields struct {
		Name        string
		Code        int
		Phone       string
		Active      bool
		Country     string
		Description string
		Identifier  FacilityIdentifierInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all fields with correct value",
			fields: fields{
				Name:        "test name",
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Country:     "KE",
				Description: "test description",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: false,
		},

		{
			name: "invalid: short name len",
			fields: fields{
				Name:        "te",
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Country:     "KE",
				Description: "test description",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: long name len",
			fields: fields{
				Name:        longWord,
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Country:     "KE",
				Description: "test description",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: short description",
			fields: fields{
				Name:        "test name",
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Country:     "KE",
				Description: "te",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: very long description",
			fields: fields{
				Name:        "test name",
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Country:     "KE",
				Description: veryLongWord,
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: missing name",
			fields: fields{
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Country:     "KE",
				Description: "test description",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: missing country",
			fields: fields{
				Name:        "test name",
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Description: "test description",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: missing description",
			fields: fields{
				Name:    "test name",
				Code:    22344,
				Phone:   interserviceclient.TestUserPhoneNumber,
				Active:  true,
				Country: "KE",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing phone",
			fields: fields{
				Name:    "test name",
				Code:    22344,
				Active:  true,
				Country: "KE",
				Identifier: FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "11111",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FacilityInput{
				Name:        tt.fields.Name,
				Phone:       tt.fields.Phone,
				Active:      tt.fields.Active,
				Country:     enums.Country(tt.fields.Country),
				Description: tt.fields.Description,
				Identifier:  tt.fields.Identifier,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("FacilityInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaginationsInput_Validate(t *testing.T) {
	type fields struct {
		Limit       int
		CurrentPage int
		Sort        SortsInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				Limit:       1,
				CurrentPage: 1,
				Sort: SortsInput{
					Direction: enums.SortDataTypeAsc,
					Field:     enums.FilterSortDataTypeActive,
				},
			},
			wantErr: false,
		},
		{
			name: "valid: all params passed",
			fields: fields{
				CurrentPage: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid: required field not passed",
			fields: fields{
				Limit: 1,
				Sort: SortsInput{
					Direction: enums.SortDataTypeAsc,
					Field:     enums.FilterSortDataTypeActive,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PaginationsInput{
				Limit:       tt.fields.Limit,
				CurrentPage: tt.fields.CurrentPage,
				Sort:        tt.fields.Sort,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PaginationsInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoginInput_Validate(t *testing.T) {
	testPIN := "0000"

	type fields struct {
		Username string
		PIN      string
		Flavour  feedlib.Flavour
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				Username: gofakeit.Username(),
				PIN:      testPIN,
				Flavour:  feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "invalid : missing pin",
			fields: fields{
				Username: gofakeit.Username(),
				Flavour:  feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing flavour",
			fields: fields{
				Username: gofakeit.Username(),
				PIN:      testPIN,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing username",
			fields: fields{
				PIN:     testPIN,
				Flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LoginInput{
				Username: tt.fields.Username,
				PIN:      tt.fields.PIN,
				Flavour:  tt.fields.Flavour,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("LoginInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFiltersInput_Validate(t *testing.T) {
	type fields struct {
		DataType enums.FilterSortDataType
		Value    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				DataType: enums.FilterSortDataTypeActive,
				Value:    "true",
			},
			wantErr: false,
		},
		{
			name: "invalid: missing datatype",
			fields: fields{
				Value: "true",
			},
			wantErr: true,
		},
		{
			name: "invalid : missing value",
			fields: fields{
				DataType: enums.FilterSortDataTypeActive,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FiltersInput{
				DataType: tt.fields.DataType,
				Value:    tt.fields.Value,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("FiltersInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPINInput_Validate(t *testing.T) {
	ID := ksuid.New().String()
	testPIN := "0000"
	type fields struct {
		UserID     *string
		PIN        *string
		ConfirmPIN *string
		Flavour    feedlib.Flavour
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				UserID:     &ID,
				PIN:        &testPIN,
				ConfirmPIN: &testPIN,
				Flavour:    feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing user id",
			fields: fields{
				PIN:        &testPIN,
				ConfirmPIN: &testPIN,
				Flavour:    feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid : missing pin",
			fields: fields{
				UserID:     &ID,
				ConfirmPIN: &testPIN,
				Flavour:    feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing confirm pin",
			fields: fields{
				UserID:  &ID,
				PIN:     &testPIN,
				Flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},

		{
			name: "invalid: missing flavour",
			fields: fields{
				UserID:     &ID,
				PIN:        &testPIN,
				ConfirmPIN: &testPIN,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &PINInput{
				UserID:     tt.fields.UserID,
				PIN:        tt.fields.PIN,
				ConfirmPIN: tt.fields.ConfirmPIN,
				Flavour:    tt.fields.Flavour,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("PINInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecurityQuestionResponseInput_Validate(t *testing.T) {
	type fields struct {
		UserID             string
		SecurityQuestionID string
		Response           string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				UserID:             "123",
				SecurityQuestionID: "123",
				Response:           "123",
			},
		},
		{
			name: "invalid: missing user id",
			fields: fields{
				SecurityQuestionID: "123",
				Response:           "123",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing security question id",
			fields: fields{
				UserID:   "123",
				Response: "123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &SecurityQuestionResponseInput{
				UserID:             tt.fields.UserID,
				SecurityQuestionID: tt.fields.SecurityQuestionID,
				Response:           tt.fields.Response,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SecurityQuestionResponseInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifySecurityQuestionInput_Validate(t *testing.T) {
	type fields struct {
		QuestionID string
		Flavour    feedlib.Flavour
		Response   string
		Username   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				QuestionID: "123",
				Flavour:    feedlib.FlavourConsumer,
				Response:   "123",
				Username:   gofakeit.Word(),
			},
			wantErr: false,
		},
		{
			name:    "invalid: missing params",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &VerifySecurityQuestionInput{
				QuestionID: tt.fields.QuestionID,
				Flavour:    tt.fields.Flavour,
				Response:   tt.fields.Response,
				Username:   gofakeit.Word(),
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("VerifySecurityQuestionInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserRespondedSecurityQuestionsInput_Validate(t *testing.T) {
	type fields struct {
		Username string
		Flavour  feedlib.Flavour
		OTP      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				Username: gofakeit.Word(),
				Flavour:  feedlib.FlavourConsumer,
				OTP:      "1234",
			},
			wantErr: false,
		},
		{
			name: "invalid: missing params",
			fields: fields{
				Flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &GetUserRespondedSecurityQuestionsInput{
				Username: tt.fields.Username,
				Flavour:  tt.fields.Flavour,
				OTP:      tt.fields.OTP,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GetUserRespondedSecurityQuestionsInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserResetPinInput_Validate(t *testing.T) {
	type fields struct {
		Username string
		Flavour  feedlib.Flavour
		PIN      string
		OTP      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				Username: gofakeit.Word(),
				Flavour:  feedlib.FlavourConsumer,
				PIN:      "1234",
				OTP:      "1234",
			},
		},
		{
			name: "invalid: missing params",
			fields: fields{
				Flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &UserResetPinInput{
				Username: gofakeit.Word(),
				Flavour:  tt.fields.Flavour,
				PIN:      tt.fields.PIN,
				OTP:      tt.fields.OTP,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserResetPinInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShareContentInput_Validate(t *testing.T) {
	type fields struct {
		ClientID  string
		ContentID int
		Channel   string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				ClientID:  "123",
				ContentID: 123,
				Channel:   "123",
			},
		},
		{
			name: "invalid: missing params",
			fields: fields{
				ClientID:  "123",
				ContentID: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ShareContentInput{
				ClientID:  tt.fields.ClientID,
				ContentID: tt.fields.ContentID,
				Channel:   tt.fields.Channel,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ShareContentInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStaffRegistrationInput_Validate(t *testing.T) {
	type fields struct {
		Username    string
		Facility    string
		StaffName   string
		Gender      enumutils.Gender
		DateOfBirth scalarutils.Date
		PhoneNumber string
		IDNumber    string
		StaffNumber string
		InviteStaff bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				Username:  gofakeit.Username(),
				Facility:  "123",
				StaffName: "123",
				Gender:    enumutils.GenderMale,
				DateOfBirth: scalarutils.Date{
					Year:  1992,
					Month: 2,
					Day:   12,
				},
				PhoneNumber: "+254098759039",
				IDNumber:    "12121212121",
				StaffNumber: "s212121",
				InviteStaff: false,
			},
		},
		{
			name: "invalid: empty id number",
			fields: fields{
				Facility:  "123",
				StaffName: "123",
				Gender:    enumutils.GenderMale,
				DateOfBirth: scalarutils.Date{
					Year:  1992,
					Month: 2,
					Day:   12,
				},
				PhoneNumber: "+254098759039",
				IDNumber:    "",
				StaffNumber: "s212121",
				InviteStaff: false,
			},
			wantErr: true,
		},
		{
			name: "invalid: INVALID id number",
			fields: fields{
				Facility:  "123",
				StaffName: "123",
				Gender:    enumutils.GenderMale,
				DateOfBirth: scalarutils.Date{
					Year:  1992,
					Month: 2,
					Day:   12,
				},
				PhoneNumber: "+254098759039",
				IDNumber:    "e12121212121",
				StaffNumber: "s212121",
				InviteStaff: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StaffRegistrationInput{
				Username:    tt.fields.Username,
				Facility:    tt.fields.Facility,
				StaffName:   tt.fields.StaffName,
				Gender:      tt.fields.Gender,
				DateOfBirth: tt.fields.DateOfBirth,
				PhoneNumber: tt.fields.PhoneNumber,
				IDNumber:    tt.fields.IDNumber,
				StaffNumber: tt.fields.StaffNumber,
				InviteStaff: tt.fields.InviteStaff,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("StaffRegistrationInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestionnaireInput_Validate(t *testing.T) {
	choice1 := "yes"
	choice2 := "no"
	type fields struct {
		Name        string
		Description string
		Questions   []*QuestionInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Happy case: valid questionnaire",
			fields: fields{
				Name:        gofakeit.BeerBlg(),
				Description: gofakeit.BeerBlg(),
				Questions: []*QuestionInput{
					{
						Text:              gofakeit.BeerBlg(),
						QuestionType:      enums.QuestionTypeCloseEnded,
						ResponseValueType: enums.QuestionResponseValueTypeBoolean,
						Required:          true,
						SelectMultiple:    false,
						Sequence:          1,
						Choices: []QuestionInputChoiceInput{
							{
								Choice: &choice1,
								Value:  "true",
								Score:  1,
							},
							{
								Choice: &choice2,
								Value:  "false",
								Score:  0,
							},
						},
					},
				},
			},
		},
		{
			name: "Sad case: open ended question with select multiple",
			fields: fields{
				Name:        gofakeit.BeerBlg(),
				Description: gofakeit.BeerBlg(),
				Questions: []*QuestionInput{
					{
						Text:              gofakeit.BeerBlg(),
						QuestionType:      enums.QuestionTypeOpenEnded,
						ResponseValueType: enums.QuestionResponseValueTypeString,
						Required:          true,
						SelectMultiple:    true,
						Sequence:          1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: open ended question with choices",
			fields: fields{
				Name:        gofakeit.BeerBlg(),
				Description: gofakeit.BeerBlg(),
				Questions: []*QuestionInput{
					{
						Text:              gofakeit.BeerBlg(),
						QuestionType:      enums.QuestionTypeOpenEnded,
						ResponseValueType: enums.QuestionResponseValueTypeString,
						Required:          true,
						Sequence:          1,
						Choices: []QuestionInputChoiceInput{
							{
								Choice: &choice1,
								Value:  "true",
								Score:  1,
							},
							{
								Choice: &choice2,
								Value:  "false",
								Score:  0,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: close ended question with less than 2 choices",
			fields: fields{
				Name:        gofakeit.BeerBlg(),
				Description: gofakeit.BeerBlg(),
				Questions: []*QuestionInput{
					{
						Text:              gofakeit.BeerBlg(),
						QuestionType:      enums.QuestionTypeCloseEnded,
						ResponseValueType: enums.QuestionResponseValueTypeString,
						Required:          true,
						SelectMultiple:    false,
						Sequence:          1,
						Choices: []QuestionInputChoiceInput{
							{
								Choice: &choice1,
								Value:  gofakeit.BeerAlcohol(),
								Score:  1,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid choice value for boolean",
			fields: fields{
				Name:        gofakeit.BeerBlg(),
				Description: gofakeit.BeerBlg(),
				Questions: []*QuestionInput{
					{
						Text:              gofakeit.BeerBlg(),
						QuestionType:      enums.QuestionTypeCloseEnded,
						ResponseValueType: enums.QuestionResponseValueTypeBoolean,
						Required:          true,
						SelectMultiple:    false,
						Sequence:          1,
						Choices: []QuestionInputChoiceInput{
							{
								Choice: &choice1,
								Value:  "invalid",
								Score:  1,
							},
							{
								Choice: &choice2,
								Value:  "false",
								Score:  0,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := QuestionnaireInput{
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Questions:   tt.fields.Questions,
			}
			if err := q.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("QuestionnaireInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestionInput_Validate(t *testing.T) {
	choice1 := "yes"
	choice2 := "no"
	type fields struct {
		Text              string
		QuestionType      enums.QuestionType
		ResponseValueType enums.QuestionResponseValueType
		Required          bool
		SelectMultiple    bool
		Sequence          int
		Choices           []QuestionInputChoiceInput
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Happy case: valid question",
			fields: fields{
				Text:              gofakeit.BeerBlg(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeString,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          1,
				Choices: []QuestionInputChoiceInput{
					{
						Choice: &choice1,
						Value:  gofakeit.BeerBlg(),
						Score:  1,
					},
					{
						Choice: &choice2,
						Value:  gofakeit.BeerBlg(),
						Score:  0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: close ended question expecting a number",
			fields: fields{
				Text:              gofakeit.BeerBlg(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeNumber,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          2,
				Choices: []QuestionInputChoiceInput{
					{
						Choice: &choice1,
						Value:  "true",
						Score:  1,
					},
					{
						Choice: &choice2,
						Value:  "false",
						Score:  0,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: close ended question expecting a boolean",
			fields: fields{
				Text:              gofakeit.BeerBlg(),
				QuestionType:      enums.QuestionTypeCloseEnded,
				ResponseValueType: enums.QuestionResponseValueTypeBoolean,
				Required:          true,
				SelectMultiple:    false,
				Sequence:          3,
				Choices: []QuestionInputChoiceInput{
					{
						Choice: &choice1,
						Value:  "1",
						Score:  1,
					},
					{
						Choice: &choice2,
						Value:  "2",
						Score:  0,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := QuestionInput{
				Text:              tt.fields.Text,
				QuestionType:      tt.fields.QuestionType,
				ResponseValueType: tt.fields.ResponseValueType,
				Required:          tt.fields.Required,
				SelectMultiple:    tt.fields.SelectMultiple,
				Sequence:          tt.fields.Sequence,
				Choices:           tt.fields.Choices,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("QuestionInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyOTPInput_Validate(t *testing.T) {
	type fields struct {
		PhoneNumber string
		Username    string
		OTP         string
		Flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Happy case: valid input",
			fields: fields{
				PhoneNumber: "+254799999999",
				Username:    "test",
				OTP:         "0000",
				Flavour:     feedlib.FlavourPro,
			},
			wantErr: false,
		},

		{
			name: "Happy case: valid input 2",
			fields: fields{
				PhoneNumber: "0799999999",
				Username:    "test",
				OTP:         "0000",
				Flavour:     feedlib.FlavourPro,
			},
			wantErr: false,
		},

		{
			name: "Happy case: valid input, non-kenyan phone number",
			fields: fields{
				PhoneNumber: "+1799999999",
				Username:    "test",
				OTP:         "0000",
				Flavour:     feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid phone",
			fields: fields{
				PhoneNumber: "799999999",
				Username:    "test",
				OTP:         "0000",
				Flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},

		{
			name: "Sad case: empty input",
			fields: fields{
				PhoneNumber: "+1799999999",
				Username:    "",
				OTP:         "0000",
				Flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing input",
			fields: fields{
				PhoneNumber: "+1799999999",
				OTP:         "0000",
				Flavour:     feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid flavour",
			fields: fields{
				PhoneNumber: "+1799999999",
				Username:    "test",
				OTP:         "0000",
				Flavour:     feedlib.Flavour("invalid"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &VerifyOTPInput{
				PhoneNumber: tt.fields.PhoneNumber,
				Username:    tt.fields.Username,
				OTP:         tt.fields.OTP,
				Flavour:     tt.fields.Flavour,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("VerifyOTPInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
