package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *mutationResolver) CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
	return r.mycarehub.Community.CreateCommunity(ctx, input)
}

func (r *queryResolver) ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error) {
	return r.mycarehub.Community.ListGetStreamUsers(ctx, input)
}

func (r *queryResolver) InviteMembersToCommunity(ctx context.Context, communityID string, userIDS []string) (bool, error) {
	return r.mycarehub.Community.InviteMembers(ctx, communityID, userIDS)
}

func (r *queryResolver) ListGetStreamChannels(ctx context.Context, input *domain.QueryOption) (*domain.QueryChannelsResponse, error) {
	return r.mycarehub.Community.ListGetStreamChannels(ctx, input)
}

func (r *queryResolver) ListCommunityMembers(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
	return r.mycarehub.Community.ListCommunityMembers(ctx, communityID)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
