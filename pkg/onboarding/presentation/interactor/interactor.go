// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
)

// Usecases combines all onboarding usecases including external
type Usecases interface {
	usecases.Usecases
}

// Interactor combines all onboarding interractors including external
type Interactor struct {
	usecases.Usecases
}

// NewUseCasesInteractor initializes all onboarding usecases including external
func NewUseCasesInteractor(infrastructure *infrastructure.Infrastructure) Usecases {
	onboardingUsecases := usecases.NewUseCaseInteractor(*infrastructure)
	return &Interactor{
		onboardingUsecases,
	}
}
