package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestServiceRequestType_Number(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceRequestType
		want string
	}{
		{
			name: "HEALTH_DIARY_ENTRY",
			e:    ServiceRequestTypeHealthDiaryEntry,
			want: "HEALTH_DIARY_ENTRY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ServiceRequestType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequestType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ServiceRequestType
		want bool
	}{
		{
			name: "valid type",
			e:    ServiceRequestTypeHealthDiaryEntry,
			want: true,
		},
		{
			name: "invalid type",
			e:    ServiceRequestType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ServiceRequestType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceRequestType_UnmarshalGQL(t *testing.T) {
	value := ServiceRequestTypeRedFlag
	invalid := ServiceRequestType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *ServiceRequestType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "RED_FLAG",
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
				t.Errorf("ServiceRequestType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceRequestType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     ServiceRequestType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     ServiceRequestTypeHealthDiaryEntry,
			b:     w,
			wantW: strconv.Quote("HEALTH_DIARY_ENTRY"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ServiceRequestType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
