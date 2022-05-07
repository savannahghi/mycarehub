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
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities/mock"
	notificationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
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
		{
			name: "Sad case - fail to get staff profile by logged in user id",
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewCommunityUsecaseMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			c := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case - cannot create channel in the database" {
				fakeDB.MockCreateCommunityFn = func(ctx context.Context, community *dto.CommunityInput) (*domain.Community, error) {
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

			if tt.name == "Sad case - fail to get staff profile by logged in user id" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, uid string) (*domain.StaffProfile, error) {
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

func TestUseCasesCommunitiesImpl_ListMembers(t *testing.T) {
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
			name: "Happy Case - Successfully list stream users",
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
			name: "Sad Case - Fail to list stream users",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad Case - Fail to list stream users" {
				fakeGetStream.MockListGetStreamUsersFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("failed to get users")
				}
			}

			got, err := communities.ListMembers(tt.args.ctx, tt.args.input)
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
		memberIDs   []string
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
				memberIDs: []string{
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
				memberIDs: []string{
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
				memberIDs: []string{
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
				memberIDs: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to get community",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				memberIDs: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to retrieve getstream user",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				memberIDs: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to retrieve user profile",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				memberIDs: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to notify user",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				memberIDs: []string{
					uuid.NewString(),
					uuid.NewString(),
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: fail to get community" {
				fakeDB.MockGetCommunityByIDFn = func(ctx context.Context, communityID string) (*domain.Community, error) {
					return nil, fmt.Errorf("failed to retrieve community")
				}
			}

			if tt.name == "sad case: fail to retrieve getstream user" {
				fakeGetStream.MockGetStreamUserFn = func(ctx context.Context, id string) (*stream.User, error) {
					return nil, fmt.Errorf("failed to retrieve getstream user")
				}
			}

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
				fakeGetStream.MockInviteMembersFn = func(ctx context.Context, memberIDs []string, channelID string, message *stream.Message) (*stream.Response, error) {
					return nil, fmt.Errorf("failed to invite members")
				}
			}

			if tt.name == "sad case: fail to retrieve user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad case: fail to notify user" {
				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify user")
				}
			}

			got, err := communities.InviteMembers(tt.args.ctx, tt.args.communityID, tt.args.memberIDs)
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

func TestUseCasesCommunitiesImpl_ListCommunities(t *testing.T) {
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
			name: "Happy Case - Successfully list stream channels",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
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
			name: "Happy Case - Successfully list stream channels, with limit",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
					Limit: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to list stream channels",
			args: args{
				ctx: context.Background(),
				input: &stream.QueryOption{
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad Case - Fail to list stream channels" {
				fakeGetStream.MockListGetStreamChannelsFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
					return nil, fmt.Errorf("failed to get channels")
				}
			}

			got, err := communities.ListCommunities(tt.args.ctx, tt.args.input)
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
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

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

func TestUseCasesCommunitiesImpl_DeleteChannels(t *testing.T) {

	type args struct {
		ctx context.Context

		communityIDs []string

		hardDelete bool
	}

	tests := []struct {
		name string

		args args

		want bool

		wantErr bool
	}{

		{

			name: "Happy Case - Successfully delete channels",

			args: args{

				ctx: context.Background(),

				communityIDs: []string{uuid.NewString()},

				hardDelete: false,
			},

			want: true,

			wantErr: false,
		},

		{

			name: "Sad Case - Fail to delete channels",

			args: args{

				ctx: context.Background(),

				communityIDs: []string{uuid.NewString()},

				hardDelete: false,
			},

			want: false,

			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			fakeGetStream := getStreamMock.NewGetStreamServiceMock()

			fakeExtension := extensionMock.NewFakeExtension()

			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad Case - Fail to delete channels" {

				fakeGetStream.MockDeleteChannelsFn = func(ctx context.Context, communityIDs []string, hardDelete bool) (*stream.AsyncTaskResponse, error) {

					return nil, fmt.Errorf("failed to delete channels")

				}

			}

			got, err := communities.DeleteCommunities(tt.args.ctx, tt.args.communityIDs, tt.args.hardDelete)

			if (err != nil) != tt.wantErr {

				t.Errorf("UseCasesCommunitiesImpl.DeleteCommunities() error = %v, wantErr %v", err, tt.wantErr)

				return

			}

			if got != tt.want {

				t.Errorf("UseCasesCommunitiesImpl.DeleteCommunities() = %v, want %v", got, tt.want)

			}

		})

	}

}

func TestUseCasesCommunitiesImpl_RejectInvite(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		userID    string
		channelID string
		message   string
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
				ctx:       ctx,
				userID:    uuid.New().String(),
				channelID: uuid.New().String(),
				message:   uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				channelID: uuid.New().String(),
				message:   uuid.New().String(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case" {
				fakeGetStream.MockRejectInviteFn = func(ctx context.Context, userID, channelID string, message *stream.Message) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.RejectInvite(tt.args.ctx, tt.args.userID, tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.RejectInvite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.RejectInvite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_AcceptInvite(t *testing.T) {
	type args struct {
		ctx       context.Context
		userID    string
		channelID string
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
				ctx:       context.Background(),
				userID:    uuid.New().String(),
				channelID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: userID is empty",
			args: args{
				ctx:       context.Background(),
				channelID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: channelID is empty",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to accept invite",
			args: args{
				ctx:       context.Background(),
				userID:    uuid.New().String(),
				channelID: uuid.New().String(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case: failed to accept invite" {
				fakeGetStream.MockAcceptInviteFn = func(ctx context.Context, userID string, channelID string, message *stream.Message) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.AcceptInvite(tt.args.ctx, tt.args.userID, tt.args.channelID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.AcceptInvite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.AcceptInvite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_AddMembersToCommunity(t *testing.T) {
	ctx := context.Background()

	userID := uuid.New().String()
	communityID := uuid.New().String()

	type args struct {
		ctx         context.Context
		memberIDs   []string
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
				memberIDs:   []string{userID},
				communityID: communityID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				memberIDs:   []string{userID},
				communityID: communityID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no community ID",
			args: args{
				ctx:         ctx,
				memberIDs:   []string{userID},
				communityID: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no user ID(s)",
			args: args{
				ctx:         ctx,
				memberIDs:   nil,
				communityID: communityID,
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case" {
				fakeGetStream.MockAddMembersToCommunityFn = func(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no community ID" {
				fakeGetStream.MockAddMembersToCommunityFn = func(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no user ID" {
				fakeGetStream.MockAddMembersToCommunityFn = func(ctx context.Context, memberIDs []string, channelID string) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.AddMembersToCommunity(tt.args.ctx, tt.args.memberIDs, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.AddMembersToCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.AddMembersToCommunity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_RemoveMembers(t *testing.T) {
	type args struct {
		ctx       context.Context
		channelID string
		memberIDs []string
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
				ctx:       context.Background(),
				channelID: uuid.New().String(),
				memberIDs: []string{uuid.New().String()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: channelID is empty",
			args: args{
				ctx:       context.Background(),
				memberIDs: []string{uuid.New().String()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: memberIDs is empty",
			args: args{
				ctx:       context.Background(),
				channelID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to remove members",
			args: args{
				ctx:       context.Background(),
				channelID: uuid.New().String(),
				memberIDs: []string{uuid.New().String()},
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case: failed to remove members" {
				fakeGetStream.MockRemoveMembersFn = func(ctx context.Context, channelID string, memberIDs []string, message *stream.Message) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.RemoveMembersFromCommunity(tt.args.ctx, tt.args.channelID, tt.args.memberIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.RemoveMembersFromCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.RemoveMembersFromCommunity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_AddModeratorsWithMessage(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New().String()

	type args struct {
		ctx         context.Context
		userIDs     []string
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
				userIDs:     []string{userID},
				communityID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				userIDs:     []string{userID},
				communityID: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - Unable to get community",
			args: args{
				ctx:         ctx,
				userIDs:     []string{userID},
				communityID: uuid.New().String(),
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case" {
				fakeGetStream.MockAddModeratorsWithMessageFn = func(ctx context.Context, userIDs []string, communityID string, message *stream.Message) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - Unable to get community" {
				fakeDB.MockGetCommunityByIDFn = func(ctx context.Context, communityID string) (*domain.Community, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := communities.AddModeratorsWithMessage(tt.args.ctx, tt.args.userIDs, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.AddModeratorsWithMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.AddModeratorsWithMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_DemoteModerators(t *testing.T) {
	type args struct {
		ctx         context.Context
		communityID string
		memberIDs   []string
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
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				memberIDs:   []string{uuid.New().String()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: communityID is empty",
			args: args{
				ctx:       context.Background(),
				memberIDs: []string{uuid.New().String()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: memberIDs is empty",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to demote moderators",
			args: args{
				ctx:         context.Background(),
				communityID: uuid.New().String(),
				memberIDs:   []string{uuid.New().String()},
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
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad case: failed to demote moderators" {
				fakeGetStream.MockDemoteModeratorsFn = func(ctx context.Context, channelID string, memberIDs []string) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.DemoteModerators(tt.args.ctx, tt.args.communityID, tt.args.memberIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.DemoteModerators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.DemoteModerators() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ListPendingInvites(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx      context.Context
		memberID string
		input    *stream.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Community
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      ctx,
				memberID: uuid.New().String(),
				input: &stream.QueryOption{
					Filter: map[string]interface{}{
						"invite": "pending",
					},
					UserID: uuid.New().String(),
					Limit:  10,
					Offset: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				memberID: "",
				input:    &stream.QueryOption{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			_, err := communities.ListPendingInvites(tt.args.ctx, tt.args.memberID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListPendingInvites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_RecommendedCommunities(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
		limit    int
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Community
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				limit:    10,
			},
			wantErr: false,
		},
		{
			name: "sad case: missing clientID",
			args: args{
				ctx:      context.Background(),
				clientID: "",
				limit:    10,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed client profile by client ID",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				limit:    10,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get user profile by user ID",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				limit:    10,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get recommended channels",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				limit:    10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: failed client profile by client ID" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by client ID")
				}
			}

			if tt.name == "sad case: failed to get user profile by user ID" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user ID")
				}
			}

			if tt.name == "sad case: failed to get recommended channels" {
				fakeGetStream.MockListGetStreamChannelsFn = func(ctx context.Context, input *stream.QueryOption) (*stream.QueryChannelsResponse, error) {
					return nil, fmt.Errorf("failed to get recommended channels")
				}
			}

			got, err := communities.RecommendedCommunities(tt.args.ctx, tt.args.clientID, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.RecommendedCommunities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesCommunitiesImpl.RecommendedCommunities() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ListCommunityBannedMembers(t *testing.T) {
	ctx := context.Background()

	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

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
				communityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				communityID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty community ID",
			args: args{
				ctx:         ctx,
				communityID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGetStream.MockListCommunityBannedMembersFn = func(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty community ID" {
				fakeGetStream.MockListCommunityBannedMembersFn = func(ctx context.Context, communityID string) (*stream.QueryBannedUsersResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := communities.ListCommunityBannedMembers(tt.args.ctx, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListCommunityBannedMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_BanUser(t *testing.T) {
	ctx := context.Background()

	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

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
			name: "Happy case",
			args: args{
				ctx:            ctx,
				targetMemberID: uuid.New().String(),
				bannedBy:       uuid.New().String(),
				communityID:    uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:            ctx,
				targetMemberID: uuid.New().String(),
				bannedBy:       uuid.New().String(),
				communityID:    uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case",
			args: args{
				ctx:            ctx,
				targetMemberID: uuid.New().String(),
				bannedBy:       uuid.New().String(),
				communityID:    uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty target ID",
			args: args{
				ctx:            ctx,
				targetMemberID: "",
				bannedBy:       uuid.New().String(),
				communityID:    uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGetStream.MockBanUserFn = func(ctx context.Context, targetMemberID, bannedBy, communityID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty target ID" {
				fakeGetStream.MockBanUserFn = func(ctx context.Context, targetMemberID, bannedBy, communityID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := communities.BanUser(tt.args.ctx, tt.args.targetMemberID, tt.args.bannedBy, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.BanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.BanUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_UnBanUser(t *testing.T) {
	ctx := context.Background()

	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

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
				targetID:    uuid.New().String(),
				communityID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				targetID:    uuid.New().String(),
				communityID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no community ID",
			args: args{
				ctx:         ctx,
				targetID:    uuid.New().String(),
				communityID: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - user does no belong to the community",
			args: args{
				ctx:         ctx,
				targetID:    uuid.New().String(),
				communityID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGetStream.MockUnBanUserFn = func(ctx context.Context, targetID, communityID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no community ID" {
				fakeGetStream.MockUnBanUserFn = func(ctx context.Context, targetID, communityID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := communities.UnBanUser(tt.args.ctx, tt.args.targetID, tt.args.communityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.UnBanUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.UnBanUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ListFlaggedMessages(t *testing.T) {
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

	ctx := context.Background()
	communityID := uuid.New().String()
	userID := uuid.New().String()
	userIIDs := []*string{&userID}

	type args struct {
		ctx          context.Context
		communityCID *string
		memberIDs    []*string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []*domain.MessageFlag
	}{
		{
			name: "Happy case",
			args: args{
				ctx:          ctx,
				communityCID: &communityID,
				memberIDs:    userIIDs,
			},
			wantErr: false,
		},
		{
			name: "Happy case - empty community ID",
			args: args{
				ctx:          ctx,
				communityCID: nil,
				memberIDs:    userIIDs,
			},
			wantErr: false,
		},
		{
			name: "Happy case - empty isReviewed",
			args: args{
				ctx:          ctx,
				communityCID: &communityID,
				memberIDs:    userIIDs,
			},
			wantErr: false,
		},
		{
			name: "Happy case - empty params",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case - failed to get message flags",
			args: args{
				ctx:          ctx,
				communityCID: &communityID,
				memberIDs:    userIIDs,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case - failed to get message flags" {
				fakeGetStream.MockListFlaggedMessagesFn = func(ctx context.Context, input *stream.QueryOption) (*stream.QueryMessageFlagsResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.ListFlaggedMessages(tt.args.ctx, tt.args.communityCID, tt.args.memberIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListFlaggedMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesCommunitiesImpl.ListFlaggedMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_DeleteCommunityMessage(t *testing.T) {
	fakeGetStream := getStreamMock.NewGetStreamServiceMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()
	communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)
	type args struct {
		ctx       context.Context
		messageID string
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
				ctx:       context.Background(),
				messageID: uuid.New().String(),
			},
			want: true,
		},
		{
			name: "Sad case - failed to delete message",
			args: args{
				ctx:       context.Background(),
				messageID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case - failed to delete message" {
				fakeGetStream.MockDeleteMessageFn = func(ctx context.Context, messageID string) (*stream.Response, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := communities.DeleteCommunityMessage(tt.args.ctx, tt.args.messageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.DeleteCommunityMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.DeleteCommunityMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_VerifyWebhook(t *testing.T) {
	ctx := context.Background()
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
			name: "Happy Case - Successfully verify a webhook request",
			args: args{
				ctx:       ctx,
				body:      []byte("random word"),
				signature: "SIGNATURE",
			},
			want: true,
		},
		{
			name: "Sad Case - Fail to verify webhook",
			args: args{
				ctx: ctx,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad Case - Fail to verify webhook" {
				fakeGetStream.MockValidateGetStreamRequestFn = func(ctx context.Context, body []byte, signature string) bool {
					return false
				}
			}

			if got := communities.ValidateGetStreamRequest(tt.args.ctx, tt.args.body, tt.args.signature); got != tt.want {
				t.Errorf("UseCasesCommunitiesImpl.ValidateGetStreamRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_ProcessGetstreamEvents(t *testing.T) {
	type args struct {
		ctx   context.Context
		event *dto.GetStreamEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully process getstream events",
			args: args{
				ctx: context.Background(),
				event: &dto.GetStreamEvent{
					CID:          "messaging:1234",
					Type:         "message.new",
					Message:      &stream.Message{},
					Reaction:     &stream.Reaction{},
					Channel:      &stream.Channel{},
					Member:       &stream.ChannelMember{},
					Members:      []*stream.ChannelMember{},
					User:         &stream.User{},
					UserID:       "",
					OwnUser:      &stream.User{},
					WatcherCount: 0,
					ExtraData:    map[string]interface{}{},
					ChannelID:    "12345",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to process getstream event",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream, fakeExtension, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "Sad Case - Fail to process getstream event" {
				fakePubsub.MockNotifyGetStreamEventFn = func(ctx context.Context, event *dto.GetStreamEvent) error {
					return fmt.Errorf("failed to process getstream event")
				}
			}

			if err := communities.ProcessGetstreamEvents(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ProcessGetstreamEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
