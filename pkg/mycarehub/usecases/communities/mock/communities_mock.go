package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CommunityUsecaseMock contains the community usecase mock methods
type CommunityUsecaseMock struct {
	MockListGetStreamUsersFn func(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error)
}

// NewCommunityUsecaseMock initializes a new instance of the Community usecase happy cases
func NewCommunityUsecaseMock() *CommunityUsecaseMock {
	return &CommunityUsecaseMock{
		MockListGetStreamUsersFn: func(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error) {
			return &domain.QueryUsersResponse{
				Users: []*domain.GetStreamUser{
					{
						ID:   uuid.NewString(),
						Role: "user",
					},
				},
			}, nil
		},
	}
}

// ListGetStreamUsers mocks the implementation for listing getstream users
func (c CommunityUsecaseMock) ListGetStreamUsers(ctx context.Context, input *domain.QueryOption) (*domain.QueryUsersResponse, error) {
	return c.MockListGetStreamUsersFn(ctx, input)
}
