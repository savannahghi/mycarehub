package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) AcceptTerms(ctx context.Context, userID string, termsID int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Terms.AcceptTerms(ctx, &userID, &termsID)
}

func (r *mutationResolver) SetNickName(ctx context.Context, userID string, nickname string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.User.SetNickName(ctx, &userID, &nickname)
}

func (r *mutationResolver) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.User.CompleteOnboardingTour(ctx, userID, flavour)
}

func (r *mutationResolver) ForgetMe(ctx context.Context, userID string, pin string, flavour feedlib.Flavour) (bool, error) {
	r.checkPreconditions()

	return r.mycarehub.User.ForgetMe(ctx, userID, pin, flavour)
}

func (r *queryResolver) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	r.checkPreconditions()
	return r.mycarehub.Terms.GetCurrentTerms(ctx)
}

func (r *queryResolver) VerifyPin(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
	return r.mycarehub.User.VerifyPIN(ctx, userID, flavour, pin)
}
