package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// ProfileUseCase represents all the profile business logic
type ProfileUseCase interface {
	// UserProfile
	UserProfile(ctx context.Context) (*base.UserProfile, error)
	GetProfile(ctx context.Context, uid string) (*base.UserProfile, error)
	GetProfileByID(ctx context.Context, id string) (*base.UserProfile, error)
	// updates the user profile of the currently logged in user
	UpdateUserProfile(ctx context.Context, input *domain.UserProfileInput) (*domain.UserResponse, error)

	SuspendUserProfile(ctx context.Context, phone string) (bool, error)

	RecordPostVisitSurvey(ctx context.Context, input domain.PostVisitSurveyInput) (bool, error)
}

// ProfileUseCaseImpl represents usecase implementation object
type ProfileUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
}

// NewProfileUseCase returns a new a onboarding usecase
func NewProfileUseCase(r repository.OnboardingRepository) *ProfileUseCaseImpl {
	return &ProfileUseCaseImpl{r}
}

// UserProfile retrieves the profile of the logged in user, if they have one
func (p *ProfileUseCaseImpl) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}
	return p.onboardingRepository.GetUserProfileByUID(ctx, uid)
}

// GetProfileByID returns the profile identified by the indicated ID
func (p *ProfileUseCaseImpl) GetProfileByID(ctx context.Context, id string) (*base.UserProfile, error) {
	return p.onboardingRepository.GetUserProfileByID(ctx, id)
}

// UpdatePrimaryPhoneNumber updates the primary phone number of a specific user profile
func (p *ProfileUseCaseImpl) UpdatePrimaryPhoneNumber(ctx context.Context, phoneNumber string) error {
	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return err
	}
	profile, err := p.onboardingRepository.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return err
	}
	return profile.UpdateProfilePrimaryPhoneNumber(ctx, p.onboardingRepository, phoneNumber)
}
