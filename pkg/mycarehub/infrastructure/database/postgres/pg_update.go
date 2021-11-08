package postgres

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// ReactivateFacility changes the status of an active facility from false to true
func (d *MyCareHubDb) ReactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.ReactivateFacility(ctx, mflCode)
}

// InactivateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) InactivateFacility(ctx context.Context, mflCode *string) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.InactivateFacility(ctx, mflCode)
}

// UpdateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) UpdateFacility(ctx context.Context, id *string, facilityInput *dto.FacilityInput) (bool, error) {
	if id == nil {
		return false, fmt.Errorf("facility's id Code cannot be empty")
	}
	return d.update.UpdateFacility(ctx, id, facilityInput)
}
