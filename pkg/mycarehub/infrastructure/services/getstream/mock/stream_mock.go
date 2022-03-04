package mock

import (
	"context"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

// GetStreamServiceMock mocks the GetStream service library implementations
type GetStreamServiceMock struct {
	MockCreateGetStreamUserTokenFn func(ctx context.Context, userID string) (string, error)
	MockCreateGetStreamUserFn      func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	MockListGetStreamUsersFn       func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error)
	MockCreateChannelFn            func(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error)
	MockDeleteChannelsFn           func(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error)
	MockInviteMembersFn            func(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error)
	MockListGetStreamChannelsFn    func(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error)
	MockListChannelMembersFn       func(ctx context.Context, channelID string, q *stream.QueryOption, sorters ...*stream.SortOption) ([]*stream.ChannelMember, error)
	MockGetChannel                 func(ctx context.Context, channelID string) (*stream.Channel, error)
	MockAddMembersToCommunityFn    func(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error)
	MockRejectInviteFn             func(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	MockAcceptInviteFn             func(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	MockRemoveMembersFn            func(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error)
	MockAddModeratorsWithMessageFn func(ctx context.Context, userIDs []string, communityID string, message *stream.Message) (*stream.Response, error)
	MockDemoteModeratorsFn         func(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error)
	MockRevokeGetStreamUserTokenFn func(ctx context.Context, userID string) (*stream.Response, error)
}

// NewGetStreamServiceMock initializes the mock service
func NewGetStreamServiceMock() *GetStreamServiceMock {
	var now = time.Now()
	return &GetStreamServiceMock{
		MockCreateGetStreamUserTokenFn: func(ctx context.Context, userID string) (string, error) {
			return uuid.New().String(), nil
		},
		MockCreateGetStreamUserFn: func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
			return &stream.UpsertUserResponse{
				User: &stream.User{
					ID:   uuid.New().String(),
					Name: gofakeit.Name(),
				},
			}, nil
		},
		MockListGetStreamUsersFn: func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
			return &stream.QueryUsersResponse{
				Users: []*stream.User{
					{
						ID:   uuid.NewString(),
						Name: gofakeit.Name(),
					},
				},
			}, nil
		},
		MockCreateChannelFn: func(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error) {
			return &stream.CreateChannelResponse{
				Channel: &stream.Channel{
					CreatedBy: &stream.User{
						ID: uuid.New().String(),
					},
				},
				Response: &stream.Response{},
			}, nil
		},
		MockInviteMembersFn: func(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
			return &stream.Response{RateLimitInfo: &stream.RateLimitInfo{
				Limit: 100,
			}}, nil
		},
		MockAddMembersToCommunityFn: func(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error) {
			return &stream.Response{RateLimitInfo: &stream.RateLimitInfo{
				Limit: 100,
			}}, nil
		},
		MockListGetStreamChannelsFn: func(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error) {

			createdBy := &stream.User{
				ID:        uuid.NewString(),
				Name:      gofakeit.Name(),
				Role:      gofakeit.Name(),
				CreatedAt: &now,
				UpdatedAt: &now,
			}
			channelMembers := []*stream.ChannelMember{
				{
					UserID: uuid.NewString(),
					User: &stream.User{
						ID:                       uuid.NewString(),
						Name:                     gofakeit.Name(),
						Image:                    gofakeit.URL(),
						Role:                     gofakeit.Name(),
						Teams:                    []string{gofakeit.Name()},
						Online:                   false,
						Invisible:                false,
						CreatedAt:                &now,
						UpdatedAt:                &now,
						LastActive:               &now,
						Mutes:                    nil,
						ChannelMutes:             nil,
						ExtraData:                map[string]interface{}{},
						RevokeTokensIssuedBefore: nil,
					},
					IsModerator:      false,
					Invited:          false,
					InviteAcceptedAt: &now,
					InviteRejectedAt: nil,
					Role:             gofakeit.Name(),
					CreatedAt:        now,
					UpdatedAt:        now,
				},
			}

			return &stream.QueryChannelsResponse{
				Channels: []*stream.Channel{
					{
						ID:            uuid.NewString(),
						Type:          gofakeit.Name(),
						CID:           uuid.NewString(),
						Team:          uuid.NewString(),
						CreatedBy:     createdBy,
						Disabled:      false,
						Frozen:        false,
						MemberCount:   1,
						Members:       channelMembers,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						LastMessageAt: time.Now(),
					},
				},
			}, nil
		},
		MockDeleteChannelsFn: func(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
			return &stream.AsyncTaskResponse{
				Response: stream.Response{
					RateLimitInfo: &stream.RateLimitInfo{
						Limit:     100,
						Remaining: 100,
						Reset:     10,
					},
				},
			}, nil
		},
		MockRejectInviteFn: func(ctx context.Context, userID, channelID string, message *stream.Message) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil
		},
		MockAcceptInviteFn: func(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil

		},
		MockRemoveMembersFn: func(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil
		},
		MockAddModeratorsWithMessageFn: func(ctx context.Context, userIDs []string, communityID string, message *stream.Message) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil
		},
		MockDemoteModeratorsFn: func(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil
		},
		MockRevokeGetStreamUserTokenFn: func(ctx context.Context, userID string) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil
		},
	}
}

// CreateGetStreamUserToken mocks creating a getstream user token
func (g GetStreamServiceMock) CreateGetStreamUserToken(ctx context.Context, userID string) (string, error) {
	return g.MockCreateGetStreamUserTokenFn(ctx, userID)
}

// CreateGetStreamUser mocks creating a getstream user
func (g GetStreamServiceMock) CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return g.MockCreateGetStreamUserFn(ctx, user)
}

// ListGetStreamUsers mocks the implementation for listing getstream users
func (g GetStreamServiceMock) ListGetStreamUsers(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
	return g.MockListGetStreamUsersFn(ctx, queryOptions)
}

// CreateChannel mocks the implementation of creating a channel
func (g GetStreamServiceMock) CreateChannel(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error) {
	return g.MockCreateChannelFn(ctx, chanType, chanID, userID, data)
}

// DeleteChannels mocks the implementation of deleting channel asynchronously
func (g GetStreamServiceMock) DeleteChannels(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
	return g.MockDeleteChannelsFn(ctx, chanIDs, hardDelete)
}

// InviteMembers mocks the implementation for inviting members to a specified channel
func (g GetStreamServiceMock) InviteMembers(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
	return g.MockInviteMembersFn(ctx, memberIDs, channelID, message)
}

// ListGetStreamChannels mocks the implementation for listing getstream channels
func (g GetStreamServiceMock) ListGetStreamChannels(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
	return g.MockListGetStreamChannelsFn(ctx, input)
}

// ListChannelMembers mocks implementation for listing channel members
func (g GetStreamServiceMock) ListChannelMembers(ctx context.Context, channelID string, query *stream.QueryOption, sorters ...*stream.SortOption) ([]*stream.ChannelMember, error) {
	return g.MockListChannelMembersFn(ctx, channelID, query, sorters...)
}

// GetChannel mocks implementation for retrieving a channel
func (g GetStreamServiceMock) GetChannel(ctx context.Context, channelID string) (*stream.Channel, error) {
	return g.MockGetChannel(ctx, channelID)
}

// AddMembersToCommunity mocks the implementation for adding members(s) to a community
func (g GetStreamServiceMock) AddMembersToCommunity(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error) {
	return g.MockAddMembersToCommunityFn(ctx, memberIDs, channelID)
}

// RejectInvite mocks the implementation for rejecting invite into a community
func (g GetStreamServiceMock) RejectInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {
	return g.MockRejectInviteFn(ctx, userID, channelID, message)
}

// AcceptInvite mocks the implementation for accepting invite into a community
func (g GetStreamServiceMock) AcceptInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {
	return g.MockAcceptInviteFn(ctx, userID, channelID, message)
}

// RemoveMembersFromCommunity mocks the implementation for removing members from a community
func (g GetStreamServiceMock) RemoveMembersFromCommunity(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error) {
	return g.MockRemoveMembersFn(ctx, channelID, memberIDs, message)
}

// AddModeratorsWithMessage mocks the implementation of adding a moderator to a community with a message
func (g GetStreamServiceMock) AddModeratorsWithMessage(ctx context.Context, userIDs []string, communityID string, message *stream.Message) (*stream.Response, error) {
	return g.MockAddModeratorsWithMessageFn(ctx, userIDs, communityID, message)
}

// DemoteModerators mocks the implementation for demoting moderators from a community
func (g GetStreamServiceMock) DemoteModerators(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error) {
	return g.MockDemoteModeratorsFn(ctx, channelID, memberIDs)
}

// RevokeGetStreamUserToken mocks the implementation for revoking a getstream user token
func (g GetStreamServiceMock) RevokeGetStreamUserToken(ctx context.Context, userID string) (*stream.Response, error) {
	return g.MockRevokeGetStreamUserTokenFn(ctx, userID)
}
