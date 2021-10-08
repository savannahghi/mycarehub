package gorm

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error)
}

// CreateFacility ...
func (db *PGInstance) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	err := db.DB.Create(facility).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}
	facilityResp := &domain.Facility{
		Name:        facility.Name,
		Code:        facility.Code,
		Active:      facility.Active,
		County:      facility.County,
		Description: facility.Description,
	}

	return facilityResp, nil
}
