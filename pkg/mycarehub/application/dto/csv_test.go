package dto

import (
	"reflect"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestFacilityCSVOutput_ValidateLabels(t *testing.T) {
	type fields struct {
		Code            string
		IdentifierType  enums.FacilityIdentifierType
		Name            string
		Level           string
		FacilityType    string
		OwnerType       string
		RegulatoryBody  string
		Country         string
		County          string
		OperationStatus string
		Contact         string
	}
	type args struct {
		labels []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Valid labels",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "contact",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid label length",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "contact", "extraLabel",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid label",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "invalidContactLabel",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FacilityCSVOutput{
				Code:            tt.fields.Code,
				IdentifierType:  tt.fields.IdentifierType,
				Name:            tt.fields.Name,
				Level:           tt.fields.Level,
				FacilityType:    tt.fields.FacilityType,
				OwnerType:       tt.fields.OwnerType,
				RegulatoryBody:  tt.fields.RegulatoryBody,
				Country:         tt.fields.Country,
				County:          tt.fields.County,
				OperationStatus: tt.fields.OperationStatus,
				Contact:         tt.fields.Contact,
			}
			if err := f.ValidateLabels(tt.args.labels); (err != nil) != tt.wantErr {
				t.Errorf("FacilityCSVOutput.ValidateLabels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFacilityCSVOutput_ParseValues(t *testing.T) {
	type fields struct {
		Code            string
		IdentifierType  enums.FacilityIdentifierType
		Name            string
		Level           string
		FacilityType    string
		OwnerType       string
		RegulatoryBody  string
		Country         string
		County          string
		OperationStatus string
		Contact         string
	}
	type args struct {
		labels []string
		values []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *FacilityCSVOutput
		wantErr bool
	}{
		{
			name: "Happy Case: parse values",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "contact",
				},
				values: []string{"25582", "MFL_CODE", "The Nairobi Hospital (Capital Centre)",
					"Level 2", "Medical Clinic", "Private Practice", "Kenya MPDB",
					"Kenya", "Nairobi", "Operational", "0202845000"},
			},
			want: &FacilityCSVOutput{
				Code:            "25582",
				IdentifierType:  "MFL_CODE",
				Name:            "The Nairobi Hospital (Capital Centre)",
				Level:           "Level 2",
				FacilityType:    "Medical Clinic",
				OwnerType:       "Private Practice",
				RegulatoryBody:  "Kenya MPDB",
				Country:         "Kenya",
				County:          "Nairobi",
				OperationStatus: "Operational",
				Contact:         "0202845000",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid contact",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "contact",
				},
				values: []string{"25582", "MFL_CODE", "The Nairobi Hospital (Capital Centre)",
					"Level 2", "Medical Clinic", "Private Practice", "Kenya MPDB",
					"Kenya", "Nairobi", "Operational", "invalidPhone"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad Case: invalid identifier",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "contact",
				},
				values: []string{"25582", "INVALID_IDENTIFIER", "The Nairobi Hospital (Capital Centre)",
					"Level 2", "Medical Clinic", "Private Practice", "Kenya MPDB",
					"Kenya", "Nairobi", "Operational", "0202845000"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad Case: invalid values length",
			fields: fields{
				Code:            "",
				IdentifierType:  "",
				Name:            "",
				Level:           "",
				FacilityType:    "",
				OwnerType:       "",
				RegulatoryBody:  "",
				Country:         "",
				County:          "",
				OperationStatus: "",
				Contact:         "",
			},
			args: args{
				labels: []string{"code", "identifierType", "name", "level",
					"facilityType", "ownerType", "regulatoryBody", "country",
					"county", "operationStatus", "contact",
				},
				values: []string{"25582", "MFL_CODE", "The Nairobi Hospital (Capital Centre)",
					"Level 2", "Medical Clinic", "Private Practice", "Kenya MPDB",
					"Kenya", "Nairobi", "Operational", "0202845000", "extraField"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FacilityCSVOutput{
				Code:            tt.fields.Code,
				IdentifierType:  tt.fields.IdentifierType,
				Name:            tt.fields.Name,
				Level:           tt.fields.Level,
				FacilityType:    tt.fields.FacilityType,
				OwnerType:       tt.fields.OwnerType,
				RegulatoryBody:  tt.fields.RegulatoryBody,
				Country:         tt.fields.Country,
				County:          tt.fields.County,
				OperationStatus: tt.fields.OperationStatus,
				Contact:         tt.fields.Contact,
			}
			got, err := f.ParseValues(tt.args.labels, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("FacilityCSVOutput.ParseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FacilityCSVOutput.ParseValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
