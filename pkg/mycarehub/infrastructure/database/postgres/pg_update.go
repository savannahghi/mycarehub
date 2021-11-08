package postgres

import (
	"context"
	"fmt"
)

// InactivateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) InactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.InactivateFacility(ctx, mflCode)
}

// ReactivateFacility changes the status of an active facility from false to true
func (d *MyCareHubDb) ReactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be nil")
	}
	return d.update.ReactivateFacility(ctx, mflCode)
}
