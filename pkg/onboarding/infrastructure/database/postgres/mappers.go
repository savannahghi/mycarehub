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
		ID:          *facilityObject.FacilityID,
		Name:        facilityObject.Name,
		Code:        facilityObject.Code,
		Active:      active,
		County:      facilityObject.County,
		Description: facilityObject.Description,
	}
}

// mapMetricObjectToDomain maps the db metrics to a domain model.
// It searches the database to fetch items specific to the metrics
func (d *OnboardingDb) mapMetricObjectToDomain(metricObject *gorm.Metric) *domain.Metric {
	if metricObject == nil {
		return nil
	}

	return &domain.Metric{
		MetricID:  *metricObject.MetricID,
		Type:      metricObject.Type,
		Payload:   metricObject.Payload,
		Timestamp: metricObject.Timestamp,
		UID:       metricObject.UID,
	}
}
