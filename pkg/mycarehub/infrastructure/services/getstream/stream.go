package getstream

import (
	"context"
	"fmt"
	"log"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/serverutils"
)

var (
	getStreamAPIKey          = serverutils.MustGetEnvVar("GET_STREAM_KEY")
	getStreamAPISecret       = serverutils.MustGetEnvVar("GET_STREAM_SECRET")
	getStreamTokenExpiryTime = time.Now().UTC().Add(time.Hour * 12)
)

// ServiceGetStream represents the various Getstream usecases
type ServiceGetStream interface {
	CreateGetStreamUserToken(ctx context.Context, userID string) (string, error)
	CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	ListGetStreamUsers(ctx context.Context, input *stream.QueryOption) (*stream.QueryUsersResponse, error)
	CreateChannel(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error)
	DeleteChannels(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error)
	InviteMembers(ctx context.Context, userIDs []string, channelID string, message *stream.Message) (*stream.Response, error)
	ListGetStreamChannels(ctr context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error)
	ListChannelMembers(ctx context.Context, channelID string, q *stream.QueryOption, sorters ...*stream.SortOption) ([]*stream.ChannelMember, error)
	GetChannel(ctx context.Context, channelID string) (*stream.Channel, error)
	AddMembersToCommunity(ctx context.Context, userIDs []string, channelID string) (*stream.Response, error)
	RejectInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	AcceptInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
}

// ChatClient is the service's struct implementation
type ChatClient struct {
	client *stream.Client
}

// NewServiceGetStream initializes a new getstream service
func NewServiceGetStream() ServiceGetStream {
	client, err := stream.NewClient(getStreamAPIKey, getStreamAPISecret)
	if err != nil {
		log.Fatalf("failed to start getstream client: %v", err)
	}

	return &ChatClient{
		client: client,
	}
}

// CreateGetStreamUserToken creates a new token for a user with optional expire time. This token is handed
// to the client side during login. It allows the client side to connect to the chat API for that user.
func (c *ChatClient) CreateGetStreamUserToken(ctx context.Context, userID string) (string, error) {
	return c.client.CreateToken(userID, getStreamTokenExpiryTime)
}

// CreateGetStreamUser creates or updates a user
func (c *ChatClient) CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return c.client.UpsertUser(ctx, user)
}

// ListGetStreamUsers returns list of users that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in the order of sort options.
func (c *ChatClient) ListGetStreamUsers(ctx context.Context, input *stream.QueryOption) (*stream.QueryUsersResponse, error) {
	user, err := c.client.QueryUsers(ctx, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateChannel creates new channel of given type and id or returns already created one.
func (c *ChatClient) CreateChannel(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error) {
	return c.client.CreateChannel(ctx, chanType, chanID, userID, data)
}

// DeleteChannels deletes channels asynchronously.
// Channels and messages will be hard deleted if hardDelete is true.
// It returns an AsyncTaskResponse object which contains the task ID, the status of the task can be check with client.GetTask method.
func (c *ChatClient) DeleteChannels(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
	return c.client.DeleteChannels(ctx, chanIDs, hardDelete)
}

// InviteMembers invites users with given IDs to the channel while at the same time composing a message to show the users
func (c *ChatClient) InviteMembers(ctx context.Context, userIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
	return c.client.Channel("messaging", channelID).InviteMembersWithMessage(ctx, userIDs, message)
}

// ListGetStreamChannels returns list of channels that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in oder of sort options.
func (c *ChatClient) ListGetStreamChannels(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
	return c.client.QueryChannels(ctx, input)
}

// ListChannelMembers returns list of channel members
func (c *ChatClient) ListChannelMembers(ctx context.Context, channelID string, query *stream.QueryOption, sorters ...*stream.SortOption) ([]*stream.ChannelMember, error) {
	return nil, nil
}

// GetChannel retrieves a channel from Getstream using the channel id
func (c *ChatClient) GetChannel(ctx context.Context, channelID string) (*stream.Channel, error) {

	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]interface{}{
				"$eq": channelID,
			},
		},
	}

	resp, err := c.client.QueryChannels(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(resp.Channels) != 1 {
		return nil, fmt.Errorf("expected a single getstream channel, got: %v", len(resp.Channels))
	}

	return resp.Channels[0], nil
}

// AddMembersToCommunity adds the specified clients/staffs to a community
func (c *ChatClient) AddMembersToCommunity(ctx context.Context, userIDs []string, channelID string) (*stream.Response, error) {
	return c.client.Channel("messaging", channelID).AddMembers(ctx, userIDs)
}

// RejectInvite rejects invitation to a getstream channel
func (c *ChatClient) RejectInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {
	return c.client.Channel("messaging", channelID).RejectInvite(ctx, userID, message)
}

// AcceptInvite accepts invitation to a getstream channel
func (c *ChatClient) AcceptInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {
	return c.client.Channel("messaging", channelID).AcceptInvite(ctx, userID, message)
}
