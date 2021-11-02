package gorm

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error)
	RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
	GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error)
	ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination domain.FacilityPage) (*domain.FacilityPage, error)
	GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error)
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*Facility, error) {
	var facility Facility
	active := strconv.FormatBool(isActive)
	err := db.DB.Where(&Facility{FacilityID: id, Active: active}).First(&facility).Error
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

// GetFacilities fetches all the healthcare facilities in the platform.
func (db *PGInstance) GetFacilities(ctx context.Context) ([]Facility, error) {
	var facility []Facility
	err := db.DB.Find(&facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query all facilities %v", err)
	}
	return facility, nil
}

// GetUserProfileByPhoneNumber retrieves a user profile using their phonenumber
func (db *PGInstance) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*User, error) {
	var user User
	if err := db.DB.Preload("Contacts", db.DB.Where(&Contact{Contact: phoneNumber})).Find(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user by phonenumber %v: %v", phoneNumber, err)
	}
	return &user, nil
}

// ListFacilities lists all facilities, the results returned are
// from search, and provided filters. They are also paginated
func (db *PGInstance) ListFacilities(ctx context.Context, searchTerm *string, filter []*domain.FiltersParam, pagination domain.FacilityPage) (*domain.FacilityPage, error) {
	var facilities []Facility
	// this will keep track of the results for pagination
	// Count query is unreliable for this since it is returning the count for all rows instead of results
	var resultCount int64

	facilitiesOutput := []domain.Facility{}

	for _, f := range filter {
		err := f.Validate()
		if err != nil {
			return nil, fmt.Errorf("failed to validate filter %v: %v", f.Value, err)
		}
	}

	paginatedFacilities := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:        pagination.Pagination.Limit,
			CurrentPage:  pagination.Pagination.CurrentPage,
			Count:        pagination.Pagination.Count,
			TotalPages:   pagination.Pagination.TotalPages,
			NextPage:     pagination.Pagination.NextPage,
			PreviousPage: pagination.Pagination.PreviousPage,
		},
		Facilities: pagination.Facilities,
	}

	mappedFilterParams := filterParamsToMap(filter)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failed to initialize filter facilities transaction %v", err)
	}

	tx.Where(mappedFilterParams).Find(&facilities)
	resultCount = int64(len(facilities))

	tx.Scopes(paginate(facilities, &paginatedFacilities.Pagination, resultCount, db.DB)).Where(mappedFilterParams).Find(&facilities)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction list facilities transaction%v", err)
	}

	for _, f := range facilities {
		active, err := strconv.ParseBool(f.Active)
		if err != nil {
			return nil, fmt.Errorf("failed to format %s to bool: %v", f.Active, err)
		}
		facility := domain.Facility{
			ID:          f.FacilityID,
			Name:        f.Name,
			Code:        f.Code,
			Active:      active,
			County:      f.County,
			Description: f.Description,
		}
		facilitiesOutput = append(facilitiesOutput, facility)
	}

	pagination.Pagination.Count = paginatedFacilities.Pagination.Count
	pagination.Pagination.TotalPages = paginatedFacilities.Pagination.TotalPages
	pagination.Pagination.Limit = paginatedFacilities.Pagination.Limit
	pagination.Facilities = facilitiesOutput
	pagination.Pagination.NextPage = paginatedFacilities.Pagination.NextPage

	pagination.Pagination.PreviousPage = paginatedFacilities.Pagination.PreviousPage

	return &pagination, nil
}

// GetUserPINByUserID fetches a user's pin using the user ID
func (db *PGInstance) GetUserPINByUserID(ctx context.Context, userID string) (*PINData, error) {
	var pin PINData
	if err := db.DB.Where(&PINData{UserID: userID, IsValid: true}).First(&pin).Error; err != nil {
		return nil, fmt.Errorf("failed to get pin: %v", err)
	}
	return &pin, nil
}
