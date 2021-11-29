package enums

import (
	"testing"
)

func TestMood_IsValid(t *testing.T) {
	tests := []struct {
		name string
		m    Mood
		want bool
	}{
		{
			name: "Happy Case - Valid type",
			m:    MoodHappy,
			want: true,
		},
		{
			name: "Sad Case - Invalid type",
			m:    Mood("Not so happy"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsValid(); got != tt.want {
				t.Errorf("Mood.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMood_String(t *testing.T) {
	tests := []struct {
		name string
		m    Mood
		want string
	}{
		{
			name: "Happy Case",
			m:    MoodHappy,
			want: MoodHappy.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Mood.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMood_UnmarshalGQL(t *testing.T) {
	validValue := MoodNeutral
	invalidType := Mood("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		m       *Mood
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Valid type",
			args: args{
				v: MoodNeutral.String(),
			},
			m:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - Invalid type",
			args: args{
				v: "invalid type",
			},
			m:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - Invalid type(int)",
			args: args{
				v: 45,
			},
			m:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Mood.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
