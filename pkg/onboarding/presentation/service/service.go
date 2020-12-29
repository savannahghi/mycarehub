package service

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

// Service represents an assemble of all use cases into a single object that can be instantiated anywhere
type Service struct {
	Onboarding *usecases.OnboardingUseCaseImpl
	Signup     *usecases.SignUpUseCasesImpl
}

// NewService returns a new instance of Service
func NewService() (*Service, error) {

	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	uc := usecases.NewOnboardingUseCase(fr)
	if uc == nil {
		return nil, fmt.Errorf("can't instantiate onboarding usecases")
	}

	su := usecases.NewSignUpUseCases(fr)
	if su == nil {
		return nil, fmt.Errorf("can't instantiate signup usecases")
	}
	return &Service{
		Onboarding: uc,
		Signup:     su,
	}, nil
}
