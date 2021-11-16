package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) RecordSecurityQuestionResponses(ctx context.Context, input []*dto.SecurityQuestionResponseInput) ([]*domain.RecordSecurityQuestionResponse, error) {
	r.checkPreconditions()
	return r.mycarehub.SecurityQuestions.RecordSecurityQuestionResponses(ctx, input)
}

func (r *queryResolver) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	r.checkPreconditions()
	return r.mycarehub.SecurityQuestions.GetSecurityQuestions(ctx, flavour)
}
