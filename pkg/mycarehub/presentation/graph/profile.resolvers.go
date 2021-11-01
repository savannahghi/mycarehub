package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

func (r *mutationResolver) SetUserPin(ctx context.Context, input *dto.PinInput) (bool, error) {
	r.checkPreconditions()
	return r.interactor.UserUseCase.SetUserPIN(ctx, input)
}
