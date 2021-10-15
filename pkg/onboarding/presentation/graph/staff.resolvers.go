package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

func (r *mutationResolver) RegisterStaffProfile(ctx context.Context, userInput dto.UserInput, staffProfileInput dto.StaffProfileInput) (*domain.StaffUserProfileOutput, error) {
	return r.interactor.StaffProfileUsecase.RegisterStaffUser(ctx, userInput, staffProfileInput)
}
