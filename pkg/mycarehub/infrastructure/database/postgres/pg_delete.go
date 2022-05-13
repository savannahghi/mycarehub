package postgres

import (
	"context"
	"fmt"
)

// DeleteFacility does the actual delete of a facility from the database.
func (d *MyCareHubDb) DeleteFacility(ctx context.Context, id int) (bool, error) {
	if id == 0 {
		return false, fmt.Errorf("an error occurred")
	}
	return d.delete.DeleteFacility(ctx, id)
}

// DeleteClientProfile method is used to delete client from the application
func (d *MyCareHubDb) DeleteClientProfile(ctx context.Context, clientID string) (bool, error) {
	return d.delete.DeleteClientProfile(ctx, clientID)
}

// DeleteUser method is used to delete a user from the system
func (d *MyCareHubDb) DeleteUser(ctx context.Context, userID string) (bool, error) {
	return d.delete.DeleteUser(ctx, userID)
}

// DeleteStaffProfile is used to delete a staff from the application
func (d *MyCareHubDb) DeleteStaffProfile(ctx context.Context, staffID string) (bool, error) {
	return d.delete.DeleteStaffProfile(ctx, staffID)
}
