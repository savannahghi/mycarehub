package usecases

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
)

func TestSurveyUseCasesImpl_RecordPostVisitSurvey(t *testing.T) {
	authenticatedContext := base.GetAuthenticatedContext(t)
	r, err := database.NewFirebaseRepository(authenticatedContext)
	if err != nil {
		t.Errorf("can't initialize a new firebase repository")
		return
	}
	rs := NewSurveyUseCases(r)

	type args struct {
		ctx   context.Context
		input domain.PostVisitSurveyInput
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
				input: domain.PostVisitSurveyInput{
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
				input: domain.PostVisitSurveyInput{
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
				input: domain.PostVisitSurveyInput{
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
			rs := rs
			got, err := rs.RecordPostVisitSurvey(tt.args.ctx, tt.args.input)
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
