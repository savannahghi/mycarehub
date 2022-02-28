package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *mutationResolver) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return r.mycarehub.Authority.AssignRoles(ctx, userID, roles)
}

func (r *mutationResolver) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return r.mycarehub.Authority.RevokeRoles(ctx, userID, roles)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
