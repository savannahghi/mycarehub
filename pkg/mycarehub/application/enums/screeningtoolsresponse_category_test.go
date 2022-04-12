package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestScreeningToolResponseCategory_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    ScreeningToolResponseCategory
		want bool
	}{
		{
			name: "valid type",
			m:    ScreeningToolResponseCategorySingleChoice,
			want: true,
		},
		{
			name: "invalid type",
			m:    ScreeningToolResponseCategory("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("ScreeningToolResponseCategory.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningToolResponseCategory_String(t *testing.T) {
	tests := []struct {
		name string
		m    ScreeningToolResponseCategory
		want string
	}{
		{
			name: "SINGLE_CHOICE",
			m:    ScreeningToolResponseCategorySingleChoice,
			want: "SINGLE_CHOICE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("ScreeningToolResponseCategory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScreeningToolResponseCategory_UnmarshalGQL(t *testing.T) {
	value := ScreeningToolResponseCategorySingleChoice
	invalid := ScreeningToolResponseCategory("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *ScreeningToolResponseCategory
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "SINGLE_CHOICE",
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
				t.Errorf("ScreeningToolResponseCategory.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScreeningToolResponseCategory_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		m     ScreeningToolResponseCategory
		b     *bytes.Buffer
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     ScreeningToolResponseCategorySingleChoice,
			b:     w,
			wantW: strconv.Quote("SINGLE_CHOICE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ScreeningToolResponseCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
