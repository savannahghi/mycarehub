package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// OnboardingUseCase represents all the profile business logic
type OnboardingUseCase interface {
	// UserProfile
	UserProfile(ctx context.Context) (*base.UserProfile, error)
	GetProfile(ctx context.Context, uid string) (*base.UserProfile, error)
	GetProfileByID(ctx context.Context, id string) (*base.UserProfile, error)
	UpdateUserProfile(ctx context.Context, input domain.UserProfileInput) (*base.UserProfile, error)
	SuspendUserProfile(ctx context.Context, phone string) (bool, error)
}

// OnboardingUseCaseImpl represents usecase implementation object
type OnboardingUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
}

// NewOnboardingUseCase returns a new a onboarding usecase
func NewOnboardingUseCase(r repository.OnboardingRepository) *OnboardingUseCaseImpl {
	return &OnboardingUseCaseImpl{r}
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (o *OnboardingUseCaseImpl) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	return o.onboardingRepository.GetUserProfile(ctx, uid)
}

// GetProfileByID returns the profile identified by the indicated ID
func (o *OnboardingUseCaseImpl) GetProfileByID(ctx context.Context, id string) (*base.UserProfile, error) {
	return o.onboardingRepository.GetUserProfileByID(ctx, id)
}

// UpdatePrimaryPhoneNumber updates the primary phone number of a specific user profile
func (o *OnboardingUseCaseImpl) UpdatePrimaryPhoneNumber(ctx context.Context, phoneNumber string) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return err
	}
	profile, err := o.onboardingRepository.GetUserProfile(ctx, uid)
	if err != nil {
		return err
	}
	return profile.UpdateProfilePrimaryPhoneNumber(ctx, o.onboardingRepository, phoneNumber)
}
