package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

func (r *entityResolver) FindUserProfileByID(ctx context.Context, id string) (*profile.UserProfile, error) {
	r.CheckDependencies()
	r.CheckUserTokenInContext(ctx)
	return r.profileService.GetProfileByID(ctx, id)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
