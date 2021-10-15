package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
)

func (r *staffProfileResolver) Roles(ctx context.Context, obj *domain.StaffProfile) ([]string, error) {
	panic(fmt.Errorf("not implemented"))
}

// StaffProfile returns generated.StaffProfileResolver implementation.
func (r *Resolver) StaffProfile() generated.StaffProfileResolver { return &staffProfileResolver{r} }

type staffProfileResolver struct{ *Resolver }
