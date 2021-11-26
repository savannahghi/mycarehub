package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *mutationResolver) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	return r.mycarehub.Content.ShareContent(ctx, input)
}

func (r *mutationResolver) BookmarkContent(ctx context.Context, userID string, contentItemID int) (bool, error) {
	return r.mycarehub.Content.BookmarkContent(ctx, userID, contentItemID)
}

func (r *mutationResolver) UnBookmarkContent(ctx context.Context, userID string, contentItemID int) (bool, error) {
	return r.mycarehub.Content.UnBookmarkContent(ctx, userID, contentItemID)
}

func (r *mutationResolver) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	r.checkPreconditions()

	return r.mycarehub.Content.LikeContent(ctx, userID, contentID)
}

func (r *mutationResolver) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	r.checkPreconditions()

	return r.mycarehub.Content.UnlikeContent(ctx, userID, contentID)
}

func (r *queryResolver) GetContent(ctx context.Context, categoryID *int, limit string) (*domain.Content, error) {
	r.checkPreconditions()
	return r.mycarehub.Content.GetContent(ctx, categoryID, limit)
}

func (r *queryResolver) ListContentCategories(ctx context.Context) ([]*domain.ContentItemCategory, error) {
	r.checkPreconditions()
	return r.mycarehub.Content.ListContentCategories(ctx)
}

func (r *queryResolver) GetUserBookmarkedContent(ctx context.Context, userID string) (*domain.Content, error) {
	r.checkPreconditions()
	return r.mycarehub.Content.GetUserBookmarkedContent(ctx, userID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
