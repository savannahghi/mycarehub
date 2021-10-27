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
	GetOrCreateStaffUser(ctx context.Context, user *User, staff *StaffProfile) (*StaffUserProfile, error)
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

// GetOrCreateStaffUser creates both the user profile and the staff profile or gets if exists.
func (db *PGInstance) GetOrCreateStaffUser(ctx context.Context, user *User, staff *StaffProfile) (*StaffUserProfile, error) {
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
		return nil, fmt.Errorf("failed to create a staff user: %v", err)
	}

	// assign userID in staff a value due to foreign keys constraint
	staff.UserID = user.UserID

	// create a staff profile, then rollback the transaction if it is unsuccessful
	if err := tx.Create(staff).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a staff profile: %v", err)
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
		return nil, fmt.Errorf("register client transaction failed: %v", err)
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create a client user: %v", err)
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
