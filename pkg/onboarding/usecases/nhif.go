package usecases

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// AddNHIFNudgeTitle is the title defined in the `engagement service`
// for the `Add NHIF` nudge
const AddNHIFNudgeTitle = "Add NHIF"

// NHIFUseCases represents all the business logic involved in NHIF
type NHIFUseCases interface {
	AddNHIFDetails(
		ctx context.Context,
		input dto.NHIFDetailsInput,
	) (*domain.NHIFDetails, error)
	NHIFDetails(ctx context.Context) (*domain.NHIFDetails, error)
}

// NHIFUseCaseImpl represents the usecase implementation object
type NHIFUseCaseImpl struct {
	onboardingRepository repository.OnboardingRepository
	baseExt              extension.BaseExtension
	profile              ProfileUseCase
	engagement           engagement.ServiceEngagement
}

// NewNHIFUseCases initializes a new NHIF usecase
func NewNHIFUseCases(
	r repository.OnboardingRepository,
	p ProfileUseCase,
	ext extension.BaseExtension,
	e engagement.ServiceEngagement,
) *NHIFUseCaseImpl {
	return &NHIFUseCaseImpl{
		onboardingRepository: r,
		profile:              p,
		baseExt:              ext,
		engagement:           e,
	}
}

// AddNHIFDetails adds NHIF details of a user
func (n NHIFUseCaseImpl) AddNHIFDetails(
	ctx context.Context,
	input dto.NHIFDetailsInput,
) (*domain.NHIFDetails, error) {
	UID, err := n.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}
	profile, err := n.onboardingRepository.GetUserProfileByUID(
		ctx,
		UID,
		false,
	)
	if err != nil {
		return nil, err
	}

	NHIF, err := n.onboardingRepository.AddNHIFDetails(
		ctx,
		input,
		profile.ID,
	)
	if err != nil {
		return nil, err
	}

	go func() {
		cons := func() error {
			return n.engagement.ResolveDefaultNudgeByTitle(
				UID,
				base.FlavourConsumer,
				AddNHIFNudgeTitle,
			)
		}
		if err := backoff.Retry(
			cons,
			backoff.NewExponentialBackOff(),
		); err != nil {
			logrus.Error(err)
		}
	}()

	return NHIF, nil
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
