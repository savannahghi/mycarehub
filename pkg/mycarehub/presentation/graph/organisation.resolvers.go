package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateOrganisation is the resolver for the createOrganisation field.
func (r *mutationResolver) CreateOrganisation(ctx context.Context, input dto.OrganisationInput) (bool, error) {
	return r.mycarehub.Organisation.CreateOrganisation(ctx, input)
}

// DeleteOrganisation is the resolver for the deleteOrganisation field.
func (r *mutationResolver) DeleteOrganisation(ctx context.Context, organisationID string) (bool, error) {
	return r.mycarehub.Organisation.DeleteOrganisation(ctx, organisationID)
}

// ListOrganisations is the resolver for the listOrganisations field.
func (r *queryResolver) ListOrganisations(ctx context.Context, paginationInput dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
	r.checkPreconditions()

	return r.mycarehub.Organisation.ListOrganisations(ctx, &paginationInput)
}

// SearchOrganisations is the resolver for the searchOrganisations field.
func (r *queryResolver) SearchOrganisations(ctx context.Context, searchParameter string) ([]*domain.Organisation, error) {
	r.checkPreconditions()

	return r.mycarehub.Organisation.SearchOrganisation(ctx, searchParameter)
}
