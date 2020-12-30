// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	Onboarding *usecases.ProfileUseCaseImpl
	Signup     *usecases.SignUpUseCasesImpl
	Otp        *usecases.OTPUseCasesImpl
	Supplier   *usecases.SupplierUseCasesImpl
	Login      *usecases.LoginUseCasesImpl
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor() (*Interactor, error) {

	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	uc := usecases.NewProfileUseCase(fr)
	if uc == nil {
		return nil, fmt.Errorf("can't instantiate onboarding usecases")
	}

	su := usecases.NewSignUpUseCases(fr)
	if su == nil {
		return nil, fmt.Errorf("can't instantiate signup usecases")
	}

	otp := usecases.NewOTPUseCasesImpl(fr)
	if otp == nil {
		return nil, fmt.Errorf("can't instantiate otp usecases")
	}

	supplier := usecases.NewSupplierUseCases(fr)
	if supplier == nil {
		return nil, fmt.Errorf("can't instantiate supplier usecases")
	}
	login := usecases.NewLoginUseCases(fr)
	if login == nil {
		return nil, fmt.Errorf("can't instantiate login usecases")
	}

	return &Interactor{
		Onboarding: uc,
		Signup:     su,
		Otp:        otp,
		Supplier:   supplier,
		Login:      login,
	}, nil
}
