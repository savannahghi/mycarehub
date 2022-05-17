package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestMetricType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    MetricType
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			m:    MetricTypeContent,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			m:    MetricType("Not a metric type"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("MetricType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricType_String(t *testing.T) {
	tests := []struct {
		name string
		m    MetricType
		want string
	}{
		{
			name: "Happy Case",
			m:    MetricTypeContent,
			want: MetricTypeContent.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("MetricType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricType_UnmarshalGQL(t *testing.T) {
	validValue := MetricTypeContent
	invalidType := MetricType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *MetricType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: MetricTypeContent.String(),
			},
			m:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			m:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(int)",
			args: args{
				v: 45,
			},
			m:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("MetricType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		m     MetricType
		wantW string
	}{
		{
			name:  "valid type enums",
			m:     MetricTypeContent,
			wantW: strconv.Quote("CONTENT"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.m.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MetricType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
