package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	domain1 "github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
)

func (r *pINResolver) User(ctx context.Context, obj *domain.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pINResolver) Pin(ctx context.Context, obj *domain.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pINResolver) ConfirmedPin(ctx context.Context, obj *domain.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pINResolver) Flavour(ctx context.Context, obj *domain.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) FailedLoginCount(ctx context.Context, obj *domain1.User) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}

// PIN returns generated.PINResolver implementation.
func (r *Resolver) PIN() generated.PINResolver { return &pINResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type pINResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
