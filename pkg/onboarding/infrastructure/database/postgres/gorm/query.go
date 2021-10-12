package gorm

import (
	"context"
	"fmt"
	"log"
)

// Query contains all the db query methods
type Query interface {
	Retrieve(id *int64) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
}

// Retrieve fetches a single facility
func (db *PGInstance) Retrieve(id *int64) (*Facility, error) {
	var facility Facility
	if err := db.DB.Where(&Facility{FacilityID: id}).First(&facility).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}

// GetFacilities fetches all the healthcare facilities in the platform.
func (db *PGInstance) GetFacilities(ctx context.Context) ([]Facility, error) {
	var facility []Facility
	facilities := db.DB.Find(&facility).Error
	log.Printf("these are the facilities %v", facilities)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to query all facilities %v", err)
	// }
	log.Printf("these are the facilities %v", facility)
	return facility, nil
}
