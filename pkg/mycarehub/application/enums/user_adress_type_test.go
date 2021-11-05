package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestAddressesType_String(t *testing.T) {
	tests := []struct {
		name string
		e    AddressesType
		want string
	}{
		{
			name: "POSTAL",
			e:    AddressesTypePostal,
			want: "POSTAL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("AddressesType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressesType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    AddressesType
		want bool
	}{
		{
			name: "valid type",
			e:    AddressesTypePostal,
			want: true,
		},
		{
			name: "invalid type",
			e:    AddressesType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("AddressesType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressesType_UnmarshalGQL(t *testing.T) {
	value := AddressesTypePostal
	invalid := AddressesType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *AddressesType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "POSTAL",
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
				t.Errorf("AddressesType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddressesType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     AddressesType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     AddressesTypePostal,
			b:     w,
			wantW: strconv.Quote("POSTAL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("AddressesType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
