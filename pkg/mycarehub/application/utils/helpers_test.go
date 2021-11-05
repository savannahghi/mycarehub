package utils

import (
	"testing"
	"time"

	"github.com/tj/assert"
)

const (
	// DefaultSaltLen is the length of generated salt for the user is 256
	DefaultSaltLen = 256
	// DefaultKeyLen is the length of encoded key in PBKDF2 function is 512
	DefaultKeyLen = 512
)

func TestGetHourMinuteSecond(t *testing.T) {
	type args struct {
		hour   time.Duration
		minute time.Duration
		second time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "Happy case",
			args: args{
				hour:   20,
				minute: 1,
				second: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetHourMinuteSecond(tt.args.hour, tt.args.minute, tt.args.second)
			assert.NotNil(t, got)
		})
	}
}

func TestNextAllowedLoginTime(t *testing.T) {
	type args struct {
		trials int
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{
			name: "Happy case",
			args: args{
				trials: 3,
			},
			wantError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextAllowedLoginTime(tt.args.trials)
			assert.NotNil(t, got)
		})
	}
}

func TestValidatePINLength(t *testing.T) {
	type args struct {
		pin string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully validate pin length",
			args: args{
				pin: "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid Pin",
			args: args{
				pin: "123456789",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePINLength(tt.args.pin); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePINLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePIN(t *testing.T) {
	type args struct {
		pin string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully validate pin",
			args: args{
				pin: "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid pin length",
			args: args{
				pin: "12",
			},
			wantErr: true,
		},
		{
			name: "Sad Case - invalid pin",
			args: args{
				pin: "asdf",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePIN(tt.args.pin); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
