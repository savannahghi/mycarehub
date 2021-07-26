package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
)

func (r *verifiedIdentifierResolver) Timestamp(ctx context.Context, obj *profileutils.VerifiedIdentifier) (*scalarutils.Date, error) {
	return &scalarutils.Date{
		Year:  obj.Timestamp.Year(),
		Day:   obj.Timestamp.Day(),
		Month: int(obj.Timestamp.Month()),
	}, nil
}

// VerifiedIdentifier returns generated.VerifiedIdentifierResolver implementation.
func (r *Resolver) VerifiedIdentifier() generated.VerifiedIdentifierResolver {
	return &verifiedIdentifierResolver{r}
}

type verifiedIdentifierResolver struct{ *Resolver }
