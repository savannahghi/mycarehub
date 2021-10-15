package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error)
	// RegisterStaffUser(user User, profile StaffProfile) (*User, *StaffProfile, error)
	SetUserPIN(ctx context.Context, pinData *PINData) (bool, error)
}

// GetOrCreateFacility ...
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	err := db.DB.Create(facility).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return facility, nil
}

// SetUserPIN does the actual saving of the users PIN in the database
func (db *PGInstance) SetUserPIN(ctx context.Context, pinData *PINData) (bool, error) {
	err := db.DB.Create(pinData).Error

	if err != nil {
		return false, fmt.Errorf("failed to save pin data: %v", err)
	}

	return true, nil
}

// CollectMetrics takes the collected metrics and saves them in the database.
func (db *PGInstance) CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error) {
	err := db.DB.Create(metrics).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return metrics, nil
}
