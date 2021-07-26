package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/profileutils"
)

func (r *entityResolver) FindMicroserviceByID(ctx context.Context, id string) (*domain.Microservice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindPageInfoByHasNextPage(ctx context.Context, hasNextPage bool) (*firebasetools.PageInfo, error) {
	return nil, nil
}

func (r *entityResolver) FindUserProfileByID(ctx context.Context, id string) (*profileutils.UserProfile, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	return r.interactor.Onboarding.GetProfileByID(ctx, &id)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
