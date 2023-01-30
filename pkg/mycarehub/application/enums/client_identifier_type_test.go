package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestClientIdentifierType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    ClientIdentifierType
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    ClientIdentifierTypeCCC,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    ClientIdentifierType("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("ClientIdentifierType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientIdentifierType_String(t *testing.T) {
	tests := []struct {
		name string
		c    ClientIdentifierType
		want string
	}{
		{
			name: "Happy Case",
			c:    ClientIdentifierTypeCCC,
			want: ClientIdentifierTypeCCC.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("ClientIdentifierType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientIdentifierType_UnmarshalGQL(t *testing.T) {
	validValue := ClientIdentifierTypeCCC
	invalidType := ClientIdentifierType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *ClientIdentifierType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: ClientIdentifierTypeCCC.String(),
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
				t.Errorf("ClientIdentifierType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientIdentifierType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     ClientIdentifierType
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     ClientIdentifierTypeCCC,
			wantW: strconv.Quote("CCC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ClientIdentifierType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
