package getstream_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	mockGetstream "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
)

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
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy Case - Successfully list get stream users" {
				fakeClient.MockQueryUsersFn = func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
					return &stream.QueryUsersResponse{}, nil
				}
			}

			if tt.name == "Sad Case - Fail to get users" {
				fakeClient.MockQueryUsersFn = func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("fail to query users")
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
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy Case - Successfully list get stream channels" {
				fakeClient.MockQueryChannelsFn = func(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error) {
					return &stream.QueryChannelsResponse{
						Channels: []*stream.Channel{{ID: gofakeit.UUID()}},
					}, nil
				}
			}

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
			name: "Happy case - retrieve channel by id",
			args: args{
				ctx:       context.Background(),
				channelID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case - multiple channels found",
			args: args{
				ctx:       context.Background(),
				channelID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - non existent channel",
			args: args{
				ctx:       context.Background(),
				channelID: "no-existent-channel",
			},
			wantErr: true,
		},
		{
			name: "Sad case - error retrieving channel",
			args: args{
				ctx:       context.Background(),
				channelID: "no-existent-channel",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Sad case - multiple channels found" {
				fakeClient.MockQueryChannelsFn = func(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error) {
					return &stream.QueryChannelsResponse{
						Channels: []*stream.Channel{{ID: gofakeit.UUID()}, {ID: gofakeit.UUID()}},
					}, nil
				}
			}

			if tt.name == "Happy case - retrieve channel by id" {
				fakeClient.MockQueryChannelsFn = func(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error) {
					return &stream.QueryChannelsResponse{
						Channels: []*stream.Channel{{ID: gofakeit.UUID()}},
					}, nil
				}
			}

			if tt.name == "Sad case - non existent channel" {
				fakeClient.MockQueryChannelsFn = func(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error) {
					return &stream.QueryChannelsResponse{
						Channels: []*stream.Channel{},
					}, nil
				}
			}

			if tt.name == "Sad case - error retrieving channel" {
				fakeClient.MockQueryChannelsFn = func(ctx context.Context, q *stream.QueryOption, sort ...*stream.SortOption) (*stream.QueryChannelsResponse, error) {
					return nil, fmt.Errorf("failed to retrieve channel")
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

func TestChatClient_RevokeGetStreamUserToken(t *testing.T) {

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
			name: "Happy case: revoke stream token",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: revoke stream token" {
				fakeClient.MockRevokeUserTokenFn = func(ctx context.Context, userID string, before *time.Time) (*stream.Response, error) {
					return &stream.Response{}, nil
				}
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

func TestChatClient_ListCommunityBannedMembers(t *testing.T) {

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
			name: "Happy case: list banned members",
			args: args{
				ctx:         context.Background(),
				communityID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: error listing banned members",
			args: args{
				ctx:         context.Background(),
				communityID: gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: list banned members" {
				fakeClient.MockQueryBannedUsersFn = func(ctx context.Context, q *stream.QueryBannedUsersOptions, sorters ...*stream.SortOption) (*stream.QueryBannedUsersResponse, error) {
					return &stream.QueryBannedUsersResponse{}, nil
				}
			}

			if tt.name == "sad case: error listing banned members" {
				fakeClient.MockQueryBannedUsersFn = func(ctx context.Context, q *stream.QueryBannedUsersOptions, sorters ...*stream.SortOption) (*stream.QueryBannedUsersResponse, error) {
					return nil, fmt.Errorf("fail to query banned users")
				}
			}

			_, err := g.ListCommunityBannedMembers(tt.args.ctx, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.ListCommunityBannedMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChatClient_UpsertUser(t *testing.T) {
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
			name: "Happy case: upsert user",
			args: args{
				ctx: context.Background(),
				user: &stream.User{
					ID:   gofakeit.UUID(),
					Name: gofakeit.Name(),
					Role: "moderator",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to upsert user",
			args: args{
				ctx: context.Background(),
				user: &stream.User{
					ID:   gofakeit.UUID(),
					Name: gofakeit.Name(),
					Role: "moderator",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: upsert user" {
				fakeClient.MockUpsertUserFn = func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
					return &stream.UpsertUserResponse{User: &stream.User{ID: gofakeit.UUID()}}, nil
				}
			}

			if tt.name == "Sad case: fail to upsert user" {
				fakeClient.MockUpsertUserFn = func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
					return nil, fmt.Errorf("failed to upsert user")
				}
			}
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

func TestChatClient_DeleteUsers(t *testing.T) {

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
			name: "Happy case: success delete user",
			args: args{
				ctx:     context.Background(),
				userIDs: []string{gofakeit.UUID()},
				options: stream.DeleteUserOptions{
					User:     stream.HardDelete,
					Messages: stream.HardDelete,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: error deleting user",
			args: args{
				ctx:     context.Background(),
				userIDs: []string{gofakeit.UUID()},
				options: stream.DeleteUserOptions{
					User:     stream.HardDelete,
					Messages: stream.HardDelete,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: success delete user" {
				fakeClient.MockDeleteUsersFn = func(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error) {
					return &stream.AsyncTaskResponse{}, nil
				}
			}

			if tt.name == "Sad Case: error deleting user" {
				fakeClient.MockDeleteUsersFn = func(ctx context.Context, userIDs []string, options stream.DeleteUserOptions) (*stream.AsyncTaskResponse, error) {
					return nil, fmt.Errorf("failed to delete user")
				}

			}

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

func TestChatClient_ListFlaggedMessages(t *testing.T) {

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
			name: "Happy case: list flagged messages",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
					Filter: map[string]interface{}{
						"channel_cid": gofakeit.UUID(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: error listing messages",
			args: args{
				ctx: context.Background(),
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
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: list flagged messages" {
				fakeClient.MockQueryMessageFlagsFn = func(ctx context.Context, q *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
					return &stream.QueryMessageFlagsResponse{}, nil
				}
			}

			if tt.name == "Sad case: error listing messages" {
				fakeClient.MockQueryMessageFlagsFn = func(ctx context.Context, q *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
					return nil, fmt.Errorf("failed to query flagged messages")
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

func TestChatClient_DeleteMessage(t *testing.T) {
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
			name: "Happy case: delete a message",
			args: args{
				ctx:       context.Background(),
				messageID: gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to delete message",
			args: args{
				ctx:       context.Background(),
				messageID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: delete a message" {
				fakeClient.MockDeleteMessageFn = func(ctx context.Context, msgID string) (*stream.Response, error) {
					return &stream.Response{}, nil
				}
			}

			if tt.name == "Sad case: fail to delete message" {
				fakeClient.MockDeleteMessageFn = func(ctx context.Context, msgID string) (*stream.Response, error) {
					return nil, fmt.Errorf("failed to delete messages")
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

func TestChatClient_CreateGetStreamUser(t *testing.T) {
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
			name: "Happy case: create user",
			args: args{
				ctx: context.Background(),
				user: &stream.User{
					ID:   uuid.New().String(),
					Name: gofakeit.Name(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to upsert user",
			args: args{
				ctx: context.Background(),
				user: &stream.User{
					ID:   uuid.New().String(),
					Name: gofakeit.Name(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			g := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "Happy case: create user" {
				fakeClient.MockUpsertUserFn = func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
					return &stream.UpsertUserResponse{User: &stream.User{ID: gofakeit.UUID()}}, nil
				}
			}

			if tt.name == "Sad case: fail to upsert user" {
				fakeClient.MockUpsertUserFn = func(ctx context.Context, user *stream.User) (*stream.UpsertUserResponse, error) {
					return nil, fmt.Errorf("failed to upsert user")
				}
			}

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

func TestChatClient_CreateGetStreamUserToken(t *testing.T) {
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
			name: "happy Case - successfully generate a user token",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			want:    "none",
			wantErr: false,
		},
		{
			name: "sad Case - error generating a user token",
			args: args{
				ctx:    context.Background(),
				userID: gofakeit.UUID(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			c := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "happy Case - successfully generate a user token" {
				fakeClient.MockCreateTokenFn = func(userID string, expire time.Time, issuedAt ...time.Time) (string, error) {
					return "none", nil
				}
			}

			if tt.name == "sad Case - error generating a user token" {
				fakeClient.MockCreateTokenFn = func(userID string, expire time.Time, issuedAt ...time.Time) (string, error) {
					return "", fmt.Errorf("fail to generate token")
				}
			}

			got, err := c.CreateGetStreamUserToken(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.CreateGetStreamUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChatClient.CreateGetStreamUserToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatClient_DeleteChannels(t *testing.T) {
	type args struct {
		ctx        context.Context
		chanIDs    []string
		hardDelete bool
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.AsyncTaskResponse
		wantErr bool
	}{
		{
			name: "happy case: delete channels",
			args: args{
				ctx:        context.Background(),
				chanIDs:    []string{gofakeit.UUID()},
				hardDelete: true,
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to delete channels",
			args: args{
				ctx:        context.Background(),
				chanIDs:    []string{gofakeit.UUID()},
				hardDelete: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			c := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "happy case: delete channels" {
				fakeClient.MockDeleteChannelsFn = func(ctx context.Context, cids []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
					return &stream.AsyncTaskResponse{}, nil
				}
			}

			if tt.name == "sad case: fail to delete channels" {
				fakeClient.MockDeleteChannelsFn = func(ctx context.Context, cids []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {
					return nil, fmt.Errorf("fail to delete channels")
				}
			}

			got, err := c.DeleteChannels(tt.args.ctx, tt.args.chanIDs, tt.args.hardDelete)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.DeleteChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_GetStreamUser(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *stream.User
		wantErr bool
	}{
		{
			name: "happy case: retrieve stream user",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: false,
		},
		{
			name: "sad case: multiple stream users",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: no stream user",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving stream user",
			args: args{
				ctx: context.Background(),
				id:  gofakeit.UUID(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			c := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "happy case: retrieve stream user" {
				fakeClient.MockQueryUsersFn = func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
					return &stream.QueryUsersResponse{Users: []*stream.User{{ID: gofakeit.UUID()}}}, nil
				}
			}

			if tt.name == "sad case: multiple stream users" {
				fakeClient.MockQueryUsersFn = func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
					return &stream.QueryUsersResponse{Users: []*stream.User{{ID: gofakeit.UUID()}, {ID: gofakeit.UUID()}}}, nil
				}
			}

			if tt.name == "sad case: no stream user" {
				fakeClient.MockQueryUsersFn = func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
					return &stream.QueryUsersResponse{Users: []*stream.User{}}, nil
				}
			}

			if tt.name == "sad case: error retrieving stream user" {
				fakeClient.MockQueryUsersFn = func(ctx context.Context, q *stream.QueryOption, sorters ...*stream.SortOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("failed to query users")
				}
			}

			got, err := c.GetStreamUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChatClient.GetStreamUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestChatClient_ValidateGetStreamRequest(t *testing.T) {
	type args struct {
		ctx       context.Context
		body      []byte
		signature string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case: valid stream hook request",
			args: args{
				ctx:       context.Background(),
				body:      make([]byte, 10),
				signature: gofakeit.HackerPhrase(),
			},
			want: true,
		},
		{
			name: "sad case: invalid stream hook request",
			args: args{
				ctx:       context.Background(),
				body:      make([]byte, 10),
				signature: gofakeit.HackerPhrase(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := mockGetstream.NewStreamClientMock()
			c := getstream.NewServiceGetStream(fakeClient)

			if tt.name == "happy case: valid stream hook request" {
				fakeClient.MockVerifyWebhookFn = func(body, signature []byte) (valid bool) {
					return true
				}
			}

			if tt.name == "sad case: invalid stream hook request" {
				fakeClient.MockVerifyWebhookFn = func(body, signature []byte) (valid bool) {
					return false
				}
			}

			if got := c.ValidateGetStreamRequest(tt.args.ctx, tt.args.body, tt.args.signature); got != tt.want {
				t.Errorf("ChatClient.ValidateGetStreamRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
