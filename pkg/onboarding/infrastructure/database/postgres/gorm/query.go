package gorm

import (
	"fmt"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// Query contains all the db query methods
type Query interface {
	Retrieve(id *int64) (*domain.Facility, error)
}

// Retrieve ...
func (db *PGInstance) Retrieve(id *int64) (*domain.Facility, error) {
	var facility domain.Facility
	if err := db.DB.Where(&domain.Facility{FacilityID: id}).First(&facility).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}
