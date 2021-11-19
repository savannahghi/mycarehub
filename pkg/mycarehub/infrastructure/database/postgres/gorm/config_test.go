package gorm_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

var testingDB *gorm.PGInstance

var (
	termsID  = 50005
	orgID    = uuid.New().String()
	pastTime = time.Now().AddDate(0, 0, -1)
)

func TestMain(m *testing.M) {
	log.Println("setting up test database")
	var err error
	testingDB, err = gorm.NewPGInstance()
	if err != nil {
		os.Exit(1)
	}
	// add organization
	createOrganization()

	//create terms
	createTermsOfService()

	log.Printf("Running tests ...")
	os.Exit(m.Run())

	// teardown
	// remove organization
	log.Println("tearing down test database")

	testingDB.DB.Unscoped().Delete(gorm.Organisation{OrganisationID: &orgID})
	testingDB.DB.Unscoped().Delete(gorm.TermsOfService{TermsID: &termsID})
}

func createOrganization() {
	organisation := &gorm.Organisation{
		OrganisationID:   &orgID,
		Active:           true,
		OrgCode:          gofakeit.Name(),
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

func createTermsOfService() {
	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	terms := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
	}

	testingDB.DB.Create(terms)
}
