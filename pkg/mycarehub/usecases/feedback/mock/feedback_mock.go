package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// FeedbackUsecaseMock contains the mock of feedback usecase methods
type FeedbackUsecaseMock struct {
	MockSendFeedbackFn func(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error)
}

// NewFeedbackUsecaseMock instantiates all the feedback usecase mock methods
func NewFeedbackUsecaseMock() *FeedbackUsecaseMock {

	return &FeedbackUsecaseMock{
		MockSendFeedbackFn: func(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error) {
			return true, nil
		},
	}
}

//SendFeedback mocks the implementation sending feedback
func (fm *FeedbackUsecaseMock) SendFeedback(ctx context.Context, payload *dto.FeedbackResponseInput) (bool, error) {
	return fm.MockSendFeedbackFn(ctx, payload)
}
