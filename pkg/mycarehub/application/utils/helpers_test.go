package utils

import (
	"context"
	"reflect"
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

func TestMakeRequest(t *testing.T) {
	type args struct {
		ctx    context.Context
		method string
		path   string
		body   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case GET",
			args: args{
				ctx:    context.Background(),
				method: "GET",
				path:   "https://google.com/",
				body:   nil,
			},
			wantErr: false,
		},
		{
			name: "Happy case POST",
			args: args{
				ctx:    context.Background(),
				method: "POST",
				path:   "https://google.com/",
				body:   nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeRequest(tt.args.ctx, tt.args.method, tt.args.path, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected facility not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestFormatFilterParamsHelper(t *testing.T) {
	type args struct {
		a map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "test",
			args: args{
				a: map[string]interface{}{
					"status": map[string]interface{}{
						"in": []string{"pending", "open", "new"},
					},
					"members": map[string]interface{}{
						"in": []string{"thierry"},
					},
					"member_count": 2,
				},
			},
			want: map[string]interface{}{
				"status": map[string]interface{}{
					"$in": []string{"pending", "open", "new"},
				},
				"members": map[string]interface{}{
					"$in": []string{"thierry"},
				},
				"member_count": 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFilterParamsHelper(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatFilterParamsHelper() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateAge(t *testing.T) {
	type args struct {
		birthday time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Happy case",
			args: args{
				birthday: time.Now(),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateAge(tt.args.birthday); got != tt.want {
				t.Errorf("CalculateAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
