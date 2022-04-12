package enums

import (
	"bytes"
	"testing"
)

func TestFeedbackType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		f    FeedbackType
		want bool
	}{
		{
			name: "valid type",
			f:    GeneralFeedbackType,
			want: true,
		},
		{
			name: "invalid type",
			f:    FeedbackType("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.IsValid(); got != tt.want {
				t.Errorf("FeedbackType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeedbackType_String(t *testing.T) {
	tests := []struct {
		name string
		f    FeedbackType
		want string
	}{
		{
			name: "GENERAL_FEEDBACK",
			f:    GeneralFeedbackType,
			want: "GENERAL_FEEDBACK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("FeedbackType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFeedbackType_UnmarshalGQL(t *testing.T) {
	validValue := GeneralFeedbackType
	invalidValue := FeedbackType("invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		f       *FeedbackType
		args    args
		wantErr bool
	}{
		{
			name: "valid type",
			f:    &validValue,
			args: args{
				v: "GENERAL_FEEDBACK",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			f:    &invalidValue,
			args: args{
				v: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("FeedbackType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFeedbackType_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		f     FeedbackType
		wantW string
	}{
		{
			name:  "GENERAL_FEEDBACK",
			f:     GeneralFeedbackType,
			wantW: `"GENERAL_FEEDBACK"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.f.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("FeedbackType.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
