package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// CreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *OnboardingDb) CreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	if facility.Code == "" {
		return nil, fmt.Errorf("`code` should be defined")
	}

	facilityObj := &gorm.Facility{
		Name:        facility.Name,
		Code:        facility.Code,
		Active:      strconv.FormatBool(facility.Active),
		County:      facility.Code,
		Description: facility.Description,
	}

	facilitySession, err := d.create.CreateFacility(ctx, facilityObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// CollectMetrics is responsible for cretating a representation of metrics data.
func (d *OnboardingDb) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {

	if metric.Type == "" {
		return nil, fmt.Errorf("metric type must be specified")
	}

	metricObj := &gorm.Metric{
		Type:      metric.Type,
		Payload:   metric.Payload,
		Timestamp: metric.Timestamp,
		UID:       metric.UID,
	}

	metricSession, err := d.create.CollectMetrics(ctx, metricObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapMetricObjectToDomain(metricSession), nil
}
