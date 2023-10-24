package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestClientType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ClientType
		want bool
	}{
		{
			name: "Happy Case - Valid client",
			e:    ClientTypeDreams,
			want: true,
		},
		{
			name: "Sad Case - Invalid client",
			e:    ClientType("Invalid Client"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("ClientType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientType_String(t *testing.T) {
	tests := []struct {
		name string
		e    ClientType
		want string
	}{
		{
			name: "Happy Case - Valid string",
			e:    ClientTypeHighRisk,
			want: ClientTypeHighRisk.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("ClientType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientType_UnmarshalGQL(t *testing.T) {
	validValue := ClientTypeHvl
	invalidType := ClientType("Invalid")

	type args struct {
		v interface{}
	}

	tests := []struct {
		name    string
		e       *ClientType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid Type",
			args: args{
				v: ClientTypeHvl.String(),
			},
			e:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - invalid Type",
			args: args{
				v: "invalid type",
			},
			e:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(int)",
			args: args{
				v: 45,
			},
			e:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("ClientType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     ClientType
		wantW string
	}{
		{
			name:  "Valid type enums",
			e:     ClientTypeOtz,
			wantW: strconv.Quote("OTZ"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ClientType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
