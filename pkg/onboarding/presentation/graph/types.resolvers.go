package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
)

func (r *facilityResolver) FacilityID(ctx context.Context, obj *domain.Facility) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

// Facility returns generated.FacilityResolver implementation.
func (r *Resolver) Facility() generated.FacilityResolver { return &facilityResolver{r} }

type facilityResolver struct{ *Resolver }
