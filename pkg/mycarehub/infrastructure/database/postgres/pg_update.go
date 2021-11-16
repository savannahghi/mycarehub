package postgres

import (
	"context"
	"fmt"
	"time"
)

// ReactivateFacility changes the status of an active facility from false to true
func (d *MyCareHubDb) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.ReactivateFacility(ctx, mflCode)
}

// InactivateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.InactivateFacility(ctx, mflCode)
}

// AcceptTerms can be used to accept or review terms of service
func (d *MyCareHubDb) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	if userID == nil || termsID == nil {
		return false, fmt.Errorf("userID or termsID cannot be nil")
	}
	return d.update.AcceptTerms(ctx, userID, termsID)
}

// UpdateUserFailedLoginCount increments a user's failed login count in the event where they fail to
// log into the app e.g when an invalid pin is passed
func (d *MyCareHubDb) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserFailedLoginCount(ctx, userID, failedLoginAttempts)
}

// UpdateUserLastFailedLoginTime updates the failed login time for a user in case an error occurs while logging in
func (d *MyCareHubDb) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserLastFailedLoginTime(ctx, userID)
}

// UpdateUserNextAllowedLoginTime updates the user's next allowed login time. This field is used to check whether we can
// allow a user to log in immediately or wait for some time before retrying the login process.
func (d *MyCareHubDb) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserNextAllowedLoginTime(ctx, userID, nextAllowedLoginTime)
}

// UpdateUserLastSuccessfulLoginTime updates the user's last successful login time to the current time in case a user
// successfully logs into the app
func (d *MyCareHubDb) UpdateUserLastSuccessfulLoginTime(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserLastSuccessfulLoginTime(ctx, userID)
}

// SetNickName is used to set the user's nickname
func (d *MyCareHubDb) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	if userID == nil || nickname == nil {
		return false, fmt.Errorf("userID or nickname cannot be empty ")
	}

	return d.update.SetNickName(ctx, userID, nickname)
}
