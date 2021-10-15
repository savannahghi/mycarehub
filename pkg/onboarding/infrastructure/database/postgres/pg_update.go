package postgres

import (
	"context"
	"time"
)

// UpdateUserLastSuccessfulLogin update the user with the last login time
func (d *OnboardingDb) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
	return d.update.UpdateUserLastSuccessfulLogin(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin ...
func (d *OnboardingDb) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
	return d.update.UpdateUserLastFailedLogin(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount ...
func (d *OnboardingDb) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error {
	return d.update.UpdateUserFailedLoginCount(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin ...
func (d *OnboardingDb) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
	return d.update.UpdateUserNextAllowedLogin(ctx, userID, nextAllowedLoginTime, flavour)
}
