package communities

import (
	"context"
	"fmt"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
)

// ListUsers is an interface that is used to list getstream users
type ListUsers interface {
	ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error)
}

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
	ListUsers
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
	GetStream getstream.ServiceGetStream
}

// NewUseCaseCommunitiesImpl initializes a new communities service
func NewUseCaseCommunitiesImpl(
	getstream getstream.ServiceGetStream,
) *UseCasesCommunitiesImpl {
	return &UseCasesCommunitiesImpl{
		GetStream: getstream,
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

	getStreamUserResponse, err := us.GetStream.ListGetStreamUsers(ctx, query)
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
