package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateFacility is the resolver for the createFacility field.
func (r *mutationResolver) CreateFacility(ctx context.Context, input dto.FacilityInput) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.GetOrCreateFacility(ctx, &input)
}

// DeleteFacility is the resolver for the deleteFacility field.
func (r *mutationResolver) DeleteFacility(ctx context.Context, mflCode int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.DeleteFacility(ctx, mflCode)
}

// ReactivateFacility is the resolver for the reactivateFacility field.
func (r *mutationResolver) ReactivateFacility(ctx context.Context, mflCode int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.ReactivateFacility(ctx, &mflCode)
}

// InactivateFacility is the resolver for the inactivateFacility field.
func (r *mutationResolver) InactivateFacility(ctx context.Context, mflCode int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.InactivateFacility(ctx, &mflCode)
}

// AddFacilityContact is the resolver for the addFacilityContact field.
func (r *mutationResolver) AddFacilityContact(ctx context.Context, facilityID string, contact string) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.AddFacilityContact(ctx, facilityID, contact)
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

// RetrieveFacilityByMFLCode is the resolver for the retrieveFacilityByMFLCode field.
func (r *queryResolver) RetrieveFacilityByMFLCode(ctx context.Context, mflCode int, isActive bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacilityByMFLCode(ctx, mflCode, isActive)
}

// ListFacilities is the resolver for the listFacilities field.
func (r *queryResolver) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	return r.mycarehub.Facility.ListFacilities(ctx, searchTerm, filterInput, &paginationInput)
}
