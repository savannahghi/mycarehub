package dto

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestFacilityInput_Validate(t *testing.T) {
	longWord := gofakeit.Sentence(100)
	veryLongWord := gofakeit.Sentence(500)

	type fields struct {
		Name        string
		Code        string
		Active      bool
		County      enums.CountyType
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
				Code:        "22344",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: "test description",
			},
			wantErr: false,
		},

		{
			name: "invalid: short name len",
			fields: fields{
				Name:        "te",
				Code:        "22344",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: long name len",
			fields: fields{
				Name:        longWord,
				Code:        "22344",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: short description",
			fields: fields{
				Name:        "test name",
				Code:        "22344",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: "te",
			},
			wantErr: true,
		},
		{
			name: "invalid: very long description",
			fields: fields{
				Name:        "test name",
				Code:        "22344",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: veryLongWord,
			},
			wantErr: true,
		},
		{
			name: "invalid: missing name",
			fields: fields{
				Code:        "22344",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing code",
			fields: fields{
				Name:        "test name",
				Active:      true,
				County:      enums.CountyTypeNairobi,
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing county",
			fields: fields{
				Name:        "test name",
				Code:        "22344",
				Active:      true,
				Description: "test description",
			},
			wantErr: true,
		},
		{
			name: "invalid: missing description",
			fields: fields{
				Name:   "test name",
				Code:   "22344",
				Active: true,
				County: enums.CountyTypeNairobi,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FacilityInput{
				Name:        tt.fields.Name,
				Code:        tt.fields.Code,
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
		Flavour     enums.Flavour
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
				Flavour:     enums.CONSUMER,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing phone number",
			fields: fields{
				PIN:     &testPIN,
				Flavour: enums.CONSUMER,
			},
			wantErr: true,
		},
		{
			name: "invalid : missing pin",
			fields: fields{
				PhoneNumber: &testPhone,
				Flavour:     enums.CONSUMER,
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
