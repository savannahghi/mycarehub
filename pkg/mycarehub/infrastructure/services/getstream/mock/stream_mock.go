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
	MockInviteMembersFn            func(ctx context.Context, userIDs []string, channelID string, message *stream.Message) (*stream.Response, error)
	MockListGetStreamChannelsFn    func(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error)
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
				Channel:  &stream.Channel{},
				Response: &stream.Response{},
			}, nil
		},
		MockInviteMembersFn: func(ctx context.Context, userIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
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
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						LastMessageAt: time.Now(),
					},
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
func (g GetStreamServiceMock) InviteMembers(ctx context.Context, userIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
	return g.MockInviteMembersFn(ctx, userIDs, channelID, message)
}

// ListGetStreamChannels mocks the implementation for listing getstream channels
func (g GetStreamServiceMock) ListGetStreamChannels(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
	return g.MockListGetStreamChannelsFn(ctx, input)
}
