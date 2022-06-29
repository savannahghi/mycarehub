package mock

import (
	"context"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
)

// NewStreamClientMock initializes the mock client
func NewStreamClientMock() *StreamClientMock {
	return &StreamClientMock{}
}

// StreamClientMock mocks the GetStream client service library implementations
type StreamClientMock struct {
	MockCreateTokenFn       func(userID string, expire time.Time, issuedAt ...time.Time) (string, error)
	MockRevokeUserTokenFn   func(ctx context.Context, userID string, before *time.Time) (*stream.Response, error)
	MockUpsertUserFn        func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	MockQueryUsersFn        func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error)
	MockCreateChannelFn     func(ctx context.Context, chanType, chanID, userID string, data *stream.ChannelRequest) (*stream.CreateChannelResponse, error)
	MockDeleteChannelsFn    func(ctx context.Context, cids []string, hardDelete bool) (*stream.AsyncTaskResponse, error)
	MockQueryChannelsFn     func(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error)
	MockChannelFn           func(channelType, channelID string) *stream.Channel
	MockDeleteUsersFn       func(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error)
	MockQueryMessageFlagsFn func(ctx context.Context, q *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error)
	MockDeleteMessageFn     func(ctx context.Context, msgID string) (*stream.Response, error)
	MockVerifyWebhookFn     func(body, signature []byte) (valid bool)
	MockQueryBannedUsersFn  func(ctx context.Context, q *stream.QueryBannedUsersOptions, sorters ...*stream.SortOption) (*stream.QueryBannedUsersResponse, error)
}

// CreateToken ...
func (c StreamClientMock) CreateToken(userID string, expire time.Time, issuedAt ...time.Time) (string, error) {
	return c.MockCreateTokenFn(userID, expire, issuedAt...)
}

// RevokeUserToken ...
func (c StreamClientMock) RevokeUserToken(ctx context.Context, userID string, before *time.Time) (*stream.Response, error) {
	return c.MockRevokeUserTokenFn(ctx, userID, before)
}

// UpsertUser ...
func (c StreamClientMock) UpsertUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return c.MockUpsertUserFn(ctx, user)
}

// QueryUsers ...
func (c StreamClientMock) QueryUsers(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
	return c.MockQueryUsersFn(ctx, q, sorters...)
}

// CreateChannel ...
func (c StreamClientMock) CreateChannel(ctx context.Context, chanType, chanID, userID string, data *stream.ChannelRequest) (*stream.CreateChannelResponse, error) {
	return c.MockCreateChannelFn(ctx, chanType, chanID, userID, data)
}

// DeleteChannels ...
func (c StreamClientMock) DeleteChannels(ctx context.Context, cids []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
	return c.MockDeleteChannelsFn(ctx, cids, hardDelete)
}

// QueryChannels ...
func (c StreamClientMock) QueryChannels(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error) {
	return c.MockQueryChannelsFn(ctx, q, sort...)
}

// Channel ...
func (c *StreamClientMock) Channel(channelType, channelID string) *stream.Channel {
	return c.MockChannelFn(channelType, channelID)
}

// DeleteUsers ...
func (c StreamClientMock) DeleteUsers(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error) {
	return c.MockDeleteUsersFn(ctx, userIDs, options)
}

// QueryMessageFlags ...
func (c StreamClientMock) QueryMessageFlags(ctx context.Context, q *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
	return c.MockQueryMessageFlagsFn(ctx, q)
}

// DeleteMessage ...
func (c StreamClientMock) DeleteMessage(ctx context.Context, msgID string) (*stream.Response, error) {
	return c.MockDeleteMessageFn(ctx, msgID)
}

// VerifyWebhook ...
func (c StreamClientMock) VerifyWebhook(body, signature []byte) (valid bool) {
	return c.MockVerifyWebhookFn(body, signature)
}

// QueryBannedUsers ...
func (c StreamClientMock) QueryBannedUsers(ctx context.Context, q *stream.QueryBannedUsersOptions, sorters ...*stream.SortOption) (*stream.QueryBannedUsersResponse, error) {
	return c.MockQueryBannedUsersFn(ctx, q, sorters...)
}
