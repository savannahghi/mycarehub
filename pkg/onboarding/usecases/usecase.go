package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*facility.UseCaseFacilityImpl
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {
	facility := facility.NewFacilityUsecase(infrastructure)

	return Interactor{
		facility,
	}
}
