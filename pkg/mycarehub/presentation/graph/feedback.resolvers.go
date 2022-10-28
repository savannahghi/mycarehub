package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// SendFeedback is the resolver for the sendFeedback field.
func (r *mutationResolver) SendFeedback(ctx context.Context, input dto.FeedbackResponseInput) (bool, error) {
	r.checkPreconditions()

	return r.mycarehub.Feedback.SendFeedback(ctx, &input)
}
