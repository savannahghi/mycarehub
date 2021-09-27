package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	engagementSvc "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*UseCaseSignUpImpl
	*engagementSvc.ServiceEngagementImpl
	*UsecaseProfileImpl
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {
	signup := NewSignUpUseCase(infrastructure)
	var engagement *engagementSvc.ServiceEngagementImpl
	profile := NewProfileUseCase(infrastructure)

	return Interactor{
		signup,
		engagement,
		profile,
	}
}
