package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// NHIFUseCases represents all the business logic involved in NHIF
type NHIFUseCases interface {
	AddNHIFDetails(ctx context.Context, input resources.NHIFDetailsInput) (*domain.NHIFDetails, error)
	NHIFDetails(ctx context.Context) (*domain.NHIFDetails, error)
}

// NHIFUseCaseImpl represents the usecase implementation object
type NHIFUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
	baseExt              extension.BaseExtension
	profile              ProfileUseCase
}

// NewNHIFUseCases initializes a new NHIF usecase
func NewNHIFUseCases(
	r repository.OnboardingRepository,
	p ProfileUseCase,
	ext extension.BaseExtension,
) *NHIFUseCaseImpl {
	return &NHIFUseCaseImpl{
		onboardingRepository: r,
		profile:              p,
		baseExt:              ext,
	}
}

// AddNHIFDetails adds NHIF details of a user
func (n NHIFUseCaseImpl) AddNHIFDetails(
	ctx context.Context,
	input resources.NHIFDetailsInput,
) (*domain.NHIFDetails, error) {
	profile, err := n.profile.UserProfile(ctx)
	if err != nil {
		return nil, err
	}

	return n.onboardingRepository.AddNHIFDetails(
		ctx,
		input,
		profile.ID,
	)
}

// NHIFDetails returns NHIF details of a user
func (n NHIFUseCaseImpl) NHIFDetails(
	ctx context.Context,
) (*domain.NHIFDetails, error) {
	profile, err := n.profile.UserProfile(ctx)
	if err != nil {
		return nil, err
	}

	return n.onboardingRepository.GetNHIFDetailsByProfileID(
		ctx,
		profile.ID,
	)
}
