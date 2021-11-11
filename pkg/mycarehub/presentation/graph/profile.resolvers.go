package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) InviteUser(ctx context.Context, userID string) (bool, error) {
	return r.interactor.UserUsecase.Invite(ctx, userID)
}
