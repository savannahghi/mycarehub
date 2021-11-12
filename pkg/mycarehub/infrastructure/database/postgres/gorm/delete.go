package gorm

import (
	"context"
	"fmt"
)

// Delete represents all `delete` ops to the database
type Delete interface {
	DeleteFacility(ctx context.Context, mflcode int) (bool, error)
}

// DeleteFacility will do the actual deletion of a facility from the database
// This operation perform HARD deletion
func (db *PGInstance) DeleteFacility(ctx context.Context, mflcode int) (bool, error) {
	if mflcode == 0 {
		return false, fmt.Errorf("MFL code cannot be empty")
	}
	err := db.DB.Where("mfl_code", mflcode).Delete(&Facility{}).Error
	if err != nil {
		return false, fmt.Errorf("an error occurred while deleting: %v", err)
	}

	return true, nil
}
