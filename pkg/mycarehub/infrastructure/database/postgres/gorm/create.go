package gorm

import (
	"context"
	"fmt"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	SaveTemporaryUserPin(ctx context.Context, pinPayload *PINData) (bool, error)
	SavePin(ctx context.Context, pinData *PINData) (bool, error)
	SaveOTP(ctx context.Context, otpInput *UserOTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse *SecurityQuestionResponse) error
}

// GetOrCreateFacility is used to get or create a facility
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	if facility == nil {
		return nil, fmt.Errorf("facility must be provided")
	}
	err := db.DB.Create(facility).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create a facility: %v", err)
	}
	return facility, nil
}

// SaveTemporaryUserPin is used to save a temporary user pin
func (db *PGInstance) SaveTemporaryUserPin(ctx context.Context, pinPayload *PINData) (bool, error) {
	if pinPayload == nil {
		return false, fmt.Errorf("pinPayload must be provided")
	}
	err := db.DB.Create(pinPayload).Error
	if err != nil {
		return false, fmt.Errorf("failed to save a pin: %v", err)
	}
	return true, nil
}

// SavePin saves the pin to the database
func (db *PGInstance) SavePin(ctx context.Context, pinData *PINData) (bool, error) {
	err := db.DB.Create(pinData).Error

	if err != nil {
		return false, fmt.Errorf("failed to save pin data: %v", err)
	}

	return true, nil
}

// SaveOTP saves the generated otp to the database
func (db *PGInstance) SaveOTP(ctx context.Context, otpInput *UserOTP) error {
	//Invalidate other OTPs before saving the new OTP by setting valid to false
	if otpInput.PhoneNumber == "" || !otpInput.Flavour.IsValid() {
		return fmt.Errorf("phone number cannot be empty")
	}

	err := db.DB.Model(&UserOTP{}).Where(&UserOTP{PhoneNumber: otpInput.PhoneNumber, Flavour: otpInput.Flavour}).
		Updates(map[string]interface{}{"is_valid": false}).Error
	if err != nil {
		return fmt.Errorf("failed to update OTP data: %v", err)
	}

	//Save the OTP by setting valid to true
	err = db.DB.Create(otpInput).Error
	if err != nil {
		return fmt.Errorf("failed to save otp data")
	}
	return nil
}

// SaveSecurityQuestionResponse saves the security question response to the database
func (db *PGInstance) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse *SecurityQuestionResponse) error {
	err := db.DB.Create(securityQuestionResponse).Error
	if err != nil {
		return fmt.Errorf("failed to save security question response data")
	}
	return nil
}
