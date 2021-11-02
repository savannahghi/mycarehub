package facility_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/firebasetools"
	onboardingExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

var (
	testInteractor interactor.Interactor

	// fakeCreate usecaseMock.CreateMock
	// fakeQuery  usecaseMock.QueryMock
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	interactor := InitializeTestInteractor(ctx)

	purgeRecords := func() {

	}

	purgeRecords()

	testInteractor = interactor

	// run the tests
	log.Printf("about to run tests\n")
	code := m.Run()
	log.Printf("finished running tests\n")

	// cleanup here
	os.Exit(code)
}

func InitializeTestInteractor(ctx context.Context) interactor.Interactor {
	//osinfra := openSourceInfra.NewInfrastructureInteractor()
	pgInstance, err := gorm.NewPGInstance()
	if err != nil {
		log.Fatal(err)
	}

	fc := &firebasetools.FirebaseClient{}
	engagementISC := onboardingExtension.NewInterServiceClient("engagement")
	baseExt := extension.NewBaseExtensionImpl(fc)
	engagement := engagement.NewServiceEngagementImpl(engagementISC, baseExt)
	onboardingExt := onboardingExtension.NewOnboardingLibImpl()

	db := postgres.NewMyCareHubDb(pgInstance, pgInstance, pgInstance)

	//infra := infrastructure.NewInteractor()
	facilityUsecase := facility.NewFacilityUsecase(db, db, db)
	clientUseCase := client.NewUseCasesClientImpl(db, db, db)
	userUsecase := user.NewUseCasesUserImpl(db, db, db, onboardingExt, engagement)

	i := interactor.NewMyCareHubInteractor(facilityUsecase, clientUseCase, userUsecase)

	return *i
}

// func InitializeFakeTestlInteractor(ctx context.Context) (usecases.Interactor, error) {

// 	var create infrastructure.Create = &fakeCreate
// 	var query infrastructure.Query = &fakeQuery

// 	infra := func() infrastructure.Interactor {
// 		return infrastructure.Interactor{
// 			Create: create,
// 			Query:  query,
// 		}
// 	}()

// 	i := usecases.NewUsecasesInteractor(infra)

// 	return i, nil
// }
