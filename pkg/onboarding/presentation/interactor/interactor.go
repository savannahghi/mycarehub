// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	pg "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	OpenSourceInfra infrastructure.Infrastructure
	database        pg.OnboardingDb
	FacilityUsecase facility.UseCasesFacility
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	openSourceInfra infrastructure.Infrastructure,
	database pg.OnboardingDb,
	facilityUseCase facility.UseCasesFacility,
) *Interactor {
	return &Interactor{
		OpenSourceInfra: openSourceInfra,
		database:        database,
		FacilityUsecase: facilityUseCase,
	}
}
