package communities

import (
	"context"
	"fmt"
	"strings"
	"time"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/savannahghi/enumutils"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
)

const (
	inviteMessage             = "%v invited %v to this community"
	promoteToModeratorMessage = "You have been promoted to be a moderator in %v community"
)

// ICreateCommunity is an interface that is used to create communities
type ICreateCommunity interface {
	CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error)
}

// IListUsers is an interface that is used to list getstream users
type IListUsers interface {
	ListCommunityMembers(ctx context.Context, communityID string, input *stream.QueryOption) ([]*domain.CommunityMember, error)
	ListMembers(ctx context.Context, input *stream.QueryOption) ([]*domain.Member, error)
}

// IInviteMembers interface holds methods that are used to send member invites
type IInviteMembers interface {
	InviteMembers(ctx context.Context, communityID string, memberIDs []string) (bool, error)
}

// IListCommunities is an interface that is used to list getstream channels
type IListCommunities interface {
	ListCommunities(ctx context.Context, input *stream.QueryOption) ([]*domain.Community, error)
	ListPendingInvites(ctx context.Context, memberID string, input *stream.QueryOption) ([]*domain.Community, error)
}

// IDeleteCommunities is an interface that is used to delete channels
type IDeleteCommunities interface {
	DeleteCommunities(ctx context.Context, communityIDs []string, hardDelete bool) (bool, error)
}

// ICommunityInvites is an interface that is used to manage community invites
type ICommunityInvites interface {
	AcceptInvite(ctx context.Context, memberID string, channelID string) (bool, error)
	RejectInvite(ctx context.Context, memberID string, channelID string) (bool, error)
}

// IManageMembers is an interface that is used to manage community members
type IManageMembers interface {
	RemoveMembersFromCommunity(ctx context.Context, channelID string, memberIDs []string) (bool, error)
	AddMembersToCommunity(ctx context.Context, memberIDs []string, communityID string) (bool, error)
}

// IModeration interface contains all the moderation functions
type IModeration interface {
	AddModeratorsWithMessage(ctx context.Context, memberIDs []string, communityID string) (bool, error)
	DemoteModerators(ctx context.Context, communityID string, memberIDs []string) (bool, error)
	BanUser(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error)
	UnBanUser(ctx context.Context, targetID string, communityID string) (bool, error)
	ListCommunityBannedMembers(ctx context.Context, communityID string) ([]*domain.Member, error)
	ListFlaggedMessages(ctx context.Context, communityCID *string, memberIDs []*string) ([]*domain.MessageFlag, error)
}

// IRecommendations interface contains all the recommendation functions
type IRecommendations interface {
	RecommendedCommunities(ctx context.Context, clientID string, limit int) ([]*domain.Community, error)
}

// IMessage interface is used to contain all the message related methods
type IMessage interface {
	DeleteCommunityMessage(ctx context.Context, messageID string) (bool, error)
}

// IValidateRequest specifies a method that is used to verify a webhook request
type IValidateRequest interface {
	ValidateGetStreamRequest(ctx context.Context, body []byte, signature string) bool
}

// IEvents specifies the methods that revolve around getstream events. Events allow the client to stay up to date with changes to the chat
type IEvents interface {
	ProcessGetstreamEvents(ctx context.Context, data *dto.GetStreamEvent) error
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
	IModeration
	IRecommendations
	IMessage
	IValidateRequest
	IEvents
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
	GetstreamService streamService.ServiceGetStream
	Create           infrastructure.Create
	ExternalExt      extension.ExternalMethodsExtension
	Query            infrastructure.Query
	Pubsub           pubsubmessaging.ServicePubsub
	Notification     notification.UseCaseNotification
}

// NewUseCaseCommunitiesImpl initializes a new communities service
func NewUseCaseCommunitiesImpl(
	getstream streamService.ServiceGetStream,
	ext extension.ExternalMethodsExtension,
	create infrastructure.Create,
	query infrastructure.Query,
	pubsub pubsubmessaging.ServicePubsub,
	notification notification.UseCaseNotification,
) *UseCasesCommunitiesImpl {
	return &UseCasesCommunitiesImpl{
		GetstreamService: getstream,
		Create:           create,
		ExternalExt:      ext,
		Query:            query,
		Pubsub:           pubsub,
		Notification:     notification,
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
			Filter:       utils.FormatFilterParamsHelper(input.Filter),
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
		var metadata domain.MemberMetadata
		err := mapstructure.Decode(user.ExtraData, &metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %v", err)
		}

		Users := domain.Member{
			ID:            user.ID,
			Name:          user.Name,
			Role:          user.Role,
			UserID:        metadata.UserID,
			UserType:      metadata.UserType,
			Username:      metadata.Username,
			Gender:        enumutils.Gender(metadata.Gender),
			AgeUpperBound: metadata.AgeUpperBound,
			AgeLowerBound: metadata.AgeLowerBound,
			ClientTypes:   metadata.ClientTypes,
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
		return nil, err
	}

	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	channelMetadata := domain.CommunityMetadata{
		MinimumAge:  channelResponse.AgeRange.LowerBound,
		MaximumAge:  channelResponse.AgeRange.UpperBound,
		Gender:      channelResponse.Gender,
		ClientType:  channelResponse.ClientType,
		InviteOnly:  channelResponse.InviteOnly,
		Name:        channelResponse.Name,
		Description: channelResponse.Description,
	}

	var extraData map[string]interface{}
	err = mapstructure.Decode(channelMetadata, &extraData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	staff, err := us.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetLoggedInUserUIDErr(err)
	}

	channel, err := us.GetstreamService.CreateChannel(ctx, "messaging", channelResponse.ID, *staff.ID, &stream.ChannelRequest{ExtraData: extraData})
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
	community, err := us.Query.GetCommunityByID(ctx, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to retrieve community with provided id: %w", err)
	}

	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get staff profile")
	}

	var invitees []domain.User
	for _, memberID := range memberIDs {
		user, err := us.GetstreamService.GetStreamUser(ctx, memberID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to retrieve getstream user: %w", err)
		}

		var metadata domain.MemberMetadata
		err = mapstructure.Decode(user.ExtraData, &metadata)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		userProfile, err := us.Query.GetUserProfileByUserID(ctx, metadata.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to retrieve getstream user's user profile: %w", err)
		}

		invitees = append(invitees, *userProfile)
	}

	var message *stream.Message
	for _, invitedUser := range invitees {
		message = &stream.Message{
			ID:   uuid.New().String(),
			Text: fmt.Sprintf(inviteMessage, staffProfile.User.Name, invitedUser.Username),
			User: &stream.User{
				ID: *staffProfile.ID,
			},
		}
	}

	_, err = us.GetstreamService.InviteMembers(ctx, memberIDs, communityID, message)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to invite members to a community: %v", err)
	}

	notificationInput := notification.ClientNotificationInput{
		Community: community,
		Inviter:   staffProfile.User,
	}
	notificationMessage := notification.ComposeClientNotification(
		enums.NotificationTypeCommunities,
		notificationInput,
	)

	for _, invitee := range invitees {
		err = us.Notification.NotifyUser(ctx, &invitee, notificationMessage)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}
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

		var metaData domain.CommunityMetadata
		err := mapstructure.Decode(channel.ExtraData, &metaData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %v", err)
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
			Description: metaData.Description,
			Gender:      metaData.Gender,
			ClientType:  metaData.ClientType,
			InviteOnly:  metaData.InviteOnly,
			Name:        metaData.Name,
			AgeRange: &domain.AgeRange{
				LowerBound: metaData.MinimumAge,
				UpperBound: metaData.MaximumAge,
			},
		})
	}

	return channelResponse, nil
}

// ListCommunityMembers retrieves the members of a community
func (us *UseCasesCommunitiesImpl) ListCommunityMembers(ctx context.Context, communityID string, input *stream.QueryOption) ([]*domain.CommunityMember, error) {
	members := []*domain.CommunityMember{}
	var query *stream.QueryOption
	var sorters []*stream.SortOption

	if input == nil {
		sorters = []*stream.SortOption{{Field: "name", Direction: 1}}
		query = &stream.QueryOption{
			Filter: map[string]interface{}{
				"banned": false,
				"joined": true,
			},
		}
	} else {
		sorters = input.Sort
		query = &stream.QueryOption{
			Filter:       utils.FormatFilterParamsHelper(input.Filter),
			UserID:       input.UserID,
			Limit:        input.Limit,
			Offset:       input.Offset,
			MessageLimit: input.MessageLimit,
			MemberLimit:  input.MemberLimit,
		}
	}

	channelMembersResponse, err := us.GetstreamService.QueryChannelMembers(ctx, communityID, query, sorters...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve members of a community: %w", err)
	}

	for _, member := range channelMembersResponse.Members {
		var metadata domain.MemberMetadata
		err := mapstructure.Decode(member.User.ExtraData, &metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %w", err)
		}

		user := domain.Member{
			ID:        member.User.ID,
			Name:      member.User.Name,
			Role:      member.User.Role,
			UserID:    metadata.UserID,
			Username:  metadata.Username,
			ExtraData: member.User.ExtraData,
		}

		commMem := &domain.CommunityMember{
			UserID:           metadata.UserID,
			User:             user,
			Role:             member.Role,
			IsModerator:      member.IsModerator,
			UserType:         metadata.UserType,
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

// RejectInvite rejects an invitation into a community
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

// AddModeratorsWithMessage adds moderators with given IDs to the channel and produces a message.
func (us *UseCasesCommunitiesImpl) AddModeratorsWithMessage(ctx context.Context, memberIDs []string, communityID string) (bool, error) {
	community, err := us.Query.GetCommunityByID(ctx, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get staff profile %w", err)
	}

	for _, v := range memberIDs {
		message := &stream.Message{
			ID:   uuid.New().String(),
			Text: fmt.Sprintf(promoteToModeratorMessage, community.Name),
			User: &stream.User{
				ID: v,
			},
		}

		_, err = us.GetstreamService.AddModeratorsWithMessage(ctx, []string{v}, communityID, message)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to add moderator(s)")
		}

		userProfile, err := us.GetUserProfileByMemberID(ctx, v)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		notificationInput := notification.ClientNotificationInput{
			Community: community,
			Promoter:  staffProfile.User,
		}

		clientNotification := notification.ComposeClientNotification(
			enums.NotificationTypePromoteToModerator,
			notificationInput,
		)

		err = us.Notification.NotifyUser(ctx, userProfile, clientNotification)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

	}
	return true, nil
}

// DemoteModerators demotes moderators in a community
// memberID can either be a Client or a Staff ID
func (us *UseCasesCommunitiesImpl) DemoteModerators(ctx context.Context, communityID string, memberIDs []string) (bool, error) {
	_, err := us.GetstreamService.DemoteModerators(ctx, communityID, memberIDs)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to demote moderators in a community: %v", err)
	}

	loggedInUserID, err := us.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staffProfile, err := us.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		return false, fmt.Errorf("failed to get staff profile %w", err)
	}

	community, err := us.Query.GetCommunityByID(ctx, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to retrieve community with provided id: %w", err)
	}

	var communityMembers []domain.User
	for _, memberID := range memberIDs {
		userProfile, err := us.GetUserProfileByMemberID(ctx, memberID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		communityMembers = append(communityMembers, *userProfile)
	}

	notificationInput := notification.ClientNotificationInput{
		Community: community,
		Demoter:   staffProfile.User,
	}

	clientNotificationMessage := notification.ComposeClientNotification(
		enums.NotificationTypeDemoteModerator,
		notificationInput,
	)

	for _, member := range communityMembers {
		err = us.Notification.NotifyUser(ctx, &member, clientNotificationMessage)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}

	return true, nil
}

// GetUserProfileByMemberID retrieves the user profile of a community member
func (us *UseCasesCommunitiesImpl) GetUserProfileByMemberID(ctx context.Context, memberID string) (*domain.User, error) {
	streamUser, err := us.GetstreamService.GetStreamUser(ctx, memberID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	var memberMetadata domain.MemberMetadata
	err = mapstructure.Decode(streamUser.ExtraData, &memberMetadata)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	userProfile, err := us.Query.GetUserProfileByUserID(ctx, memberMetadata.UserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}
	return userProfile, nil
}

// ListPendingInvites lists all communities that a user has been invited into
func (us *UseCasesCommunitiesImpl) ListPendingInvites(ctx context.Context, memberID string, input *stream.QueryOption) ([]*domain.Community, error) {
	if memberID == "" {
		return nil, fmt.Errorf("memberID cannot be empty")
	}
	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"invite": "pending",
		},
		UserID: memberID,
		Limit:  input.Limit,
		Offset: input.Offset,
	}

	streamChannelsResponse, err := us.GetstreamService.ListGetStreamChannels(ctx, query)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("an error occurred: %v", err)
	}

	var communities []*domain.Community
	for _, channel := range streamChannelsResponse.Channels {
		var metaData domain.CommunityMetadata

		err := mapstructure.Decode(channel.ExtraData, &metaData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %v", err)
		}

		community := &domain.Community{
			ID:          channel.ID,
			Name:        metaData.Name,
			Description: metaData.Description,
			MemberCount: channel.MemberCount,
			CreatedAt:   channel.CreatedAt,
			Gender:      metaData.Gender,
			ClientType:  metaData.ClientType,
			InviteOnly:  metaData.InviteOnly,
			AgeRange: &domain.AgeRange{
				LowerBound: metaData.MinimumAge,
				UpperBound: metaData.MaximumAge,
			},
		}

		communities = append(communities, community)
	}

	return communities, nil
}

// UnBanUser unbans a user from the specified channel
func (us *UseCasesCommunitiesImpl) UnBanUser(ctx context.Context, targetID string, communityID string) (bool, error) {
	return us.GetstreamService.UnBanUser(ctx, targetID, communityID)
}

// RecommendedCommunities returns a list of communities that have been recommended to the user
//the recommendations are based on the channel metadata and the client data; client type, age range, gender
// e.g. if a channel's age range value is 25-30, and the client's age is between this range,
// this channel will be recommended to the user if they have not joined
func (us *UseCasesCommunitiesImpl) RecommendedCommunities(ctx context.Context, clientID string, limit int) ([]*domain.Community, error) {
	if clientID == "" {
		return nil, fmt.Errorf("memberID cannot be empty")
	}

	joinedChannelsMap := make(map[string]bool)

	clientProfile, err := us.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	userID := clientProfile.UserID

	clientUserProfile, err := us.Query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}
	var clientTypes []string
	for _, k := range clientProfile.ClientTypes {
		clientTypes = append(clientTypes, k.String())
	}

	joinedChannels, err := us.GetstreamService.ListGetStreamChannels(ctx, &stream.QueryOption{
		Filter: map[string]interface{}{
			"members": map[string]interface{}{"$in": []string{clientID}},
		},
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	for _, channel := range joinedChannels.Channels {
		joinedChannelsMap[channel.ID] = true
	}

	var dob = time.Now()
	if clientUserProfile.DateOfBirth != nil {
		dob = *clientUserProfile.DateOfBirth
	}

	age := utils.CalculateAge(dob)

	clientGender := clientUserProfile.Gender.String()
	clientGender = strings.ToLower(clientGender)

	query := &stream.QueryOption{
		Filter: map[string]interface{}{
			"clientType": clientTypes,
			"gender":     []string{clientGender},
			"minimumAge": map[string]interface{}{"$lte": age},
			"maximumAge": map[string]interface{}{"$gte": age},
			"inviteOnly": false,
		},
		Limit: limit,
	}

	streamChannelsResponse, err := us.GetstreamService.ListGetStreamChannels(ctx, query)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("an error occurred: %v", err)
	}
	var communities []*domain.Community
	for _, channel := range streamChannelsResponse.Channels {
		if _, ok := joinedChannelsMap[channel.ID]; ok {
			continue
		}

		var metaData domain.CommunityMetadata

		err := mapstructure.Decode(channel.ExtraData, &metaData)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %v", err)
		}

		// TODO:check if user has joined the community
		// TODO: nullable filters for client type and age

		community := &domain.Community{
			ID:          channel.ID,
			Name:        metaData.Name,
			Description: metaData.Description,
			MemberCount: channel.MemberCount,
			CreatedAt:   channel.CreatedAt,
			ClientType:  metaData.ClientType,
			Gender:      metaData.Gender,
			AgeRange: &domain.AgeRange{
				LowerBound: metaData.MinimumAge,
				UpperBound: metaData.MaximumAge,
			},
			InviteOnly: metaData.InviteOnly,
		}

		communities = append(communities, community)
	}

	return communities, nil
}

// ListCommunityBannedMembers is used to list members banned from a channel.
func (us *UseCasesCommunitiesImpl) ListCommunityBannedMembers(ctx context.Context, communityID string) ([]*domain.Member, error) {
	if communityID == "" {
		return nil, fmt.Errorf("communityID cannot be empty")
	}

	response, err := us.GetstreamService.ListCommunityBannedMembers(ctx, communityID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var members []*domain.Member
	for _, v := range response.Bans {
		var metadata domain.MemberMetadata
		err := mapstructure.Decode(v.User.ExtraData, &metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %v", err)
		}

		member := &domain.Member{
			ID:       v.User.ID,
			UserID:   metadata.UserID,
			Name:     v.User.Name,
			Role:     v.User.Role,
			UserType: metadata.UserType,
		}

		members = append(members, member)
	}
	return members, nil
}

// BanUser is used to ban user from a specified channel
func (us *UseCasesCommunitiesImpl) BanUser(ctx context.Context, targetMemberID string, bannedBy string, communityID string) (bool, error) {
	if targetMemberID == "" {
		return false, fmt.Errorf("target member ID cannot be empty")
	}
	return us.GetstreamService.BanUser(ctx, targetMemberID, bannedBy, communityID)
}

// ListFlaggedMessages returns a list of flaged messages
// passing in a communityID only will return all the flagged messages in that community
// passing in memberIDs only will return all the flagged messages for members in all channels
// passing in both a communityID and memberIDs will return all the flagged messages for members in that channel
func (us *UseCasesCommunitiesImpl) ListFlaggedMessages(ctx context.Context, communityCID *string, memberIDs []*string) ([]*domain.MessageFlag, error) {
	var (
		newMemberIDs    []string
		newCommunityCID string
	)
	messageFilter := map[string]interface{}{}

	for _, memberID := range memberIDs {
		if memberID != nil {
			newMemberIDs = append(newMemberIDs, *memberID)
		}
	}
	if communityCID != nil {
		newCommunityCID = *communityCID
		messageFilter["channel_cid"] = newCommunityCID
	}

	if len(newMemberIDs) > 0 {
		messageFilter["user_id"] = map[string][]string{"$in": newMemberIDs}
	}

	query := &stream.QueryOption{
		Filter: messageFilter,
	}
	messageFlagsResponse, err := us.GetstreamService.ListFlaggedMessages(ctx, query)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to get flagged messages: %v", err)
	}

	messageFlags := []*domain.MessageFlag{}
	for _, messageFlag := range messageFlagsResponse.Flags {

		newMessageFlag := &domain.MessageFlag{}
		err := mapstructure.Decode(messageFlag, newMessageFlag)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %v", err)
		}

		if messageFlag.Message.DeletedAt != nil {
			continue
		}

		messageFlags = append(messageFlags, newMessageFlag)
	}
	return messageFlags, nil
}

// DeleteCommunityMessage is used to delete a message from a channel
func (us *UseCasesCommunitiesImpl) DeleteCommunityMessage(ctx context.Context, messageID string) (bool, error) {
	_, err := us.GetstreamService.DeleteMessage(ctx, messageID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to delete message: %v", err)
	}
	return true, nil
}

// ValidateGetStreamRequest verifies that a request is coming from Stream and has not been tampered by a 3rd party
// The requests are made on an unauthenticated endpoint hence it can easily be tampered with.
func (us *UseCasesCommunitiesImpl) ValidateGetStreamRequest(ctx context.Context, body []byte, signature string) bool {
	return us.GetstreamService.ValidateGetStreamRequest(ctx, body, signature)
}

// ProcessGetstreamEvents published the event payload to a pubsub topic where it will be processed. This makes it
// easy to leverage asynchronicity
func (us *UseCasesCommunitiesImpl) ProcessGetstreamEvents(ctx context.Context, event *dto.GetStreamEvent) error {
	return us.Pubsub.NotifyGetStreamEvent(ctx, event)
}
