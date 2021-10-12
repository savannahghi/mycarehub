package facility_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

var (
	testInfrastructureInteractor infrastructure.Interactor
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	infra, err := InitializeTestInfrastructure(ctx)
	if err != nil {
		log.Printf("failed to initialize infrastructure: %v", err)
	}

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

// InitializeTestInfrastructure initializes the test infrastructure.
func InitializeTestInfrastructure(ctx context.Context) (infrastructure.Interactor, error) {
	interactor := infrastructure.NewInteractor()
	return interactor, nil
}
