package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/profile/domain"
	"gitlab.slade360emr.com/go/profile/pkg/profile/repository"
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

// NewOnboardingUseCase returns a new authentication service
func NewOnboardingUseCase(r repository.OnboardingRepository) *OnboardingUseCaseImpl {
	return &OnboardingUseCaseImpl{r}
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (o OnboardingUseCaseImpl) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	return o.onboardingRepository.GetUserProfile(ctx, uid)
}

// GetProfileByID returns the profile identified by the indicated ID
func (o OnboardingUseCaseImpl) GetProfileByID(ctx context.Context, id string) (*base.UserProfile, error) {
	return o.onboardingRepository.GetUserProfileByID(ctx, id)
}
