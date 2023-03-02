package authorization

import (
	"bytes"
	"strconv"
	"testing"
)

func TestDefaultRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    DefaultRole
		want bool
	}{
		{
			name: "valid type",
			m:    DefaultRoleAdmin,
			want: true,
		},
		{
			name: "invalid type",
			m:    DefaultRole("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("DefaultRole.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultRole_String(t *testing.T) {
	tests := []struct {
		name string
		m    DefaultRole
		want string
	}{
		{
			name: "Default Admin",
			m:    DefaultRoleAdmin,
			want: "Default Admin",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("DefaultRole.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultRole_UnmarshalGQL(t *testing.T) {
	value := DefaultRoleAdmin
	invalid := DefaultRole("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *DefaultRole
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "Default Admin",
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
				t.Errorf("DefaultRole.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultRole_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		m     DefaultRole
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     DefaultRoleAdmin,
			wantW: strconv.Quote("Default Admin"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("DefaultRole.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
