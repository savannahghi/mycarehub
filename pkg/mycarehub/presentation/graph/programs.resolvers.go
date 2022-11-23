package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// CreateProgram is the resolver for the createProgram field.
func (r *mutationResolver) CreateProgram(ctx context.Context, input dto.ProgramInput) (bool, error) {
	return r.mycarehub.Programs.CreateProgram(ctx, &input)
}
