package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestCountyType_String(t *testing.T) {
	tests := []struct {
		name string
		e    CountyType
		want string
	}{
		{
			name: "Mombasa",
			e:    CountyTypeMombasa,
			want: "Mombasa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("CountyType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountyType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    CountyType
		want bool
	}{
		{
			name: "valid type",
			e:    CountyTypeMombasa,
			want: true,
		},
		{
			name: "invalid type",
			e:    CountyType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("CountyType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountyType_UnmarshalGQL(t *testing.T) {
	pmtc := CountyTypeMombasa
	invalid := CountyType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *CountyType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &pmtc,
			args: args{
				v: "Mombasa",
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
				t.Errorf("CountyType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCountyType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     CountyType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     CountyTypeMombasa,
			b:     w,
			wantW: strconv.Quote("Mombasa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("CountyType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
