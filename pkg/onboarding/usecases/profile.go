package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	libProfileUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// UsecaseProfile represent the open sourced profile usecase
type UsecaseProfile interface {
	libProfileUsecase.ProfileUseCase
}

// UsecaseProfileImpl is the open sourced usecase implememtation
type UsecaseProfileImpl struct {
	infrastructure infrastructure.Interactor
}

// NewProfileUseCase instantiates profile usecases
func NewProfileUseCase(
	infrastructure infrastructure.Interactor,
) *UsecaseProfileImpl {
	return &UsecaseProfileImpl{
		infrastructure: infrastructure,
	}

}
