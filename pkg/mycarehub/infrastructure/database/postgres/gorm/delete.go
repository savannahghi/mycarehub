package gorm

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
)

// Delete represents all `delete` ops to the database
type Delete interface {
	DeleteFacility(ctx context.Context, mflcode int) (bool, error)
	DeleteClientProfile(ctx context.Context, clientID string) (bool, error)
	DeleteUser(ctx context.Context, userID string) (bool, error)
	DeleteStaffProfile(ctx context.Context, staffID string) (bool, error)
}

// DeleteFacility will do the actual deletion of a facility from the database
// This operation perform HARD deletion
func (db *PGInstance) DeleteFacility(ctx context.Context, mflcode int) (bool, error) {
	if mflcode == 0 {
		return false, fmt.Errorf("MFL code cannot be empty")
	}

	err := db.DB.Where("mfl_code", mflcode).First(&Facility{}).Delete(&Facility{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred while deleting: %v", err)
	}

	return true, nil
}

// DeleteClientProfile will do the actual deletion of a client from the database
func (db *PGInstance) DeleteClientProfile(ctx context.Context, clientID string) (bool, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failed to initialize client delete transaction")
	}

	// Get client contacts
	var clientContacts []ClientContacts
	err := tx.Model(&ClientContacts{}).Where(&ClientContacts{ClientID: &clientID}).Find(&clientContacts).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while fetching client contacts: %v", err)
	}

	for _, clientContact := range clientContacts {
		err := tx.Model(&ClientContacts{}).Unscoped().Where(&ClientContacts{ClientID: clientContact.ClientID}).Delete(&ClientContacts{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("an error occurred while deleting client contacts: %v", err)
		}

		err = tx.Model(&Contact{}).Unscoped().Where(&Contact{ContactID: clientContact.ContactID, Flavour: feedlib.FlavourConsumer}).First(&Contact{}).Delete(&Contact{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("an error occurred while deleting contact: %v", err)
		}
	}

	// Get client identifiers
	var clientIdentifiers []ClientIdentifiers
	err = tx.Model(&ClientIdentifiers{}).Where(&ClientIdentifiers{ClientID: &clientID}).Find(&clientIdentifiers).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while fetching client identifiers: %v", err)
	}

	for _, clientIdentifier := range clientIdentifiers {
		err = tx.Model(&ClientIdentifiers{}).Unscoped().Where("identifier_id", clientIdentifier.IdentifierID).First(&ClientIdentifiers{}).Delete(&ClientIdentifiers{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("an error occurred while deleting client identifiers: %v", err)
		}

		err = tx.Model(&Identifier{}).Unscoped().Where("id", clientIdentifier.IdentifierID).First(&Identifier{}).Delete(&Identifier{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("an error occurred while deleting identifier: %v", err)
		}
	}

	err = tx.Model(&Client{}).Unscoped().Where("id", clientID).First(&Client{}).Delete(&Client{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred while deleting client profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to delete client profile failed: %v", err)
	}

	return true, nil
}

// DeleteStaffProfile will do the actual deletion of a staff profile from the database
func (db *PGInstance) DeleteStaffProfile(ctx context.Context, staffID string) (bool, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failed to initialize staff profile deletion transaction")
	}

	// Get staff identifier
	var staffIdentifiers []StaffIdentifiers
	err := tx.Model(&StaffIdentifiers{}).Where(&StaffIdentifiers{StaffID: &staffID}).Find(&staffIdentifiers).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while fetching staff identifiers: %v", err)
	}

	for _, staffIdentifier := range staffIdentifiers {
		err := tx.Model(&StaffIdentifiers{}).Unscoped().Where(&StaffIdentifiers{StaffID: staffIdentifier.StaffID}).Delete(&StaffIdentifiers{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("an error occurred while deleting staff identifiers: %v", err)
		}
	}

	// Get staff contact
	var staffContacts []StaffContacts
	err = tx.Model(&StaffContacts{}).Where(&StaffContacts{StaffID: &staffID}).Find(&staffContacts).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while fetching staff contacts: %v", err)
	}

	for _, staffContact := range staffContacts {
		err = tx.Model(&StaffContacts{}).Unscoped().Where(&StaffContacts{StaffID: staffContact.StaffID}).Delete(&StaffContacts{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("an error occurred while deleting staff contacts: %v", err)
		}

		err = tx.Model(&Contact{}).Unscoped().Where(&Contact{ContactID: staffContact.ContactID, Flavour: feedlib.FlavourPro}).First(&Contact{}).Delete(&Contact{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("an error occurred while deleting contacts: %v", err)
		}
	}

	// Get staff facilities
	var staffFacilities []StaffFacilities
	err = tx.Model(&StaffFacilities{}).Where(&StaffFacilities{StaffID: &staffID}).Find(&staffFacilities).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while fetching staff facilities: %v", err)
	}

	for _, staffFacility := range staffFacilities {
		err = tx.Model(&StaffFacilities{}).Unscoped().Where(&StaffFacilities{StaffID: staffFacility.StaffID}).First(&StaffFacilities{}).Delete(&StaffFacilities{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("an error occurred while deleting staff facilities: %v", err)
		}
	}

	err = tx.Model(&StaffProfile{}).Unscoped().Where("id", staffID).First(&StaffProfile{}).Delete(&StaffProfile{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while deleting staff profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to delete staff failed: %v", err)
	}

	return true, nil
}

// DeleteUser will do the actual deletion of a user from the database
func (db *PGInstance) DeleteUser(ctx context.Context, userID string) (bool, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false, fmt.Errorf("failed to initialize user deletion transaction")
	}
	err := tx.Model(&UserGroups{}).Unscoped().Where(&UserGroups{UserID: &userID}).Delete(&UserGroups{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while deleting user groups: %v", err)
	}

	err = tx.Model(&UserAuthToken{}).Unscoped().Where(&UserAuthToken{UserID: &userID}).Delete(&UserAuthToken{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while deleting user token: %v", err)
	}

	err = tx.Model(&UserPermisssions{}).Unscoped().Where(&UserPermisssions{UserID: &userID}).Delete(&UserPermisssions{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("an error occurred while  user permission: %v", err)
	}

	err = tx.Model(&User{}).Unscoped().Where("id", userID).First(&User{}).Delete(&User{}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred while deleting user profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to delete user profile failed: %v", err)
	}

	return true, nil
}
