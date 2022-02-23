package communities_test

import (
	"context"
	"fmt"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities/mock"
)

func TestUseCaseStreamImpl_CreateCommunity(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		input dto.CommunityInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Community
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				input: dto.CommunityInput{
					Name:        "test",
					Description: "test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 10,
						UpperBound: 20,
					},
					Gender:     []*enumutils.Gender{&enumutils.AllGender[0]},
					ClientType: []*enums.ClientType{&enums.AllClientType[0]},
					InviteOnly: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case - cannot create channel in the database",
			args: args{
				ctx: ctx,
				input: dto.CommunityInput{
					Name:        "test",
					Description: "test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:     []*enumutils.Gender{&enumutils.AllGender[0]},
					ClientType: []*enums.ClientType{&enums.AllClientType[0]},
					InviteOnly: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to get logged in user",
			args: args{
				ctx: ctx,
				input: dto.CommunityInput{
					Name:        "test",
					Description: "test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:     []*enumutils.Gender{&enumutils.AllGender[0]},
					ClientType: []*enums.ClientType{&enums.AllClientType[0]},
					InviteOnly: false,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to create streams channel",
			args: args{
				ctx: ctx,
				input: dto.CommunityInput{
					Name:        "test",
					Description: "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewCommunityUsecaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()

			c := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB)

			if tt.name == "Sad case - cannot create channel in the database" {
				fakeDB.MockCreateChannelFn = func(ctx context.Context, community *dto.CommunityInput) (*domain.Community, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - fail to create streams channel" {
				fakeGetStream.MockCreateChannelFn = func(ctx context.Context, chanType, chanID, userID string, data map[string]interface{}) (*stream.CreateChannelResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.CreateCommunity(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseStreamImpl.CreateCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ListGetStreamUsers(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *domain.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list stream users",
			args: args{
				ctx: context.Background(),
				input: &domain.QueryOption{
					Filter: map[string]interface{}{
						"role": "user",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to list stream users",
			args: args{
				ctx: context.Background(),
				input: &domain.QueryOption{
					Filter: map[string]interface{}{
						"role": "user",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB)

			if tt.name == "Sad Case - Fail to list stream users" {
				fakeGetStream.MockListGetStreamUsersFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("failed to get users")
				}
			}

			got, err := communities.ListGetStreamUsers(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListGetStreamUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_InviteMembers(t *testing.T) {
	type args struct {
		ctx         context.Context
		communityID string
		userIDS     []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully invite members",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				userIDS: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get logged in user",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				userIDS: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff user",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				userIDS: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to invite members",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				userIDS: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB)

			if tt.name == "Sad Case - Fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "Sad Case - Fail to get staff user" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}

			if tt.name == "Sad Case - Fail to invite members" {
				fakeGetStream.MockInviteMembersFn = func(ctx context.Context, userIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
					return nil, fmt.Errorf("failed to invite members")
				}
			}

			got, err := communities.InviteMembers(tt.args.ctx, tt.args.communityID, tt.args.userIDS)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.InviteMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.InviteMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ListGetStreamChannels(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *domain.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list stream channels",
			args: args{
				ctx: context.Background(),
				input: &domain.QueryOption{
					Filter: map[string]interface{}{
						"type": "channel",
					},
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully list stream channels, no params",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to list stream channels",
			args: args{
				ctx: context.Background(),
				input: &domain.QueryOption{
					Filter: map[string]interface{}{
						"type": "channel",
					},
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB)

			if tt.name == "Sad Case - Fail to list stream channels" {
				fakeGetStream.MockListGetStreamChannelsFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
					return nil, fmt.Errorf("failed to get channels")
				}
			}

			got, err := communities.ListGetStreamChannels(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListGetStreamChannels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ListCommunityMembers(t *testing.T) {
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB)

	type args struct {
		ctx         context.Context
		communityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - success list community members",
			args: args{
				ctx:         context.Background(),
				communityID: "test-community",
			},
			wantErr: false,
		},
		{
			name: "Sad case - fail invalid community id",
			args: args{
				ctx:         context.Background(),
				communityID: "test-community",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case - success list community members" {
				fakeGetStream.MockGetChannel = func(ctx context.Context, channelID string) (*stream.Channel, error) {

					user := &stream.User{
						ID:   uuid.NewString(),
						Name: "john doe",
						Role: "user",
						ExtraData: map[string]interface{}{
							"userType": "CLIENT",
							"userID":   uuid.NewString(),
						},
					}

					return &stream.Channel{
						Members: []*stream.ChannelMember{
							{
								User:        user,
								Role:        "member",
								IsModerator: false,
							},
						},
					}, nil
				}
			}

			if tt.name == "Sad case - fail invalid community id" {
				fakeGetStream.MockGetChannel = func(ctx context.Context, channelID string) (*stream.Channel, error) {
					return nil, fmt.Errorf("channel does not exist")
				}
			}

			_, err := communities.ListCommunityMembers(tt.args.ctx, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListCommunityMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
