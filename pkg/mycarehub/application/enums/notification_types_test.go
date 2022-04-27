package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestNotificationType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    NotificationType
		want bool
	}{
		{
			name: "valid type",
			m:    NotificationTypeAppointment,
			want: true,
		},
		{
			name: "invalid type",
			m:    NotificationType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("NotificationType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationType_String(t *testing.T) {
	tests := []struct {
		name string
		m    NotificationType
		want string
	}{
		{
			name: "APPOINTMENT",
			m:    NotificationTypeAppointment,
			want: "APPOINTMENT",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("NotificationType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotificationType_UnmarshalGQL(t *testing.T) {
	value := NotificationTypeAppointment
	invalid := NotificationType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *NotificationType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			m:    &value,
			args: args{
				v: "APPOINTMENT",
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
				t.Errorf("NotificationType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotificationType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		m     NotificationType
		b     *bytes.Buffer
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     NotificationTypeAppointment,
			b:     w,
			wantW: strconv.Quote("APPOINTMENT"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("NotificationType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
