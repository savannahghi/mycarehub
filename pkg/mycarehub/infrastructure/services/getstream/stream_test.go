package getstream_test

import (
	"context"
	"fmt"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()

			if tt.name == "Sad Case - Fail to get users" {
				fakeGetStream.MockListGetStreamUsersFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("failed to get users")
				}
			}

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

	//TODO: Delete the created test channel from stream server
	_, err := g.DeleteChannels(ctx, []string{channelID}, true)
	if err != nil {
		t.Errorf("ChatClient.DeleteChannel() error = %v", err)
	}
}
