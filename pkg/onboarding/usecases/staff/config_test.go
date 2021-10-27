package staff_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/firebasetools"
	onboardingExtension "github.com/savannahghi/onboarding-service/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/client"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/metric"
	usecaseMock "github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/mock"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/staff"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/user"
	baseExt "github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	openSourceInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	libOnboardingUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

var (
	testInfrastructureInteractor     infrastructure.Interactor
	testInteractor                   interactor.Interactor
	testFakeInfrastructureInteractor usecases.Interactor

	fakeCreate usecaseMock.CreateMock
	fakeQuery  usecaseMock.QueryMock
	fakeUpdate usecaseMock.UpdateMock
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("failed to initialize infrastructure: %v", err)
	}

	interactor := InitializeTestInteractor(ctx)

	testInteractor = interactor

	fakeInfra, err := InitializeFakeTestlInteractor(ctx)
	if err != nil {
		log.Printf("failed to initialize fake usecase interractor: %v", err)
	}
	testFakeInfrastructureInteractor = fakeInfra

	purgeRecords := func() {

	}

	purgeRecords()

	testInfrastructureInteractor = infra

	// run the tests
	log.Printf("about to run tests\n")
	code := m.Run()
	log.Printf("finished running tests\n")

	// cleanup here
	os.Exit(code)
}

func InitializeTestInfrastructure(ctx context.Context) (infrastructure.Interactor, error) {
	i := infrastructure.NewInteractor()
	return i, nil
}

func InitializeFakeTestlInteractor(ctx context.Context) (usecases.Interactor, error) {

	var create infrastructure.Create = &fakeCreate
	var query infrastructure.Query = &fakeQuery
	var update infrastructure.Update = &fakeUpdate

	infra := func() infrastructure.Interactor {
		return infrastructure.Interactor{
			Create: create,
			Query:  query,
			Update: update,
		}
	}()

	i := usecases.NewUsecasesInteractor(infra)

	return i, nil
}

func InitializeTestInteractor(ctx context.Context) interactor.Interactor {
	osinfra := openSourceInfra.NewInfrastructureInteractor()
	pgInstance, err := gorm.NewPGInstance()
	if err != nil {
		log.Fatal(err)
	}
	// Initialize ISC clients
	engagementISC := onboardingExtension.NewInterServiceClient("engagement")

	infra := infrastructure.NewInteractor()
	facilityUsecase := facility.NewFacilityUsecase(infra)
	metricUsecase := metric.NewMetricUsecase(infra)
	db := postgres.NewOnboardingDb(pgInstance, pgInstance, pgInstance, pgInstance)
	var fc firebasetools.IFirebaseClient
	baseExtension := baseExt.NewBaseExtensionImpl(fc)
	pinExtension := baseExt.NewPINExtensionImpl()
	libUsecasee := libOnboardingUsecase.NewUsecasesInteractor(osinfra, baseExtension, pinExtension)
	onboardingExtension := onboardingExtension.NewOnboardingLibImpl()
	engagement := engagement.NewServiceEngagementImpl(engagementISC, baseExtension)
	userUsecase := user.NewUseCasesUserImpl(infra, onboardingExtension, engagement)
	staff := staff.NewUsecasesStaffProfileImpl(infra)
	client := client.NewUseCasesClientImpl(infra)
	i := interactor.NewOnboardingInteractor(osinfra, *db, libUsecasee, facilityUsecase, metricUsecase, userUsecase, staff, client)

	return *i
}
