package usecases

import (
	libSignUpUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// UseCaseSignUp represent the open sourced sign up usecase
type UseCaseSignUp interface {
	libSignUpUsecase.SignUpUseCases
}

// UseCaseSignUpImpl is the open sourced usecase implememtation
type UseCaseSignUpImpl struct {
	LibSignUp libSignUpUsecase.SignUpUseCases
}

// NewSignUpUseCase instantiates signup usecases
func NewSignUpUseCase(
	libSignUp libSignUpUsecase.SignUpUseCases,
) *UseCaseSignUpImpl {
	return &UseCaseSignUpImpl{
		LibSignUp: libSignUp,
	}
}
