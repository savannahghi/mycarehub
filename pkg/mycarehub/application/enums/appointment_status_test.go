package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestAppointmentStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    AppointmentStatus
		want bool
	}{
		{
			name: "valid type",
			m:    AppointmentStatusScheduled,
			want: true,
		},
		{
			name: "invalid type",
			m:    AppointmentStatus("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("AppointmentStatus.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentStatus_String(t *testing.T) {
	tests := []struct {
		name string
		m    AppointmentStatus
		want string
	}{
		{
			name: "SCHEDULED",
			m:    AppointmentStatusScheduled,
			want: "SCHEDULED",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("AppointmentStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppointmentStatus_UnmarshalGQL(t *testing.T) {
	value := AppointmentStatusScheduled
	invalid := AppointmentStatus("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *AppointmentStatus
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "SCHEDULED",
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
				t.Errorf("AppointmentStatus.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppointmentStatus_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		m     AppointmentStatus
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     AppointmentStatusScheduled,
			wantW: strconv.Quote("SCHEDULED"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AppointmentStatus.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
