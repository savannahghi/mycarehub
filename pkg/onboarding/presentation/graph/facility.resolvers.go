package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
)

func (r *mutationResolver) CreateFacility(ctx context.Context, input dto.FacilityInput) (*domain.Facility, error) {
	log.Printf("\n\n\n\n\n\n\n\n\n\nthis is where we are %v", input)
	return r.interactor.FacilityUsecase.CreateFacility(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
