package enums

import (
	"testing"
)

func TestQuestionnaireQuestionTypeChoices_IsValid(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionnaireQuestionTypeChoices
		want bool
	}{
		{
			name: "OpenEnded Case - Valid type",
			q:    OpenEnded,
			want: true,
		},
		{
			name: "OpenEnded Case - invalid type",
			q:    QuestionnaireQuestionTypeChoices("ANY QUESTION"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsValid(); got != tt.want {
				t.Errorf("QuestionnaireQuestionTypeChoices.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireQuestionTypeChoices_String(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionnaireQuestionTypeChoices
		want string
	}{
		{
			name: "OpenEnded Case",
			q:    OpenEnded,
			want: OpenEnded.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("QuestionnaireQuestionTypeChoices.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireQuestionTypeChoices_UnmarshalGQL(t *testing.T) {
	validValue := OpenEnded
	invalidType := QuestionnaireQuestionTypeChoices("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		q       *QuestionnaireQuestionTypeChoices
		args    args
		wantErr bool
	}{
		{
			name: "OpenEnded Case - Valid type",
			args: args{
				v: OpenEnded.String(),
			},
			q:       &validValue,
			wantErr: false,
		},
		{
			name: "Sad Case - invalid type",
			args: args{
				v: "invalid type",
			},
			q:       &invalidType,
			wantErr: true,
		},
		{
			name: "Sad Case - invalid type",
			args: args{
				v: 45,
			},
			q:       &validValue,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("QuestionnaireQuestionTypeChoices.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestionnaireResponseValueChoices_IsValid(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionnaireResponseValueChoices
		want bool
	}{
		{
			name: "Number Choice Case - Valid type",
			q:    NumberResponseValue,
			want: true,
		},
		{
			name: "Number Choice Case - invalid type",
			q:    QuestionnaireResponseValueChoices("ANY QUESTION"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsValid(); got != tt.want {
				t.Errorf("QuestionnaireResponseValueChoices.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireResponseValueChoices_String(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionnaireResponseValueChoices
		want string
	}{
		{
			name: "Number Choice Case",
			q:    NumberResponseValue,
			want: NumberResponseValue.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("QuestionnaireResponseValueChoices.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireResponseValueChoices_UnmarshalGQL(t *testing.T) {
	validChoice := NumberResponseValue
	invalidType := QuestionnaireResponseValueChoices("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		q       *QuestionnaireResponseValueChoices
		args    args
		wantErr bool
	}{
		{
			name: "Number Choice Case - Valid type",
			args: args{
				v: NumberResponseValue.String(),
			},
			q:       &validChoice,
			wantErr: false,
		},
		{
			name: "Number Choice Case - invalid type",
			args: args{
				v: invalidType.String(),
			},
			q:       &validChoice,
			wantErr: true,
		},
		{
			name: "Number Choice Case - invalid type",
			args: args{
				v: 45,
			},
			q:       &validChoice,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("QuestionnaireResponseValueChoices.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
