package postgres

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// DeleteFacility does the actual delete of a facility from the database.
func (d *MyCareHubDb) DeleteFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	identifierObj := &gorm.FacilityIdentifier{
		FacilityID: identifier.FacilityID,
	}
	return d.delete.DeleteFacility(ctx, identifierObj)
}

// DeleteStaffProfile is used to delete a staff from the application
func (d *MyCareHubDb) DeleteStaffProfile(ctx context.Context, staffID string) error {
	return d.delete.DeleteStaffProfile(ctx, staffID)
}

// DeleteCommunity deletes the specified community from the database
func (d *MyCareHubDb) DeleteCommunity(ctx context.Context, communityID string) error {
	return d.delete.DeleteCommunity(ctx, communityID)
}

// RemoveFacilitiesFromClientProfile updates the client profile and removes the specified facilities in their assigned facilities
func (d *MyCareHubDb) RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) error {
	return d.delete.RemoveFacilitiesFromClientProfile(ctx, clientID, facilities)
}

// RemoveFacilitiesFromStaffProfile updates the staff profile and removes the specified facilities in their assigned facilities
func (d *MyCareHubDb) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	return d.delete.RemoveFacilitiesFromStaffProfile(ctx, staffID, facilities)
}

// DeleteOrganisation deletes the specified organisation from the database
func (d *MyCareHubDb) DeleteOrganisation(ctx context.Context, organisation *domain.Organisation) error {
	org := &gorm.Organisation{
		ID: &organisation.ID,
	}
	return d.delete.DeleteOrganisation(ctx, org)
}

// DeleteAccessToken retrieves an access token using the signature
func (d *MyCareHubDb) DeleteAccessToken(ctx context.Context, signature string) error {
	return d.delete.DeleteAccessToken(ctx, signature)
}

// DeleteRefreshToken retrieves a refresh token using the signature
func (d *MyCareHubDb) DeleteRefreshToken(ctx context.Context, signature string) error {
	return d.delete.DeleteRefreshToken(ctx, signature)
}

// DeleteClientProfile method is used to delete a client user from the system
func (d *MyCareHubDb) DeleteClientProfile(ctx context.Context, clientID string, userID *string) error {
	return d.delete.DeleteClientProfile(ctx, clientID, userID)
}
