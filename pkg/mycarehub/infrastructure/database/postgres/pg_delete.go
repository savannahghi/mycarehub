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
