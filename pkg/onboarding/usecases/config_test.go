package usecases_test

import (
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database"
	mockRepo "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/mock"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	extMock "github.com/savannahghi/onboarding/pkg/onboarding/application/extension/mock"
	libInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	sharelib "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

var fakeBaseExt extMock.FakeBaseExtensionImpl
var fakePinExt extMock.PINExtensionImpl
var fakeRepo mockRepo.FakeOnboardingRepository

func InitializeFakeInfrastructure() infrastructure.Infrastructure {
	var r database.Repository = &fakeRepo

	type InfrastructureMock struct {
		database.Repository
	}
	infra := func() infrastructure.Infrastructure {
		return &InfrastructureMock{
			r,
		}
	}()

	return infra
}

// InitializeFakeOnboardingInteractor represents a fakeonboarding interactor
func InitializeFakeOnboardingInteractor() (interactor.Interactor, error) {
	var ext extension.BaseExtension = &fakeBaseExt
	var pinExt extension.PINExtension = &fakePinExt

	// infrastricture from library
	libOpenInfra := libInfra.NewInfrastructureInteractor()
	// internally declared infrastructure
	svcInfra, err := infrastructure.NewInfrastructureInteractor()
	if err != nil {
		return interactor.Interactor{}, fmt.Errorf("failed to initialize new infrastructure interractor: %v", err)
	}
	sharedLibinteractor := sharelib.NewUsecasesInteractor(
		libOpenInfra,
		ext,
		pinExt,
	)
	svcUsecaseInteractor := usecases.NewUsecasesInteractor(svcInfra, sharedLibinteractor, ext, pinExt)

	i, err := interactor.NewOnboardingInteractor(libOpenInfra, sharedLibinteractor, *svcUsecaseInteractor)
	if err != nil {
		return interactor.Interactor{}, fmt.Errorf("failed to initialize new usecases interractor: %v", err)
	}
	return *i, nil

}
