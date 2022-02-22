package communities_test

import (
	"context"
	"fmt"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
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

			c := communities.NewUseCaseCommunities(fakeGetStream, fakeDB, fakeExtension)

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
