package mock

import (
	"context"
	"net/http"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

// Client mocks the implementation of a getstream client
type Client struct {
	BaseURL string
	HTTP    *http.Client `json:"-"`

	apiKey       string
	apiSecret    []byte
	authToken    string
	ClientID     string
	ClientSecret string
	Client       stream.Client
}

// NewClient mocks the implementation of creating a new getstream client
func (g Client) NewClient() *Client {
	return &Client{
		BaseURL:   "",
		HTTP:      &http.Client{},
		apiKey:    "",
		apiSecret: []byte("ew"),
		authToken: "",
	}
}

// GetStreamServiceMock mocks the GetStream service library implementations
type GetStreamServiceMock struct {
	MockCreateGetStreamUserTokenFn   func(ctx context.Context, userID string) (string, error)
	MockCreateGetStreamUserFn        func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	MockListGetStreamUsersFn         func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error)
	MockCreateChannelFn              func(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error)
	MockDeleteChannelsFn             func(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error)
	MockInviteMembersFn              func(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error)
	MockListGetStreamChannelsFn      func(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error)
	MockGetChannel                   func(ctx context.Context, channelID string) (*stream.Channel, error)
	MockAddMembersToCommunityFn      func(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error)
	MockRejectInviteFn               func(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	MockAcceptInviteFn               func(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	MockRemoveMembersFn              func(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error)
	MockAddModeratorsWithMessageFn   func(ctx context.Context, userIDs []string, communityID string, message *stream.Message) (*stream.Response, error)
	MockDemoteModeratorsFn           func(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error)
	MockRevokeGetStreamUserTokenFn   func(ctx context.Context, userID string) (*stream.Response, error)
	MockDeleteUsersFn                func(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error)
	MockBanUserFn                    func(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error)
	MockUnBanUserFn                  func(ctx context.Context, targetID string, communityID string) (bool, error)
	MockListCommunityBannedMembersFn func(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error)
	MockUpsertUserFn                 func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	MockListFlaggedMessagesFn        func(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error)
	MockDeleteMessageFn              func(ctx context.Context, messageID string) (*stream.Response, error)
	MockValidateGetStreamRequestFn   func(ctx context.Context, body []byte, signature string) bool
	MockGetStreamUserFn              func(ctx context.Context, id string) (*stream.User, error)
}

// NewGetStreamServiceMock initializes the mock service
func NewGetStreamServiceMock() *GetStreamServiceMock {
	var now = time.Now()
	return &GetStreamServiceMock{
		MockGetStreamUserFn: func(ctx context.Context, id string) (*stream.User, error) {
			return &stream.User{
				ID:   uuid.NewString(),
				Name: gofakeit.Name(),
				ExtraData: map[string]interface{}{
					"userID":   gofakeit.UUID(),
					"userType": "STAFF",
					"nickName": gofakeit.Name(),
				},
			}, nil
		},
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

		MockListCommunityBannedMembersFn: func(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error) {
			return &stream.QueryBannedUsersResponse{
				Bans: []*stream.Ban{
					{
						Channel:   &stream.Channel{},
						User:      &stream.User{},
						Expires:   &time.Time{},
						Reason:    "",
						Shadow:    false,
						BannedBy:  &stream.User{},
						CreatedAt: time.Now(),
					},
				},
				Response: stream.Response{
					RateLimitInfo: &stream.RateLimitInfo{},
				},
			}, nil
		},

		MockDeleteUsersFn: func(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error) {
			return &stream.AsyncTaskResponse{
				Response: stream.Response{
					RateLimitInfo: &stream.RateLimitInfo{
						Limit:     0,
						Remaining: 0,
						Reset:     0,
					},
				},
			}, nil
		},
		MockBanUserFn: func(ctx context.Context, targetMemberID, bannedBy, communityID string) (bool, error) {
			return true, nil
		},
		MockUnBanUserFn: func(ctx context.Context, targetID, communityID string) (bool, error) {
			return true, nil
		},

		MockGetChannel: func(ctx context.Context, channelID string) (*stream.Channel, error) {
			return &stream.Channel{
				ID:     channelID,
				Type:   "",
				CID:    channelID,
				Team:   "",
				Config: stream.ChannelConfig{},
			}, nil
		},
		MockUpsertUserFn: func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
			return &stream.UpsertUserResponse{
				Response: stream.Response{
					RateLimitInfo: &stream.RateLimitInfo{
						Limit:     0,
						Remaining: 0,
						Reset:     0,
					},
				},
			}, nil
		},
		MockListFlaggedMessagesFn: func(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
			return &stream.QueryMessageFlagsResponse{
				Response: stream.Response{
					RateLimitInfo: &stream.RateLimitInfo{
						Limit:     0,
						Remaining: 0,
						Reset:     0,
					},
				},
			}, nil
		},
		MockDeleteMessageFn: func(ctx context.Context, messageID string) (*stream.Response, error) {
			return &stream.Response{
				RateLimitInfo: &stream.RateLimitInfo{
					Limit:     0,
					Remaining: 0,
					Reset:     0,
				},
			}, nil
		},
		MockValidateGetStreamRequestFn: func(ctx context.Context, body []byte, signature string) bool {
			return true
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

// DeleteUsers mocks the implementation for deleting users
func (g GetStreamServiceMock) DeleteUsers(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error) {
	return g.MockDeleteUsersFn(ctx, userIDs, options)
}

// BanUser mocks the implementation banning a user from a specified channel
func (g GetStreamServiceMock) BanUser(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error) {
	return g.MockBanUserFn(ctx, targetMemberID, bannedBy, communityID)
}

// UnBanUser mocks the implementation of unbanning a user from a specified channel
func (g GetStreamServiceMock) UnBanUser(ctx context.Context, targetID string, communityID string) (bool, error) {
	return g.MockUnBanUserFn(ctx, targetID, communityID)
}

// ListCommunityBannedMembers mocks the implementation of listing the community members
func (g GetStreamServiceMock) ListCommunityBannedMembers(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error) {
	return g.MockListCommunityBannedMembersFn(ctx, communityID)
}

// UpsertUser mocks the implementation of upserting a user
func (g GetStreamServiceMock) UpsertUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return g.MockUpsertUserFn(ctx, user)
}

// ListFlaggedMessages mocks the implementation for querying message flags
func (g GetStreamServiceMock) ListFlaggedMessages(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
	return g.MockListFlaggedMessagesFn(ctx, input)
}

// DeleteMessage mocks the implementation for deleting messages
func (g GetStreamServiceMock) DeleteMessage(ctx context.Context, messageID string) (*stream.Response, error) {
	return g.MockDeleteMessageFn(ctx, messageID)
}

// ValidateGetStreamRequest mocks the implementation of verifying the webhook request
func (g GetStreamServiceMock) ValidateGetStreamRequest(ctx context.Context, body []byte, signature string) bool {
	return g.MockValidateGetStreamRequestFn(ctx, body, signature)
}

// GetStreamUser retrieves a getstream user given the ID
func (g GetStreamServiceMock) GetStreamUser(ctx context.Context, id string) (*stream.User, error) {
	return g.MockGetStreamUserFn(ctx, id)
}
