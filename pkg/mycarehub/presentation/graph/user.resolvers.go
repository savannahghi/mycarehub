package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
)

func (r *mutationResolver) GetCurrentTerms(ctx context.Context) (string, error) {
	r.checkPreconditions()
	return r.interactor.TermsUsecase.GetCurrentTerms(ctx)
}
