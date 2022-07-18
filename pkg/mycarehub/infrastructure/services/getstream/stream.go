package getstream

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/serverutils"
)

var (
	getStreamTokenExpiryInDays = serverutils.MustGetEnvVar("GET_STREAM_TOKEN_EXPIRY_DAYS")
)

// ServiceGetStream represents the various Getstream usecases
type ServiceGetStream interface {
	CreateGetStreamUserToken(ctx context.Context, userID string) (string, error)
	RevokeGetStreamUserToken(ctx context.Context, userID string) (*stream.Response, error)
	CreateGetStreamUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	ListGetStreamUsers(ctx context.Context, input *stream.QueryOption) (*stream.QueryUsersResponse, error)
	CreateChannel(ctx context.Context, chanType, chanID, userID string, data *stream.ChannelRequest) (*stream.CreateChannelResponse, error)
	DeleteChannels(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error)
	InviteMembers(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error)
	ListGetStreamChannels(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error)
	GetChannel(ctx context.Context, channelID string) (*stream.Channel, error)
	AddMembersToCommunity(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error)
	RejectInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	AcceptInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error)
	RemoveMembersFromCommunity(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error)
	AddModeratorsWithMessage(ctx context.Context, userIDs []string, communityID string, msg *stream.Message) (*stream.Response, error)
	DemoteModerators(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error)
	DeleteUsers(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error)
	BanUser(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error)
	UnBanUser(ctx context.Context, targetID string, communityID string) (bool, error)
	ListCommunityBannedMembers(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error)
	UpsertUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	ListFlaggedMessages(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error)
	DeleteMessage(ctx context.Context, messageID string) (*stream.Response, error)
	ValidateGetStreamRequest(ctx context.Context, body []byte, signature string) bool
	GetStreamUser(ctx context.Context, id string) (*stream.User, error)
	QueryChannelMembers(ctx context.Context, channelID string, input *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryMembersResponse, error)
}

// IStreamClient defines the methods we consume from getstream library
type IStreamClient interface {
	CreateToken(userID string, expire time.Time, issuedAt ...time.Time) (string, error)
	RevokeUserToken(ctx context.Context, userID string, before *time.Time) (*stream.Response, error)
	UpsertUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error)
	QueryUsers(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error)
	CreateChannel(ctx context.Context, chanType, chanID, userID string, data *stream.ChannelRequest) (*stream.CreateChannelResponse, error)
	DeleteChannels(ctx context.Context, cids []string, hardDelete bool) (*stream.AsyncTaskResponse, error)
	QueryChannels(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error)
	Channel(channelType, channelID string) *stream.Channel
	DeleteUsers(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error)
	QueryMessageFlags(ctx context.Context, q *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error)
	DeleteMessage(ctx context.Context, msgID string) (*stream.Response, error)
	VerifyWebhook(body, signature []byte) (valid bool)
	QueryBannedUsers(ctx context.Context, q *stream.QueryBannedUsersOptions, sorters ...*stream.SortOption) (*stream.QueryBannedUsersResponse, error)
}

// ChatClient is the service's struct implementation
type ChatClient struct {
	client IStreamClient
}

// NewServiceGetStream initializes a new getstream service
func NewServiceGetStream(client IStreamClient) *ChatClient {
	return &ChatClient{
		client: client,
	}
}

// CreateGetStreamUserToken creates a new token for a user with optional expire time. This token is handed
// to the client side during login. It allows the client side to connect to the chat API for that user.
func (c *ChatClient) CreateGetStreamUserToken(ctx context.Context, userID string) (string, error) {
	tokenExpiry, err := strconv.Atoi(getStreamTokenExpiryInDays)
	if err != nil {
		return "", fmt.Errorf("failed to convert expiry days to integer: %v", err)
	}
	getStreamTokenExpiryTime := time.Now().UTC().AddDate(0, 0, tokenExpiry)

	return c.client.CreateToken(userID, getStreamTokenExpiryTime, time.Now())
}

// RevokeGetStreamUserToken expires a users token. It sets a `revoke_tokens_issued_before` time which implies
// that any token issued before this time will be considered expired and fail to authenticate.
func (c *ChatClient) RevokeGetStreamUserToken(ctx context.Context, userID string) (*stream.Response, error) {
	revokeTime := time.Now()
	return c.client.RevokeUserToken(ctx, userID, &revokeTime)
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
func (c *ChatClient) CreateChannel(ctx context.Context, chanType, chanID, userID string, data *stream.ChannelRequest) (*stream.CreateChannelResponse, error) {
	response, err := c.client.CreateChannel(ctx, chanType, chanID, userID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	_, err = c.AddModeratorsWithMessage(ctx, []string{userID}, chanID, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteChannels deletes channels asynchronously.
// Channels and messages will be hard deleted if hardDelete is true.
// It returns an AsyncTaskResponse object which contains the task ID, the status of the task can be check with client.GetTask method.
func (c *ChatClient) DeleteChannels(ctx context.Context, chanIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
	return c.client.DeleteChannels(ctx, chanIDs, hardDelete)
}

// InviteMembers invites users with given IDs to the channel while at the same time composing a message to show the users
func (c *ChatClient) InviteMembers(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {

	channel := c.client.Channel("messaging", channelID)

	return channel.InviteMembersWithMessage(ctx, memberIDs, message)
}

// ListGetStreamChannels returns list of channels that match QueryOption.
// If any number of SortOption are set, result will be sorted by field and direction in oder of sort options.
func (c *ChatClient) ListGetStreamChannels(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
	return c.client.QueryChannels(ctx, input)
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
func (c *ChatClient) AddMembersToCommunity(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error) {

	channel := c.client.Channel("messaging", channelID)

	return channel.AddMembers(ctx, memberIDs)
}

// RejectInvite rejects invitation to a getstream channel
func (c *ChatClient) RejectInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {

	channel := c.client.Channel("messaging", channelID)

	return channel.RejectInvite(ctx, userID, message)
}

// AcceptInvite accepts invitation to a getstream channel
func (c *ChatClient) AcceptInvite(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {

	channel := c.client.Channel("messaging", channelID)

	return channel.AcceptInvite(ctx, userID, message)
}

// RemoveMembersFromCommunity deletes members from a community
func (c *ChatClient) RemoveMembersFromCommunity(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error) {

	channel := c.client.Channel("messaging", channelID)

	return channel.RemoveMembers(ctx, memberIDs, message)
}

// AddModeratorsWithMessage adds moderators with given IDs to the channel and produces a message.
func (c *ChatClient) AddModeratorsWithMessage(ctx context.Context, userIDs []string, communityID string, message *stream.Message) (*stream.Response, error) {

	channel := c.client.Channel("messaging", communityID)

	return channel.AddModeratorsWithMessage(ctx, userIDs, message)
}

// DemoteModerators demotes moderators to members
func (c *ChatClient) DemoteModerators(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error) {

	channel := c.client.Channel("messaging", channelID)

	return channel.DemoteModerators(ctx, memberIDs...)
}

// DeleteUsers deletes users from the platform with the specified options.
// Users and messages will be hard deleted if hardDelete is true.
func (c *ChatClient) DeleteUsers(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error) {
	return c.client.DeleteUsers(ctx, userIDs, options)
}

// ListCommunityBannedMembers is used to list members banned from a channel.
func (c *ChatClient) ListCommunityBannedMembers(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error) {
	options := &stream.QueryBannedUsersOptions{
		QueryOption: &stream.QueryOption{Filter: map[string]interface{}{
			"channel_cid": "messaging:" + communityID,
		}},
	}

	response, err := c.client.QueryBannedUsers(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("an error occurred: %w", err)
	}

	return response, nil
}

// BanUser bans a user from a specified channel
func (c *ChatClient) BanUser(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error) {
	similar := strings.EqualFold(targetMemberID, bannedBy)
	if similar {
		return false, fmt.Errorf("users cannot ban themselves from a channel")
	}

	channel := c.client.Channel("messaging", communityID)

	_, err := channel.BanUser(ctx, targetMemberID, bannedBy)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to ban user %s from channel %s: %w", targetMemberID, communityID, err)
	}

	return true, nil
}

// UnBanUser unbans a user who was banned in the specified channel
func (c *ChatClient) UnBanUser(ctx context.Context, targetID string, communityID string) (bool, error) {

	channel := c.client.Channel("messaging", communityID)

	if targetID == "" || communityID == "" {
		return false, fmt.Errorf("neither targetID nor communityID can be empty")
	}

	options := &stream.QueryBannedUsersOptions{
		QueryOption: &stream.QueryOption{Filter: map[string]interface{}{
			"channel_cid": "messaging:" + communityID,
		}},
	}
	response, err := c.client.QueryBannedUsers(ctx, options)
	if err != nil {
		return false, fmt.Errorf("an error occurred: %w", err)
	}

	for _, v := range response.Bans {
		if v.User.ID == targetID {
			_, err = channel.UnBanUser(ctx, targetID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, fmt.Errorf("unable to unban user: %w", err)
			}

			return true, nil
		}
	}

	return false, nil
}

// UpsertUser updates a user's details
func (c *ChatClient) UpsertUser(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
	return c.client.UpsertUser(ctx, user)
}

// ListFlaggedMessages returns a list of message flags for the given user and channel.
func (c *ChatClient) ListFlaggedMessages(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
	return c.client.QueryMessageFlags(ctx, input)
}

// DeleteMessage deletes messages from the platform with the specified options.
func (c *ChatClient) DeleteMessage(ctx context.Context, messageID string) (*stream.Response, error) {
	return c.client.DeleteMessage(ctx, messageID)
}

// ValidateGetStreamRequest verifies the request as coming from getstream and not tampered by a 3rd party
func (c *ChatClient) ValidateGetStreamRequest(ctx context.Context, body []byte, signature string) bool {
	return c.client.VerifyWebhook(body, []byte(signature))
}

// GetStreamUser retrieves a getstream user given the ID
func (c *ChatClient) GetStreamUser(ctx context.Context, id string) (*stream.User, error) {
	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"id": map[string]interface{}{
				"$eq": id,
			},
		},
	}
	response, err := c.client.QueryUsers(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(response.Users) != 1 {
		return nil, fmt.Errorf("expected a single user got: %v", len(response.Users))
	}

	return response.Users[0], nil
}

// QueryChannelMembers returns a list of members for the given community
func (c *ChatClient) QueryChannelMembers(ctx context.Context, channelID string, input *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryMembersResponse, error) {
	channel := c.client.Channel("messaging", channelID)
	return channel.QueryMembers(ctx, input, sorters...)
}
