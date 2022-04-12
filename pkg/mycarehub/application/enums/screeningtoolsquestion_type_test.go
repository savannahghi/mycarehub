package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestScreeningToolType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    ScreeningToolType
		want bool
	}{
		{
			name: "valid type",
			m:    ScreeningToolTypeTB,
			want: true,
		},
		{
			name: "invalid type",
			m:    ScreeningToolType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("ScreeningToolType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningToolType_String(t *testing.T) {
	tests := []struct {
		name string
		m    ScreeningToolType
		want string
	}{
		{
			name: "TB_ASSESSMENT",
			m:    ScreeningToolTypeTB,
			want: "TB_ASSESSMENT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("ScreeningToolType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningToolType_UnmarshalGQL(t *testing.T) {
	value := ScreeningToolTypeTB
	invalid := ScreeningToolType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *ScreeningToolType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "TB_ASSESSMENT",
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
				t.Errorf("ScreeningToolType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScreeningToolType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		m     ScreeningToolType
		b     *bytes.Buffer
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     ScreeningToolTypeTB,
			b:     w,
			wantW: strconv.Quote("TB_ASSESSMENT"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ScreeningToolType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
