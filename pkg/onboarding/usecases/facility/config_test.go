package facility_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	usecaseMock "github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/mock"
	openSourceInfra "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
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
	infra := infrastructure.NewInteractor()
	facilityUsecase := facility.NewFacilityUsecase(infra)
	db := postgres.NewOnboardingDb(pgInstance, pgInstance, pgInstance)
	i := interactor.NewOnboardingInteractor(osinfra, *db, facilityUsecase)

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
