package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestTerminologies_String(t *testing.T) {
	tests := []struct {
		name string
		e    Terminologies
		want string
	}{
		{
			name: "CIEL",
			e:    TerminologiesCIEL,
			want: "CIEL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Terminologies.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTerminologies_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Terminologies
		want bool
	}{
		{
			name: "valid type",
			e:    TerminologiesCIEL,
			want: true,
		},
		{
			name: "invalid type",
			e:    Terminologies("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Terminologies.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTerminologies_UnmarshalGQL(t *testing.T) {
	value := TerminologiesCIEL
	invalid := Terminologies("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *Terminologies
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "CIEL",
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
				t.Errorf("Terminologies.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTerminologies_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     Terminologies
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     TerminologiesCIEL,
			b:     w,
			wantW: strconv.Quote("CIEL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Terminologies.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
