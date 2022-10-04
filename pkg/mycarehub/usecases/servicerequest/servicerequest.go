package servicerequest

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"gorm.io/gorm"
)

// Service requests are tasks for the healthcare staff on the platform. Some examples are:
// red flags raised by content posted to a group or entries into health diaries
// appointment reschedule requests
// These tasks will be presented on a list and notified (e.g., via push notifications).
// Each task will have a status. When first created, the tasks will be marked as “PENDING”.
// Once the relevant actions are taken, it will be possible to mark them as “IN PROGRESS”, “RESOLVED” and add relevant notes.
// In order to ensure that a task is not addressed by multiple people at the same time, each task will be updated with a record of the user and timestamp each time the status is changed.

const (
	callCenterNumber = "0790 360 360"
)

// ICreateServiceRequest is an interface that holds the method signature for creating a service request
type ICreateServiceRequest interface {
	CreateServiceRequest(ctx context.Context, input *dto.ServiceRequestInput) (bool, error)
	CreatePinResetServiceRequest(
		ctx context.Context,
		phoneNumber string,
		cccNumber string,
		flavour feedlib.Flavour,
	) (bool, error)
}

// ISetInProgresssBy is an interface that contains the method signature for assigning the staff currently working on a request
type ISetInProgresssBy interface {
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
}

// IGetServiceRequests is an interface that holds the method signature for getting service requests
type IGetServiceRequests interface {
	GetServiceRequests(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) (*dto.RedFlagServiceRequestResponse, error)
	GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error)
	SearchServiceRequests(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
}

// IResolveServiceRequest is an interface that holds the method signature for resolving a service request
type IResolveServiceRequest interface {
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error)
	VerifyClientPinResetServiceRequest(
		ctx context.Context,
		clientID string,
		serviceRequestID string,
		cccNumber string,
		phoneNumber string,
		physicalIdentityVerified bool,
		state string,
	) (bool, error)
	VerifyStaffPinResetServiceRequest(ctx context.Context, phoneNumber string, serviceRequestID string, verificationStatus string) (bool, error)
}

// IUpdateServiceRequest is the interface holding the method signature for updating service requests.
type IUpdateServiceRequest interface {
	UpdateServiceRequestsFromKenyaEMR(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error)
}

// UseCaseServiceRequest holds all the interfaces that represent the service request business logic
type UseCaseServiceRequest interface {
	ICreateServiceRequest
	IGetServiceRequests
	ISetInProgresssBy
	IResolveServiceRequest
	IUpdateServiceRequest
}

// UseCasesServiceRequestImpl embeds the service request logic
type UseCasesServiceRequestImpl struct {
	Create       infrastructure.Create
	Query        infrastructure.Query
	Update       infrastructure.Update
	ExternalExt  extension.ExternalMethodsExtension
	User         user.UseCasesUser
	Notification notification.UseCaseNotification
}

// NewUseCaseServiceRequestImpl creates a new service request instance
func NewUseCaseServiceRequestImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
	ext extension.ExternalMethodsExtension,
	user user.UseCasesUser,
	notification notification.UseCaseNotification,
) *UseCasesServiceRequestImpl {
	return &UseCasesServiceRequestImpl{
		Create:       create,
		Query:        query,
		Update:       update,
		ExternalExt:  ext,
		User:         user,
		Notification: notification,
	}
}

// CreateServiceRequest creates a service request
func (u *UseCasesServiceRequestImpl) CreateServiceRequest(ctx context.Context, input *dto.ServiceRequestInput) (bool, error) {
	switch input.Flavour {
	case feedlib.FlavourConsumer:
		if input.ClientID == "" {
			return false, fmt.Errorf("client ID is required")
		}
		clientProfile, err := u.Query.GetClientProfileByClientID(ctx, input.ClientID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ClientProfileNotFoundErr(err)
		}
		serviceRequestInput := &dto.ServiceRequestInput{
			Active:      true,
			RequestType: input.RequestType,
			Request:     input.Request,
			Status:      "PENDING",
			ClientID:    input.ClientID,
			FacilityID:  clientProfile.DefaultFacilityID,
			Meta:        input.Meta,
		}
		err = u.Create.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create client's service request: %v", err)
		}

		clientNotification := &domain.Notification{
			Title:   "Your service request has been created",
			Body:    "",
			Flavour: feedlib.FlavourConsumer,
			Type:    enums.NotificationTypeServiceRequest,
		}
		err = u.Notification.NotifyUser(ctx, clientProfile.User, clientNotification)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

		facility, err := u.Query.RetrieveFacility(ctx, &clientProfile.DefaultFacilityID, true)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

		requestType := enums.ServiceRequestType(input.RequestType)
		notificationInput := notification.StaffNotificationArgs{
			Subject:            clientProfile.User,
			ServiceRequestType: &requestType,
		}
		staffNotification := notification.ComposeStaffNotification(
			enums.NotificationTypeServiceRequest,
			notificationInput,
		)
		err = u.Notification.NotifyFacilityStaffs(ctx, facility, staffNotification)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

		return true, nil

	case feedlib.FlavourPro:
		if input.StaffID == "" {
			return false, fmt.Errorf("staff ID cannot be empty")
		}
		staffProfile, err := u.Query.GetStaffProfileByStaffID(ctx, input.StaffID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.StaffProfileNotFoundErr(err)
		}
		serviceRequestInput := &dto.ServiceRequestInput{
			Active:      true,
			RequestType: input.RequestType,
			Request:     input.Request,
			Status:      "PENDING",
			StaffID:     input.StaffID,
			FacilityID:  staffProfile.DefaultFacilityID,
		}
		err = u.Create.CreateStaffServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create staff's service request: %v", err)
		}
		return true, nil
	default:
		return false, fmt.Errorf("invalid flavour defined: %v", input.Flavour)
	}

}

// SetInProgressBy assigns to a service request, staff currently working on the service request
func (u *UseCasesServiceRequestImpl) SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error) {
	if requestID == "" || staffID == "" {
		return false, fmt.Errorf("request ID or staff ID cannot be empty")
	}
	return u.Update.SetInProgressBy(ctx, requestID, staffID)
}

// GetServiceRequests gets service requests based on the parameters provided
func (u *UseCasesServiceRequestImpl) GetServiceRequests(
	ctx context.Context,
	requestType *string,
	requestStatus *string,
	facilityID string,
	flavour feedlib.Flavour,
) ([]*domain.ServiceRequest, error) {
	if requestType != nil {
		if !enums.ServiceRequestType(*requestType).IsValid() {
			return nil, fmt.Errorf("invalid request type: %v", *requestType)
		}
	}
	if requestStatus != nil {
		if !enums.ServiceRequestStatus(*requestStatus).IsValid() {
			return nil, fmt.Errorf("invalid request status: %v", *requestStatus)
		}
	}

	return u.Query.GetServiceRequests(ctx, requestType, requestStatus, facilityID, flavour)
}

// GetServiceRequestsForKenyaEMR fetches all the most recent service requests  that have not been
// synced to KenyaEMR.
func (u *UseCasesServiceRequestImpl) GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) (*dto.RedFlagServiceRequestResponse, error) {
	serviceRequests, err := u.Query.GetServiceRequestsForKenyaEMR(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &dto.RedFlagServiceRequestResponse{
		RedFlagServiceRequests: serviceRequests,
	}, nil
}

// GetPendingServiceRequestsCount gets the total number of service requests
func (u *UseCasesServiceRequestImpl) GetPendingServiceRequestsCount(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
	return u.Query.GetPendingServiceRequestsCount(ctx, facilityID)
}

// ResolveServiceRequest resolves a service request
func (u *UseCasesServiceRequestImpl) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error) {
	if staffID == nil {
		return false, fmt.Errorf("staff ID is required")
	}
	if serviceRequestID == nil {
		return false, fmt.Errorf("service request ID is required")
	}
	serviceRequest, err := u.Query.GetServiceRequestByID(ctx, *serviceRequestID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get service request: %v", err)
	}
	clientProfile, err := u.Query.GetClientProfileByClientID(ctx, serviceRequest.ClientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get client profile: %v", err)
	}
	if serviceRequest.RequestType == enums.ServiceRequestTypePinReset.String() {
		userProfile := &domain.User{
			ID: &clientProfile.UserID,
		}
		updatePayload := map[string]interface{}{
			"next_allowed_login":    time.Now(),
			"failed_login_count":    0,
			"failed_security_count": 0,
		}
		err := u.Update.UpdateUser(ctx, userProfile, updatePayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to update user: %v", err)
		}

		err = u.Update.UpdateFailedSecurityQuestionsAnsweringAttempts(ctx, clientProfile.UserID, 0)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to reset client's failed security answering attempts: %v", err)
		}
	}

	resolveErr := u.Update.ResolveServiceRequest(ctx, staffID, serviceRequestID, enums.ServiceRequestStatusResolved.String(), action, comment)
	if resolveErr != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to update service request: %v", err)
	}

	return true, nil
}

// UpdateServiceRequestsFromKenyaEMR is used to update service requests from KenyaEMR to MyCareHub service.
func (u *UseCasesServiceRequestImpl) UpdateServiceRequestsFromKenyaEMR(ctx context.Context, payload *dto.UpdateServiceRequestsPayload) (bool, error) {
	var serviceRequests []domain.ServiceRequest

	for _, request := range payload.ServiceRequests {
		serviceRequest, err := u.Query.GetServiceRequestByID(ctx, request.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		if request.RequestType == enums.ServiceRequestTypeAppointments.String() &&
			request.Status == enums.ServiceRequestStatusResolved.String() {
			client, err := u.Query.GetClientProfileByClientID(ctx, serviceRequest.ClientID)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}

			// the appointment service request is not valid
			appointmentID, exists := serviceRequest.Meta["appointmentID"]
			if !exists {
				continue
			}

			filter := domain.Appointment{ID: appointmentID.(string)}
			appointment, err := u.Query.GetAppointment(ctx, filter)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}

			updates := map[string]interface{}{
				"has_rescheduled_appointment": false,
			}
			updatedAppointment, err := u.Update.UpdateAppointment(ctx, appointment, updates)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}

			notificationMessage := notification.ComposeClientNotification(
				enums.NotificationTypeAppointment,
				notification.ClientNotificationInput{
					Appointment:   updatedAppointment,
					IsRescheduled: true,
				},
			)

			err = u.Notification.NotifyUser(ctx, client.User, notificationMessage)
			if err != nil {
				helpers.ReportErrorToSentry(err)
			}

		}

		mapped := &domain.ServiceRequest{
			ID:           request.ID,
			RequestType:  request.RequestType,
			Status:       request.Status,
			InProgressAt: &request.InProgressAt,
			InProgressBy: &request.InProgressBy,
			ResolvedAt:   &request.ResolvedAt,
			ResolvedBy:   &request.ResolvedBy,
		}

		serviceRequests = append(serviceRequests, *mapped)
	}

	serviceReq := &domain.UpdateServiceRequestsPayload{
		ServiceRequests: serviceRequests,
	}

	return u.Update.UpdateServiceRequests(ctx, serviceReq)
}

// CreatePinResetServiceRequest creates a PIN_RESET service request. This occurs when a user attempts to change
// their pin but they don't succeed.
func (u *UseCasesServiceRequestImpl) CreatePinResetServiceRequest(ctx context.Context, phoneNumber string, cccNumber string, flavour feedlib.Flavour) (bool, error) {
	switch flavour {
	case feedlib.FlavourConsumer:
		if cccNumber == "" {
			return false, fmt.Errorf("ccc number cannot be empty")
		}
		if phoneNumber == "" {
			return false, fmt.Errorf("phone number cannot be empty")
		}

		var meta = map[string]interface{}{}
		meta["ccc_number"] = cccNumber

		// TODO: Check if the service request exists before creating a new one
		userProfile, err := u.Query.GetUserProfileByPhoneNumber(ctx, phoneNumber, feedlib.FlavourConsumer)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ProfileNotFoundErr(err)
		}

		clientProfile, err := u.Query.GetClientProfileByUserID(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ClientProfileNotFoundErr(err)
		}

		meta["is_ccc_number_valid"] = true
		_, err = u.Query.GetClientProfileByCCCNumber(ctx, cccNumber)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				meta["is_ccc_number_valid"] = false
			} else {
				helpers.ReportErrorToSentry(err)
				return false, exceptions.GetError(err)
			}
		}

		serviceRequestInput := &dto.ServiceRequestInput{
			Active:      true,
			RequestType: enums.ServiceRequestTypePinReset.String(),
			Request:     "Change PIN Request",
			ClientID:    *clientProfile.ID,
			FacilityID:  clientProfile.DefaultFacilityID,
			Flavour:     feedlib.FlavourConsumer,
			Meta:        meta,
		}

		_, err = u.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		return true, nil

	case feedlib.FlavourPro:
		userProfile, err := u.Query.GetUserProfileByPhoneNumber(ctx, phoneNumber, feedlib.FlavourPro)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ProfileNotFoundErr(err)
		}

		staffProfile, err := u.Query.GetStaffProfileByUserID(ctx, *userProfile.ID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.StaffProfileNotFoundErr(err)
		}

		serviceRequestInput := &dto.ServiceRequestInput{
			Active:      true,
			RequestType: enums.ServiceRequestTypeStaffPinReset.String(),
			Request:     "Change PIN Request",
			StaffID:     *staffProfile.ID,
			FacilityID:  staffProfile.DefaultFacilityID,
			Flavour:     feedlib.FlavourPro,
		}

		_, err = u.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		return true, nil

	default:
		return false, nil
	}

}

// VerifyStaffPinResetServiceRequest is used to approve and/or reject staff's service reset request.
// This is used by the admin to reset the login credentials of a staff who has been unable to sign into the portal
// and has requested for help from the admin through a service request.
func (u *UseCasesServiceRequestImpl) VerifyStaffPinResetServiceRequest(ctx context.Context, phoneNumber string, serviceRequestID string, verificationStatus string) (bool, error) {
	if phoneNumber == "" || serviceRequestID == "" || verificationStatus == "" {
		return false, fmt.Errorf("neither phoneNumber, serviceRequestID nor verificationStatus can be empty")
	}
	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}
	loggedInStaffProfile, err := u.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}
	userProfile, err := u.Query.GetUserProfileByPhoneNumber(ctx, phoneNumber, feedlib.FlavourPro)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}
	_, err = u.Query.GetStaffProfileByUserID(ctx, *userProfile.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}

	err = u.Update.UpdateUser(ctx, &domain.User{ID: userProfile.ID}, map[string]interface{}{
		"pin_update_required": true,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(err)
	}

	return u.VerifyServiceRequestResponse(ctx, verificationStatus, phoneNumber, serviceRequestID, userProfile, loggedInStaffProfile, feedlib.FlavourPro)

}

// VerifyClientPinResetServiceRequest is used to approve/reject a pin reset service request. This is used by the
// healthcare worker to reset the login credentials of a user who failed to login and requested for help from
// the health care worker.
//
// The basic workflow is
// 1. Get the logged in user ID - This will be used to identify the staff who resolved the request
// 2. Verify that the patient was physically verified by the healthcare worker and that the provided
// ccc number matches the one on their profile
// 3. Mark the service request as IN_PROGRESS
// 4. Send a fresh invite to the user and invalidate the previous pins
// 5. Update the field `pin_change_required` to true and mark the service request as resolved
func (u *UseCasesServiceRequestImpl) VerifyClientPinResetServiceRequest(
	ctx context.Context,
	clientID string,
	serviceRequestID string,
	cccNumber string,
	phoneNumber string,
	physicalIdentityVerified bool,
	state string,
) (bool, error) {
	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	staff, err := u.Query.GetStaffProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}

	userProfile, err := u.Query.GetUserProfileByPhoneNumber(ctx, phoneNumber, feedlib.FlavourConsumer)
	if err != nil {
		return false, err
	}

	if !physicalIdentityVerified {
		return false, fmt.Errorf("the patient has not been physically verified by the healthcare worker")
	}

	err = u.Update.UpdateUser(ctx, &domain.User{ID: userProfile.ID}, map[string]interface{}{
		"pin_update_required": true,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(err)
	}

	return u.VerifyServiceRequestResponse(ctx, state, phoneNumber, serviceRequestID, userProfile, staff, feedlib.FlavourConsumer)

}

// VerifyServiceRequestResponse returns the boolean response indicating whether the processing of a service request is successful or not.
func (u *UseCasesServiceRequestImpl) VerifyServiceRequestResponse(
	ctx context.Context, state, phoneNumber, serviceRequestID string,
	user *domain.User,
	staff *domain.StaffProfile,
	flavour feedlib.Flavour,
) (bool, error) {
	switch state {
	case enums.VerifyServiceRequestStateRejected.String():
		text := fmt.Sprintf(
			"Dear %s, your request to reset your pin has been rejected. "+
				"For enquiries call us on %s.", user.Name, callCenterNumber,
		)

		_, err := u.ExternalExt.SendSMS(ctx, phoneNumber, text, enumutils.SenderIDBewell)
		if err != nil {
			return false, err
		}

		switch flavour {
		case feedlib.FlavourPro:
			_, err = u.Update.ResolveStaffServiceRequest(ctx, staff.ID, &serviceRequestID, enums.ServiceRequestStatusRejected.String())
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}
		case feedlib.FlavourConsumer:
			err := u.Update.ResolveServiceRequest(ctx, staff.ID, &serviceRequestID, enums.ServiceRequestStatusRejected.String(), []string{}, nil)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}
		}

		return true, nil

	case enums.VerifyServiceRequestStateApproved.String():
		_, err := u.SetInProgressBy(ctx, serviceRequestID, *staff.ID)
		if err != nil {
			return false, err
		}

		tempPin, err := u.User.GenerateTemporaryPin(ctx, *user.ID, flavour)
		if err != nil {
			return false, err
		}

		text := fmt.Sprintf(
			"Dear %s, your request to reset your pin has been accepted. "+
				"Your One Time PIN is %s.", user.Name, tempPin,
		)

		_, err = u.ExternalExt.SendSMS(ctx, phoneNumber, text, enumutils.SenderIDBewell)
		if err != nil {
			return false, err
		}

		updatePayload := map[string]interface{}{
			"next_allowed_login":    time.Now(),
			"failed_login_count":    0,
			"failed_security_count": 0,
			"pin_update_required":   true,
		}
		err = u.Update.UpdateUser(ctx, user, updatePayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to update user: %v", err)
		}

		switch flavour {
		case feedlib.FlavourPro:
			_, err = u.Update.ResolveStaffServiceRequest(ctx, staff.ID, &serviceRequestID, enums.ServiceRequestStatusResolved.String())
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return false, err
			}
		case feedlib.FlavourConsumer:
			err := u.Update.ResolveServiceRequest(ctx, staff.ID, &serviceRequestID, enums.ServiceRequestStatusResolved.String(), []string{}, nil)
			if err != nil {
				return false, err
			}
		}

		return true, nil

	default:
		return false, fmt.Errorf("unknown state provided")
	}
}

// SearchServiceRequests searches for service requests based on the provided search term. Can be username, phone or request type.
func (u *UseCasesServiceRequestImpl) SearchServiceRequests(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error) {
	switch flavour {
	case feedlib.FlavourConsumer:
		return u.Query.SearchClientServiceRequests(ctx, searchTerm, requestType, facilityID)

	case feedlib.FlavourPro:
		return u.Query.SearchStaffServiceRequests(ctx, searchTerm, requestType, facilityID)

	default:
		return nil, fmt.Errorf("unknown flavour provided")
	}
}
