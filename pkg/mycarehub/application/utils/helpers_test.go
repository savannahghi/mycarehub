package utils

import (
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

func TestConvertJsonStringToMap(t *testing.T) {
	type args struct {
		jsonString string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				jsonString: `{"key":"value"}`,
			},
			want: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name: "invalid json",
			args: args{
				jsonString: `{"key":"value`,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty json",
			args: args{
				jsonString: ``,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "non map json",
			args: args{
				jsonString: `["yes","no"]`,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertJSONStringToMap(tt.args.jsonString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertJSONStringToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertJSONStringToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
