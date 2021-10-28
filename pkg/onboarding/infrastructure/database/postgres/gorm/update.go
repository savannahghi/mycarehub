package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
)

// Update represents all `update` ops to the database
type Update interface {
	UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error
	UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateStaffUserProfile(ctx context.Context, userID string, user *User, staff *StaffProfile) (bool, error)
	TransferClient(
		ctx context.Context,
		clientID string,
		originFacilityID string,
		destinationFacilityID string,
		reason enums.TransferReason,
		notes string,
	) (bool, error)
	InvalidatePIN(ctx context.Context, userID string) error
}

// UpdateUserLastSuccessfulLogin updates users last successful login time
func (db *PGInstance) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{LastSuccessfulLogin: &lastLoginTime}).Error
}

// UpdateUserLastFailedLogin updates user's last failed login time
func (db *PGInstance) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{LastFailedLogin: &lastFailedLoginTime}).Error
}

// UpdateUserFailedLoginCount updates users failed login count
func (db *PGInstance) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).
		Updates(&User{FailedLoginCount: failedLoginCount}).Error
}

// UpdateUserNextAllowedLogin updates the users next allowed login time
func (db *PGInstance) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, flavour)
	if err != nil {
		return fmt.Errorf("unable to get user profile by userID when updating: %v", err)
	}

	return db.DB.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{NextAllowedLogin: &nextAllowedLoginTime}).Error
}

// UpdateStaffUserProfile updates the staff user
func (db *PGInstance) UpdateStaffUserProfile(ctx context.Context, userID string, user *User, staff *StaffProfile) (bool, error) {
	userProfile, err := db.GetUserProfileByUserID(ctx, userID, user.Flavour)
	if err != nil {
		return false, fmt.Errorf("unable to get user profile by userID: %v", err)
	}

	staffProfile, err := db.GetStaffProfile(ctx, staff.StaffNumber)
	if err != nil {
		return false, fmt.Errorf("unable to get staff profile by staff number: %v", err)
	}

	// Initialize a database transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failied initialize database transaction %v", err)
	}

	for _, c := range user.Contacts {
		if err := tx.Model(&Contact{}).Where(&Contact{UserID: &userID}).Updates(&Contact{
			Type:    c.Type,
			Contact: c.Contact,
			Active:  c.Active,
			OptedIn: c.OptedIn,
		}).Error; err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to update a user data: %v", err)
		}
	}

	for _, a := range staff.Addresses {
		if err := tx.Model(&Addresses{}).Where(&Addresses{StaffProfileID: staffProfile.StaffProfileID}).Updates(&Addresses{
			Type:       a.Type,
			Text:       a.Text,
			Country:    a.Country,
			PostalCode: a.PostalCode,
			County:     a.County,
			Active:     a.Active,
		}).Error; err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to update a staff profile: %v", err)
		}
	}

	// Update a user profile, then rollback the transaction if it is unsuccessful
	if err := tx.Model(&User{}).Where(&User{UserID: userProfile.UserID, Flavour: userProfile.Flavour}).Updates(&User{
		Languages: user.Languages,
	}).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update a user data: %v", err)
	}

	// update a staff profile, then rollback the transaction if it is unsuccessful
	if err := tx.Model(&StaffProfile{}).Where(&StaffProfile{UserID: userProfile.UserID}).Updates(&StaffProfile{
		DefaultFacilityID: staff.DefaultFacilityID,
	}).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update a staff data: %v", err)
	}

	// try to commit the transactions
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to update a staff profile failed: %v", err)
	}

	return true, nil
}

// TransferClient handles transfer of a client
func (db *PGInstance) TransferClient(
	ctx context.Context,
	clientID string,
	originFacilityID string,
	destinationFacilityID string,
	reason enums.TransferReason,
	notes string,
) (bool, error) {
	if clientID == "" {
		return false, fmt.Errorf("clientID cannot be empty")
	}

	clientProfile, err := db.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		return false, fmt.Errorf("unable to get client profile with the provided client ID: %v", err)
	}

	transferRecord := &ClientTransfer{
		ClientID:              clientID,
		OriginFacilityID:      originFacilityID,
		DestinationFacilityID: destinationFacilityID,
		Reason:                reason,
		Notes:                 notes,
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("unable to initialize database transaction: %v", err)
	}

	//Update client profile with the new facility. Rollback if transaction fails
	if err := tx.Model(&ClientProfile{}).Where(&ClientProfile{UserID: clientProfile.UserID}).Updates(&ClientProfile{
		FacilityID: destinationFacilityID,
	}).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update client profile: %v", err)
	}

	//Create transfer record for the client
	if err := tx.Create(transferRecord).Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to create client transfer record: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to commit the transaction: %v", err)
	}

	return true, nil
}

// InvalidatePIN toggles the valid field of a pin from true to false
func (db *PGInstance) InvalidatePIN(ctx context.Context, userID string) error {
	return db.DB.Model(&PINData{}).Where(&PINData{UserID: userID, IsValid: true}).Select("IsValid").Updates(PINData{IsValid: false}).Error
}
