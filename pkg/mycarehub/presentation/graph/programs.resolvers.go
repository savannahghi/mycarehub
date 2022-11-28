package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateProgram is the resolver for the createProgram field.
func (r *mutationResolver) CreateProgram(ctx context.Context, input dto.ProgramInput) (bool, error) {
	return r.mycarehub.Programs.CreateProgram(ctx, &input)
}

// SetCurrentProgram is the resolver for the setCurrentProgram field.
func (r *mutationResolver) SetCurrentProgram(ctx context.Context, id string) (bool, error) {
	return r.mycarehub.Programs.SetCurrentProgram(ctx, id)
}

// ListUserPrograms is the resolver for the listUserPrograms field.
func (r *queryResolver) ListUserPrograms(ctx context.Context, userID string) ([]*domain.Program, error) {
	return r.mycarehub.Programs.ListUserPrograms(ctx, userID)
}
