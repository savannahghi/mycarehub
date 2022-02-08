package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestPermissionType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    PermissionType
		want bool
	}{
		{
			name: "valid type",
			m:    PermissionTypeCanResetUserPassword,
			want: true,
		},
		{
			name: "invalid type",
			m:    PermissionType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("PermissionType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionType_String(t *testing.T) {
	tests := []struct {
		name string
		m    PermissionType
		want string
	}{
		{
			name: "CAN_RESET_USER_PASSWORD",
			m:    PermissionTypeCanResetUserPassword,
			want: "CAN_RESET_USER_PASSWORD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("PermissionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionType_UnmarshalGQL(t *testing.T) {
	value := PermissionTypeCanResetUserPassword
	invalid := PermissionType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *PermissionType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "CAN_RESET_USER_PASSWORD",
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
				t.Errorf("PermissionType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		m     PermissionType
		b     *bytes.Buffer
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     PermissionTypeCanResetUserPassword,
			b:     w,
			wantW: strconv.Quote("CAN_RESET_USER_PASSWORD"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PermissionType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
