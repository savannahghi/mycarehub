package infrastructure

import (
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	repo "github.com/savannahghi/onboarding-service/pkg/onboarding/repository"
)

// Infrastructure is an implementation of the infrastructure interface
// It combines each individual service implementation
type Infrastructure struct {
	Database repo.FacilityRepository
}

// NewInfrastructureInteractor initializes a new infrastructure Infrastructure
func NewInfrastructureInteractor() (*Infrastructure, error) {

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate repository in infrastructure: %v", err)
	}

	db := postgres.NewOnboardingDb(pg, pg)

	return &Infrastructure{
		db,
	}, nil
}
