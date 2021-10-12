package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
)

func (r *mutationResolver) CreateFacility(ctx context.Context, input dto.FacilityInput) (*domain.Facility, error) {
	return r.interactor.FacilityUsecase.CreateFacility(ctx, input)
}

func (r *mutationResolver) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return r.interactor.FacilityUsecase.DeleteFacility(ctx, id)
}

func (r *queryResolver) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return r.interactor.FacilityUsecase.FetchFacilities(ctx)
}

func (r *queryResolver) RetrieveFacility(ctx context.Context, id int) (*domain.Facility, error) {
	intID := int64(id)
	return r.interactor.FacilityUsecase.RetrieveFacility(ctx, &intID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
