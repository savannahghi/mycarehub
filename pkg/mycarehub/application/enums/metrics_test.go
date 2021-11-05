package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestMetricType_String(t *testing.T) {
	tests := []struct {
		name string
		e    MetricType
		want string
	}{
		{
			name: "Engagement",
			e:    EngagementMetrics,
			want: "Engagement",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MetricType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MetricType
		want bool
	}{
		{
			name: "valid type",
			e:    EngagementMetrics,
			want: true,
		},
		{
			name: "invalid type",
			e:    MetricType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MetricType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricType_UnmarshalGQL(t *testing.T) {
	value := EngagementMetrics
	invalid := MetricType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MetricType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "Engagement",
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
				t.Errorf("MetricType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetricType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     MetricType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     EngagementMetrics,
			b:     w,
			wantW: strconv.Quote("Engagement"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MetricType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
