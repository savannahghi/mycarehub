package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

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

	facilitySession, err := d.create.GetOrCreateFacility(ctx, facilityObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// SetUserPIN does the actual saving of the users PIN in the database
func (d *OnboardingDb) SetUserPIN(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	if pinData.UserID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	pinObj := &gorm.PINData{
		Base:      gorm.Base{},
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		IsValid:   pinData.IsValid,
		Flavour:   pinData.Flavour,
		Salt:      pinData.Salt,
	}

	_, err := d.create.SetUserPIN(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to set user pin: %v", err)
	}

	return true, nil
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

// RegisterStaffUser creates both the user profile and the staff profile.
func (d *OnboardingDb) RegisterStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	if staff.DefaultFacilityID == nil {
		return nil, fmt.Errorf("expected default facility ID to be provided")
	}
	userObject := &gorm.User{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		Flavour:     user.Flavour,
	}

	staffObject := &gorm.StaffProfile{
		StaffNumber:       staff.StaffNumber,
		DefaultFacilityID: staff.DefaultFacilityID,
	}

	userStaffProfile, err := d.create.RegisterStaffUser(ctx, userObject, staffObject)
	if err != nil {

		return nil, fmt.Errorf("failed to create user session %v", err)
	}

	return d.mapRegisterStaffObjectToDomain(userStaffProfile), nil

}
