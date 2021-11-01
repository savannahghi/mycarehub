// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	pg "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	OpenSourceInfra infrastructure.Infrastructure
	database        pg.OnboardingDb
	FacilityUsecase facility.UseCasesFacility
	ClientUseCase   client.UseCasesClientProfile
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	openSourceInfra infrastructure.Infrastructure,
	database pg.OnboardingDb,
	facilityUseCase facility.UseCasesFacility,
	clientUseCase client.UseCasesClientProfile,
) *Interactor {
	return &Interactor{
		OpenSourceInfra: openSourceInfra,
		database:        database,
		FacilityUsecase: facilityUseCase,
		ClientUseCase:   clientUseCase,
	}
}
