package enums

import (
	"bytes"
	"strconv"
	"testing"
)

func TestVisibility_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Visibility
		want bool
	}{
		{
			name: "Happy Case - Valid visibility",
			e:    PrivateVisibility,
			want: true,
		},
		{
			name: "Sad Case - Invalid visibility",
			e:    Visibility("Invalid Visibility"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Visibility.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibility_String(t *testing.T) {
	tests := []struct {
		name string
		e    Visibility
		want string
	}{
		{
			name: "Happy Case - Valid string",
			e:    PrivateVisibility,
			want: PrivateVisibility.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Visibility.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVisibility_UnmarshalGQL(t *testing.T) {
	validValue := PrivateVisibility
	invalidType := Visibility("Invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *Visibility
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid Type",
			args: args{
				v: PrivateVisibility.String(),
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
				t.Errorf("Visibility.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVisibility_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     Visibility
		wantW string
	}{
		{
			name:  "Valid type enums",
			e:     PublicVisibility,
			wantW: strconv.Quote("public"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Visibility.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPreset_IsValid(t *testing.T) {
	tests := []struct {
		name string
		e    Preset
		want bool
	}{
		{
			name: "Happy Case - Valid preset",
			e:    PresetPrivateChat,
			want: true,
		},
		{
			name: "Sad Case - Invalid preset",
			e:    Preset("Invalid preset"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.IsValid(); got != tt.want {
				t.Errorf("Preset.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreset_String(t *testing.T) {
	tests := []struct {
		name string
		e    Preset
		want string
	}{
		{
			name: "Happy Case - Valid string",
			e:    PresetPublicChat,
			want: PresetPublicChat.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Preset.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreset_UnmarshalGQL(t *testing.T) {
	validValue := PresetPublicChat
	invalidType := Preset("Invalid")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		e       *Preset
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid Type",
			args: args{
				v: PresetTrustedPrivateChat.String(),
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
				t.Errorf("Preset.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPreset_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		e     Preset
		wantW string
	}{
		{
			name:  "Valid type enums",
			e:     PresetTrustedPrivateChat,
			wantW: strconv.Quote("trusted_private_chat"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.e.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Preset.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
