package gorm

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"gorm.io/gorm"
)

// Create contains all the methods used to perform a create operation in DB
type Create interface {
	GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error)
	SaveTemporaryUserPin(ctx context.Context, pinPayload *PINData) (bool, error)
	SavePin(ctx context.Context, pinData *PINData) (bool, error)
	SaveOTP(ctx context.Context, otpInput *UserOTP) error
	SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*SecurityQuestionResponse) error
	CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *ClientHealthDiaryEntry) error
	CreateServiceRequest(ctx context.Context, serviceRequestInput *ClientServiceRequest) error
	CreateClientCaregiver(ctx context.Context, clientID string, clientCaregiver *Caregiver) error
	CreateCommunity(ctx context.Context, community *Community) (*Community, error)
	CreateRelatedPerson(ctx context.Context, person *RelatedPerson) error
	CreateContact(ctx context.Context, contact *Contact) error
	CreateAppointment(ctx context.Context, appointment *Appointment) error
	AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*ScreeningToolsResponse) error
}

// GetOrCreateFacility is used to get or create a facility
func (db *PGInstance) GetOrCreateFacility(ctx context.Context, facility *Facility) (*Facility, error) {
	if facility == nil {
		return nil, fmt.Errorf("facility must be provided")
	}
	err := db.DB.Create(facility).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to save a pin: %v", err)
	}
	return true, nil
}

// SavePin saves the pin to the database
func (db *PGInstance) SavePin(ctx context.Context, pinData *PINData) (bool, error) {
	err := db.DB.Create(pinData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
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
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to update OTP data: %v", err)
	}

	//Save the OTP by setting valid to true
	err = db.DB.Create(otpInput).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to save otp data")
	}
	return nil
}

// SaveSecurityQuestionResponse saves the security question response to the database if it does not exist,
// otherwise it updates the existing one
func (db *PGInstance) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*SecurityQuestionResponse) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed initialize database transaction %v", err)
	}
	for _, questionResponse := range securityQuestionResponse {
		SaveSecurityQuestionResponseUpdatePayload := &SecurityQuestionResponse{
			Response: questionResponse.Response,
		}
		err := tx.Model(&SecurityQuestionResponse{}).Where(&SecurityQuestionResponse{UserID: questionResponse.UserID, QuestionID: questionResponse.QuestionID}).First(&questionResponse).Error
		if err == nil {
			helpers.ReportErrorToSentry(err)
			err := tx.Model(&SecurityQuestionResponse{}).Where(&SecurityQuestionResponse{UserID: questionResponse.UserID, QuestionID: questionResponse.QuestionID}).Updates(&SaveSecurityQuestionResponseUpdatePayload).Error
			if err != nil {
				helpers.ReportErrorToSentry(err)
				tx.Rollback()
				return fmt.Errorf("failed to update security question response data: %v", err)
			}
		} else if err == gorm.ErrRecordNotFound {
			err = tx.Create(&questionResponse).Error
			if err != nil {
				helpers.ReportErrorToSentry(err)
				tx.Rollback()
				return fmt.Errorf("failed to create security question response data: %v", err)
			}
		} else {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return fmt.Errorf("failed to get security question response data: %v", err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create/update security question responses failed: %v", err)
	}

	return nil
}

// CreateHealthDiaryEntry records the health diary entries from a client. This is necessary for engagement with clients
// on a day-by-day basis
func (db *PGInstance) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *ClientHealthDiaryEntry) error {
	tx := db.DB.Begin()

	err := tx.Create(healthDiaryInput).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

// CreateServiceRequest creates a service request entry into the database. This step is reached only IF the client is
// in a VERY_BAD mood. We get this mood from the mood scale provided by the front end.
// This operation is done within a transaction to prevent a situation where a health entry can be recorded
// but a service request is not successfully created.
func (db *PGInstance) CreateServiceRequest(
	ctx context.Context,
	serviceRequestInput *ClientServiceRequest,
) error {
	tx := db.DB.Begin()

	err := tx.Create(serviceRequestInput).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// CreateClientCaregiver is used to create a caregiver
func (db *PGInstance) CreateClientCaregiver(ctx context.Context, clientID string, clientCaregiver *Caregiver) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed initialize database transaction %v", err)
	}

	err := tx.Create(clientCaregiver).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create caregiver: %v", err)
	}

	err = tx.Model(&Client{}).Where(&Client{ID: &clientID}).Updates(map[string]interface{}{"caregiver_id": clientCaregiver.CaregiverID}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update client with caregiver id: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create caregiver failed: %v", err)
	}

	return nil

}

// CreateCommunity creates a channel in the database
func (db *PGInstance) CreateCommunity(ctx context.Context, community *Community) (*Community, error) {
	err := db.DB.Create(community).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to create a community: %v", err)
	}
	return community, nil
}

// CreateRelatedPerson creates a related person in the database
func (db *PGInstance) CreateRelatedPerson(ctx context.Context, person *RelatedPerson) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(person).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create related person: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create related person failed: %v", err)
	}

	return nil
}

// CreateContact creates a person's contact in the database
func (db *PGInstance) CreateContact(ctx context.Context, contact *Contact) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(contact).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create contact: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create contact failed: %v", err)
	}

	return nil
}

// CreateAppointment creates an appointment in the database
func (db *PGInstance) CreateAppointment(ctx context.Context, appointment *Appointment) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(appointment).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create an appointment: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create an appointment failed: %v", err)
	}

	return nil
}

// AnswerScreeningToolQuestions answers the screening tool questions
func (db *PGInstance) AnswerScreeningToolQuestions(ctx context.Context, screeningToolResponses []*ScreeningToolsResponse) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, response := range screeningToolResponses {
		err := tx.Create(response).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create screening tool response: %v", err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to create screening tool responses failed: %v", err)
	}

	return nil
}
