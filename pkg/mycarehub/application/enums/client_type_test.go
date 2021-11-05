package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestClientType_String(t *testing.T) {
	tests := []struct {
		name string
		e    ClientType
		want string
	}{
		{
			name: "PMTCT",
			e:    ClientTypePmtct,
			want: "PMTCT",
		},
		{
			name: "OVC",
			e:    ClientTypeOvc,
			want: "OVC",
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

func TestClientType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    ClientType
		want bool
	}{
		{
			name: "valid PMTCT",
			e:    ClientTypeOvc,
			want: true,
		},
		{
			name: "invalid client type",
			e:    ClientType("invalid"),
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

func TestClientType_UnmarshalGQL(t *testing.T) {
	pmtc := ClientTypePmtct
	invalid := ClientType("invalid")
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
			name: "valid Pmtct client",
			e:    &pmtc,
			args: args{
				v: "PMTCT",
			},
			wantErr: false,
		},
		{
			name: "invalid client type",
			e:    &invalid,
			args: args{
				v: "this is not a valid client type",
			},
			wantErr: true,
		},
		{
			name: "non string client type",
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
			name:  "valid Pmtct client type enums",
			e:     ClientTypePmtct,
			wantW: strconv.Quote("PMTCT"),
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
