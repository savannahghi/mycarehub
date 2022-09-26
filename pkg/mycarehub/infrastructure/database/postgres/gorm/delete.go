package gorm

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Delete represents all `delete` ops to the database
type Delete interface {
	DeleteFacility(ctx context.Context, mflcode int) (bool, error)
	DeleteUser(ctx context.Context, userID string, clientID *string, staffID *string, flavour feedlib.Flavour) error
	DeleteStaffProfile(ctx context.Context, staffID string) error
	DeleteCommunity(ctx context.Context, communityID string) error
	RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) error
	RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) error
}

// DeleteFacility will do the actual deletion of a facility from the database
// This operation perform HARD deletion
func (db *PGInstance) DeleteFacility(ctx context.Context, mflcode int) (bool, error) {
	if mflcode == 0 {
		return false, fmt.Errorf("MFL code cannot be empty")
	}

	err := db.DB.Where("mfl_code", mflcode).First(&Facility{}).Delete(&Facility{}).Error
	if err != nil {
		return false, fmt.Errorf("an error occurred while deleting: %v", err)
	}

	return true, nil
}

// DeleteStaffProfile will do the actual deletion of a staff profile from the database
func (db *PGInstance) DeleteStaffProfile(ctx context.Context, staffID string) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to initialize staff profile deletion transaction")
	}

	// Get staff identifier
	var staffIdentifiers []StaffIdentifiers
	err := tx.Model(&StaffIdentifiers{}).Where(&StaffIdentifiers{StaffID: &staffID}).Find(&staffIdentifiers).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("an error occurred while fetching staff identifiers: %v", err)
	}

	for _, staffIdentifier := range staffIdentifiers {
		err := tx.Model(&StaffIdentifiers{}).Unscoped().Where(&StaffIdentifiers{StaffID: staffIdentifier.StaffID}).Delete(&StaffIdentifiers{}).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("an error occurred while deleting staff identifiers: %v", err)
		}
	}

	// Get staff contact
	var staffContacts []StaffContacts
	err = tx.Model(&StaffContacts{}).Where(&StaffContacts{StaffID: &staffID}).Find(&staffContacts).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("an error occurred while fetching staff contacts: %v", err)
	}

	for _, staffContact := range staffContacts {
		err = tx.Model(&StaffContacts{}).Unscoped().Where(&StaffContacts{StaffID: staffContact.StaffID}).Delete(&StaffContacts{}).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("an error occurred while deleting staff contacts: %v", err)
		}

		err = tx.Model(&Contact{}).Unscoped().Where(&Contact{ContactID: staffContact.ContactID, Flavour: feedlib.FlavourPro}).First(&Contact{}).Delete(&Contact{}).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("an error occurred while deleting contacts: %v", err)
		}
	}

	// Get staff facilities
	var staffFacilities []StaffFacilities
	err = tx.Model(&StaffFacilities{}).Where(&StaffFacilities{StaffID: &staffID}).Find(&staffFacilities).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("an error occurred while fetching staff facilities: %v", err)
	}

	for _, staffFacility := range staffFacilities {
		err = tx.Model(&StaffFacilities{}).Unscoped().Where(&StaffFacilities{StaffID: staffFacility.StaffID}).First(&StaffFacilities{}).Delete(&StaffFacilities{}).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("an error occurred while deleting staff facilities: %v", err)
		}
	}

	err = tx.Model(&StaffProfile{}).Unscoped().Where("id", staffID).First(&StaffProfile{}).Delete(&StaffProfile{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("an error occurred while deleting staff profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to delete staff failed: %v", err)
	}

	return nil
}

// DeleteUser will do the actual deletion of a user from the database
func (db *PGInstance) DeleteUser(
	ctx context.Context,
	userID string,
	clientID *string,
	staffID *string,
	flavour feedlib.Flavour,
) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to initialize user deletion transaction")
	}

	switch flavour {
	case feedlib.FlavourConsumer:
		err := tx.Unscoped().Preload(clause.Associations).Delete(&Client{ID: clientID}).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete client profile: %w", err)
		}
	case feedlib.FlavourPro:
		err := db.DeleteStaffProfile(ctx, *staffID)
		if err != nil {
			return err
		}
	}

	err := tx.Unscoped().Delete(&User{UserID: &userID}).Error
	if err != nil {
		return fmt.Errorf("an error occurred while deleting user profile: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to delete user profile failed: %w", err)
	}

	return nil
}

// DeleteCommunity deletes the specified community from the database
func (db *PGInstance) DeleteCommunity(ctx context.Context, communityID string) error {
	err := db.DB.Where("id = ?", communityID).Delete(&Community{}).Error
	if err != nil {
		// skip error if not found
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fmt.Errorf("an error occurred while deleting community: %v", err)
	}
	return nil
}

// RemoveFacilitiesFromClientProfile updates the client profile and removes the specified facilities
func (db *PGInstance) RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) error {
	clientFacilities := ClientFacilities{}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, facilityID := range facilities {
		err := tx.Where(&ClientFacilities{
			ClientID:   &clientID,
			FacilityID: &facilityID,
		}).Delete(&clientFacilities).Error
		if err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed transaction commit to update client profile: %w", err)
	}
	return nil
}

// RemoveFacilitiesFromStaffProfile updates the client profile and removes the specified facilities
func (db *PGInstance) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	staffFacilities := StaffFacilities{}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, facilityID := range facilities {
		err := tx.Where(&StaffFacilities{
			StaffID:    &staffID,
			FacilityID: &facilityID,
		}).Delete(&staffFacilities).Error
		if err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed transaction commit to update staff profile: %w", err)
	}
	return nil
}
