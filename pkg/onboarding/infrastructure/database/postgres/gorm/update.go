package gorm

import (
	"context"
	"fmt"
	"time"
)

// Update represents all `update` ops to the database
type Update interface {
	UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error
	UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error
	UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error
	UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error
}

// UpdateUserLastSuccessfulLogin ...
func (db *PGInstance) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour string) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{LastSuccessfulLogin: &lastLoginTime}).Error
}

// UpdateUserLastFailedLogin ...
func (db *PGInstance) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour string) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{LastFailedLogin: &lastFailedLoginTime}).Error
}

// UpdateUserFailedLoginCount ...
func (db *PGInstance) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour string) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).
		Updates(&User{FailedLoginCount: failedLoginCount}).Error
}

// UpdateUserNextAllowedLogin ...
func (db *PGInstance) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour string) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{NextAllowedLogin: &nextAllowedLoginTime}).Error
}
