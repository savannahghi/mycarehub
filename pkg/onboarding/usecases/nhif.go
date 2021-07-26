package usecases

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/sirupsen/logrus"
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
	ctx, span := tracer.Start(ctx, "AddNHIFDetails")
	defer span.End()

	UID, err := n.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.UserNotFoundError(err)
	}
	profile, err := n.onboardingRepository.GetUserProfileByUID(
		ctx,
		UID,
		false,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	NHIF, err := n.onboardingRepository.AddNHIFDetails(
		ctx,
		input,
		profile.ID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	go func() {
		cons := func() error {
			return n.engagement.ResolveDefaultNudgeByTitle(
				ctx,
				UID,
				feedlib.FlavourConsumer,
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
	ctx, span := tracer.Start(ctx, "NHIFDetails")
	defer span.End()

	profile, err := n.profile.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	return n.onboardingRepository.GetNHIFDetailsByProfileID(
		ctx,
		profile.ID,
	)
}
