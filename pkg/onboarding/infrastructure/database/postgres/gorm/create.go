package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error)
	RegisterStaffUser(ctx context.Context, user User, profile StaffProfile) (*StaffUserProfileTable, error)
}

// GetOrCreateFacility ...
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	err := db.DB.Create(facility).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return facility, nil
}

// CollectMetrics takes the collected metrics and saves them in the database.
func (db *PGInstance) CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error) {
	err := db.DB.Create(metrics).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return metrics, nil
}

// RegisterStaffUser creates a staff user profile
func (db *PGInstance) RegisterStaffUser(ctx context.Context, user User, profile StaffProfile) (*StaffUserProfileTable, error) {
	// Initialize a database transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failied initialize database transaction %v", err)
	}

	// create a user profile, then rollback the transaction if it is unsuccessful
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a user %v", err)
	}

	// create a staff profile, then rollback the transaction if it is unsuccessful
	if err := tx.Create(profile).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a staff profile %v", err)
	}

	// TODO: assign a facility for the first time registration

	// try to commit the successful transactions
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("transaction commit to create a staff profile failed: %v", err)
	}

	return &StaffUserProfileTable{
		User:         user,
		StaffProfile: profile,
	}, nil
}
