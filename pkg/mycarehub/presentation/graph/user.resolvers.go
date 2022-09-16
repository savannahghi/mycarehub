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
	return r.mycarehub.User.Consent(ctx, phoneNumber, flavour)
}

func (r *mutationResolver) SetPushToken(ctx context.Context, token string) (bool, error) {
	return r.mycarehub.User.RegisterPushToken(ctx, token)
}

func (r *mutationResolver) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite *bool) (bool, error) {
	return r.mycarehub.User.InviteUser(ctx, userID, phoneNumber, flavour, *reinvite)
}

func (r *mutationResolver) SetUserPin(ctx context.Context, input *dto.PINInput) (bool, error) {
	return r.mycarehub.User.SetUserPIN(ctx, *input)
}

func (r *mutationResolver) TransferClientToFacility(ctx context.Context, clientID string, facilityID string) (bool, error) {
	return r.mycarehub.User.TransferClientToFacility(ctx, &clientID, &facilityID)
}

func (r *mutationResolver) SetStaffDefaultFacility(ctx context.Context, userID string, facilityID string) (bool, error) {
	return r.mycarehub.User.SetStaffDefaultFacility(ctx, userID, facilityID)
}

func (r *mutationResolver) SetClientDefaultFacility(ctx context.Context, userID string, facilityID string) (bool, error) {
	return r.mycarehub.User.SetClientDefaultFacility(ctx, userID, facilityID)
}

func (r *mutationResolver) AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error) {
	return r.mycarehub.User.AddFacilitiesToStaffProfile(ctx, staffID, facilities)
}

func (r *mutationResolver) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error) {
	return r.mycarehub.User.RemoveFacilitiesFromStaffProfile(ctx, staffID, facilities)
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

func (r *queryResolver) SearchClientUser(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
	return r.mycarehub.User.SearchClientUser(ctx, searchParameter)
}

func (r *queryResolver) SearchStaffUser(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
	return r.mycarehub.User.SearchStaffUser(ctx, searchParameter)
}

func (r *queryResolver) GetClientProfileByCCCNumber(ctx context.Context, cCCNumber string) (*domain.ClientProfile, error) {
	return r.mycarehub.User.GetClientProfileByCCCNumber(ctx, cCCNumber)
}

func (r *queryResolver) GetUserLinkedFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return r.mycarehub.User.GetUserLinkedFacilities(ctx)
}
