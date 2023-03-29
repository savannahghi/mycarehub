package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestUserIdentifierType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    UserIdentifierType
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    UserIdentifierTypeCCC,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    UserIdentifierType("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("UserIdentifierType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserIdentifierType_String(t *testing.T) {
	tests := []struct {
		name string
		c    UserIdentifierType
		want string
	}{
		{
			name: "Happy Case",
			c:    UserIdentifierTypeCCC,
			want: UserIdentifierTypeCCC.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("UserIdentifierType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserIdentifierType_UnmarshalGQL(t *testing.T) {
	validValue := UserIdentifierTypeCCC
	invalidType := UserIdentifierType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *UserIdentifierType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: UserIdentifierTypeCCC.String(),
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
				t.Errorf("UserIdentifierType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserIdentifierType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     UserIdentifierType
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     UserIdentifierTypeCCC,
			wantW: strconv.Quote("CCC"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UserIdentifierType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
