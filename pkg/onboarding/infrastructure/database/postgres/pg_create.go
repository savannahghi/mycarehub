package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// GetOrCreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *OnboardingDb) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	if facility.Code == "" {
		return nil, fmt.Errorf("`code` should be defined")
	}

	facilityObj := &gorm.Facility{
		Name:        facility.Name,
		Code:        facility.Code,
		Active:      strconv.FormatBool(facility.Active),
		County:      facility.County,
		Description: facility.Description,
	}

	facilityData, err := d.create.GetOrCreateFacility(ctx, facilityObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilityData), nil
}

// CollectMetrics is responsible for creating a representation of metrics data.
func (d *OnboardingDb) CollectMetrics(ctx context.Context, metric *dto.MetricInput) (*domain.Metric, error) {
	// TODO: make this checks helpers in a different place
	if metric.Type == "" {
		return nil, fmt.Errorf("metric type must be specified")
	}

	metricObj := &gorm.Metric{
		Type:      metric.Type,
		Payload:   metric.Payload,
		Timestamp: metric.Timestamp,
		UID:       metric.UID,
	}

	metricData, err := d.create.CollectMetrics(ctx, metricObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapMetricObjectToDomain(metricData), nil
}

// CreateUser is responsible from creating a representation of a user
func (d *OnboardingDb) CreateUser(
	ctx context.Context,
	input *dto.CreateUserInput,
) (*domain.User, error) {
	// TODO: make this checks helpers in a different place
	if *input.UserType == "" {
		return nil, fmt.Errorf("user type must be specified")
	}
	if *input.Gender == "" {
		return nil, fmt.Errorf("gender must be specified")
	}

	userObject := &gorm.User{
		Username:   *input.FirstName,
		LastName:   *input.LastName,
		MiddleName: input.MiddleName,
		Gender:     *input.Gender,
		UserType:   *input.UserType,
		Active:     input.Active,
	}

	userData, err := d.create.CreateUser(ctx, userObject)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	return d.mapUserObjectToDomain(userData), nil
}
