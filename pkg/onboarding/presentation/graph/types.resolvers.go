package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
)

func (r *facilityResolver) ID(ctx context.Context, obj *domain.Facility) (string, error) {
	return obj.ID.String(), nil
}

// Facility returns generated.FacilityResolver implementation.
func (r *Resolver) Facility() generated.FacilityResolver { return &facilityResolver{r} }

type facilityResolver struct{ *Resolver }
