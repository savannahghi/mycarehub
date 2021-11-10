package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestFlavour_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		m     Flavour
		wantW string
	}{
		{
			name:  "valid enums",
			m:     PRO,
			wantW: strconv.Quote("PRO"),
		},
		{
			name:  "valid enums",
			m:     CONSUMER,
			wantW: strconv.Quote("CONSUMER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Flavour.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestFlavour_UnmarshalGQL(t *testing.T) {
	value := PRO
	invalid := Flavour("invalid")

	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *Flavour
		args    args
		wantErr bool
	}{
		{
			name: "valid enum",
			m:    &value,
			args: args{
				v: "PRO",
			},
			wantErr: false,
		},
		{
			name: "invalid enum",
			m:    &invalid,
			args: args{
				v: "INVALID",
			},
			wantErr: true,
		},
		{
			name: "non string enum",
			m:    &invalid,
			args: args{
				v: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Flavour.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFlavour_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    Flavour
		want bool
	}{
		{
			name: "valid type",
			m:    PRO,
			want: true,
		},
		{
			name: "invalid type",
			m:    Flavour("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("Flavour.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlavour_String(t *testing.T) {
	tests := []struct {
		name string
		m    Flavour
		want string
	}{
		{
			name: "CONSUMER",
			m:    CONSUMER,
			want: "CONSUMER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Flavour.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
