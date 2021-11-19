package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
)

func (r *queryResolver) SendOtp(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
	r.checkPreconditions()
	return r.mycarehub.OTP.GenerateAndSendOTP(ctx, phoneNumber, flavour)
}
