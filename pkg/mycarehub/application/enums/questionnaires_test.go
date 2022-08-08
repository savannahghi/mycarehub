package enums

import (
	"testing"
)

func TestQuestionnaireQuestionTypeChoices_IsValid(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionType
		want bool
	}{
		{
			name: "QuestionTypeOpenEnded Case - Valid type",
			q:    QuestionTypeOpenEnded,
			want: true,
		},
		{
			name: "QuestionTypeOpenEnded Case - invalid type",
			q:    QuestionType("ANY QUESTION"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsValid(); got != tt.want {
				t.Errorf("QuestionType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireQuestionTypeChoices_String(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionType
		want string
	}{
		{
			name: "QuestionTypeOpenEnded Case",
			q:    QuestionTypeOpenEnded,
			want: QuestionTypeOpenEnded.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("QuestionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireQuestionTypeChoices_UnmarshalGQL(t *testing.T) {
	validValue := QuestionTypeOpenEnded
	invalidType := QuestionType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		q       *QuestionType
		args    args
		wantErr bool
	}{
		{
			name: "QuestionTypeOpenEnded Case - Valid type",
			args: args{
				v: QuestionTypeOpenEnded.String(),
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
				t.Errorf("QuestionType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestionnaireResponseValueChoices_IsValid(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionResponseValueType
		want bool
	}{
		{
			name: "Number Choice Case - Valid type",
			q:    QuestionResponseValueTypeNumber,
			want: true,
		},
		{
			name: "Number Choice Case - invalid type",
			q:    QuestionResponseValueType("ANY QUESTION"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsValid(); got != tt.want {
				t.Errorf("QuestionResponseValueType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireResponseValueChoices_String(t *testing.T) {
	tests := []struct {
		name string
		q    QuestionResponseValueType
		want string
	}{
		{
			name: "Number Choice Case",
			q:    QuestionResponseValueTypeNumber,
			want: QuestionResponseValueTypeNumber.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("QuestionResponseValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestionnaireResponseValueChoices_UnmarshalGQL(t *testing.T) {
	validChoice := QuestionResponseValueTypeNumber
	invalidType := QuestionResponseValueType("INVALID")
	type args struct {
		v interface{}
	}
	tests := []struct {
		name    string
		q       *QuestionResponseValueType
		args    args
		wantErr bool
	}{
		{
			name: "Number Choice Case - Valid type",
			args: args{
				v: QuestionResponseValueTypeNumber.String(),
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
				t.Errorf("QuestionResponseValueType.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
