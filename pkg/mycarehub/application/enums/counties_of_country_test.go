package enums

import (
	"reflect"
	"testing"
)

func TestValidateCountiesOfCountries(t *testing.T) {
	type args struct {
		country CountryType
		county  CountyType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			args: args{
				country: CountryTypeKenya,
				county:  CountyTypeNakuru,
			},
			wantErr: false,
		},
		{
			name: "invalid contry type",
			args: args{
				country: CountryType("invalid"),
				county:  CountyTypeNakuru,
			},
			wantErr: true,
		},
		{
			name: "invalid county type",
			args: args{
				country: CountryTypeKenya,
				county:  CountyType("invalid"),
			},
			wantErr: true,
		},
		// Todo: add test case to validate a county belongs to a country once another country is added.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCountiesOfCountries(tt.args.country, tt.args.county); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCountiesOfCountries() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findSelectedCountryCounties(t *testing.T) {
	countriesCounties := []CountiesOfCountry{
		{
			Country:  CountryTypeKenya,
			Counties: KenyanCounties,
		},
	}
	want1 := CountiesOfCountry{
		Country:  CountryTypeKenya,
		Counties: KenyanCounties,
	}

	type args struct {
		countriesCounties []CountiesOfCountry
		countryInput      CountryType
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 *CountiesOfCountry
	}{
		{
			name: "valid type",
			args: args{
				countriesCounties: countriesCounties,
				countryInput:      CountryTypeKenya,
			},
			want:  true,
			want1: &want1,
		},
		{
			name: "invalid country type",
			args: args{
				countriesCounties: countriesCounties,
				countryInput:      CountryType("invalid"),
			},
			want:  false,
			want1: nil,
		},
		{
			name: "invalid country list type",
			args: args{
				countriesCounties: []CountiesOfCountry{},
				countryInput:      CountryTypeKenya,
			},
			want:  false,
			want1: nil,
		},
		{
			name:  "empty args",
			args:  args{},
			want:  false,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findSelectedCountryCounties(tt.args.countriesCounties, tt.args.countryInput)
			if got != tt.want {
				t.Errorf("findSelectedCountryCounties() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("findSelectedCountryCounties() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_findCounty(t *testing.T) {
	type args struct {
		counties    []CountyType
		countyInput CountyType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			args: args{
				counties:    KenyanCounties,
				countyInput: CountyTypeNakuru,
			},
			wantErr: false,
		},
		{
			name: "invalid county list type",
			args: args{
				counties:    []CountyType{},
				countyInput: CountyTypeNakuru,
			},
			wantErr: true,
		},
		{
			name: "invalid county type",
			args: args{
				counties:    KenyanCounties,
				countyInput: CountyType("invalid"),
			},
			wantErr: true,
		},
		{
			name:    "empty params passed",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := findCounty(tt.args.counties, tt.args.countyInput); (err != nil) != tt.wantErr {
				t.Errorf("findCounty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
