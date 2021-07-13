package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain/model"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
)

func (r *navActionResolver) Icon(ctx context.Context, obj *base.NavAction) (*model.NewLink, error) {
	return &model.NewLink{
		ID:          obj.Icon.ID,
		Title:       obj.Icon.Title,
		URL:         obj.Icon.URL,
		Description: obj.Icon.Description,
		Thumbnail:   obj.Icon.Thumbnail,
	}, nil
}

func (r *verifiedIdentifierResolver) Timestamp(ctx context.Context, obj *base.VerifiedIdentifier) (*base.Date, error) {
	return &base.Date{
		Year:  obj.Timestamp.Year(),
		Day:   obj.Timestamp.Day(),
		Month: int(obj.Timestamp.Month()),
	}, nil
}

// NavAction returns generated.NavActionResolver implementation.
func (r *Resolver) NavAction() generated.NavActionResolver { return &navActionResolver{r} }

// VerifiedIdentifier returns generated.VerifiedIdentifierResolver implementation.
func (r *Resolver) VerifiedIdentifier() generated.VerifiedIdentifierResolver {
	return &verifiedIdentifierResolver{r}
}

type navActionResolver struct{ *Resolver }
type verifiedIdentifierResolver struct{ *Resolver }
