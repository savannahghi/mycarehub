package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *queryResolver) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
	r.checkPreconditions()
	return r.mycarehub.Content.GetContent(ctx, categoryID, limit)
}

func (r *queryResolver) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	r.checkPreconditions()

	return r.mycarehub.Content.ListContentCategories(ctx)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
