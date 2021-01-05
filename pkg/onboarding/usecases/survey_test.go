package usecases_test

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestSurveyUseCasesImpl_RecordPostVisitSurvey(t *testing.T) {
	authenticatedContext := base.GetAuthenticatedContext(t)
	s, err := InitializeTestService(authenticatedContext)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx   context.Context
		input resources.PostVisitSurveyInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: authenticatedContext,
				input: resources.PostVisitSurveyInput{
					LikelyToRecommend: 10,
					Criticism:         "very good developers",
					Suggestions:       "pay them more",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "bad case - invalid input",
			args: args{
				ctx: authenticatedContext,
				input: resources.PostVisitSurveyInput{
					LikelyToRecommend: 11,
					Criticism:         "piece of crap",
					Suggestions:       "replace it all",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "bad case - user not found",
			args: args{
				ctx: context.Background(),
				input: resources.PostVisitSurveyInput{
					LikelyToRecommend: 0,
					Criticism:         "piece of crap",
					Suggestions:       "replace it all",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := s
			got, err := rs.Survey.RecordPostVisitSurvey(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"SurveyUseCasesImpl.RecordPostVisitSurvey() error = %v,wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if got != tt.want {
				t.Errorf("SurveyUseCasesImpl.RecordPostVisitSurvey() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
