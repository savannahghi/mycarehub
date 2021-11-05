package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestRolesType_String(t *testing.T) {
	tests := []struct {
		name string
		e    RolesType
		want string
	}{
		{
			name: "CAN_REGISTER_STAFF",
			e:    RolesTypeCanRegisterStaff,
			want: "CAN_REGISTER_STAFF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("RolesType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRolesType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    RolesType
		want bool
	}{
		{
			name: "valid type",
			e:    RolesTypeCanRegisterStaff,
			want: true,
		},
		{
			name: "invalid type",
			e:    RolesType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("RolesType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRolesType_UnmarshalGQL(t *testing.T) {
	value := RolesTypeCanRegisterStaff
	invalid := RolesType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *RolesType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "CAN_REGISTER_STAFF",
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
				t.Errorf("RolesType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRolesType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     RolesType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     RolesTypeCanRegisterStaff,
			b:     w,
			wantW: strconv.Quote("CAN_REGISTER_STAFF"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("RolesType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
