package postgres

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// OnboardingDb struct implements ther service's business specific calls to the database
type OnboardingDb struct {
	create gorm.Create
	query  gorm.Query
	// update gorm.Update
	// delete gorm.Delete
}

// NewOnboardingDb initializes a new instance of the OnboardingDB struct
func NewOnboardingDb(c gorm.Create, q gorm.Query) *OnboardingDb {
	return &OnboardingDb{create: c, query: q}
}
