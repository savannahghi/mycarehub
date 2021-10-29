package postgres

import (
	"strconv"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *OnboardingDb) mapFacilityObjectToDomain(facilityObject *gorm.Facility) *domain.Facility {
	if facilityObject == nil {
		return nil
	}

	active, err := strconv.ParseBool(facilityObject.Active)
	if err != nil {
		return nil
	}

	return &domain.Facility{
		ID:          facilityObject.FacilityID,
		Name:        facilityObject.Name,
		Code:        facilityObject.Code,
		Active:      active,
		County:      facilityObject.County,
		Description: facilityObject.Description,
	}
}
