package infrastructure

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// Create represents a contract that contains all `create` ops to the database
type Create interface {
	CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}

// ServiceEngagementImpl represents engagement usecases
type ServiceCreate struct {
	Create Create
}

// NewServiceEngagementImpl returns new instance of ServiceEngagementImpl
func NewServicServiceCreateeEngagementImpl() *ServiceCreate {
	return &ServiceCreate{}
}

func (f ServiceCreate) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	return f.Create.CreateFacility(ctx, facility)
}
