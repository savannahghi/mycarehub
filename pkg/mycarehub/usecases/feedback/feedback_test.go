package feedback_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback/mock"
)

func TestUsecaseFeedbackImpl_SendFeedback(t *testing.T) {
	ctx := context.Background()

	feedbackInput := &dto.FeedbackResponseInput{
		UserID:           "test",
		Message:          "test",
		RequiresFollowUp: true,
	}
	noUserIDFeedback := &dto.FeedbackResponseInput{
		UserID:           "",
		Message:          "test",
		RequiresFollowUp: true,
	}
	noMessageFeedback := &dto.FeedbackResponseInput{
		UserID:           "user-id",
		Message:          "",
		RequiresFollowUp: true,
	}

	type args struct {
		ctx     context.Context
		payload *dto.FeedbackResponseInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				payload: feedbackInput,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				payload: feedbackInput,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no user ID",
			args: args{
				ctx:     ctx,
				payload: noUserIDFeedback,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no message",
			args: args{
				ctx:     ctx,
				payload: noMessageFeedback,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - unable to send message",
			args: args{
				ctx:     ctx,
				payload: noMessageFeedback,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeFeedback := mock.NewFeedbackUsecaseMock()
			fakeExtension := extensionMock.NewFakeExtension()

			f := feedback.NewUsecaseFeedback(fakeDB, fakeExtension)

			if tt.name == "Sad case" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			if tt.name == "Sad case - no user ID" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			if tt.name == "Sad case - no message" {
				fakeFeedback.MockSendFeedbackFn = func(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error) {
					return false, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			if tt.name == "Sad case - unable to send message" {
				fakeExtension.MockSendFeedbackFn = func(ctx context.Context, subject, feedbackMessage string) (bool, error) {
					return false, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			got, err := f.SendFeedback(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseFeedbackImpl.SendFeedback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseFeedbackImpl.SendFeedback() = %v, want %v", got, tt.want)
			}
		})
	}
}
