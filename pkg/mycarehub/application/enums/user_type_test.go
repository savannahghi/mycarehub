package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestUsersType_String(t *testing.T) {
	tests := []struct {
		name string
		e    UsersType
		want string
	}{
		{
			name: "HEALTHCAREWORKER",
			e:    HealthcareWorkerUser,
			want: "HEALTHCAREWORKER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("UsersType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsersType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    UsersType
		want bool
	}{
		{
			name: "valid type",
			e:    HealthcareWorkerUser,
			want: true,
		},
		{
			name: "invalid type",
			e:    UsersType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("UsersType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsersType_UnmarshalGQL(t *testing.T) {
	value := HealthcareWorkerUser
	invalid := UsersType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *UsersType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "HEALTHCAREWORKER",
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
				t.Errorf("UsersType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     UsersType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     HealthcareWorkerUser,
			b:     w,
			wantW: strconv.Quote("HEALTHCAREWORKER"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UsersType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
