package postgres

import "github.com/savannahghi/onboarding-service/pkg/onboarding/domain"

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *OnboardingDb) mapFacilityObjectToDomain(facilityObject *domain.Facility) *domain.Facility {
	if facilityObject == nil {
		return nil
	}

	facilityID := *facilityObject.FacilityID

	return &domain.Facility{
		FacilityID:  &facilityID,
		Name:        facilityObject.Name,
		Code:        facilityObject.Code,
		Active:      facilityObject.Active,
		County:      facilityObject.County,
		Description: facilityObject.Description,
	}
}
