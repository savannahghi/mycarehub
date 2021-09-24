// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	libOnboardingUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	OpenSourceInfra    infrastructure.Infrastructure
	OpenSourceUsecases libOnboardingUsecase.Interactor
	LoginUseCase       usecases.UseCaseLogin
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	openSourceInfra infrastructure.Infrastructure,
	openSourceUsecases libOnboardingUsecase.Interactor,
	loginUsecase usecases.UseCaseLogin,
) (*Interactor, error) {
	return &Interactor{
		OpenSourceInfra:    openSourceInfra,
		OpenSourceUsecases: openSourceUsecases,
		LoginUseCase:       loginUsecase,
	}, nil
}
