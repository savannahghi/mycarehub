package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error)
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

func (db *PGInstance) CreateUser(ctx context.Context, user *User) (*User, error) {
	err := db.DB.Create(user).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create user %v", err)
	}
	return user, nil
}
