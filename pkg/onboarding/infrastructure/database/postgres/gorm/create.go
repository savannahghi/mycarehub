package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	CreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
}

// CreateFacility ...
func (db *PGInstance) CreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	err := db.DB.Create(facility).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return facility, nil
}
