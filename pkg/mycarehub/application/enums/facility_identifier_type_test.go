package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestFacilityIdentifierType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    FacilityIdentifierType
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    FacilityIdentifierTypeMFLCode,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    FacilityIdentifierType("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("FacilityIdentifierType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacilityIdentifierType_String(t *testing.T) {
	tests := []struct {
		name string
		c    FacilityIdentifierType
		want string
	}{
		{
			name: "Happy Case",
			c:    FacilityIdentifierTypeMFLCode,
			want: FacilityIdentifierTypeMFLCode.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("FacilityIdentifierType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFacilityIdentifierType_UnmarshalGQL(t *testing.T) {
	validValue := FacilityIdentifierTypeMFLCode
	invalidType := FacilityIdentifierType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *FacilityIdentifierType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: FacilityIdentifierTypeMFLCode.String(),
			},
			c:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			c:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(float)",
			args: args{
				v: 45.1,
			},
			c:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FacilityIdentifierType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFacilityIdentifierType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     FacilityIdentifierType
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     FacilityIdentifierTypeMFLCode,
			wantW: strconv.Quote("MFL_CODE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FacilityIdentifierType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
