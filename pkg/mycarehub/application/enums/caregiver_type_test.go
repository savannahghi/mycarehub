package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestCaregiverType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    CaregiverType
		want bool
	}{
		{
			name: "valid type",
			m:    CaregiverTypeFather,
			want: true,
		},
		{
			name: "invalid type",
			m:    CaregiverType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("CaregiverType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCaregiverType_String(t *testing.T) {
	tests := []struct {
		name string
		m    CaregiverType
		want string
	}{
		{
			name: "SIBLING",
			m:    CaregiverTypeSibling,
			want: "SIBLING",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("CaregiverType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCaregiverType_UnmarshalGQL(t *testing.T) {
	value := CaregiverTypeSibling
	invalid := CaregiverType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *CaregiverType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "SIBLING",
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
				t.Errorf("CaregiverType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCaregiverType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		m     CaregiverType
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     CaregiverTypeSibling,
			wantW: strconv.Quote("SIBLING"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("CaregiverType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
