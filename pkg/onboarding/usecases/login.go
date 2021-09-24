package usecases

import (
	libLoginUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// UseCaseLogin represents open source login usecases
type UseCaseLogin interface {
	libLoginUsecase.LoginUseCases
}

// UseCaseLoginImpl is the login usecase implementation
type UseCaseLoginImpl struct {
	LibLogin libLoginUsecase.LoginUseCases
}

// NewUseCaseLogin instantiates login usecases
func NewUseCaseLogin(
	libLogin libLoginUsecase.LoginUseCases,
) *UseCaseLoginImpl {
	return &UseCaseLoginImpl{
		LibLogin: libLogin,
	}
}
