package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// ReactivateFacility changes the status of an active facility from false to true
func (d *MyCareHubDb) ReactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	identifierObj := &gorm.FacilityIdentifier{
		Type:  identifier.Type.String(),
		Value: identifier.Value,
	}
	return d.update.ReactivateFacility(ctx, identifierObj)
}

// InactivateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) InactivateFacility(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
	identifierObj := &gorm.FacilityIdentifier{
		Type:  identifier.Type.String(),
		Value: identifier.Value,
	}
	return d.update.InactivateFacility(ctx, identifierObj)
}

// AcceptTerms can be used to accept or review terms of service
func (d *MyCareHubDb) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	if userID == nil || termsID == nil {
		return false, fmt.Errorf("userID or termsID cannot be nil")
	}
	return d.update.AcceptTerms(ctx, userID, termsID)
}

// CompleteOnboardingTour updates the user's pin change required from true to false. It'll be used to
// determine the onboarding journey for a user.
func (d *MyCareHubDb) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID must be defined")
	}
	return d.update.CompleteOnboardingTour(ctx, userID, flavour)
}

// InvalidatePIN invalidates a pin that is linked to the user profile.
// This is done by toggling the IsValid field to false
func (d *MyCareHubDb) InvalidatePIN(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	return d.update.InvalidatePIN(ctx, userID)
}

// UpdateIsCorrectSecurityQuestionResponse updates the user's security question response
func (d *MyCareHubDb) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	return d.update.UpdateIsCorrectSecurityQuestionResponse(ctx, userID, isCorrectSecurityQuestionResponse)
}

// SetInProgressBy updates the the value of the staff assigned to a service request
func (d *MyCareHubDb) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	return d.update.SetInProgressBy(ctx, requestID, staffID)
}

// ResolveServiceRequest resolves a service request
func (d *MyCareHubDb) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, status string, action []string, comment *string) error {
	serviceRequest, err := d.query.GetClientServiceRequestByID(ctx, *serviceRequestID)
	if err != nil {
		return err
	}

	metadata, err := utils.ConvertJSONStringToMap(serviceRequest.Meta)
	if err != nil {
		return err
	}

	if metadata == nil {
		metadata = map[string]interface{}{
			"comment": comment,
			"action":  action,
		}
	} else {
		metadata["comment"] = comment
		metadata["action"] = action
	}

	newMetaData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	serviceRequestUpdatePayload := map[string]interface{}{
		"status":         status,
		"resolved_by_id": staffID,
		"resolved_at":    time.Now(),
		"meta":           string(newMetaData),
	}

	clientServiceRequest := &gorm.ClientServiceRequest{
		ID: serviceRequestID,
	}

	return d.update.UpdateClientServiceRequest(ctx, clientServiceRequest, serviceRequestUpdatePayload)
}

// ResolveStaffServiceRequest resolves a staff's service request
func (d *MyCareHubDb) ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error) {
	return d.update.ResolveStaffServiceRequest(ctx, staffID, serviceRequestID, verificationStatus)
}

// UpdateAppointment updates an appointment
func (d *MyCareHubDb) UpdateAppointment(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
	ap := &gorm.Appointment{
		ID:         appointment.ID,
		ExternalID: appointment.ExternalID,
	}
	updatedAppointment, err := d.update.UpdateAppointment(ctx, ap, updateData)
	if err != nil {
		return nil, err
	}

	appointmentDate, err := utils.ConvertTimeToScalarDate(updatedAppointment.Date)
	if err != nil {
		return nil, err
	}

	return &domain.Appointment{
		ID:                        updatedAppointment.ID,
		ExternalID:                updatedAppointment.ExternalID,
		Reason:                    updatedAppointment.Reason,
		Date:                      appointmentDate,
		ClientID:                  updatedAppointment.ClientID,
		FacilityID:                updatedAppointment.FacilityID,
		Provider:                  updatedAppointment.Provider,
		HasRescheduledAppointment: updatedAppointment.HasRescheduledAppointment,
	}, nil
}

// UpdateServiceRequests updates service requests
func (d *MyCareHubDb) UpdateServiceRequests(ctx context.Context, payload *domain.UpdateServiceRequestsPayload) (bool, error) {
	var serviceRequests []*gorm.ClientServiceRequest
	for _, k := range payload.ServiceRequests {
		// Update service request
		serviceRequest := &gorm.ClientServiceRequest{
			ID:             &k.ID,
			RequestType:    k.RequestType,
			Status:         k.Status,
			InProgressAt:   k.InProgressAt,
			ResolvedAt:     k.ResolvedAt,
			InProgressByID: k.InProgressBy,
			ResolvedByID:   k.ResolvedBy,
		}

		serviceRequests = append(serviceRequests, serviceRequest)
	}

	return d.update.UpdateServiceRequests(ctx, serviceRequests)
}

// UpdateUserPinChangeRequiredStatus updates a users pin_change_required status. This will
// be used to redirect a user to the change pin page on the app
func (d *MyCareHubDb) UpdateUserPinChangeRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return d.update.UpdateUserPinChangeRequiredStatus(ctx, userID, flavour, status)
}

// UpdateUserPinUpdateRequiredStatus updates a users pin_update_required status. This will
// enable to redirect a user to the change pin page on the app
func (d *MyCareHubDb) UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return d.update.UpdateUserPinUpdateRequiredStatus(ctx, userID, flavour, status)
}

// UpdateHealthDiary updates the status of the specified health diary
func (d *MyCareHubDb) UpdateHealthDiary(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
	healthDiaryEntry := &gorm.ClientHealthDiaryEntry{
		ClientHealthDiaryEntryID: clientHealthDiaryEntry.ID,
		ClientID:                 clientHealthDiaryEntry.ClientID,
	}

	return d.update.UpdateHealthDiary(ctx, healthDiaryEntry, updateData)
}

// UpdateClient updates the client details for a particular client
func (d *MyCareHubDb) UpdateClient(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
	c, err := d.update.UpdateClient(ctx, &gorm.Client{ID: client.ID}, updates)
	if err != nil {
		return nil, err
	}
	var clientList []enums.ClientType
	for _, k := range c.ClientTypes {
		clientList = append(clientList, enums.ClientType(k))
	}

	clientFacilities, _, err := d.GetClientFacilities(ctx, dto.ClientFacilityInput{ClientID: client.ID}, nil)
	if err != nil {
		return nil, err
	}
	facilitiesMap := make(map[string]string)

	for _, f := range clientFacilities {
		facilitiesMap[*f.ID] = f.Name
	}

	return &domain.ClientProfile{
		ID:                      c.ID,
		Active:                  c.Active,
		ClientTypes:             clientList,
		UserID:                  *c.UserID,
		TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
		FHIRPatientID:           c.FHIRPatientID,
		HealthRecordID:          c.HealthRecordID,
		ClientCounselled:        c.ClientCounselled,
		OrganisationID:          c.OrganisationID,
		DefaultFacility: &domain.Facility{
			ID:   &c.FacilityID,
			Name: facilitiesMap[c.FacilityID],
		},
		Facilities: clientFacilities,
	}, nil
}

// UpdateFailedSecurityQuestionsAnsweringAttempts resets the failed attempts for answered security questions
func (d *MyCareHubDb) UpdateFailedSecurityQuestionsAnsweringAttempts(ctx context.Context, userID string, failCount int) error {
	return d.update.UpdateFailedSecurityQuestionsAnsweringAttempts(ctx, userID, failCount)
}

// UpdateUser updates the user details
func (d *MyCareHubDb) UpdateUser(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
	userPayload := &gorm.User{
		UserID: user.ID,
	}
	return d.update.UpdateUser(ctx, userPayload, updateData)
}

// UpdateUserSurveys updates the user surveys. The update is performed with regard to the data passed in the survey model.
func (d *MyCareHubDb) UpdateUserSurveys(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error {
	surveyPayload := &gorm.UserSurvey{
		LinkID:    survey.LinkID,
		ProjectID: survey.ProjectID,
		FormID:    survey.FormID,
	}
	return d.update.UpdateUserSurveys(ctx, surveyPayload, updateData)
}

// UpdateFacility updates the facility with the provided facility details
func (d *MyCareHubDb) UpdateFacility(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
	gormFacility := &gorm.Facility{
		FacilityID: facility.ID,
	}

	return d.update.UpdateFacility(ctx, gormFacility, updateData)
}

// UpdateNotification updates the notification with the provided notification details
func (d *MyCareHubDb) UpdateNotification(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error {
	gormNotification := &gorm.Notification{
		ID: notification.ID,
	}

	return d.update.UpdateNotification(ctx, gormNotification, updateData)
}

// UpdateClientServiceRequest updates the service request with the provided service request details
func (d *MyCareHubDb) UpdateClientServiceRequest(ctx context.Context, clientServiceRequest *domain.ServiceRequest, updateData map[string]interface{}) error {
	gormServiceRequest := &gorm.ClientServiceRequest{
		ID: &clientServiceRequest.ID,
	}

	return d.update.UpdateClientServiceRequest(ctx, gormServiceRequest, updateData)
}

// UpdateStaff updates the staff details for a particular staff
func (d *MyCareHubDb) UpdateStaff(ctx context.Context, staff *domain.StaffProfile, updates map[string]interface{}) error {
	_, err := d.update.UpdateStaff(ctx, &gorm.StaffProfile{ID: staff.ID}, updates)
	if err != nil {
		return err
	}
	return nil
}

// AddFacilitiesToStaffProfile updates the current facility list of a client
func (d *MyCareHubDb) AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) error {
	return d.update.AddFacilitiesToStaffProfile(ctx, staffID, facilities)
}

// AddFacilitiesToClientProfile updates the current facility list of a client
func (d *MyCareHubDb) AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) error {
	return d.update.AddFacilitiesToClientProfile(ctx, clientID, facilities)
}

// UpdateCaregiverClient updates the caregiver client details for either the caregiver or client.
func (d *MyCareHubDb) UpdateCaregiverClient(ctx context.Context, caregiverClient *domain.CaregiverClient, updateData map[string]interface{}) error {
	gormCaregiverClient := &gorm.CaregiverClient{
		ClientID:    caregiverClient.ClientID,
		CaregiverID: caregiverClient.CaregiverID,
	}

	return d.update.UpdateCaregiverClient(ctx, gormCaregiverClient, updateData)
}

// UpdateCaregiver updates the caregiver profile
func (d *MyCareHubDb) UpdateCaregiver(ctx context.Context, caregiver *domain.CaregiverProfile, updates map[string]interface{}) error {
	gormCaregiver := &gorm.Caregiver{
		ID: caregiver.ID,
	}

	return d.update.UpdateCaregiver(ctx, gormCaregiver, updates)
}

// UpdateUserContact is used to updates the user's contact details
func (d *MyCareHubDb) UpdateUserContact(ctx context.Context, contact *domain.Contact, updateData map[string]interface{}) error {
	gormContact := &gorm.Contact{
		Type:   contact.ContactType,
		Value:  contact.ContactValue,
		UserID: contact.UserID,
	}
	return d.update.UpdateUserContact(ctx, gormContact, updateData)
}

// UpdateClientIdentifier updates the client identifier details for a particular client
func (d *MyCareHubDb) UpdateClientIdentifier(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error {
	return d.update.UpdateClientIdentifier(ctx, clientID, identifierType, identifierValue, programID)
}

// UpdateProgram update the details of a particular program
func (d *MyCareHubDb) UpdateProgram(ctx context.Context, program *domain.Program, updateData map[string]interface{}) error {
	gormProgram := &gorm.Program{
		ID: program.ID,
	}

	return d.update.UpdateProgram(ctx, gormProgram, updateData)
}

// UpdateAuthorizationCode updates the details of a given code
func (d *MyCareHubDb) UpdateAuthorizationCode(ctx context.Context, code *domain.AuthorizationCode, updateData map[string]interface{}) error {
	authCode := &gorm.AuthorizationCode{
		ID: code.ID,
	}

	return d.update.UpdateAuthorizationCode(ctx, authCode, updateData)
}

// UpdateAccessToken updates the details of a given access token
func (d *MyCareHubDb) UpdateAccessToken(ctx context.Context, token *domain.AccessToken, updateData map[string]interface{}) error {
	authCode := &gorm.AccessToken{
		ID: token.ID,
	}

	return d.update.UpdateAccessToken(ctx, authCode, updateData)
}

// UpdateRefreshToken updates the details of a given refresh token
func (d *MyCareHubDb) UpdateRefreshToken(ctx context.Context, token *domain.RefreshToken, updateData map[string]interface{}) error {
	authCode := &gorm.RefreshToken{
		ID: token.ID,
	}

	return d.update.UpdateRefreshToken(ctx, authCode, updateData)
}

// UpdateBooking updates the booking model given the models data and the update data
func (d *MyCareHubDb) UpdateBooking(ctx context.Context, booking *domain.Booking, updateData map[string]interface{}) error {
	updatePayload := &gorm.Booking{
		ID:               booking.ID,
		ProgramID:        booking.ProgramID,
		VerificationCode: booking.VerificationCode,
	}

	return d.update.UpdateBooking(ctx, updatePayload, updateData)
}
