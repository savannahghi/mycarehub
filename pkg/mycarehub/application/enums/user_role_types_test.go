package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestUserRoleType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    UserRoleType
		want bool
	}{
		{
			name: "valid type",
			m:    UserRoleTypeClientManagement,
			want: true,
		},
		{
			name: "invalid type",
			m:    UserRoleType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("UserRoleType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRoleType_String(t *testing.T) {
	tests := []struct {
		name string
		m    UserRoleType
		want string
	}{
		{
			name: "CLIENT_MANAGEMENT",
			m:    UserRoleTypeClientManagement,
			want: "CLIENT_MANAGEMENT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("UserRoleType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRoleType_UnmarshalGQL(t *testing.T) {
	value := UserRoleTypeClientManagement
	invalid := UserRoleType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *UserRoleType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "CLIENT_MANAGEMENT",
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
				t.Errorf("UserRoleType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRoleType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		m     UserRoleType
		b     *bytes.Buffer
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     UserRoleTypeClientManagement,
			b:     w,
			wantW: strconv.Quote("CLIENT_MANAGEMENT"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UserRoleType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
