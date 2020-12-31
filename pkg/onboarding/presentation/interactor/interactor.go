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
	Onboarding   usecases.ProfileUseCase
	Signup       usecases.SignUpUseCases
	Otp          usecases.OTPUseCases
	Supplier     usecases.SupplierUseCases
	Login        usecases.LoginUseCases
	Survey       usecases.SurveyUseCases
	UserPIN      usecases.UserPINUseCases
	ERP          usecases.ERPUseCases
	ChargeMaster usecases.ChargeMasterUseCases
	Engage       usecases.EngagementUseCases
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor() (*Interactor, error) {

	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	profile := usecases.NewProfileUseCase(fr)
	otp := usecases.NewOTPUseCasesImpl(fr)
	erp := usecases.NewERPUseCases(fr)
	chrg := usecases.NewChargeMasterUseCasesImpl(fr)
	engage := usecases.NewEngagementUseCasesImpl(fr)
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage)
	login := usecases.NewLoginUseCases(fr)
	survey := usecases.NewSurveyUseCases(fr)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier)

	return &Interactor{
		Onboarding:   profile,
		Signup:       su,
		Otp:          otp,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engage:       engage,
	}, nil
}
