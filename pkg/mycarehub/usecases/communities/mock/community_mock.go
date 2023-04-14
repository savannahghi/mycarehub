package mock

import (
	"context"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CommunityUsecaseMock is used to mock community methods
type CommunityUsecaseMock struct {
	MockCreateCommunityFn func(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error)
	MockListCommunitiesFn func(ctx context.Context) ([]string, error)
	MockSearchUsersFn     func(ctx context.Context, limit *int, searchTerm string) (*domain.MatrixUserSearchResult, error)
}

// NewCommunityUsecaseMock instantiates all the community usecase mock methods
func NewCommunityUsecaseMock() *CommunityUsecaseMock {
	return &CommunityUsecaseMock{
		MockCreateCommunityFn: func(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error) {
			return &domain.Community{
				ID:          uuid.NewString(),
				RoomID:      gofakeit.BeerName(),
				Name:        gofakeit.BeerName(),
				Description: gofakeit.BeerAlcohol(),
				AgeRange: &domain.AgeRange{
					LowerBound: 10,
					UpperBound: 20,
				},
				OrganisationID: uuid.NewString(),
				ProgramID:      uuid.NewString(),
				FacilityID:     uuid.NewString(),
			}, nil
		},
		MockListCommunitiesFn: func(ctx context.Context) ([]string, error) {
			return []string{"test"}, nil
		},
		MockSearchUsersFn: func(ctx context.Context, limit *int, searchTerm string) (*domain.MatrixUserSearchResult, error) {
			return &domain.MatrixUserSearchResult{
				Limited: false,
				Results: []domain.Result{
					{
						UserID:      gofakeit.UUID(),
						DisplayName: "test",
						AvatarURL:   "mxc://bar.com/foo",
					},
				},
			}, nil
		},
	}
}

// CreateCommunity mocks the implementation of creating a room.
func (c *CommunityUsecaseMock) CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error) {
	return c.MockCreateCommunityFn(ctx, communityInput)
}

// ListCommunities mocks the listing of communities
func (c *CommunityUsecaseMock) ListCommunities(ctx context.Context) ([]string, error) {
	return c.MockListCommunitiesFn(ctx)
}

// SearchUsers mocks the implementation od searching for Matrix users
func (c *CommunityUsecaseMock) SearchUsers(ctx context.Context, limit *int, searchTerm string) (*domain.MatrixUserSearchResult, error) {
	return c.MockSearchUsersFn(ctx, limit, searchTerm)
}
