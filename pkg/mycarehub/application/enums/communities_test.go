package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestMessageType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    MessageType
		want bool
	}{
		{
			name: "Happy Case - Valid Message",
			e:    MessageTypeRegular,
			want: true,
		},
		{
			name: "Sad Case - Invalid Message",
			e:    MessageType("Invalid Message"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("MessageType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageType_String(t *testing.T) {
	tests := []struct {
		name string
		e    MessageType
		want string
	}{
		{
			name: "Happy Case - Valid string",
			e:    MessageTypeRegular,
			want: MessageTypeRegular.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("MessageType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageType_UnmarshalGQL(t *testing.T) {
	validValue := MessageTypeRegular
	invalidType := MessageType("Invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *MessageType
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid Type",
			args: args{
				v: MessageTypeRegular.String(),
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
				t.Errorf("MessageType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     MessageType
		wantW string
	}{
		{
			name:  "Valid type enums",
			e:     MessageTypeRegular,
			wantW: strconv.Quote("regular"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("MessageType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
