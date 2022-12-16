package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/graph/generated"
)

// Flavour is the resolver for the Flavour field.
func (r *userResolver) Flavour(ctx context.Context, obj *domain.User) (*feedlib.Flavour, error) {
	panic(fmt.Errorf("not implemented: Flavour - Flavour"))
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
