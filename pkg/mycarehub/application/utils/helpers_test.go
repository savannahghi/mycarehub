package utils

import (
	"reflect"
	"testing"
	"time"

	"github.com/savannahghi/scalarutils"
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
			want:    map[string]interface{}{},
			wantErr: false,
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

func TestConvertTimeToScalarDate(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    scalarutils.Date
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				t: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
			want:    scalarutils.Date{Year: 2020, Month: 1, Day: 1},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				t: time.Time{},
			},
			want:    scalarutils.Date{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertTimeToScalarDate(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertTimeToScalarDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertTimeToScalarDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_InterfaceToInt(t *testing.T) {
	type args struct {
		n interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Happy Case:  initialize existing int",
			args: args{
				n: 130,
			},
			want: 130,
		},
		{
			name: "Sad Case:  nil input",
			args: args{
				n: nil,
			},
			want: 0,
		},
		{
			name: "Sad Case:  invalid input",
			args: args{
				n: "",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceToInt(tt.args.n); got != tt.want {
				t.Errorf("InterfaceToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interfaceToString(t *testing.T) {
	type args struct {
		n interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test_interfaceToString:  initialize existing string",
			args: args{
				n: "130",
			},
			want: "130",
		},
		{
			name: "Test_interfaceToString:  initialize string",
			args: args{
				n: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceToString(tt.args.n); got != tt.want {
				t.Errorf("InterfaceToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateOTP(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "Happy case - generate OTP",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateOTP()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOTP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTruncateMatrixUserID(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Happy case: truncates string",
			args: args{
				userID: "@abiudrn:prohealth360.org",
			},
			want: "abiudrn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateMatrixUserID(tt.args.userID); got != tt.want {
				t.Errorf("TruncateMatrixUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceToFloat64(t *testing.T) {
	type args struct {
		n interface{}
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "Happy case: convert interface to float64",
			args: args{
				n: 4.0,
			},
			want: 4.0,
		},
		{
			name: "Sad case: nil input",
			args: args{
				n: nil,
			},
			want: 0.0,
		},
		{
			name: "Sad case: invalid type passed",
			args: args{
				n: "",
			},
			want: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceToFloat64(tt.args.n); got != tt.want {
				t.Errorf("InterfaceToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}
