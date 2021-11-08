package gorm

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// Update represents all `update` operations to the database
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *string) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *string) (bool, error)
	UpdateFacility(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error)
}

// ReactivateFacility perfoms the actual re-activation of the facility in the database
func (db *PGInstance) ReactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: "false"}).
		Updates(&Facility{Active: "true"}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// InactivateFacility perfoms the actual inactivation of the facility in the database
func (db *PGInstance) InactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: "true"}).
		Updates(&Facility{Active: "false"}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

// UpdateFacility performs the actual inactivation of the facility in the database
func (db *PGInstance) UpdateFacility(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
	if id == nil {
		return false, fmt.Errorf("id cannot be empty")
	}
	if facilityInput == nil {
		return false, fmt.Errorf("facility input cannot be empty")
	}
	err := db.DB.Model(&Facility{}).Where(&Facility{FacilityID: id}).Updates(&Facility{
		Name:        facilityInput.Name,
		Code:        facilityInput.Code,
		Active:      strconv.FormatBool(facilityInput.Active),
		County:      facilityInput.County,
		Description: facilityInput.Description,
	}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
