package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
)

func (r *kYCRequestResolver) FiledTimestamp(ctx context.Context, obj *domain.KYCRequest) (*base.Date, error) {
	return &base.Date{
		Year:  obj.FiledTimestamp.Year(),
		Day:   obj.FiledTimestamp.Day(),
		Month: int(obj.FiledTimestamp.Month()),
	}, nil
}

func (r *kYCRequestResolver) ProcessedTimestamp(ctx context.Context, obj *domain.KYCRequest) (*base.Date, error) {
	return &base.Date{
		Year:  obj.ProcessedTimestamp.Year(),
		Day:   obj.ProcessedTimestamp.Day(),
		Month: int(obj.ProcessedTimestamp.Month()),
	}, nil
}

func (r *verifiedIdentifierResolver) Timestamp(ctx context.Context, obj *base.VerifiedIdentifier) (*base.Date, error) {
	return &base.Date{
		Year:  obj.Timestamp.Year(),
		Day:   obj.Timestamp.Day(),
		Month: int(obj.Timestamp.Month()),
	}, nil
}

// KYCRequest returns generated.KYCRequestResolver implementation.
func (r *Resolver) KYCRequest() generated.KYCRequestResolver { return &kYCRequestResolver{r} }

// VerifiedIdentifier returns generated.VerifiedIdentifierResolver implementation.
func (r *Resolver) VerifiedIdentifier() generated.VerifiedIdentifierResolver {
	return &verifiedIdentifierResolver{r}
}

type kYCRequestResolver struct{ *Resolver }
type verifiedIdentifierResolver struct{ *Resolver }
