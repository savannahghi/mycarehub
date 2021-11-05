package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestCountryType_String(t *testing.T) {
	tests := []struct {
		name string
		e    CountryType
		want string
	}{
		{
			name: "KENYA",
			e:    CountryTypeKenya,
			want: "KENYA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("CountryType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountryType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    CountryType
		want bool
	}{
		{
			name: "valid type",
			e:    CountryTypeKenya,
			want: true,
		},
		{
			name: "invalid type",
			e:    CountryType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("CountryType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountryType_UnmarshalGQL(t *testing.T) {
	pmtc := CountryTypeKenya
	invalid := CountryType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *CountryType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &pmtc,
			args: args{
				v: "KENYA",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			e:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
			e:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("CountryType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCountryType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     CountryType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     CountryTypeKenya,
			b:     w,
			wantW: strconv.Quote("KENYA"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("CountryType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
