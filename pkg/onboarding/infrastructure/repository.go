package infrastructure

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	pg "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	CreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
}

// ServiceCreateImpl represents create contract implementation object
type ServiceCreateImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceCreateImpl returns new instance of ServiceCreateImpl
func NewServiceCreateImpl(on pg.OnboardingDb) Create {
	return &ServiceCreateImpl{
		onboarding: on,
	}
}

// CreateFacility is responsible from creating a representation of a facility
func (f ServiceCreateImpl) CreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.onboarding.CreateFacility(ctx, &facility)
}
