package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
)

func (r *queryResolver) SendOtp(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (string, error) {
	r.checkPreconditions()
	return r.interactor.OTP.GenerateAndSendOTP(ctx, userID, phoneNumber, flavour)
}
