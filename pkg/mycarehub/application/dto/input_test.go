package dto

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
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
		County      string
		Description string
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
				County:      "Nairobi",
				Description: "test description",
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
				County:      "Nairobi",
				Description: "test description",
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
				County:      "Nairobi",
				Description: "test description",
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
				County:      "Nairobi",
				Description: "te",
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
				County:      "Nairobi",
				Description: veryLongWord,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing name",
			fields: fields{
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				County:      "Nairobi",
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing code",
			fields: fields{
				Name:        "test name",
				Active:      true,
				Phone:       interserviceclient.TestUserPhoneNumber,
				County:      "Nairobi",
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing county",
			fields: fields{
				Name:        "test name",
				Code:        22344,
				Phone:       interserviceclient.TestUserPhoneNumber,
				Active:      true,
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing description",
			fields: fields{
				Name:   "test name",
				Code:   22344,
				Phone:  interserviceclient.TestUserPhoneNumber,
				Active: true,
				County: "Nairobi",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing phone",
			fields: fields{
				Name:   "test name",
				Code:   22344,
				Active: true,
				County: "Nairobi",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FacilityInput{
				Name:        tt.fields.Name,
				Code:        tt.fields.Code,
				Phone:       tt.fields.Phone,
				Active:      tt.fields.Active,
				County:      tt.fields.County,
				Description: tt.fields.Description,
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
	testPhone := interserviceclient.TestUserPhoneNumber
	testPIN := "0000"

	type fields struct {
		PhoneNumber *string
		PIN         *string
		Flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				PhoneNumber: &testPhone,
				PIN:         &testPIN,
				Flavour:     feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing phone number",
			fields: fields{
				PIN:     &testPIN,
				Flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid : missing pin",
			fields: fields{
				PhoneNumber: &testPhone,
				Flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing flavour",
			fields: fields{
				PhoneNumber: &testPhone,
				PIN:         &testPIN,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &LoginInput{
				PhoneNumber: tt.fields.PhoneNumber,
				PIN:         tt.fields.PIN,
				Flavour:     tt.fields.Flavour,
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
		QuestionID  string
		Flavour     feedlib.Flavour
		Response    string
		PhoneNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				QuestionID:  "123",
				Flavour:     feedlib.FlavourConsumer,
				Response:    "123",
				PhoneNumber: interserviceclient.TestUserPhoneNumber,
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
				QuestionID:  tt.fields.QuestionID,
				Flavour:     tt.fields.Flavour,
				Response:    tt.fields.Response,
				PhoneNumber: tt.fields.PhoneNumber,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("VerifySecurityQuestionInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserRespondedSecurityQuestionsInput_Validate(t *testing.T) {
	type fields struct {
		PhoneNumber string
		Flavour     feedlib.Flavour
		OTP         string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				PhoneNumber: gofakeit.Phone(),
				Flavour:     feedlib.FlavourConsumer,
				OTP:         "1234",
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
				PhoneNumber: tt.fields.PhoneNumber,
				Flavour:     tt.fields.Flavour,
				OTP:         tt.fields.OTP,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GetUserRespondedSecurityQuestionsInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserResetPinInput_Validate(t *testing.T) {
	type fields struct {
		PhoneNumber string
		Flavour     feedlib.Flavour
		PIN         string
		OTP         string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: all params passed",
			fields: fields{
				PhoneNumber: gofakeit.Phone(),
				Flavour:     feedlib.FlavourConsumer,
				PIN:         "1234",
				OTP:         "1234",
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
				PhoneNumber: tt.fields.PhoneNumber,
				Flavour:     tt.fields.Flavour,
				PIN:         tt.fields.PIN,
				OTP:         tt.fields.OTP,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UserResetPinInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShareContentInput_Validate(t *testing.T) {
	type fields struct {
		UserID    string
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
				UserID:    "123",
				ContentID: 123,
				Channel:   "123",
			},
		},
		{
			name: "invalid: missing params",
			fields: fields{
				UserID:    "123",
				ContentID: 123,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &ShareContentInput{
				UserID:    tt.fields.UserID,
				ContentID: tt.fields.ContentID,
				Channel:   tt.fields.Channel,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ShareContentInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
