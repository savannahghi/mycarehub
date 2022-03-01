package getstream_test

import (
	"context"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
)

func TestGetStreamClient_CreateUserGetStreamToken(t *testing.T) {
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
				userID: uuid.New().String(),
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
			getStream := getstream.NewServiceGetStream()
			got, err := getStream.CreateGetStreamUserToken(tt.args.ctx, tt.args.userID)
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

func TestChatClient_ListGetStreamUsers(t *testing.T) {
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
			getStream := getstream.NewServiceGetStream()

			got, err := getStream.ListGetStreamUsers(tt.args.ctx, tt.args.input)
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

func TestChatClient_CreateChannel(t *testing.T) {
	g := getstream.NewServiceGetStream()

	ctx := context.Background()
	channelID := "channelJnJ"

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
			name: "Happy case - Create channel",
			args: args{
				ctx:      ctx,
				chanType: "messaging",
				chanID:   channelID,
				userID:   uuid.New().String(),
				data: map[string]interface{}{
					"age": map[string]interface{}{
						"lowerBound": 10,
						"upperBound": 20,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case - unable to create channel",
			args: args{
				ctx:      ctx,
				chanType: "test",
				chanID:   "",
				userID:   uuid.New().String(),
				data:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

	_, err := g.DeleteChannels(ctx, []string{"messaging:" + channelID}, true)
	if err != nil {
		t.Errorf("ChatClient.DeleteChannel() error = %v", err)
	}
}

func TestChatClient_ListGetStreamChannels(t *testing.T) {
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
			getStream := getstream.NewServiceGetStream()

			got, err := getStream.ListGetStreamChannels(tt.args.ctx, tt.args.input)
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

func TestChatClient_GetChannel(t *testing.T) {
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
			name: "Happy case - successfully retrieve a getstream channel",
			args: args{
				ctx:       context.Background(),
				channelID: channelID,
			},
			wantErr: false,
		},
		{
			name: "Sad case - channel does not exist",
			args: args{
				ctx:       context.Background(),
				channelID: "no-existent-channel",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid channel id",
			args: args{
				ctx:       context.Background(),
				channelID: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			_, err := c.GetChannel(tt.args.ctx, tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.GetChannel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestChatClient_RejectInvite(t *testing.T) {
	ctx := context.Background()
	invitedUserID := uuid.New().String()
	user := stream.User{
		ID:        invitedUserID,
		Name:      "test",
		Invisible: false,
	}
	streamUser, err := c.CreateGetStreamUser(ctx, &user)
	if err != nil {
		t.Errorf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	_, err = c.InviteMembers(ctx, []string{streamUser.User.ID}, ch.Channel.ID, nil)
	if err != nil {
		t.Errorf("ChatClient.InviteMembers() error = %v", err)
		return
	}
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
				userID:    invitedUserID,
				channelID: channelID,
				message:   nil,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid user id",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				channelID: channelID,
				message:   nil,
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid channel id",
			args: args{
				ctx:       ctx,
				userID:    invitedUserID,
				channelID: "",
				message:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.RejectInvite(tt.args.ctx, tt.args.userID, tt.args.channelID, tt.args.message)
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

func TestChatClient_AcceptInvite(t *testing.T) {
	ctx := context.Background()
	invitedUserID := uuid.New().String()
	user := stream.User{
		ID:        invitedUserID,
		Name:      "test user accepted invite",
		Invisible: false,
	}
	streamUser, err := c.CreateGetStreamUser(ctx, &user)
	if err != nil {
		t.Errorf("ChatClient.CreateGetStreamUser() error = %v", err)
		return
	}

	_, err = c.InviteMembers(ctx, []string{streamUser.User.ID}, ch.Channel.ID, nil)
	if err != nil {
		t.Errorf("ChatClient.InviteMembers() error = %v", err)
		return
	}
	customInviteMessage := "the user" + user.Name + "accepted the invite"
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
				userID:    invitedUserID,
				channelID: channelID,
				message: &stream.Message{
					Text: customInviteMessage,
					User: &stream.User{
						ID:   user.ID,
						Name: user.Name,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid user id",
			args: args{
				ctx:       ctx,
				userID:    uuid.New().String(),
				channelID: channelID,
				message:   nil,
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid channel id",
			args: args{
				ctx:       ctx,
				userID:    invitedUserID,
				channelID: uuid.New().String(),
				message:   &stream.Message{Text: customInviteMessage},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.AcceptInvite(tt.args.ctx, tt.args.userID, tt.args.channelID, tt.args.message)
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

func TestChatClient_RemoveMembers(t *testing.T) {
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
				memberIDs: []string{member1},
				message:   nil,
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid channel id",
			args: args{
				ctx:       ctx,
				channelID: uuid.New().String(),
				memberIDs: []string{member1},
				message:   nil,
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid user id",
			args: args{
				ctx:       ctx,
				channelID: channelID,
				memberIDs: []string{uuid.New().String()},
				message:   nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.RemoveMembersFromCommunity(tt.args.ctx, tt.args.channelID, tt.args.memberIDs, tt.args.message)
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
