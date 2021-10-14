package gorm

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// Query contains all the db query methods
type Query interface {
	RetrieveFacility(ctx context.Context, id *uuid.UUID) (*Facility, error)
	GetFacilities(ctx context.Context) ([]Facility, error)
	FindFacility(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error)
}

// RetrieveFacility fetches a single facility
func (db *PGInstance) RetrieveFacility(ctx context.Context, id *uuid.UUID) (*Facility, error) {
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

// FindFacility gets a list of facilities according to the passed filters
// Return an empty list if not found or a list with values if found and an error
func (db *PGInstance) FindFacility(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.FacilityFilterInput, sort []*dto.FacilitySortInput) (*dto.FacilityConnection, error) {
	var (
		facilities []Facility
	)

	// err := db.DB.Preload(
	// 	"Facility", func(db *gorm.DB) *gorm.DB {
	// 		return db.Offset(2).Limit(2)
	// 	},
	// ).Find(&facilities).Error
	err := db.DB.Find(&facilities).Error
	if err != nil {
		return nil, fmt.Errorf("error querying for facilities %v", err)
	}

	type apiResp struct {
		dto.PaginationListResponse

		Results []Facility `json:"results,omitempty"`
	}

	r := apiResp{}

	r.Results = facilities

	startOffset := firebasetools.CreateAndEncodeCursor(r.StartIndex)
	endOffset := firebasetools.CreateAndEncodeCursor(r.EndIndex)
	hasNextPage := r.Next != ""
	hasPreviousPage := r.Previous != ""

	edges := []*dto.FacilityEdge{}
	for pos, f := range r.Results {
		edge := dto.FacilityEdge{
			Node: &domain.Facility{
				ID:          *f.FacilityID,
				Name:        f.Name,
				Code:        f.Code,
				County:      f.County,
				Description: f.Description,
			},
			Cursor: firebasetools.CreateAndEncodeCursor(pos + 1),
		}
		edges = append(edges, &edge)
	}
	pageInfo := &firebasetools.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		StartCursor:     startOffset,
		EndCursor:       endOffset,
	}
	connection := &dto.FacilityConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
	return connection, nil

}
