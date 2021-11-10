package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func (r *mutationResolver) GetCurrentTerms(ctx context.Context, flavour enums.Flavour) (string, error) {
	r.checkPreconditions()
	return r.interactor.TermsUsecase.GetCurrentTerms(ctx, flavour)
}
