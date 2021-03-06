package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Update represents all `update` operations to the database
type Update interface {
	InactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	ReactivateFacility(ctx context.Context, mflCode *int) (bool, error)
	UpdateFacility(ctx context.Context, facility *Facility, updateData map[string]interface{}) error
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
	SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error)
	CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error)
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
	UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error
	ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificattionStatus string) (bool, error)
	AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error)
	UpdateAppointment(ctx context.Context, appointment *Appointment, updateData map[string]interface{}) (*Appointment, error)
	InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error
	UpdateServiceRequests(ctx context.Context, payload []*ClientServiceRequest) (bool, error)
	UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	UpdateClient(ctx context.Context, client *Client, updates map[string]interface{}) (*Client, error)
	UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error
	UpdateHealthDiary(ctx context.Context, clientHealthDiaryEntry *ClientHealthDiaryEntry, updateData map[string]interface{}) error
	UpdateFailedSecurityQuestionsAnsweringAttempts(ctx context.Context, userID string, failCount int) error
	UpdateUserSurveys(ctx context.Context, survey *UserSurvey, updateData map[string]interface{}) error
	UpdateUser(ctx context.Context, user *User, updateData map[string]interface{}) error
	UpdateNotification(ctx context.Context, notification *Notification, updateData map[string]interface{}) error
	UpdateClientServiceRequest(ctx context.Context, clientServiceRequest *ClientServiceRequest, updateData map[string]interface{}) error
}

// ReactivateFacility performs the actual re-activation of the facility in the database
func (db *PGInstance) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: false}).
		Updates(&Facility{Active: true}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// InactivateFacility perfoms the actual inactivation of the facility in the database
func (db *PGInstance) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("mflCode cannot be empty")
	}

	err := db.DB.Model(&Facility{}).Where(&Facility{Code: *mflCode, Active: true}).
		Updates(&Facility{Active: false}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	return true, nil
}

// AcceptTerms performs the actual modification of the users data for the terms accepted as well as the id of the terms accepted
func (db *PGInstance) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	if userID == nil || termsID == nil {
		return false, fmt.Errorf("userID or termsID cannot be nil")
	}

	if err := db.DB.Model(&User{}).Where(&User{UserID: userID}).
		Updates(&User{TermsAccepted: true, AcceptedTermsID: termsID}).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred while updating the user: %v", err)
	}

	return true, nil
}

// SetNickName is used to set the user's nickname in the database
func (db *PGInstance) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	if userID == nil || nickname == nil {
		return false, fmt.Errorf("userID or nickname cannot be nil")
	}
	err := db.DB.Model(&User{}).Where(&User{UserID: userID}).Updates(&User{Username: *nickname}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to set nickname")
	}

	return true, nil
}

// SetInProgressBy updates the staff assigned to a service request
func (db *PGInstance) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	if requestID == "" || staffID == "" {
		return false, fmt.Errorf("requestID or staffID cannot be empty")
	}
	if err := db.DB.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: &requestID}).Updates(map[string]interface{}{
		"status":            enums.ServiceRequestStatusInProgress,
		"in_progress_by_id": staffID,
		"in_progress_at":    time.Now(),
	}).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update the service request: %v", err)
	}
	return true, nil
}

// CompleteOnboardingTour updates the user's pin change required from true to false
// It also updates the phone_verified, set_pin and set_security_questions to true
// It'll be used to determine the onboarding journey for a user i.e where to redirect a user
// after they log in
func (db *PGInstance) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if !flavour.IsValid() {
		return false, fmt.Errorf("invalid flavour provided")
	}
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID, Flavour: flavour}).Updates(map[string]interface{}{
		"pin_change_required":        false,
		"has_set_pin":                true,
		"has_set_security_questions": true,
		"is_phone_verified":          true,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	return true, nil
}

// InvalidatePIN toggles the valid field of a pin from true to false
func (db *PGInstance) InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")

	}
	err := db.DB.Model(&PINData{}).Where(&PINData{UserID: userID, IsValid: true, Flavour: flavour}).Select("active").Updates(PINData{IsValid: false}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred while invalidating the pin: %v", err)
	}
	return true, nil
}

// UpdateIsCorrectSecurityQuestionResponse updates the is_correct_security_question_response field in the database
func (db *PGInstance) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")

	}
	err := db.DB.Model(&SecurityQuestionResponse{}).Where(&SecurityQuestionResponse{UserID: userID}).Updates(map[string]interface{}{
		"is_correct": isCorrectSecurityQuestionResponse,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("an error occurred while updating the is correct security question response: %v", err)
	}
	return true, nil
}

// UpdateClient updates details for a particular client
func (db *PGInstance) UpdateClient(ctx context.Context, client *Client, updates map[string]interface{}) (*Client, error) {
	updateClient := &Client{}

	if client.ID == nil {
		return nil, fmt.Errorf("client id is required")
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to initialize database transaction %v", err)
	}

	err := tx.Model(updateClient).Where(client).Updates(updates).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to update client profile: %v", err)
	}

	err = tx.First(updateClient, client.ID).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to retrieve client profile: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return nil, fmt.Errorf("failed transaction commit to update client profile: %v", err)
	}

	return updateClient, nil
}

// UpdateClientCaregiver updates the caregiver for a client
func (db *PGInstance) UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	var (
		client Client
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return fmt.Errorf("failed to initialize database transaction %v", err)
	}

	if err := tx.Model(&Client{}).Where(&Client{ID: &caregiverInput.ClientID}).First(&client).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get client: %v", err)
	}

	err := tx.Model(&Caregiver{}).Where(&Caregiver{CaregiverID: client.CaregiverID}).Updates(map[string]interface{}{
		"first_name":     caregiverInput.FirstName,
		"last_name":      caregiverInput.LastName,
		"phone_number":   caregiverInput.PhoneNumber,
		"caregiver_type": caregiverInput.CaregiverType,
	}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update caregiver: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit to update client caregiver failed: %v", err)
	}

	return nil
}

// ResolveStaffServiceRequest resolves the service request for a given staff
func (db *PGInstance) ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error) {
	currentTime := time.Now()

	err := db.DB.Model(&StaffServiceRequest{}).Where(&StaffServiceRequest{ID: serviceRequestID}).Updates(StaffServiceRequest{
		Status:       verificationStatus,
		ResolvedByID: staffID,
		ResolvedAt:   &currentTime,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update staff's service request: %v", err)
	}

	return true, nil
}

// AssignRoles assigns roles to a user
func (db *PGInstance) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	var (
		user User
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to initialize database transaction %v", err)
	}

	err := tx.Model(&User{}).Where(&User{UserID: &userID}).First(&user).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	for _, role := range roles {
		var (
			roleID string
		)

		err := tx.Raw(`SELECT id FROM authority_authorityrole WHERE name = ?`, role.String()).Row().Scan(&roleID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to get authority role: %v", err)
		}

		err = tx.Model(&AuthorityRoleUser{}).Where(&AuthorityRoleUser{UserID: user.UserID, RoleID: &roleID}).FirstOrCreate(&AuthorityRoleUser{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to assign role: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to update user roles failed: %v", err)
	}

	return true, nil
}

// RevokeRoles revokes roles from a user
func (db *PGInstance) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	var user User

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to initialize database transaction %v", err)
	}

	err := tx.Model(&User{}).Where(&User{UserID: &userID}).First(&user).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("failed to get user: %v", err)
	}

	for _, role := range roles {
		var roleID string

		err := tx.Raw(`SELECT id FROM authority_authorityrole WHERE name = ?`, role.String()).Row().Scan(&roleID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to get authority role: %v", err)
		}

		err = tx.Model(&AuthorityRoleUser{}).Where(&AuthorityRoleUser{UserID: user.UserID, RoleID: &roleID}).Delete(&AuthorityRoleUser{}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return false, fmt.Errorf("failed to revoke role: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("transaction commit to update user roles failed: %v", err)
	}

	return true, nil
}

// UpdateAppointment updates the details of an appointment requires the ID or appointment_uuid to be provided
func (db *PGInstance) UpdateAppointment(ctx context.Context, appointment *Appointment, updateData map[string]interface{}) (*Appointment, error) {
	var (
		appointmentToUpdate Appointment
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("failed to initialize database transaction %v", err)
	}

	if appointment.ID != "" {
		err := tx.Model(&Appointment{}).Where(&Appointment{ID: appointment.ID}).First(&appointmentToUpdate).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			tx.Rollback()
			return nil, fmt.Errorf("failed to get appointment: %v", err)
		}
	} else {
		tx.Rollback()
		return nil, fmt.Errorf("failed to get appointment: no ID or appointment_uuid provided")
	}

	err := tx.Model(&Appointment{}).Where(&Appointment{ID: appointmentToUpdate.ID}).Updates(updateData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return nil, fmt.Errorf("failed to update appointment: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return nil, fmt.Errorf("transaction commit to update appointment failed: %v", err)
	}

	return &appointmentToUpdate, nil
}

// InvalidateScreeningToolResponse invalidates a screening tool response
func (db *PGInstance) InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error {
	var (
		client               Client
		screeningToolRespose ScreeningToolsResponse
	)

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failed to initialize database transaction %v", err)
	}

	err := tx.Model(&Client{}).Where(&Client{ID: &clientID}).First(&client).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("failed to get client: %v", err)
	}

	err = tx.Model(&screeningToolRespose).Where(
		&ScreeningToolsResponse{
			ClientID:   clientID,
			QuestionID: questionID,
		}).Updates(map[string]interface{}{"active": false}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("failed to invalidate screening tool response: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("transaction commit to update screening tool response failed: %v", err)
	}
	return nil
}

// UpdateServiceRequests performs and update to the client service requests
func (db *PGInstance) UpdateServiceRequests(ctx context.Context, payload []*ClientServiceRequest) (bool, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("unable to initialize database transaction: %v", err)
	}

	for _, k := range payload {
		err := db.DB.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: k.ID}).Updates(map[string]interface{}{
			"status":         k.Status,
			"in_progress_at": k.InProgressAt,
			"resolved_at":    k.ResolvedAt,
		}).Error
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("unable to update client's service request: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return false, fmt.Errorf("unable to commit transaction: %v", err)
	}

	return true, nil
}

// UpdateUserPinChangeRequiredStatus updates a user's pin change required status
func (db *PGInstance) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID, Flavour: flavour}).Updates(map[string]interface{}{
		"pin_change_required": status,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}
	return nil
}

// UpdateUserPinUpdateRequiredStatus updates a user's pin update required status
func (db *PGInstance) UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	err := db.DB.Model(&User{}).Where(&User{UserID: &userID, Flavour: flavour}).Updates(map[string]interface{}{
		"pin_update_required": status,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}
	return nil
}

// UpdateHealthDiary updates the status of the specified health diary entry
func (db *PGInstance) UpdateHealthDiary(ctx context.Context, clientHealthDiaryEntry *ClientHealthDiaryEntry, updateData map[string]interface{}) error {
	err := db.DB.Model(&ClientHealthDiaryEntry{}).Where(&clientHealthDiaryEntry).Updates(updateData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("unable to update health diary shares status for client: %v", err)
	}

	return nil
}

// UpdateUserSurveys updates the user surveys. The update is performed with regard to the data passed in the survey model.
func (db *PGInstance) UpdateUserSurveys(ctx context.Context, survey *UserSurvey, updateData map[string]interface{}) error {
	if err := db.DB.Model(&UserSurvey{}).Where(&survey).Updates(updateData).Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("an error occurred while updating the user surveys: %w", err)
	}

	return nil
}

// UpdateUser updates the user model
func (db *PGInstance) UpdateUser(ctx context.Context, user *User, updateData map[string]interface{}) error {
	err := db.DB.Model(&User{}).Where(&User{UserID: user.UserID}).Updates(updateData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("unable to update user: %v", err)
	}

	return nil
}

// UpdateFacility updates the facility model
func (db *PGInstance) UpdateFacility(ctx context.Context, facility *Facility, updateData map[string]interface{}) error {
	err := db.DB.Model(&Facility{}).Where(&Facility{FacilityID: facility.FacilityID}).Updates(updateData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("unable to update facility: %v", err)
	}

	return nil
}

//UpdateNotification updates a notification with the new data
func (db *PGInstance) UpdateNotification(ctx context.Context, notification *Notification, updateData map[string]interface{}) error {
	err := db.DB.Model(&Notification{}).Where(&Notification{ID: notification.ID}).Updates(updateData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("unable to update notification: %w", err)
	}

	return nil
}

// UpdateFailedSecurityQuestionsAnsweringAttempts sets the failed security attempts
// the reset happens in an instance where:
// 1. the fail count is less than 3 and the user successfully answers the security questions correctly
// 2. the fail count is 3, the service request for resetting the pin is resolved (client), the user should set the security questions again
// 3. verification of the security questions is unsuccessful
func (db *PGInstance) UpdateFailedSecurityQuestionsAnsweringAttempts(ctx context.Context, userID string, failCount int) error {
	var user User

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("failied initialize database transaction %v", err)
	}

	err := tx.Model(&User{}).Where(&User{UserID: &userID}).First(&user).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("failed to get user: %v", err)
	}

	err = tx.Model(&User{}).Where(&User{UserID: &userID}).Updates(map[string]interface{}{
		"failed_security_count": failCount,
	}).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("failed to update user failed security count: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		helpers.ReportErrorToSentry(err)
		tx.Rollback()
		return fmt.Errorf("transaction commit to update user failed: %v", err)
	}
	return nil
}

// UpdateClientServiceRequest updates the client service request
func (db *PGInstance) UpdateClientServiceRequest(ctx context.Context, clientServiceRequest *ClientServiceRequest, updateData map[string]interface{}) error {
	err := db.DB.Model(&ClientServiceRequest{}).Where(&ClientServiceRequest{ID: clientServiceRequest.ID}).Updates(&updateData).Error
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return fmt.Errorf("unable to update client service request: %v", err)
	}

	return nil
}
