package infrastructure

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	pg "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
)

// Create represents a contract that contains all `create` ops to the database
//
// All the  contracts for create operations are assembled here
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error)
	CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error)
	CreateUser(
		ctx context.Context,
		input *dto.CreateUserInput,
	) (*domain.User, error)
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

// GetOrCreateFacility is responsible for creating a representation of a facility
func (f ServiceCreateImpl) GetOrCreateFacility(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
	return f.onboarding.GetOrCreateFacility(ctx, &facility)
}

// CollectMetrics is responsible for creating a representation of a metric
func (f ServiceCreateImpl) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
	return f.onboarding.CollectMetrics(ctx, metric)
}

// CreateUser is responsible for creating a representation of a user
func (f ServiceCreateImpl) CreateUser(
	ctx context.Context,
	input *dto.CreateUserInput,
) (*domain.User, error) {
	return f.onboarding.CreateUser(ctx, input)
}

// Query contains all query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error)
	GetFacilities(ctx context.Context) ([]*domain.Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error)
}

// ServiceQueryImpl contains implementation for the Query interface
type ServiceQueryImpl struct {
	onboarding pg.OnboardingDb
}

// NewServiceQueryImpl is the initializer for Service query
func NewServiceQueryImpl(on pg.OnboardingDb) *ServiceQueryImpl {
	return &ServiceQueryImpl{
		onboarding: on,
	}
}

// RetrieveFacility  is a repository implementation method for RetrieveFacility
func (q ServiceQueryImpl) RetrieveFacility(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error) {
	return q.onboarding.RetrieveFacility(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode  is a repository implementation method for RetrieveFacilityByMFLCode
func (q ServiceQueryImpl) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	return q.onboarding.RetrieveByFacilityMFLCode(ctx, MFLCode, isActive)
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
