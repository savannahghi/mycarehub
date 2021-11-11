package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	SavePin(ctx context.Context, pinInput *PINData) (bool, error)
}

// GetOrCreateFacility creates a facility if it doesn't exist or gets it if it exists.
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	if facility == nil {
		return nil, fmt.Errorf("facility must be provided")
	}
	err := db.DB.Create(facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}
	return facility, nil
}

// SavePin does the actual saving of the users PIN in the database
func (db *PGInstance) SavePin(ctx context.Context, pinData *PINData) (bool, error) {
	err := db.DB.Create(pinData).Error

	if err != nil {
		return false, fmt.Errorf("failed to save pin data: %v", err)
	}

	return true, nil
}
