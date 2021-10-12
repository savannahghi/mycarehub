package postgres

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// CreateFacility ...
func (d *OnboardingDb) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {

	facilityObj := &gorm.Facility{
		Name:        facility.Name,
		Code:        facility.Code,
		Active:      facility.Active,
		County:      facility.Code,
		Description: facility.Description,
	}

	facilitySession, err := d.create.CreateFacility(ctx, facilityObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}
