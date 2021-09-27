package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	libSignUpUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// UseCaseSignUp represent the open sourced sign up usecase
type UseCaseSignUp interface {
	libSignUpUsecase.SignUpUseCases
}

// UseCaseSignUpImpl is the open sourced usecase implememtation
type UseCaseSignUpImpl struct {
	infrastructure infrastructure.Interactor
}

// NewSignUpUseCase instantiates signup usecases
func NewSignUpUseCase(
	infrastructure infrastructure.Interactor,
) *UseCaseSignUpImpl {
	return &UseCaseSignUpImpl{
		infrastructure: infrastructure,
	}

}
