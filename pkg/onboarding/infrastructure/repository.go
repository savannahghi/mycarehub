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

// Delete represents all the deletion action interfaces
type Delete interface {
	DeleteFacility(ctx context.Context, id string) (bool, error)
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

// CreateFacility is responsible for creating a representation of a facility
func (f ServiceCreateImpl) CreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.onboarding.CreateFacility(ctx, &facility)
}

// Query represents a contract that contains all `get` ops to the database
//
// All the  contracts for get operations are assembled here
type Query interface {
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
}

// ServiceQueryImpl represents create contract implementation object
type ServiceQueryImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceQueryImpl returns new instance of ServiceQueryImpl
func NewServiceQueryImpl(on pg.OnboardingDb) Query {
	return &ServiceQueryImpl{
		onboarding: on,
	}
}

//GetFacilities is responsible for returning a slice of healthcare facilities in the platform.
func (q ServiceQueryImpl) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	return q.onboarding.GetFacilities(ctx)
}

// DeleteFacility is responsible for deletion of a facility from the database using the facility's id
func (f ServiceDeleteImpl) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return f.onboarding.DeleteFacility(ctx, id)
}

// ServiceDeleteImpl represents delete facility implementation object
type ServiceDeleteImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceDeleteImpl returns new instance of NewServiceDeleteImpl
func NewServiceDeleteImpl(on pg.OnboardingDb) Delete {
	return &ServiceDeleteImpl{
		onboarding: on,
	}
}
