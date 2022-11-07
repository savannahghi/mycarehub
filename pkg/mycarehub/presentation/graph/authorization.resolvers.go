package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// AddPolicy is the resolver for the addPolicy field.
func (r *mutationResolver) AddPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.AddPolicy(ctx, subject, permission)
}

// RemovePolicy is the resolver for the removePolicy field.
func (r *mutationResolver) RemovePolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.RemovePolicy(ctx, subject, permission)
}

// AddGroupingPolicy is the resolver for the addGroupingPolicy field.
func (r *mutationResolver) AddGroupingPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.AddGroupingPolicy(ctx, subject, permission)
}

// RemoveGroupingPolicy is the resolver for the removeGroupingPolicy field.
func (r *mutationResolver) RemoveGroupingPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.RemoveGroupingPolicy(ctx, subject, permission)
}

// CheckPermissions is the resolver for the checkPermissions field.
func (r *queryResolver) CheckPermissions(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.CheckPermissions(ctx, subject, permission)
}

// CheckAuthorization is the resolver for the checkAuthorization field.
func (r *queryResolver) CheckAuthorization(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.CheckAuthorization(ctx, subject, permission)
}

// IsAuthorized is the resolver for the isAuthorized field.
func (r *queryResolver) IsAuthorized(ctx context.Context, permission dto.PermissionInput) (bool, error) {
	return r.mycarehub.Authorization.IsAuthorized(ctx, permission)
}
