package utils

import (
	"testing"
	"time"

	"github.com/tj/assert"
)

func TestCalculateNextAllowedLoginTime(t *testing.T) {
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
			name: "Happy Case - Success",
			args: args{
				hour:   0,
				minute: 0,
				second: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateNextAllowedLoginTime(tt.args.hour, tt.args.minute, tt.args.second)
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
