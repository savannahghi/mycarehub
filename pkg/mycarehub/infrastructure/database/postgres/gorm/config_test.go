package gorm_test

import (
	"log"
	"os"
	"testing"

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

	log.Printf("Running tests ...")
	code := m.Run()

	defer tearDown()

	os.Exit(code)

}

func tearDown() {
	testingDB.DB.Migrator().DropTable(&gorm.Contact{})
	testingDB.DB.Migrator().DropTable(&gorm.PINData{})
	testingDB.DB.Migrator().DropTable(&gorm.User{})
	testingDB.DB.Migrator().DropTable(&gorm.Facility{})
}
