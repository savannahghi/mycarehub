package database

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

//FacilityRepository exposes the repository methods
type FacilityRepository interface {
	CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}
