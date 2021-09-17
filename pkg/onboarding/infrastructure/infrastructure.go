package infrastructure

import (
	"github.com/savannahghi/firebasetools"
	baseExt "github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	libUtils "github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	libInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	engagementSvc "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	libOnboardingUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

const (
	engagementService = "engagement"
)

// Interactor is an implementation of the infrastructure interface
// It combines each individual service implementation
type Interactor struct {
	libOnboardingUsecase.LoginUseCases
	libOnboardingUsecase.SignUpUseCases
	engagementSvc.ServiceEngagementImpl
}

// NewInteractor initializes a new infrastructure interactor
func NewInteractor() Interactor {

	i := libInfra.NewInfrastructureInteractor()
	var fc firebasetools.IFirebaseClient
	baseExtension := baseExt.NewBaseExtensionImpl(fc)
	pinExtension := baseExt.NewPINExtensionImpl()
	profile := libOnboardingUsecase.NewProfileUseCase(i, baseExtension)
	userPinUseCase := libOnboardingUsecase.NewUserPinUseCase(i, profile, baseExtension, pinExtension)
	login := libOnboardingUsecase.NewLoginUseCases(i, profile, baseExtension, pinExtension)
	signup := libOnboardingUsecase.NewSignUpUseCases(i, profile, userPinUseCase, baseExtension)
	engagementClient := libUtils.NewInterServiceClient(engagementService, baseExtension)
	engagement := engagementSvc.NewServiceEngagementImpl(engagementClient, baseExtension)

	return Interactor{
		login,
		signup,
		*engagement,
	}
}
