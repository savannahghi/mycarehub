package gorm_test

import (
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

var testingDB *gorm.PGInstance

func TestMain(m *testing.M) {
	log.Println("setting up test database")
	var err error
	testingDB, err = gorm.NewPGInstance()
	if err != nil {
		os.Exit(1)
	}
	// add organization
	createOrganization()

	log.Printf("Running tests ...")
	os.Exit(m.Run())

	// teardown
	// remove organization
	log.Println("tearing down test database")
	orgID := os.Getenv("DEFAULT_ORG_ID")
	testingDB.DB.Unscoped().Delete(gorm.Organisation{OrganisationID: &orgID})
}

func createOrganization() {
	orgID := os.Getenv("DEFAULT_ORG_ID")
	organisation := &gorm.Organisation{
		OrganisationID:   &orgID,
		Active:           true,
		Deleted:          false,
		OrgCode:          "ORG_123",
		Code:             gofakeit.Number(100, 344),
		OrganisationName: gofakeit.Name(),
		EmailAddress:     gofakeit.Email(),
		PhoneNumber:      gofakeit.Phone(),
		PostalAddress:    gofakeit.Address().Address,
		PhysicalAddress:  gofakeit.Address().City,
		DefaultCountry:   "KEN",
	}

	testingDB.DB.Create(organisation)
}
