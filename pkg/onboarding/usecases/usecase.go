package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	engagementSvc "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	// *UseCaseLoginImpl
	*UseCaseSignUpImpl
	*engagementSvc.ServiceEngagementImpl
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {
	// login := NewUseCaseLogin(infrastructure)
	signup := NewSignUpUseCase(infrastructure)
	var engagement *engagementSvc.ServiceEngagementImpl

	return Interactor{
		// login,
		signup,
		engagement,
	}
}
