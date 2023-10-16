package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestCountry_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    Country
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    CountryKenya,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    Country("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("Country.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountry_String(t *testing.T) {
	tests := []struct {
		name string
		c    Country
		want string
	}{
		{
			name: "Happy Case",
			c:    CountryKenya,
			want: CountryKenya.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Country.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountry_UnmarshalGQL(t *testing.T) {
	validValue := CountryKenya
	invalidType := Country("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *Country
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: CountryKenya.String(),
			},
			c:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			c:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(float)",
			args: args{
				v: 45.1,
			},
			c:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Country.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCountry_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     Country
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     CountryKenya,
			wantW: strconv.Quote("KE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Country.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
