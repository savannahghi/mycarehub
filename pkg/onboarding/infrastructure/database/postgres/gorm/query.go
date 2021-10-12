package gorm

import (
	"fmt"
)

// Query contains all the db query methods
type Query interface {
	Retrieve(id *int64) (*Facility, error)
}

// Retrieve ...
func (db *PGInstance) Retrieve(id *int64) (*Facility, error) {
	var facility Facility
	if err := db.DB.Where(&Facility{FacilityID: id}).First(&facility).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}
