package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestSecurityQuestionResponseType_String(t *testing.T) {
	tests := []struct {
		name string
		e    SecurityQuestionResponseType
		want string
	}{
		{
			name: "NUMBER",
			e:    NumberResponse,
			want: "NUMBER",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("SecurityQuestionResponseType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityQuestionResponseType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    SecurityQuestionResponseType
		want bool
	}{
		{
			name: "valid type",
			e:    NumberResponse,
			want: true,
		},
		{
			name: "invalid type",
			e:    SecurityQuestionResponseType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("SecurityQuestionResponseType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSecurityQuestionResponseType_UnmarshalGQL(t *testing.T) {
	value := DateResponse
	invalid := SecurityQuestionResponseType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *SecurityQuestionResponseType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			e:    &value,
			args: args{
				v: "DATE",
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
				t.Errorf("SecurityQuestionResponseType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSecurityQuestionResponseType_MarshalGQL(t *testing.T) {
	w := &bytes.Buffer{}
	tests := []struct {
		name  string
		e     SecurityQuestionResponseType
		b     *bytes.Buffer
		wantW string
		panic bool
	}{
		{
			name:  "valid type enums",
			e:     StringResponse,
			b:     w,
			wantW: strconv.Quote("STRING"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.MarshalGQL(tt.b)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SecurityQuestionResponseType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
