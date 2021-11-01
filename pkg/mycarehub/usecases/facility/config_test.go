package facility_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/firebasetools"
	onboardingExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	usecaseMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	openSourceInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

var (
	testInfrastructureInteractor     infrastructure.Interactor
	testInteractor                   interactor.Interactor
	testFakeInfrastructureInteractor usecases.Interactor

	fakeCreate usecaseMock.CreateMock
	fakeQuery  usecaseMock.QueryMock
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("failed to initialize infrastructure: %v", err)
	}

	interactor := InitializeTestInteractor(ctx)

	fakeInfra, err := InitializeFakeTestlInteractor(ctx)
	if err != nil {
		log.Printf("failed to initialize fake usecase interractor: %v", err)
	}
	testFakeInfrastructureInteractor = fakeInfra

	purgeRecords := func() {

	}

	purgeRecords()

	testInfrastructureInteractor = infra
	testInteractor = interactor

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

func InitializeTestInteractor(ctx context.Context) interactor.Interactor {
	osinfra := openSourceInfra.NewInfrastructureInteractor()
	pgInstance, err := gorm.NewPGInstance()
	if err != nil {
		log.Fatal(err)
	}
	fc := &firebasetools.FirebaseClient{}
	engagementISC := onboardingExtension.NewInterServiceClient("engagement")
	baseExt := extension.NewBaseExtensionImpl(fc)
	engagement := engagement.NewServiceEngagementImpl(engagementISC, baseExt)
	onboardingExt := onboardingExtension.NewOnboardingLibImpl()
	infra := infrastructure.NewInteractor()
	facilityUsecase := facility.NewFacilityUsecase(infra)
	clientUseCase := client.NewUseCasesClientImpl(infra)
	userUsecase := user.NewUseCasesUserImpl(infra, onboardingExt, engagement)
	db := postgres.NewOnboardingDb(pgInstance, pgInstance, pgInstance)
	i := interactor.NewOnboardingInteractor(osinfra, *db, facilityUsecase, clientUseCase, userUsecase)

	return *i
}

func InitializeFakeTestlInteractor(ctx context.Context) (usecases.Interactor, error) {

	var create infrastructure.Create = &fakeCreate
	var query infrastructure.Query = &fakeQuery

	infra := func() infrastructure.Interactor {
		return infrastructure.Interactor{
			Create: create,
			Query:  query,
		}
	}()

	i := usecases.NewUsecasesInteractor(infra)

	return i, nil
}
