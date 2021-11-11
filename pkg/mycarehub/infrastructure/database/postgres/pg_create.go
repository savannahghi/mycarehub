package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// GetOrCreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *MyCareHubDb) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	if err := facility.Validate(); err != nil {
		return nil, fmt.Errorf("facility input validation failed: %s", err)
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

// SavePin does the actual saving of the users PIN in the database
func (d *MyCareHubDb) SavePin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
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
		// Flavour:   pinData.Flavour,
		Salt: pinData.Salt,
	}

	_, err := d.create.SavePin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}
