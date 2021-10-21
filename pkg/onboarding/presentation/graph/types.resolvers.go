package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
	domain1 "github.com/savannahghi/onboarding/pkg/onboarding/domain"
)

func (r *pINResolver) User(ctx context.Context, obj *domain1.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pINResolver) Pin(ctx context.Context, obj *domain1.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pINResolver) ConfirmedPin(ctx context.Context, obj *domain1.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *pINResolver) Flavour(ctx context.Context, obj *domain1.PIN) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// PIN returns generated.PINResolver implementation.
func (r *Resolver) PIN() generated.PINResolver { return &pINResolver{r} }

type pINResolver struct{ *Resolver }
