package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) AcceptTerms(ctx context.Context, userID string, termsID int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Terms.AcceptTerms(ctx, &userID, &termsID)
}

func (r *mutationResolver) SetNickName(ctx context.Context, userID string, nickname string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.User.SetNickName(ctx, userID, nickname)
}

func (r *mutationResolver) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.User.CompleteOnboardingTour(ctx, userID, flavour)
}

func (r *mutationResolver) CreateOrUpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
	return r.mycarehub.User.CreateOrUpdateClientCaregiver(ctx, caregiverInput)
}

func (r *mutationResolver) RegisterClient(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error) {
	return r.mycarehub.User.RegisterClient(ctx, input)
}

func (r *mutationResolver) RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
	return r.mycarehub.User.RegisterStaff(ctx, input)
}

func (r *mutationResolver) OptOut(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	return r.mycarehub.User.Consent(ctx, phoneNumber, flavour, false)
}

func (r *mutationResolver) SetPushToken(ctx context.Context, token string) (bool, error) {
	return r.mycarehub.User.RegisterPushToken(ctx, token)
}

func (r *queryResolver) GetCurrentTerms(ctx context.Context, flavour feedlib.Flavour) (*domain.TermsOfService, error) {
	r.checkPreconditions()
	return r.mycarehub.Terms.GetCurrentTerms(ctx, flavour)
}

func (r *queryResolver) VerifyPin(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
	return r.mycarehub.User.VerifyPIN(ctx, userID, flavour, pin)
}

func (r *queryResolver) GetClientCaregiver(ctx context.Context, clientID string) (*domain.Caregiver, error) {
	return r.mycarehub.User.GetClientCaregiver(ctx, clientID)
}

func (r *queryResolver) SearchClientsByCCCNumber(ctx context.Context, cCCNumber string) ([]*domain.ClientProfile, error) {
	return r.mycarehub.User.SearchClientsByCCCNumber(ctx, cCCNumber)
}

func (r *queryResolver) SearchStaffByStaffNumber(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
	return r.mycarehub.User.SearchStaffByStaffNumber(ctx, staffNumber)
}

func (r *queryResolver) GetClientProfileByCCCNumber(ctx context.Context, cCCNumber string) (*domain.ClientProfile, error) {
	return r.mycarehub.User.GetClientProfileByCCCNumber(ctx, cCCNumber)
}
