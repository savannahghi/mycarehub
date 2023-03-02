package authorization

import (
	"bytes"
	"strconv"
	"testing"
)

func TestPermissionCategory_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    PermissionCategory
		want bool
	}{
		{
			name: "valid type",
			m:    PermissionCategoryAppointment,
			want: true,
		},
		{
			name: "invalid type",
			m:    PermissionCategory("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("PermissionCategory.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionCategory_String(t *testing.T) {
	tests := []struct {
		name string
		m    PermissionCategory
		want string
	}{
		{
			name: "User",
			m:    PermissionCategoryUser,
			want: "User",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("PermissionCategory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionCategory_UnmarshalGQL(t *testing.T) {
	value := PermissionCategoryUser
	invalid := PermissionCategory("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *PermissionCategory
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "User",
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
				t.Errorf("PermissionCategory.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPermissionCategory_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		m     PermissionCategory
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     PermissionCategoryUser,
			wantW: strconv.Quote("User"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("PermissionCategory.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
