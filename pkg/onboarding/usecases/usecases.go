package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	libExtension "github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	libInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	libUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// Interactor represents the login usecse interactor
type Interactor struct {
	LoginUseCases
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(
	infrastructure infrastructure.Infrastructure,
	p libUsecase.ProfileUseCase,
	ext libExtension.BaseExtension,
	pin libExtension.PINExtension,
) *Interactor {
	infra := libInfra.NewInfrastructureInteractor()
	login := libUsecase.NewLoginUseCases(
		infra,
		p,
		ext,
		pin,
	)
	return &Interactor{
		login,
	}
}
