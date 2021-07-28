package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

func TestSurveyUseCasesImpl_RecordPostVisitSurvey(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v",
			err,
		)
		return
	}

	type args struct {
		ctx   context.Context
		input dto.PostVisitSurveyInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_record_post_vist_survey",
			args: args{
				ctx: ctx,
				input: dto.PostVisitSurveyInput{
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
				ctx: ctx,
				input: dto.PostVisitSurveyInput{
					LikelyToRecommend: 11,
					Criticism:         "piece of crap",
					Suggestions:       "replace it all",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_record_survey",
			args: args{
				ctx: ctx,
				input: dto.PostVisitSurveyInput{
					LikelyToRecommend: 5,
					Criticism:         "very good developers",
					Suggestions:       "pay them more",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_user_not_found",
			args: args{
				ctx: context.Background(),
				input: dto.PostVisitSurveyInput{
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

			if tt.name == "invalid:_user_not_found" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged in user")
				}
			}

			if tt.name == "valid:_record_post_vist_survey" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "8716-7e2aead29f2c", nil
				}
				fakeRepo.RecordPostVisitSurveyFn = func(ctx context.Context, input dto.PostVisitSurveyInput, UID string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_record_survey" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "8716-7e2aead29f2c", nil
				}
				fakeRepo.RecordPostVisitSurveyFn = func(ctx context.Context, input dto.PostVisitSurveyInput, UID string) error {
					return fmt.Errorf("unable to record post visit survey")
				}
			}

			got, err := i.Survey.RecordPostVisitSurvey(tt.args.ctx, tt.args.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got != tt.want {
					t.Errorf("expected %v got %v  ", tt.want, got)
					return
				}
			}

		})
	}
}
