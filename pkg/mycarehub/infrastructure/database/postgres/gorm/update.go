package gorm

import (
	"context"
	"fmt"
)

// Update represents all `update` operations to the database
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *string) (bool, error)
}

// InactivateFacility perfoms the actual inactivation of the facility in the database
func (db *PGInstance) InactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}
	// facility, err := db.RetrieveFacilityByMFLCode(ctx, *mflCode, true)
	// if err != nil {
	// 	return false, err
	// }

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: "true"}).
		Updates(&Facility{Active: "false"}).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
