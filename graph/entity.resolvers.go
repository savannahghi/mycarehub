package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

func (r *entityResolver) FindCoverByPayerName(ctx context.Context, payerName string) (*profile.Cover, error) {
	// todo(dexter) implement this
	return nil, nil
}

func (r *entityResolver) FindPageInfoByHasNextPage(ctx context.Context, hasNextPage bool) (*base.PageInfo, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *entityResolver) FindUserProfileByUID(ctx context.Context, uid string) (*profile.UserProfile, error) {
	return r.profileService.GetProfile(ctx, uid)
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
