package feedback_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	mailMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/mail/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback/mock"
)

func TestUsecaseFeedbackImpl_SendFeedback(t *testing.T) {
	ctx := context.Background()

	feedbackInput := &dto.FeedbackResponseInput{
		UserID:            uuid.New().String(),
		FeedbackType:      enums.GeneralFeedbackType,
		SatisfactionLevel: 4,
		ServiceName:       "JOIN",
		Feedback:          "test",
		RequiresFollowUp:  true,
	}
	noMessageFeedback := &dto.FeedbackResponseInput{
		UserID:           "user-id",
		Feedback:         "",
		RequiresFollowUp: true,
	}
	invalidFeedbackType := &dto.FeedbackResponseInput{
		UserID:            uuid.New().String(),
		FeedbackType:      enums.FeedbackType("invalid"),
		SatisfactionLevel: 4,
		ServiceName:       "JOIN",
		Feedback:          "test",
		RequiresFollowUp:  true,
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
			name: "Happy case: send feedback",
			args: args{
				ctx:     ctx,
				payload: feedbackInput,
			},
			want:    true,
			wantErr: false,
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
		{
			name: "Sad case - invalid feedback type",
			args: args{
				ctx:     ctx,
				payload: invalidFeedbackType,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - unable to persist feedback",
			args: args{
				ctx:     ctx,
				payload: invalidFeedbackType,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - unable go get client profile",
			args: args{
				ctx:     ctx,
				payload: invalidFeedbackType,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeFeedback := mock.NewFeedbackUsecaseMock()
			fakeMailService := mailMock.NewMailServiceMock()

			f := feedback.NewUsecaseFeedback(fakeDB, fakeDB, fakeMailService)

			if tt.name == "Sad case - no message" {
				fakeFeedback.MockSendFeedbackFn = func(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error) {
					return false, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			if tt.name == "Sad case - unable to send message" {
				fakeMailService.MockSendFeedbackFn = func(ctx context.Context, subject, feedbackMessage string) (bool, error) {
					return false, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			if tt.name == "Sad case - invalid feedback type" {
				fakeMailService.MockSendFeedbackFn = func(ctx context.Context, subject, feedbackMessage string) (bool, error) {
					return false, fmt.Errorf("an error occurred while sending feedback")
				}
			}
			if tt.name == "Sad case - unable to persist feedback" {
				fakeDB.MockSaveFeedbackFn = func(ctx context.Context, feedback *domain.FeedbackResponse) error {
					return fmt.Errorf("an error occurred while saving feedback")
				}
			}
			if tt.name == "Sad case - unable go get client profile" {
				fakeDB.MockGetClientProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred while getting client profile")
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
