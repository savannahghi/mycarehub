package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

func (r *mutationResolver) GetOrCreateStaffUser(ctx context.Context, userInput dto.UserInput, staffInput dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	return r.interactor.StaffUsecase.GetOrCreateStaffUser(ctx, &userInput, &staffInput)
}

func (r *mutationResolver) UpdateStaffUserProfile(ctx context.Context, userID string, userInput *dto.UserInput, staffInput *dto.StaffProfileInput) (bool, error) {
	return r.interactor.StaffUsecase.UpdateStaffUserProfile(ctx, userID, userInput, staffInput)
}
