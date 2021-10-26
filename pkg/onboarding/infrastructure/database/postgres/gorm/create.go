package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error)
	SavePin(ctx context.Context, pinData *PINData) (bool, error)
	RegisterStaffUser(ctx context.Context, user *User, staff *StaffProfile) (*StaffUserProfile, error)
	RegisterClient(
		ctx context.Context,
		userInput *User,
		clientInput *ClientProfile,
	) (*ClientUserProfile, error)

	AddIdentifier(
		ctx context.Context,
		identifier *Identifier,
	) (*Identifier, error)
}

// GetOrCreateFacility ...
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	err := db.DB.Create(facility).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return facility, nil
}

// SavePin does the actual saving of the users PIN in the database
func (db *PGInstance) SavePin(ctx context.Context, pinData *PINData) (bool, error) {
	err := db.DB.Create(pinData).Error

	if err != nil {
		return false, fmt.Errorf("failed to save pin data: %v", err)
	}

	return true, nil
}

// CollectMetrics takes the collected metrics and saves them in the database.
func (db *PGInstance) CollectMetrics(ctx context.Context, metrics *Metric) (*Metric, error) {
	err := db.DB.Create(metrics).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return metrics, nil
}

// RegisterStaffUser creates both the user profile and the staff profile.
func (db *PGInstance) RegisterStaffUser(ctx context.Context, user *User, staff *StaffProfile) (*StaffUserProfile, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, fmt.Errorf("failied initialize database transaction %v", err)
	}

	// create a user profile, then rollback the transaction if it is unsuccessful
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a user %v", err)
	}

	// assign userID in staff a value due to foreign keys constraint
	staff.UserID = user.UserID

	// Save the staff facilities in a temporary variable to be used when appending associations
	tempFacilities := staff.Facilities

	// Purge the staff facilities so we wont have duplicate facilities in results (We are using tempFacilities)
	staff.Facilities = []*Facility{}

	// create a staff profile, then rollback the transaction if it is unsuccessful
	// Omit the creation of facilities since we don't want to create new facilities when creating staff
	// Associate staff to facilities and append  the temp facilities
	if err := tx.Omit("Facilities").Create(staff).Association("Facilities").Append(tempFacilities); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a staff profile %v", err)
	}

	// TODO: remove manual saving of the associations
	// Appending temp facilities is not saving data in the join table
	// hence this manual save to the join table
	for _, f := range tempFacilities {
		// Prepare the join table with the  facilities and staff IDs
		staffProfileFacility := StaffprofileFacility{
			StaffProfileID: *staff.StaffProfileID,
			FacilityID:     *f.FacilityID,
		}
		// Search for unique together fields in the join table and if they do not exist, insert the values for staffID and facilityID
		var count int64
		tx.Where("staff_profile_id = ? AND facility_id = ?", staff.StaffProfileID, f.FacilityID).Find(&staffProfileFacility).Count(&count)
		if count < 1 {
			err := tx.Create(staffProfileFacility).Error
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create a values in staffprofile - facility join table %v", err)
			}
		}
	}

	// try to commit the transactions
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("transaction commit to create a staff profile failed: %v", err)
	}

	return &StaffUserProfile{
		User:  user,
		Staff: staff,
	}, nil
}

// RegisterClient picks the clients registration details and saves them in the database
func (db *PGInstance) RegisterClient(
	ctx context.Context,
	user *User,
	clientProfile *ClientProfile,
) (*ClientUserProfile, error) {
	// begin a transaction
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a user %v", err)
	}

	clientProfile.UserID = user.UserID

	if err := tx.Create(clientProfile).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a client profile %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("transaction commit to create a staff profile failed: %v", err)
	}

	return &ClientUserProfile{
		User:   user,
		Client: clientProfile,
	}, nil

}

// AddIdentifier saves a client's identifier record to the database
func (db *PGInstance) AddIdentifier(ctx context.Context, identifier *Identifier) (*Identifier, error) {
	if err := db.DB.Create(identifier).Error; err != nil {
		return nil, fmt.Errorf("failed to create identifier: %v", err)
	}
	return identifier, nil
}
