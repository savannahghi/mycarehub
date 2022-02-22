package communities

import (
	"context"
	"fmt"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
)

// CreateCommunity is an interface that is used to create communities
type CreateCommunity interface {
	CreateCommunity(ctx context.Context, input dto.CommunityInput) (*domain.Community, error)
}

// ListUsers is an interface that is used to list getstream users
type ListUsers interface {
	ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error)
}

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
	CreateCommunity
	ListUsers
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
	GetstreamService streamService.ServiceGetStream
	Create           infrastructure.Create
	ExternalExt      extension.ExternalMethodsExtension
}

// NewUseCaseCommunities initializes a new communities service
func NewUseCaseCommunities(
	getstream streamService.ServiceGetStream,
	create infrastructure.Create,
	ext extension.ExternalMethodsExtension,
) *UseCasesCommunitiesImpl {
	return &UseCasesCommunitiesImpl{
		GetstreamService: getstream,
		Create:           create,
		ExternalExt:      ext,
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
