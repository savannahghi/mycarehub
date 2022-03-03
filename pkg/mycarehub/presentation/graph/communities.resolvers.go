package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	stream_chat "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

func (r *mutationResolver) CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
	return r.mycarehub.Community.CreateCommunity(ctx, input)
}

func (r *mutationResolver) DeleteCommunities(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error) {
	return r.mycarehub.Community.DeleteCommunities(ctx, communityIDs, hardDelete)
}

func (r *mutationResolver) RejectInvitation(ctx context.Context, memberID string, communityID string) (bool, error) {
	return r.mycarehub.Community.RejectInvite(ctx, memberID, communityID)
}

func (r *mutationResolver) AcceptInvitation(ctx context.Context, memberID string, communityID string) (bool, error) {
	return r.mycarehub.Community.AcceptInvite(ctx, memberID, communityID)
}

func (r *mutationResolver) AddMembersToCommunity(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
	return r.mycarehub.Community.AddMembersToCommunity(ctx, memberIDs, communityID)
}

func (r *mutationResolver) RemoveMembersFromCommunity(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	return r.mycarehub.Community.RemoveMembersFromCommunity(ctx, communityID, memberIDs)
}

func (r *mutationResolver) AddModerators(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
	return r.mycarehub.Community.AddModeratorsWithMessage(ctx, memberIDs, communityID)
}

func (r *mutationResolver) DemoteModerators(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	return r.mycarehub.Community.DemoteModerators(ctx, communityID, memberIDs)
}

func (r *queryResolver) ListMembers(ctx context.Context, input *stream_chat.QueryOption) ([]*domain.Member, error) {
	return r.mycarehub.Community.ListMembers(ctx, input)
}

func (r *queryResolver) InviteMembersToCommunity(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	return r.mycarehub.Community.InviteMembers(ctx, communityID, memberIDs)
}

func (r *queryResolver) ListCommunities(ctx context.Context, input *stream_chat.QueryOption) ([]*domain.Community, error) {
	return r.mycarehub.Community.ListCommunities(ctx, input)
}

func (r *queryResolver) ListCommunityMembers(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
	return r.mycarehub.Community.ListCommunityMembers(ctx, communityID)
}

func (r *queryResolver) ListPendingInvites(ctx context.Context, memberID string, input *stream_chat.QueryOption) ([]*domain.Community, error) {
	return r.mycarehub.Community.ListPendingInvites(ctx, memberID, input)
}

func (r *queryResolver) RecommendedCommunities(ctx context.Context, clientID string, limit int) ([]*domain.Community, error) {
	return r.mycarehub.Community.RecommendedCommunities(ctx, clientID, limit)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
