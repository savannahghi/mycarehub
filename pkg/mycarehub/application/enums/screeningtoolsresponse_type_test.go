package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestScreeningToolResponseType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    ScreeningToolResponseType
		want bool
	}{
		{
			name: "valid type",
			m:    ScreeningToolResponseTypeInteger,
			want: true,
		},
		{
			name: "invalid type",
			m:    ScreeningToolResponseType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("ScreeningToolResponseType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningToolResponseType_String(t *testing.T) {
	tests := []struct {
		name string
		m    ScreeningToolResponseType
		want string
	}{
		{
			name: "INTEGER",
			m:    ScreeningToolResponseTypeInteger,
			want: "INTEGER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("ScreeningToolResponseType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningToolResponseType_UnmarshalGQL(t *testing.T) {
	value := ScreeningToolResponseTypeInteger
	invalid := ScreeningToolResponseType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *ScreeningToolResponseType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "INTEGER",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			m:    &invalid,
			args: args{
				v: "this is not a valid type",
			},
			wantErr: true,
		},
		{
			name: "non string type",
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
				t.Errorf("ScreeningToolResponseType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScreeningToolResponseType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		m     ScreeningToolResponseType
		b     *bytes.Buffer
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     ScreeningToolResponseTypeInteger,
			b:     w,
			wantW: strconv.Quote("INTEGER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ScreeningToolResponseType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
