package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

func (r *mutationResolver) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	return r.interactor.UserUsecase.InviteUser(ctx, userID, phoneNumber, flavour)
}

func (r *mutationResolver) SetUserPin(ctx context.Context, input *dto.PINInput) (bool, error) {
	return r.interactor.UserUsecase.SetUserPIN(ctx, *input)
}
