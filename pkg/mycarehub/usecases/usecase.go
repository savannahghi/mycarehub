package usecases

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
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
