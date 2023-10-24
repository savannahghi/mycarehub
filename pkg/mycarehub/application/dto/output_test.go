package dto

import (
	"reflect"
	"testing"
)

func TestOrganisationOutput_ParseValues(t *testing.T) {
	type fields struct {
		Name            string
		Description     string
		EmailAddress    string
		PhoneNumber     string
		PostalAddress   string
		PhysicalAddress string
		DefaultCountry  string
	}
	type args struct {
		values []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *OrganisationInput
		wantErr bool
	}{
		{
			name: "Happy Case: parse values",
			fields: fields{
				Name:            "",
				Description:     "",
				EmailAddress:    "",
				PhoneNumber:     "",
				PostalAddress:   "",
				PhysicalAddress: "",
				DefaultCountry:  "",
			},
			args: args{
				values: []byte(`{
					"name": "test",
					"description": "test",
					"emailAddress": "test@test.org",
					"phoneNumber": "0999999999",
					"postalAddress": "test",
					"physicalAddress": "test",
					"defaultCountry": "test"
				}`),
			},
			want: &OrganisationInput{
				Name:            "test",
				Description:     "test",
				EmailAddress:    "test@test.org",
				PhoneNumber:     "0999999999",
				PostalAddress:   "test",
				PhysicalAddress: "test",
				DefaultCountry:  "test",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid payload",
			fields: fields{
				Name:            "",
				Description:     "",
				EmailAddress:    "",
				PhoneNumber:     "",
				PostalAddress:   "",
				PhysicalAddress: "",
				DefaultCountry:  "",
			},
			args: args{
				values: []byte(`{
					"invalid": "test",
					"description": "test",
					"emailAddress": "test@test.org",
					"phoneNumber": "0999999999",
					"postalAddress": "test",
					"physicalAddress": "test",
					"defaultCountry": "test"
				}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad Case: empty name",
			fields: fields{
				Name:            "",
				Description:     "",
				EmailAddress:    "",
				PhoneNumber:     "",
				PostalAddress:   "",
				PhysicalAddress: "",
				DefaultCountry:  "",
			},
			args: args{
				values: []byte(`{
					"name": "",
					"description": "test",
					"emailAddress": "test@test.org",
					"phoneNumber": "0999999999",
					"postalAddress": "test",
					"physicalAddress": "test",
					"defaultCountry": "test"
				}`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad Case: nil name",
			fields: fields{
				Name:            "",
				Description:     "",
				EmailAddress:    "",
				PhoneNumber:     "",
				PostalAddress:   "",
				PhysicalAddress: "",
				DefaultCountry:  "",
			},
			args: args{
				values: []byte(`{
					"description": "test",
					"emailAddress": "test@test.org",
					"phoneNumber": "0999999999",
					"postalAddress": "test",
					"physicalAddress": "test",
					"defaultCountry": "test"
				}`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrganisationOutput{
				Name:            tt.fields.Name,
				Description:     tt.fields.Description,
				EmailAddress:    tt.fields.EmailAddress,
				PhoneNumber:     tt.fields.PhoneNumber,
				PostalAddress:   tt.fields.PostalAddress,
				PhysicalAddress: tt.fields.PhysicalAddress,
				DefaultCountry:  tt.fields.DefaultCountry,
			}
			got, err := o.ParseValues(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrganisationOutput.ParseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrganisationOutput.ParseValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgramJSONOutput_ParseValues(t *testing.T) {
	type fields struct {
		Name string
	}
	type args struct {
		values []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ProgramInput
		wantErr bool
	}{
		{
			name: "Happy case: parse values",
			fields: fields{
				Name: "",
			},
			args: args{
				values: []byte(`
					{
						"name": "test",
						"description": "test"
					}
				`),
			},
			want: &ProgramInput{
				Name:        "test",
				Description: "test",
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid, empty description",
			fields: fields{
				Name: "",
			},
			args: args{
				values: []byte(`
					{
						"name": "test",
						"description": ""
					}
				`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: invalid, missing description",
			fields: fields{
				Name: "",
			},
			args: args{
				values: []byte(`
					{
						"name": "test",
					}
				`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: invalid payload",
			fields: fields{
				Name: "",
			},
			args: args{
				values: []byte(`
					{
						"invalid": "test"
					}
				`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProgramJSONOutput{
				Name: tt.fields.Name,
			}
			got, err := p.ParseValues(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProgramJSONOutput.ParseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProgramJSONOutput.ParseValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseValues(t *testing.T) {
	type args struct {
		values []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ProgramJSONOutput
		wantErr bool
	}{
		{
			name: "Happy case: parse values",
			args: args{
				values: []byte(`
					{
						"name": "test",
						"description": "test"
					}
				`),
			},
			want: &ProgramJSONOutput{
				Name:        "test",
				Description: "test",
			},
			wantErr: false,
		},
		{
			name: "Sad case: empty name",
			args: args{
				values: []byte(`
					{
						"name": "",
						"description": "test"
					}
				`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: nil name",
			args: args{
				values: []byte(`
					{
						"description": "test"
					}
				`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case: invalid json",
			args: args{
				values: []byte(`invalid json`),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseValues(ProgramJSONOutput{}, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
