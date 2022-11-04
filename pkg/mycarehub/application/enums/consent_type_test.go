package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestConsentState_IsValid(t *testing.T) {
	tests := []struct {
		name string
		c    ConsentState
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			c:    ConsentStateAccepted,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			c:    ConsentState("INVALID"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsValid(); got != tt.want {
				t.Errorf("ConsentState.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentState_String(t *testing.T) {
	tests := []struct {
		name string
		c    ConsentState
		want string
	}{
		{
			name: "Happy Case",
			c:    ConsentStateAccepted,
			want: ConsentStateAccepted.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("ConsentState.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConsentState_UnmarshalGQL(t *testing.T) {
	validValue := ConsentStateAccepted
	invalidType := ConsentState("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		c       *ConsentState
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: ConsentStateAccepted.String(),
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
				t.Errorf("ConsentState.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConsentState_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		c     ConsentState
		wantW string
	}{
		{
			name:  "valid type enums",
			c:     ConsentStateRejected,
			wantW: strconv.Quote("REJECTED"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.c.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("ConsentState.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
