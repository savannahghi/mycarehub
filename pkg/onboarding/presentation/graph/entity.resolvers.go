package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
)

func (r *entityResolver) FindMicroserviceByID(ctx context.Context, id string) (*domain.Microservice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindPageInfoByHasNextPage(ctx context.Context, hasNextPage bool) (*base.PageInfo, error) {
	return nil, nil
}

func (r *entityResolver) FindUserProfileByID(ctx context.Context, id string) (*base.UserProfile, error) {
	r.checkPreconditions()
	r.CheckUserTokenInContext(ctx)
	return r.interactor.Onboarding.GetProfileByID(ctx, &id)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
