package postgres

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
)

// DeleteFacility does the actual delete of a facility from the database.
func (d *MyCareHubDb) DeleteFacility(ctx context.Context, id int) (bool, error) {
	if id == 0 {
		return false, fmt.Errorf("an error occurred")
	}
	return d.delete.DeleteFacility(ctx, id)
}

// DeleteUser method is used to delete a user from the system
func (d *MyCareHubDb) DeleteUser(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error {
	return d.delete.DeleteUser(ctx, userID, clientID, staffID, flavour)
}

// DeleteStaffProfile is used to delete a staff from the application
func (d *MyCareHubDb) DeleteStaffProfile(ctx context.Context, staffID string) error {
	return d.delete.DeleteStaffProfile(ctx, staffID)
}
