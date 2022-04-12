package getstream_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	mockGetstream "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
)

var (
	fakeGetstream = mockGetstream.NewGetStreamServiceMock()
)

func TestChatClient_UnitTest_CreateUserGetStreamToken(t *testing.T) {
	streamClient := &stream.Client{}
	g := getstream.NewServiceGetStream(streamClient)

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully generate a user token",
			args: args{
				ctx:    context.Background(),
				userID: "fe9a8f7c-f8f9-4f0c-b8b1-f8b8f8b8f8b8",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to generate token",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad Case - Fail to generate token" {
				fakeGetstream.MockCreateGetStreamUserTokenFn = func(ctx context.Context, userID string) (string, error) {
					return "", fmt.Errorf("failed to generate token")
				}
			}

			got, err := g.CreateGetStreamUserToken(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStreamClient.CreateGetStreamUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == "" {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_ListGetStreamUsers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://example.com/users?api_key=&payload=%7B%22watch%22%3Afalse%2C%22state%22%3Afalse%2C%22presence%22%3Afalse%2C%22filter_conditions%22%3A%7B%22role%22%3A%22user%22%7D%7D",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)

	type args struct {
		ctx   context.Context
		input *stream.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list get stream users",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
					Filter: map[string]interface{}{
						"role": "user",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get users",
			args: args{
				ctx:   context.Background(),
				input: &stream.QueryOption{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - Fail to get users" {
				fakeGetstream.MockListGetStreamUsersFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("failed to get users")
				}
			}

			got, err := g.ListGetStreamUsers(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.ListGetStreamUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_CreateChannel(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/streamTestChannel?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)

	ctx := context.Background()

	type args struct {
		ctx      context.Context
		chanType string
		chanID   string
		userID   string
		data     map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.CreateChannelResponse
		wantErr bool
	}{
		{
			name: "Sad case - empty channel id",
			args: args{
				ctx:      ctx,
				chanType: "test",
				chanID:   "",
				userID:   userToAddToNewChannelID,
				data:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - empty channel id" {
				fakeGetstream.MockCreateChannelFn = func(ctx context.Context, chanType string, chanID string, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error) {
					return nil, fmt.Errorf("failed to create channel")
				}
			}
			got, err := g.CreateChannel(tt.args.ctx, tt.args.chanType, tt.args.chanID, tt.args.userID, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.CreateChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}

}

func TestChatClient_UnitTest_ListGetStreamChannels(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"value": "fixed",
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)

	type args struct {
		ctx   context.Context
		input *stream.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list get stream channels",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
					Filter: map[string]interface{}{
						"type": "messaging",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := g.ListGetStreamChannels(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.ListGetStreamChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_GetChannel(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)

	type args struct {
		ctx       context.Context
		channelID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad case - channel does not exist",
			args: args{
				ctx:       context.Background(),
				channelID: "no-existent-channel",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - channel does not exist" {
				fakeGetstream.MockGetChannel = func(ctx context.Context, channelID string) (*stream.Channel, error) {
					return nil, fmt.Errorf("failed to get channel")
				}
			}

			_, err := g.GetChannel(tt.args.ctx, tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.GetChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestChatClient_UnitTest_RejectInvite(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/testChannelJnJ?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		userID    string
		channelID string
		message   *stream.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.Response
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				userID:    userToRejectInviteID,
				channelID: channelID,
				message:   nil,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid channel id",
			args: args{
				ctx:       ctx,
				userID:    userToRejectInviteID,
				channelID: "",
				message:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.RejectInvite(tt.args.ctx, tt.args.userID, tt.args.channelID, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.RejectInvite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", err)
			}
		})
	}
}

func TestChatClient_UnitTest_AcceptInvite(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/testChannelJnJ?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()
	customInviteMessage := "the user " + userToAcceptInviteName + "accepted the invite"
	type args struct {
		ctx       context.Context
		userID    string
		channelID string
		message   *stream.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.Response
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				userID:    userToAcceptInviteID,
				channelID: channelID,
				message: &stream.Message{
					Text: customInviteMessage,
					User: &stream.User{
						ID:   userToAcceptInviteID,
						Name: userToAcceptInviteName,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: channel does not exist",
			args: args{
				ctx:       ctx,
				userID:    userToAcceptInviteID,
				channelID: "no-existent-channel",
				message:   &stream.Message{Text: customInviteMessage},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.AcceptInvite(tt.args.ctx, tt.args.userID, tt.args.channelID, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.AcceptInvite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ChatClient.AcceptInvite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatClient_UnitTest_RemoveMembers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/testChannelJnJ?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		channelID string
		memberIDs []string
		message   *stream.Message
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.Response
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				channelID: channelID,
				memberIDs: []string{userRemoveFromCommunityID},
				message:   nil,
			},
			wantErr: false,
		},
		{
			name: "Sad case: channel does not exist",
			args: args{
				ctx:       ctx,
				channelID: "no-existent-channel",
				memberIDs: []string{userRemoveFromCommunityID},
				message:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.RemoveMembersFromCommunity(tt.args.ctx, tt.args.channelID, tt.args.memberIDs, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.RemoveMembersFromCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", err)
			}
		})
	}
}

func TestChatClient_UnitTest_DemoteModerators(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/testChannelJnJ?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		channelID string
		memberIDs []string
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.Response
		wantErr bool
	}{
		{
			name: "Happy case: demote moderators",
			args: args{
				ctx:       ctx,
				channelID: channelID,
				memberIDs: []string{moderatorToDemoteID},
			},
			wantErr: false,
		},
		{
			name: "Sad case: non-existent channel",
			args: args{
				ctx:       ctx,
				channelID: "no-existent-channel",
				memberIDs: []string{moderatorToDemoteID},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.DemoteModerators(tt.args.ctx, tt.args.channelID, tt.args.memberIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.DemoteModerators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", err)
			}
		})
	}
}

func TestChatClient_UnitTest_RevokeGetStreamUserToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("PATCH", "/users?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.Response
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				userID: userToRevokeGetstreamTokenID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				userID: "5ea9dc51-a67e-4e5d-aaba-c590c9a66b67",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				httpmock.RegisterResponder("PATCH", "/users?api_key=",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(404, map[string]interface{}{
							"channel": map[string]interface{}{},
						})
						return resp, err
					},
				)
			}
			got, err := g.RevokeGetStreamUserToken(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.RevokeGetStreamUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", err)
			}
		})
	}
}

func TestChatClient_UnitTest_BanUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/testChannelJnJ?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()
	type args struct {
		ctx            context.Context
		targetMemberID string
		bannedBy       string
		communityID    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Sad Case",
			args: args{
				ctx:            ctx,
				targetMemberID: "",
				bannedBy:       "",
				communityID:    "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - same targetMemberID",
			args: args{
				ctx:            ctx,
				targetMemberID: userToBanID,
				bannedBy:       defaultModeratorID,
				communityID:    channelCID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.BanUser(tt.args.ctx, tt.args.targetMemberID, tt.args.bannedBy, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.BanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChatClient.BanUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatClient_UnitTest_UnBanUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/query_banned_users?api_key=&payload=%7B%22watch%22%3Afalse%2C%22state%22%3Afalse%2C%22presence%22%3Afalse%2C%22filter_conditions%22%3A%7B%22channel_cid%22%3A%22messaging%3AtestChannelJnJ%22%7D%7D",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx         context.Context
		targetID    string
		communityID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				targetID:    userToUnbanID,
				communityID: channelID,
			},
			want:    false, // no user banned
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				targetID:    "6d743db7-d2ea-4364-a581-b15bad19ada7",
				communityID: channelID,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.UnBanUser(tt.args.ctx, tt.args.targetID, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.UnBanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChatClient.UnBanUser() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestChatClient_UnitTest_ListCommunityBannedMembers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/query_banned_users?api_key=&payload=%7B%22watch%22%3Afalse%2C%22state%22%3Afalse%2C%22presence%22%3Afalse%2C%22filter_conditions%22%3A%7B%22channel_cid%22%3A%22messaging%3AtestChannelJnJ%22%7D%7D",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx         context.Context
		communityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Member
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				communityID: channelID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := g.ListCommunityBannedMembers(tt.args.ctx, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.ListCommunityBannedMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_UpsertUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/users?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	type args struct {
		ctx  context.Context
		user *stream.User
	}
	tests := []struct {
		name string

		args    args
		want    *stream.UpsertUserResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: context.Background(),
				user: &stream.User{
					ID:   userToUpsertID,
					Name: "Test",
					Role: "moderator",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.UpsertUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.UpsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ChatClient.UpsertUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatClient_UnitTest_DeleteUsers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/users/delete?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userIDs []string
		options stream.DeleteUserOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.AsyncTaskResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				userIDs: []string{userToDeleteID},
				options: stream.DeleteUserOptions{
					User:     stream.HardDelete,
					Messages: stream.HardDelete,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.DeleteUsers(tt.args.ctx, tt.args.userIDs, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.DeleteUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ChatClient.DeleteUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatClient_UnitTest_ListFlaggedMessages(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/moderation/flags/message?api_key=&payload=%7B%22watch%22%3Afalse%2C%22state%22%3Afalse%2C%22presence%22%3Afalse%2C%22filter_conditions%22%3A%7B%22channel_cid%22%3A%22testChannelJnJ%22%7D%7D",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		input *stream.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.QueryMessageFlagsResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: &stream.QueryOption{
					Filter: map[string]interface{}{
						"channel_cid": channelID,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				input: &stream.QueryOption{
					Filter: map[string]interface{}{
						"channel_cid": uuid.New().String(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGetstream.MockListFlaggedMessagesFn = func(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
					return nil, fmt.Errorf("error")
				}
			}
			got, err := g.ListFlaggedMessages(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.ListFlaggedMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_DeleteMessage(t *testing.T) {
	messageID := uuid.New().String()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("DELETE", "/messages/"+messageID+"?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)

	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		messageID string
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.Response
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				messageID: messageID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				messageID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGetstream.MockDeleteMessageFn = func(ctx context.Context, messageID string) (*stream.Response, error) {
					return nil, fmt.Errorf("error")
				}
			}
			got, err := g.DeleteMessage(tt.args.ctx, tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.DeleteMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_CreateGetStreamUser(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/users?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)
	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	type args struct {
		ctx  context.Context
		user *stream.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: context.Background(),
				user: &stream.User{
					ID:   uuid.New().String(),
					Name: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := g.CreateGetStreamUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.CreateGetStreamUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_UnitTest_AddModeratorsWithMessage(t *testing.T) {
	communityID := uuid.New().String()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/channels/messaging/"+communityID+"?api_key=",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"channel": map[string]interface{}{},
			})
			return resp, err
		},
	)
	streamClient := &stream.Client{
		BaseURL: "https://example.com",
		HTTP:    &http.Client{},
	}
	g := getstream.NewServiceGetStream(streamClient)
	type args struct {
		ctx         context.Context
		userIDs     []string
		communityID string
		message     *stream.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         context.Background(),
				userIDs:     []string{uuid.NewString()},
				communityID: communityID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.AddModeratorsWithMessage(tt.args.ctx, tt.args.userIDs, tt.args.communityID, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.AddModeratorsWithMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
