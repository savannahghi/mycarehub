package postgres

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// OnboardingDb struct implements the service's business specific calls to the database
type OnboardingDb struct {
	create gorm.Create
	query  gorm.Query
	delete gorm.Delete
}

// NewOnboardingDb initializes a new instance of the OnboardingDB struct
func NewOnboardingDb(c gorm.Create, q gorm.Query, d gorm.Delete) *OnboardingDb {
	return &OnboardingDb{create: c, query: q, delete: d}
}
