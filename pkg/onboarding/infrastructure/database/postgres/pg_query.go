package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

//GetFacilities returns a slice of healthcare facilities in the platform.
func (d *OnboardingDb) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	var facility []*domain.Facility
	facilities, err := d.query.GetFacilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	if len(facilities) == 0 {
		return facility, nil
	}
	for _, m := range facilities {
		singleFacility := domain.Facility{
			ID:          *m.FacilityID,
			Name:        m.Name,
			Code:        m.Code,
			Active:      m.Active,
			County:      m.County,
			Description: m.Description,
		}

		facility = append(facility, &singleFacility)
	}

	return facility, nil
}

// RetrieveFacility gets a facility by ID from the database
func (d *OnboardingDb) RetrieveFacility(ctx context.Context, id *uuid.UUID) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacility(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve one facility: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// / FindFacility gets a facility by ID from the database
func (d *OnboardingDb) FindFacility(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error) {
	facilitySession, err := d.query.FindFacility(ctx, pagination, filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to find facility: %v", err)
	}

	return facilitySession, nil
}
