package gorm

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/savannahghi/feedlib"
)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
	GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*User, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error)
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error) {
	var facility Facility
	active := strconv.FormatBool(isActive)
	err := db.DB.Where(&Facility{FacilityID: id, Active: active}).Find(&facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get facility by ID %v: %v", id, err)
	}
	return &facility, nil
}

// RetrieveFacilityByMFLCode fetches a single facility using MFL Code
func (db *PGInstance) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*Facility, error) {
	var facility Facility
	active := strconv.FormatBool(isActive)
	if err := db.DB.Where(&Facility{Code: MFLCode, Active: active}).First(&facility).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by MFL Code %v and status %v: %v", MFLCode, active, err)
	}
	return &facility, nil
}

// GetUserProfileByUserID fetches a user profile facility using the user ID
func (db *PGInstance) GetUserProfileByUserID(ctx context.Context, userID string, flavour string) (*User, error) {
	var user User
	if err := db.DB.Where(&User{UserID: &userID, Flavour: feedlib.Flavour(flavour)}).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by userID %v: %v", userID, err)
	}
	return &user, nil
}

// GetUserPINByUserID fetches a user profile facility using the user ID
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID}).First(&pin).Error; err != nil {
		return nil, fmt.Errorf("failed to get facility by MFL Code %v: %v", userID, err)
	}
	return &pin, nil
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
