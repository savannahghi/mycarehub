package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *queryResolver) ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error) {
	return r.mycarehub.Community.ListGetStreamUsers(ctx, input)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
