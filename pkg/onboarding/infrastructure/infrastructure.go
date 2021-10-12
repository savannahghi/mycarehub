package infrastructure

import (
	"log"

	"github.com/savannahghi/firebasetools"
	pg "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	baseExt "github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	libUtils "github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	libInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	engagementSvc "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	libOnboardingUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

const (
	engagementService = "engagement"
)

// Infrastructure is an implementation of the infrastructure interface
// It combines each individual service implementation
type Infrastructure struct {
	Create
	Query
	libOnboardingUsecase.LoginUseCases
	libOnboardingUsecase.SignUpUseCases
	engagementSvc.ServiceEngagementImpl
	libOnboardingUsecase.ProfileUseCase
}

// Interactor is an implementation of the infrastructure interface
// It combines each individual service implementation
type Interactor struct {
	Create
	Query
	libOnboardingUsecase.LoginUseCases
	libOnboardingUsecase.SignUpUseCases
	engagementSvc.ServiceEngagementImpl
	libOnboardingUsecase.ProfileUseCase
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
	postgres, err := gorm.NewPGInstance()
	if err != nil {
		log.Fatal(err)
	}
	db := pg.NewOnboardingDb(postgres, postgres)
	create := NewServiceCreateImpl(*db)
	query := NewServiceQueryImpl(*db)

	return Interactor{
		create,
		query,
		login,
		signup,
		*engagement,
		profile,
	}
}
