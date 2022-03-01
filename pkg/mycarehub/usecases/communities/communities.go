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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
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
	ListCommunityMembers(ctx context.Context, communityID string) ([]*domain.CommunityMember, error)
	ListMembers(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error)
}

// IInviteMembers interface holds methods that are used to send member invites
type IInviteMembers interface {
	InviteMembers(ctx context.Context, communityID string, memberIDs []string) (bool, error)
}

// IListCommunities is an interface that is used to list getstream channels
type IListCommunities interface {
	ListCommunities(ctx context.Context, input *stream.QueryOption) ([]*domain.Community, error)
}

// IDeleteCommunities is an interface that is used to delete channels
type IDeleteCommunities interface {
	DeleteCommunities(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error)
}

// ICommunityInvites is an interface that is used to manage community invites
type ICommunityInvites interface {
	AcceptInvite(ctx context.Context, userID string, channelID string) (bool, error)
	RejectInvite(ctx context.Context, userID string, channelID string) (bool, error)
}

// IManageMembers is an interface that is used to manage community members
type IManageMembers interface {
	RemoveMembersFromCommunity(ctx context.Context, channelID string, memberIDs []string) (bool, error)
	AddMembersToCommunity(ctx context.Context, memberIDs []string, communityID string) (bool, error)
}

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
	ICreateCommunity
	IInviteMembers
	IListUsers
	IListCommunities
	IDeleteCommunities
	ICommunityInvites
	IManageMembers
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

// ListMembers returns list of the members that match QueryOption that's passed as the input
func (us *UseCasesCommunitiesImpl) ListMembers(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error) {
	var query *stream.QueryOption

	if input == nil {
		query = &stream.QueryOption{
			Filter: map[string]interface{}{
				"role": "user",
			},
		}
	} else {
		query = &stream.QueryOption{
			Filter:       input.Filter,
			UserID:       input.UserID,
			Limit:        input.Limit,
			Offset:       input.Offset,
			MessageLimit: input.MessageLimit,
			MemberLimit:  input.MemberLimit,
		}
	}

	userResponse := []*domain.Member{}

	getStreamUserResponse, err := us.GetstreamService.ListGetStreamUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to return getstream users :%v", err)
	}

	for _, user := range getStreamUserResponse.Users {
		var userID string
		if val, ok := user.ExtraData["userID"]; ok {
			userID = val.(string)
		}

		Users := domain.Member{
			ID:     user.ID,
			Name:   user.Name,
			Role:   user.Role,
			UserID: userID,
		}
		userResponse = append(userResponse, &Users)
	}

	return userResponse, nil
}

// CreateCommunity creates channel with the GetStream chat service
func (us *UseCasesCommunitiesImpl) CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error) {
	channelResponse, err := us.Create.CreateCommunity(ctx, &input)
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

	staff, err := us.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	channel, err := us.GetstreamService.CreateChannel(ctx, "messaging", channelResponse.ID, *staff.ID, data)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to create channel: %v", err)
	}

	_, err = us.GetstreamService.AddMembersToCommunity(ctx, []string{channel.Channel.CreatedBy.ID}, channel.Channel.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to add member to channel: %v", err)
	}

	return &domain.Community{
		ID:          channelResponse.ID,
		CID:         channelResponse.CID,
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
func (us *UseCasesCommunitiesImpl) InviteMembers(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
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

	_, err = us.GetstreamService.InviteMembers(ctx, memberIDs, communityID, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to invite members to a community: %v", err)
	}

	return true, nil
}

// ListCommunities returns list of the communities that match QueryOption that's passed as the input
func (us *UseCasesCommunitiesImpl) ListCommunities(ctx context.Context, input *stream.QueryOption) ([]*domain.Community, error) {
	channelResponse := []*domain.Community{}

	input.Filter = utils.FormatFilterParamsHelper(input.Filter)

	query := stream.QueryOption{
		Filter: input.Filter,
		Limit:  input.Limit,
		Offset: input.Offset,
		Sort:   input.Sort,
	}

	getStreamChannelResponse, err := us.GetstreamService.ListGetStreamChannels(ctx, &query)
	if err != nil {
		return nil, fmt.Errorf("failed to return getstream channels :%v", err)
	}

	for _, channel := range getStreamChannelResponse.Channels {

		createdBy := &domain.Member{
			ID:   channel.CreatedBy.ID,
			Name: channel.CreatedBy.Name,
			Role: channel.CreatedBy.Role,
		}

		channelResponse = append(channelResponse, &domain.Community{
			ID:          channel.ID,
			CID:         channel.CID,
			CreatedBy:   createdBy,
			Disabled:    channel.Disabled,
			Frozen:      channel.Frozen,
			MemberCount: channel.MemberCount,
			CreatedAt:   channel.CreatedAt,
			UpdatedAt:   channel.UpdatedAt,
		})
	}

	return channelResponse, nil
}

// ListCommunityMembers retrieves the members of a community
func (us *UseCasesCommunitiesImpl) ListCommunityMembers(ctx context.Context, communityID string) ([]*domain.CommunityMember, error) {
	members := []*domain.CommunityMember{}

	channel, err := us.GetstreamService.GetChannel(ctx, communityID)
	if err != nil {
		return nil, err
	}

	for _, member := range channel.Members {
		var userType string
		var userID string

		if val, ok := member.User.ExtraData["userType"]; ok {
			userType = val.(string)
		}

		if val, ok := member.User.ExtraData["userID"]; ok {
			userID = val.(string)
		}

		user := domain.Member{
			ID:     member.User.ID,
			Name:   member.User.Name,
			Role:   member.User.Role,
			UserID: userID,
		}

		commMem := &domain.CommunityMember{
			UserID:           userID,
			User:             user,
			Role:             member.Role,
			IsModerator:      member.IsModerator,
			UserType:         userType,
			Invited:          member.Invited,
			InviteAcceptedAt: member.InviteAcceptedAt,
			InviteRejectedAt: member.InviteRejectedAt,
		}

		members = append(members, commMem)

	}

	return members, nil
}

// DeleteCommunities deletes the specified communities by provided [cid]
func (us *UseCasesCommunitiesImpl) DeleteCommunities(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error) {
	_, err := us.GetstreamService.DeleteChannels(ctx, communityIDs, hardDelete)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to delete channels: %v", err)
	}

	return true, nil
}

// AddMembersToCommunity adds a user (staff/client) to a community
// memberID can either be a Client or a Staff ID
func (us *UseCasesCommunitiesImpl) AddMembersToCommunity(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
	if len(memberIDs) == 0 {
		return false, fmt.Errorf("memberIDs cannot be empty")
	}
	if communityID == "" {
		return false, fmt.Errorf("communityID cannot be empty")
	}

	community, err := us.Query.GetCommunityByID(ctx, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get community by ID: %v", err)
	}

	if community.InviteOnly {
		return false, fmt.Errorf("group is invite only")
	}

	_, err = us.GetstreamService.AddMembersToCommunity(ctx, memberIDs, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to add member(s) to a community: %v", err)
	}

	return true, nil
}

// RejectInvite rejects an invite into a community
func (us *UseCasesCommunitiesImpl) RejectInvite(ctx context.Context, userID string, channelID string) (bool, error) {
	message := &stream.Message{
		ID: uuid.New().String(),
		User: &stream.User{
			ID: userID,
		},
	}

	_, err := us.GetstreamService.RejectInvite(ctx, userID, channelID, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed reject invite into a community")
	}

	return true, nil
}

// AcceptInvite accepts an invite into a community
func (us *UseCasesCommunitiesImpl) AcceptInvite(ctx context.Context, userID string, channelID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("user id cannot be empty")
	}
	if channelID == "" {
		return false, fmt.Errorf("channel id cannot be empty")
	}
	message := &stream.Message{
		ID: uuid.New().String(),
		User: &stream.User{
			ID: userID,
		},
	}

	_, err := us.GetstreamService.AcceptInvite(ctx, userID, channelID, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to accept invite into a community")
	}

	return true, nil
}

// RemoveMembersFromCommunity removes members from a community
// memberID can either be a Client or a Staff ID
func (us *UseCasesCommunitiesImpl) RemoveMembersFromCommunity(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	if len(memberIDs) == 0 {
		return false, fmt.Errorf("user id cannot be empty")
	}
	if communityID == "" {
		return false, fmt.Errorf("channel id cannot be empty")
	}
	message := &stream.Message{
		ID: uuid.New().String(),
		User: &stream.User{
			ID: memberIDs[0],
		},
	}
	_, err := us.GetstreamService.RemoveMembersFromCommunity(ctx, communityID, memberIDs, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to remove members from a community: %v", err)
	}

	return true, nil
}
