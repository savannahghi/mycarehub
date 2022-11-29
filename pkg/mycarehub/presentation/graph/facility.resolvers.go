package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateFacility is the resolver for the createFacility field.
func (r *mutationResolver) CreateFacility(ctx context.Context, facility dto.FacilityInput, identifier dto.FacilityIdentifierInput) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.GetOrCreateFacility(ctx, &facility, &identifier)
}

// DeleteFacility is the resolver for the deleteFacility field.
func (r *mutationResolver) DeleteFacility(ctx context.Context, identifier dto.FacilityIdentifierInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.DeleteFacility(ctx, &identifier)
}

// ReactivateFacility is the resolver for the reactivateFacility field.
func (r *mutationResolver) ReactivateFacility(ctx context.Context, identifier dto.FacilityIdentifierInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.ReactivateFacility(ctx, &identifier)
}

// InactivateFacility is the resolver for the inactivateFacility field.
func (r *mutationResolver) InactivateFacility(ctx context.Context, identifier dto.FacilityIdentifierInput) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.InactivateFacility(ctx, &identifier)
}

// AddFacilityContact is the resolver for the addFacilityContact field.
func (r *mutationResolver) AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.AddFacilityContact(ctx, facilityID, contact)
}

// AddFacilityToProgram is the resolver for the addFacilityToProgram field.
func (r *mutationResolver) AddFacilityToProgram(ctx context.Context, facilityIDs []string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.AddFacilityToProgram(ctx, facilityIDs)
}

// SearchFacility is the resolver for the searchFacility field.
func (r *queryResolver) SearchFacility(ctx context.Context, searchParameter *string) ([]*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.SearchFacility(ctx, searchParameter)
}

// RetrieveFacility is the resolver for the retrieveFacility field.
func (r *queryResolver) RetrieveFacility(ctx context.Context, id string, active bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacility(ctx, &id, active)
}

// RetrieveFacilityByIdentifier is the resolver for the retrieveFacilityByIdentifier field.
func (r *queryResolver) RetrieveFacilityByIdentifier(ctx context.Context, identifier dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacilityByIdentifier(ctx, &identifier, isActive)
}

// ListFacilities is the resolver for the listFacilities field.
func (r *queryResolver) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	return r.mycarehub.Facility.ListFacilities(ctx, searchTerm, filterInput, &paginationInput)
}
