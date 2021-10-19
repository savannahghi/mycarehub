package staff_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases"
	usecaseMock "github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/mock"
)

var (
	testInfrastructureInteractor     infrastructure.Interactor
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

	infra := func() infrastructure.Interactor {
		return infrastructure.Interactor{
			Create: create,
			Query:  query,
		}
	}()

	i := usecases.NewUsecasesInteractor(infra)

	return i, nil
}
