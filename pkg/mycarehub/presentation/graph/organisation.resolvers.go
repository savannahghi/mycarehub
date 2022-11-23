package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// CreateOrganisation is the resolver for the createOrganisation field.
func (r *mutationResolver) CreateOrganisation(ctx context.Context, input dto.OrganisationInput) (bool, error) {
	return r.mycarehub.Organisation.CreateOrganisation(ctx, input)
}

// DeleteOrganisation is the resolver for the deleteOrganisation field.
func (r *mutationResolver) DeleteOrganisation(ctx context.Context, organisationID string) (bool, error) {
	return r.mycarehub.Organisation.DeleteOrganisation(ctx, organisationID)
}
