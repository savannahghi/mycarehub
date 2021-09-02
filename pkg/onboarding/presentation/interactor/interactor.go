// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	sharelib "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	OpenSourceInfra    infrastructure.Infrastructure
	OpenSourceUsecases sharelib.Interactor
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	openSourceInfra infrastructure.Infrastructure,
	openSourceUsecases sharelib.Interactor,
) (*Interactor, error) {
	return &Interactor{
		OpenSourceInfra:    openSourceInfra,
		OpenSourceUsecases: openSourceUsecases,
	}, nil
}
