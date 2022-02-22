package communities

import (
	"context"
	"fmt"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
)

const (
	inviteMessage = "%v has invited you to join this community"
)

// ICreateCommunity is an interface that is used to create communities
type ICreateCommunity interface {
	CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error)
}

// IListUsers is an interface that is used to list getstream users
type IListUsers interface {
	ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error)
}

// IInviteMembers interface holds methods that are used to send member invites
type IInviteMembers interface {
	InviteMembers(ctx context.Context, communityID string, userIDS []string) (bool, error)
}

// IListChannels is an interface that is used to list getstream channels
type IListChannels interface {
	ListGetStreamChannels(ctx context.Context, input *domain.QueryOption) (*domain.QueryChannelsResponse, error)
}

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
	ICreateCommunity
	IInviteMembers
	IListUsers
	IListChannels
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
	GetstreamService streamService.ServiceGetStream
	Create           infrastructure.Create
	ExternalExt      extension.ExternalMethodsExtension
	Query            infrastructure.Query
}

// NewUseCaseCommunitiesImpl initializes a new communities service
func NewUseCaseCommunitiesImpl(
	getstream streamService.ServiceGetStream,
	ext extension.ExternalMethodsExtension,
	create infrastructure.Create,
	query infrastructure.Query,
) *UseCasesCommunitiesImpl {
	return &UseCasesCommunitiesImpl{
		GetstreamService: getstream,
		Create:           create,
		ExternalExt:      ext,
		Query:            query,
	}
}

// ListGetStreamUsers returns list of users that match QueryOption that's passed as the input
func (us *UseCasesCommunitiesImpl) ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error) {
	userResponse := []*domain.GetStreamUser{}

	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"role": "user",
		},
		UserID:       input.UserID,
		Limit:        input.Limit,
		Offset:       input.Offset,
		MessageLimit: input.MessageLimit,
		MemberLimit:  input.MemberLimit,
	}

	getStreamUserResponse, err := us.GetstreamService.ListGetStreamUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to return getstream users :%v", err)
	}

	for _, user := range getStreamUserResponse.Users {
		Users := domain.GetStreamUser{
			ID:        user.ID,
			Name:      user.Name,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		userResponse = append(userResponse, &Users)
	}

	return &domain.QueryUsersResponse{
		Users: userResponse,
	}, nil
}

// CreateCommunity creates channel with the GetStream chat service
func (us *UseCasesCommunitiesImpl) CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
	channelResponse, err := us.Create.CreateChannel(ctx, &input)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	data := map[string]interface{}{
		"minimumAge": channelResponse.AgeRange.LowerBound,
		"maximumAge": channelResponse.AgeRange.UpperBound,
		"gender":     channelResponse.Gender,
		"clientType": channelResponse.ClientType,
		"inviteOnly": channelResponse.InviteOnly,
		"name":       channelResponse.Name,
	}

	_, err = us.GetstreamService.CreateChannel(ctx, "messaging", channelResponse.ID, loggedInUserID, data)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create channel: %v", err)
	}

	return &domain.Community{
		ID:          channelResponse.ID,
		Name:        channelResponse.Name,
		Description: channelResponse.Description,
		AgeRange: &domain.AgeRange{
			LowerBound: channelResponse.AgeRange.LowerBound,
			UpperBound: channelResponse.AgeRange.UpperBound,
		},
		Gender:     channelResponse.Gender,
		ClientType: channelResponse.ClientType,
		InviteOnly: channelResponse.InviteOnly,
	}, nil
}

// InviteMembers invites specified members to a community
func (us *UseCasesCommunitiesImpl) InviteMembers(ctx context.Context, communityID string, userIDS []string) (bool, error) {
	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get staff profile")
	}

	// TODO: Fetch the channel to get the channel name and pass it as part of the message
	message := &stream.Message{
		ID:   uuid.New().String(),
		Text: fmt.Sprintf(inviteMessage, staffProfile.User.Name),
		User: &stream.User{
			ID: *staffProfile.ID,
		},
	}

	_, err = us.GetstreamService.InviteMembers(ctx, userIDS, communityID, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to invite members to a community: %v", err)
	}

	return true, nil
}

// ListGetStreamChannels returns list of channels that match QueryOption that's passed as the input
func (us *UseCasesCommunitiesImpl) ListGetStreamChannels(ctx context.Context, input *domain.QueryOption) (*domain.QueryChannelsResponse, error) {
	channelResponse := []*domain.GetStreamChannel{}

	var query stream.QueryOption
	// set default offset and limit if none is passed
	if input == nil {
		query = stream.QueryOption{
			Filter: map[string]interface{}{},
			Limit:  10,
			Offset: 0,
		}
	} else {
		query = stream.QueryOption{
			Filter: input.Filter,
			Limit:  input.Limit,
			Offset: input.Offset,
		}
	}

	getStreamChannelResponse, err := us.GetstreamService.ListGetStreamChannels(ctx, &query)
	if err != nil {
		return nil, fmt.Errorf("failed to return getstream channels :%v", err)
	}

	for _, channel := range getStreamChannelResponse.Channels {

		createdBy := &domain.GetStreamUser{
			ID:        channel.CreatedBy.ID,
			Name:      channel.CreatedBy.Name,
			Role:      channel.CreatedBy.Role,
			CreatedAt: channel.CreatedBy.CreatedAt,
			UpdatedAt: channel.CreatedBy.UpdatedAt,
		}

		channelResponse = append(channelResponse, &domain.GetStreamChannel{
			ID:            channel.ID,
			Type:          channel.Type,
			CID:           channel.CID,
			Team:          channel.Team,
			CreatedBy:     createdBy,
			Disabled:      channel.Disabled,
			Frozen:        channel.Frozen,
			MemberCount:   channel.MemberCount,
			CreatedAt:     channel.CreatedAt,
			UpdatedAt:     channel.UpdatedAt,
			LastMessageAt: channel.UpdatedAt,
		})
	}

	return &domain.QueryChannelsResponse{
		Channels: channelResponse,
	}, nil
}
