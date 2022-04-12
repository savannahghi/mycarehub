package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) AssignOrRevokeRoles(ctx context.Context, userID string, roles []*enums.UserRoleType) (bool, error) {
	return r.mycarehub.Authority.AssignOrRevokeRoles(ctx, userID, roles)
}

func (r *queryResolver) GetUserRoles(ctx context.Context, userID string) ([]*domain.AuthorityRole, error) {
	return r.mycarehub.Authority.GetUserRoles(ctx, userID)
}

func (r *queryResolver) GetAllAuthorityRoles(ctx context.Context) ([]*domain.AuthorityRole, error) {
	return r.mycarehub.Authority.GetAllRoles(ctx)
}
