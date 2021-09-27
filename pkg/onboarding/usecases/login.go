package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	libLoginUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// UseCaseLogin represents open source login usecases
type UseCaseLogin interface {
	libLoginUsecase.LoginUseCases
}

// UseCaseLoginImpl is the login usecase implementation
type UseCaseLoginImpl struct {
	infrastructure infrastructure.Interactor
}

// NewUseCaseLogin instantiates login usecases
func NewUseCaseLogin(
	infrastructure infrastructure.Interactor,
) *UseCaseLoginImpl {
	return &UseCaseLoginImpl{
		infrastructure: infrastructure,
	}
}
