package postgres

import (
	"context"
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
func (d *MyCareHubDb) ReactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.ReactivateFacility(ctx, mflCode)
}

// InactivateFacility changes the status of an active facility from true to false
func (d *MyCareHubDb) InactivateFacility(ctx context.Context, mflCode *int) (bool, error) {
	if mflCode == nil {
		return false, fmt.Errorf("facility's MFL Code cannot be empty")
	}
	return d.update.InactivateFacility(ctx, mflCode)
}

// AcceptTerms can be used to accept or review terms of service
func (d *MyCareHubDb) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	if userID == nil || termsID == nil {
		return false, fmt.Errorf("userID or termsID cannot be nil")
	}
	return d.update.AcceptTerms(ctx, userID, termsID)
}

// UpdateUserFailedLoginCount increments a user's failed login count in the event where they fail to
// log into the app e.g when an invalid pin is passed
func (d *MyCareHubDb) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginAttempts int) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserFailedLoginCount(ctx, userID, failedLoginAttempts)
}

// UpdateUserLastFailedLoginTime updates the failed login time for a user in case an error occurs while logging in
func (d *MyCareHubDb) UpdateUserLastFailedLoginTime(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserLastFailedLoginTime(ctx, userID)
}

// UpdateUserNextAllowedLoginTime updates the user's next allowed login time. This field is used to check whether we can
// allow a user to log in immediately or wait for some time before retrying the login process.
func (d *MyCareHubDb) UpdateUserNextAllowedLoginTime(ctx context.Context, userID string, nextAllowedLoginTime time.Time) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserNextAllowedLoginTime(ctx, userID, nextAllowedLoginTime)
}

// UpdateUserProfileAfterLoginSuccess updates the user's last successful login time to the current time in case a user
// successfully logs into the app
func (d *MyCareHubDb) UpdateUserProfileAfterLoginSuccess(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("userID must be defined")
	}
	return d.update.UpdateUserProfileAfterLoginSuccess(ctx, userID)
}

// SetNickName is used to set the user's nickname
func (d *MyCareHubDb) SetNickName(ctx context.Context, userID *string, nickname *string) (bool, error) {
	if userID == nil || nickname == nil {
		return false, fmt.Errorf("userID or nickname cannot be empty ")
	}

	return d.update.SetNickName(ctx, userID, nickname)
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
func (d *MyCareHubDb) InvalidatePIN(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	return d.update.InvalidatePIN(ctx, userID, flavour)
}

// UpdateIsCorrectSecurityQuestionResponse updates the user's security question response
func (d *MyCareHubDb) UpdateIsCorrectSecurityQuestionResponse(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
	if userID == "" {
		return false, fmt.Errorf("userID cannot be empty")
	}
	return d.update.UpdateIsCorrectSecurityQuestionResponse(ctx, userID, isCorrectSecurityQuestionResponse)
}

// ShareContent updates content share count
func (d *MyCareHubDb) ShareContent(ctx context.Context, input dto.ShareContentInput) (bool, error) {
	if input.Validate() != nil {
		return false, fmt.Errorf("input cannot be empty")
	}
	return d.update.ShareContent(ctx, input)
}

//BookmarkContent updates the user's bookmark status for a content
func (d *MyCareHubDb) BookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if contentID == 0 || userID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}
	return d.update.BookmarkContent(ctx, userID, contentID)
}

// UnBookmarkContent removes the bookmark for a given user
func (d *MyCareHubDb) UnBookmarkContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if contentID == 0 || userID == "" {
		return false, fmt.Errorf("contentID or userID cannot be nil")
	}
	return d.update.UnBookmarkContent(ctx, userID, contentID)
}

// LikeContent increments the number of likes for a particular content
func (d *MyCareHubDb) LikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID or contentID cannot be empty")
	}

	return d.update.LikeContent(ctx, userID, contentID)
}

// UnlikeContent decrements the number of likes for a particular content
func (d *MyCareHubDb) UnlikeContent(ctx context.Context, userID string, contentID int) (bool, error) {
	if userID == "" || contentID == 0 {
		return false, fmt.Errorf("userID or contentID cannot be empty")
	}

	return d.update.UnlikeContent(ctx, userID, contentID)
}

// SetInProgressBy updates the the value of the staff assigned to a service request
func (d *MyCareHubDb) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	return d.update.SetInProgressBy(ctx, requestID, staffID)
}

// ViewContent gets a content item and updates the view count
func (d *MyCareHubDb) ViewContent(ctx context.Context, userID string, contentID int) (bool, error) {
	return d.update.ViewContent(ctx, userID, contentID)
}

// UpdateClientCaregiver updates the caregiver for a client
func (d *MyCareHubDb) UpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) error {
	return d.update.UpdateClientCaregiver(ctx, caregiverInput)
}

// ResolveServiceRequest resolves a service request
func (d *MyCareHubDb) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, status string) (bool, error) {
	return d.update.ResolveServiceRequest(ctx, staffID, serviceRequestID, status)
}

// ResolveStaffServiceRequest resolves a staff's service request
func (d *MyCareHubDb) ResolveStaffServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, verificationStatus string) (bool, error) {
	return d.update.ResolveStaffServiceRequest(ctx, staffID, serviceRequestID, verificationStatus)
}

// AssignRoles assigns roles to a user
func (d *MyCareHubDb) AssignRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return d.update.AssignRoles(ctx, userID, roles)
}

// RevokeRoles revokes roles from a user
func (d *MyCareHubDb) RevokeRoles(ctx context.Context, userID string, roles []enums.UserRoleType) (bool, error) {
	return d.update.RevokeRoles(ctx, userID, roles)
}

// InvalidateScreeningToolResponse invalidates a screening tool response
func (d *MyCareHubDb) InvalidateScreeningToolResponse(ctx context.Context, clientID string, questionID string) error {
	return d.update.InvalidateScreeningToolResponse(ctx, clientID, questionID)
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

// UpdateUserActiveStatus updates a user's `active` status. It will be used to opt out/in a user
func (d *MyCareHubDb) UpdateUserActiveStatus(ctx context.Context, userID string, flavour feedlib.Flavour, active bool) error {
	return d.update.UpdateUserActiveStatus(ctx, userID, flavour, active)
}

// UpdateUserPinUpdateRequiredStatus updates a users pin_update_required status. This will
// enable to redirect a user to the change pin page on the app
func (d *MyCareHubDb) UpdateUserPinUpdateRequiredStatus(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
	return d.update.UpdateUserPinUpdateRequiredStatus(ctx, userID, flavour, status)
}

// UpdateHealthDiary updates the status of the specified health diary
func (d *MyCareHubDb) UpdateHealthDiary(ctx context.Context, payload *gorm.ClientHealthDiaryEntry) (bool, error) {
	return d.update.UpdateHealthDiary(ctx, payload)
}

// UpdateClient updates the client details for a particular client
func (d *MyCareHubDb) UpdateClient(ctx context.Context, client *domain.ClientProfile, updates map[string]interface{}) (*domain.ClientProfile, error) {
	c, err := d.update.UpdateClient(ctx, &gorm.Client{ID: client.ID}, updates)
	if err != nil {
		return nil, err
	}

	return &domain.ClientProfile{
		ID:                      c.ID,
		Active:                  c.Active,
		ClientType:              c.ClientType,
		UserID:                  *c.UserID,
		TreatmentEnrollmentDate: c.TreatmentEnrollmentDate,
		FHIRPatientID:           c.FHIRPatientID,
		HealthRecordID:          c.HealthRecordID,
		TreatmentBuddy:          c.TreatmentBuddy,
		ClientCounselled:        c.ClientCounselled,
		OrganisationID:          c.OrganisationID,
		FacilityID:              c.FacilityID,
		CHVUserID:               c.CHVUserID,
		CaregiverID:             c.CaregiverID,
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

// UpdateFacility updates the facility with the provided facility details
func (d *MyCareHubDb) UpdateFacility(ctx context.Context, facility *domain.Facility, updateData map[string]interface{}) error {
	gormFacility := &gorm.Facility{
		FacilityID: facility.ID,
	}

	return d.update.UpdateFacility(ctx, gormFacility, updateData)
}
