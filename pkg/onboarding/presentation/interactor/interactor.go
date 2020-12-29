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
	Onboarding usecases.ProfileUseCase
	Signup     usecases.SignUpUseCases
	Otp        usecases.OTPUseCases
	Supplier   usecases.SupplierUseCases
	Login      usecases.LoginUseCases
	Survey     usecases.SurveyUseCases
	UserPIN    usecases.UserPINUseCases
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor() (*Interactor, error) {

	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	uc := usecases.NewProfileUseCase(fr)
	su := usecases.NewSignUpUseCases(fr, uc)
	otp := usecases.NewOTPUseCasesImpl(fr)
	supplier := usecases.NewSupplierUseCases(fr)
	login := usecases.NewLoginUseCases(fr)
	survey := usecases.NewSurveyUseCases(fr)
	userpin := usecases.NewUserPinUseCase(fr)

	return &Interactor{
		Onboarding: uc,
		Signup:     su,
		Otp:        otp,
		Supplier:   supplier,
		Login:      login,
		Survey:     survey,
		UserPIN:    userpin,
	}, nil
}
