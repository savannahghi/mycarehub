package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	RegisterClient(
		ctx context.Context,
		userInput *User,
		clientInput *ClientProfile,
	) (*ClientUserProfile, error)
	SavePin(ctx context.Context, pinData *PINData) (bool, error)
}

// GetOrCreateFacility ...
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	err := db.DB.Create(facility).Error

	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}

	return facility, nil
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

// SavePin saves the pin to the database
func (db *PGInstance) SavePin(ctx context.Context, pinData *PINData) (bool, error) {
	if pinData == nil {
		return false, fmt.Errorf("nil pin data provided")
	}
	err := db.DB.Create(pinData).Error
	if err != nil {
		return false, fmt.Errorf("failed to save pin data: %v", err)
	}
	return true, nil
}
