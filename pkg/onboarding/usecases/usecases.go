package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
)

// Usecases combines all onboarding usecases
type Usecases interface {
	facility.UseCasesFacility
}

// Interactor combines all onboarding interractors
type Interactor struct {
	facility.UseCasesFacility
}

// NewUseCaseInteractor initializes all onboarding usecases
func NewUseCaseInteractor(infrastructure infrastructure.Infrastructure) Usecases {
	fcl := facility.NewFacilityUsecase(infrastructure)

	return &Interactor{
		fcl,
	}
}
