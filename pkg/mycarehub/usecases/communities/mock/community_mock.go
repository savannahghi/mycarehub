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
	}
}

// CreateCommunity mocks the implementation of creating a room.
func (c *CommunityUsecaseMock) CreateCommunity(ctx context.Context, communityInput *dto.CommunityInput) (*domain.Community, error) {
	return c.MockCreateCommunityFn(ctx, communityInput)
}
