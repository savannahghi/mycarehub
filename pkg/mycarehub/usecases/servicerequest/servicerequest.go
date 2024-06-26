package servicerequest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/healthcrm"
	serviceSMS "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms"
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
		username string,
		cccNumber string,
		flavour feedlib.Flavour,
	) (bool, error)
	CompleteVisit(ctx context.Context, staffID string, serviceRequestID string, bookingID string, notes string) (bool, error)
}

// ISetInProgresssBy is an interface that contains the method signature for assigning the staff currently working on a request
type ISetInProgresssBy interface {
	SetInProgressBy(ctx context.Context, requestID string, staffID string) (bool, error)
}

// IGetServiceRequests is an interface that holds the method signature for getting service requests
type IGetServiceRequests interface {
	GetServiceRequests(ctx context.Context, requestType string, requestStatus *string, facilityID string, flavour feedlib.Flavour, pagination *dto.PaginationsInput) (*domain.ServiceRequestPage, error)
	GetServiceRequestsForKenyaEMR(ctx context.Context, payload *dto.ServiceRequestPayload) (*dto.RedFlagServiceRequestResponse, error)
	GetPendingServiceRequestsCount(ctx context.Context) (*domain.ServiceRequestsCountResponse, error)
	SearchServiceRequests(ctx context.Context, searchTerm string, flavour feedlib.Flavour, requestType string, facilityID string) ([]*domain.ServiceRequest, error)
}

// IResolveServiceRequest is an interface that holds the method signature for resolving a service request
type IResolveServiceRequest interface {
	ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error)
	VerifyClientPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus, physicalIdentityVerified bool) (bool, error)
	VerifyStaffPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus) (bool, error)
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
	SMS          serviceSMS.IServiceSMS
	HealthCRM    healthcrm.IHealthCRMService
}

// NewUseCaseServiceRequestImpl creates a new service request instance
func NewUseCaseServiceRequestImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
	ext extension.ExternalMethodsExtension,
	user user.UseCasesUser,
	notification notification.UseCaseNotification,
	sms serviceSMS.IServiceSMS,
	healthCRM healthcrm.IHealthCRMService,
) *UseCasesServiceRequestImpl {
	return &UseCasesServiceRequestImpl{
		Create:       create,
		Query:        query,
		Update:       update,
		ExternalExt:  ext,
		User:         user,
		Notification: notification,
		SMS:          sms,
		HealthCRM:    healthCRM,
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
			Active:         true,
			RequestType:    input.RequestType,
			Request:        input.Request,
			Status:         "PENDING",
			ClientID:       input.ClientID,
			FacilityID:     *clientProfile.DefaultFacility.ID,
			Meta:           input.Meta,
			ProgramID:      clientProfile.User.CurrentProgramID,
			OrganisationID: clientProfile.User.CurrentOrganizationID,
			CaregiverID:    input.CaregiverID,
		}
		err = u.Create.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create client's service request: %v", err)
		}

		clientNotification := &domain.Notification{
			Title:          "Your service request has been created",
			Body:           "",
			UserID:         clientProfile.User.ID,
			Flavour:        feedlib.FlavourConsumer,
			Type:           enums.NotificationTypeServiceRequest,
			ProgramID:      clientProfile.User.CurrentProgramID,
			OrganisationID: clientProfile.User.CurrentOrganizationID,
		}
		err = u.Notification.NotifyUser(ctx, clientProfile.User, clientNotification)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

		facility, err := u.Query.RetrieveFacility(ctx, clientProfile.DefaultFacility.ID, true)
		if err != nil {
			helpers.ReportErrorToSentry(err)
		}

		requestType := enums.ServiceRequestType(input.RequestType)
		notificationInput := notification.StaffNotificationArgs{
			Subject:            clientProfile.User,
			ServiceRequestType: &requestType,
		}

		notificationType := enums.NotificationType(input.RequestType)
		if input.RequestType != enums.ServiceRequestBooking.String() {
			notificationType = enums.NotificationTypeServiceRequest
		}

		staffNotification := notification.ComposeStaffNotification(notificationType, notificationInput)

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
			Active:         true,
			RequestType:    input.RequestType,
			Request:        input.Request,
			Status:         "PENDING",
			StaffID:        input.StaffID,
			FacilityID:     *staffProfile.DefaultFacility.ID,
			ProgramID:      staffProfile.User.CurrentProgramID,
			OrganisationID: staffProfile.User.CurrentOrganizationID,
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
	requestType string,
	requestStatus *string,
	facilityID string,
	flavour feedlib.Flavour,
	pagination *dto.PaginationsInput,
) (*domain.ServiceRequestPage, error) {
	if requestType != "" {
		if !enums.ServiceRequestType(requestType).IsValid() {
			return nil, fmt.Errorf("invalid request type: %v", requestType)
		}
	}
	if requestStatus != nil {
		if !enums.ServiceRequestStatus(*requestStatus).IsValid() {
			return nil, fmt.Errorf("invalid request status: %v", *requestStatus)
		}
	}

	page := &domain.Pagination{
		Limit:       pagination.Limit,
		CurrentPage: pagination.CurrentPage,
	}

	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, err
	}

	userProfile, err := u.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		return nil, err
	}

	exists, err := u.Query.CheckIfFacilityExistsInProgram(ctx, userProfile.CurrentProgramID, facilityID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("facility %v does not exist in program %v", facilityID, userProfile.CurrentProgramID)
	}

	results, page, err := u.Query.GetServiceRequests(ctx, &requestType, requestStatus, facilityID, userProfile.CurrentProgramID, flavour, page)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	if requestType == enums.ServiceRequestBooking.String() {
		var wg sync.WaitGroup

		resultsChannel := make([]chan []domain.FacilityService, len(results))

		for i, result := range results {
			wg.Add(1)

			resultsChannel[i] = make(chan []domain.FacilityService)

			// Goroutine for each result
			go func(serviceRequest *domain.ServiceRequest, resultChan chan []domain.FacilityService) {
				defer wg.Done()

				serviceList, ok := serviceRequest.Meta["serviceIDs"].([]interface{})
				if !ok {
					helpers.ReportErrorToSentry(fmt.Errorf("a service should have at least one service id"))
					resultChan <- []domain.FacilityService{}
				}

				facilityServices := u.fetchServices(ctx, serviceList)

				resultChan <- facilityServices

			}(result, resultsChannel[i])
		}

		go func() {
			wg.Wait()
			for _, resultChan := range resultsChannel {
				close(resultChan)
			}
		}()

		for i, resultChan := range resultsChannel {
			results[i].Services = <-resultChan
		}
	}

	return &domain.ServiceRequestPage{
		Results:    results,
		Pagination: *page,
	}, nil
}

// fetchServices is a helper function to fetch services from health CRM
func (u *UseCasesServiceRequestImpl) fetchServices(ctx context.Context, serviceList []interface{}) []domain.FacilityService {
	var facilityServices []domain.FacilityService

	for _, id := range serviceList {
		serviceID := id.(string)

		service, err := u.HealthCRM.GetServiceByID(ctx, serviceID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			continue
		}

		facilityServices = append(facilityServices, *service)
	}

	return facilityServices
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
func (u *UseCasesServiceRequestImpl) GetPendingServiceRequestsCount(ctx context.Context) (*domain.ServiceRequestsCountResponse, error) {
	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	userProfile, err := u.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	staffProfile, err := u.Query.GetStaffProfile(ctx, *userProfile.ID, userProfile.CurrentProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	return u.Query.GetPendingServiceRequestsCount(ctx, *staffProfile.DefaultFacility.ID, staffProfile.ProgramID)
}

// ResolveServiceRequest resolves a service request
func (u *UseCasesServiceRequestImpl) ResolveServiceRequest(ctx context.Context, staffID *string, serviceRequestID *string, action []string, comment *string) (bool, error) {
	if staffID == nil {
		return false, fmt.Errorf("staff ID is required")
	}
	if serviceRequestID == nil {
		return false, fmt.Errorf("service request ID is required")
	}

	serviceRequest, err := u.Query.GetClientServiceRequestByID(ctx, *serviceRequestID)
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
		serviceRequest, err := u.Query.GetClientServiceRequestByID(ctx, request.ID)
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
func (u *UseCasesServiceRequestImpl) CreatePinResetServiceRequest(ctx context.Context, username string, cccNumber string, flavour feedlib.Flavour) (bool, error) {
	switch flavour {
	case feedlib.FlavourConsumer:
		if cccNumber == "" {
			return false, fmt.Errorf("ccc number cannot be empty")
		}
		if username == "" {
			return false, fmt.Errorf("username cannot be empty")
		}

		var meta = map[string]interface{}{}
		meta["ccc_number"] = cccNumber

		// TODO: Check if the service request exists before creating a new one
		userProfile, err := u.Query.GetUserProfileByUsername(ctx, username)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ProfileNotFoundErr(err)
		}

		clientProfile, err := u.Query.GetClientProfile(ctx, *userProfile.ID, userProfile.CurrentProgramID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ClientProfileNotFoundErr(err)
		}

		meta["is_ccc_number_valid"] = true
		_, err = u.Query.GetProgramClientProfileByIdentifier(ctx, clientProfile.ProgramID, enums.UserIdentifierTypeCCC.String(), cccNumber)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				meta["is_ccc_number_valid"] = false
			} else {
				helpers.ReportErrorToSentry(err)
				return false, exceptions.GetError(err)
			}
		}

		serviceRequestInput := &dto.ServiceRequestInput{
			Active:         true,
			RequestType:    enums.ServiceRequestTypePinReset.String(),
			Request:        "Change PIN Request",
			ClientID:       *clientProfile.ID,
			FacilityID:     *clientProfile.DefaultFacility.ID,
			Flavour:        feedlib.FlavourConsumer,
			Meta:           meta,
			ProgramID:      clientProfile.User.CurrentProgramID,
			OrganisationID: clientProfile.User.CurrentOrganizationID,
		}

		_, err = u.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, err
		}

		return true, nil

	case feedlib.FlavourPro:
		userProfile, err := u.Query.GetUserProfileByUsername(ctx, username)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.ProfileNotFoundErr(err)
		}

		staffProfile, err := u.Query.GetStaffProfile(ctx, *userProfile.ID, userProfile.CurrentProgramID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, exceptions.StaffProfileNotFoundErr(err)
		}

		serviceRequestInput := &dto.ServiceRequestInput{
			Active:         true,
			RequestType:    enums.ServiceRequestTypeStaffPinReset.String(),
			Request:        "Change PIN Request",
			StaffID:        *staffProfile.ID,
			FacilityID:     *staffProfile.DefaultFacility.ID,
			Flavour:        feedlib.FlavourPro,
			ProgramID:      staffProfile.User.CurrentProgramID,
			OrganisationID: staffProfile.User.CurrentOrganizationID,
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
func (u *UseCasesServiceRequestImpl) VerifyStaffPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus) (bool, error) {
	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	loggedInUserProfile, err := u.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.ProfileNotFoundErr(err)
	}

	loggedInStaffProfile, err := u.Query.GetStaffProfile(ctx, loggedInUserID, loggedInUserProfile.CurrentProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}

	serviceRequest, err := u.Query.GetStaffServiceRequestByID(ctx, serviceRequestID)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("failed to get staff service request by ID %s: %w", serviceRequestID, err))
		return false, fmt.Errorf("failed to get staff service request by ID %s: %w", serviceRequestID, err)
	}

	staffProfile, err := u.Query.GetStaffProfileByStaffID(ctx, serviceRequest.StaffID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	err = u.Update.UpdateUser(ctx, &domain.User{ID: &staffProfile.UserID}, map[string]interface{}{
		"pin_update_required": true,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(err)
	}

	phoneNumber := staffProfile.User.Contacts.ContactValue

	return u.VerifyServiceRequestResponse(ctx, status.String(), phoneNumber, serviceRequestID, staffProfile.User, loggedInStaffProfile, feedlib.FlavourPro)

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
func (u *UseCasesServiceRequestImpl) VerifyClientPinResetServiceRequest(ctx context.Context, serviceRequestID string, status enums.PINResetVerificationStatus, physicalIdentityVerified bool) (bool, error) {
	loggedInUserID, err := u.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}

	loggedInUserProfile, err := u.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	loggedInStaffProfile, err := u.Query.GetStaffProfile(ctx, loggedInUserID, loggedInUserProfile.CurrentProgramID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.StaffProfileNotFoundErr(err)
	}

	serviceRequest, err := u.Query.GetClientServiceRequestByID(ctx, serviceRequestID)
	if err != nil {
		helpers.ReportErrorToSentry(fmt.Errorf("failed to get staff service request by ID %s: %w", serviceRequestID, err))
		return false, fmt.Errorf("failed to get staff service request by ID %s: %w", serviceRequestID, err)
	}

	clientProfile, err := u.Query.GetClientProfileByClientID(ctx, serviceRequest.ClientID)
	if err != nil {
		return false, err
	}

	if !physicalIdentityVerified {
		return false, fmt.Errorf("the patient has not been physically verified by the healthcare worker")
	}

	err = u.Update.UpdateUser(ctx, &domain.User{ID: &clientProfile.UserID}, map[string]interface{}{
		"pin_update_required": true,
	})
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.UpdateProfileErr(err)
	}

	phoneNumber := clientProfile.User.Contacts.ContactValue

	return u.VerifyServiceRequestResponse(ctx, status.String(), phoneNumber, serviceRequestID, clientProfile.User, loggedInStaffProfile, feedlib.FlavourConsumer)

}

// VerifyServiceRequestResponse returns the boolean response indicating whether the processing of a service request is successful or not.
func (u *UseCasesServiceRequestImpl) VerifyServiceRequestResponse(
	ctx context.Context, state, phoneNumber, serviceRequestID string,
	user *domain.User,
	staff *domain.StaffProfile,
	flavour feedlib.Flavour,
) (bool, error) {
	switch state {
	case enums.PINResetVerificationStatusRejected.String():
		text := fmt.Sprintf(
			"Dear %s, your request to reset your pin has been rejected. "+
				"For enquiries call us on %s.", user.Name, callCenterNumber,
		)

		_, err := u.SMS.SendSMS(ctx, text, []string{phoneNumber})
		if err != nil {
			helpers.ReportErrorToSentry(err)
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

	case enums.PINResetVerificationStatusApproved.String():
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

		_, sendErr := u.SMS.SendSMS(ctx, text, []string{phoneNumber})
		if sendErr != nil {
			helpers.ReportErrorToSentry(err)
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

// CompleteVisit is used to complete/mark a booking service as attended to when the client visits the service provider
func (u *UseCasesServiceRequestImpl) CompleteVisit(ctx context.Context, staffID string, serviceRequestID string, bookingID string, notes string) (bool, error) {
	if staffID == "" {
		return false, fmt.Errorf("staff ID is required")
	}
	if serviceRequestID == "" {
		return false, fmt.Errorf("service request ID is required")
	}

	serviceRequest, err := u.Query.GetClientServiceRequestByID(ctx, serviceRequestID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("failed to get service request: %v", err)
	}

	err = u.Update.ResolveServiceRequest(ctx, &staffID, &serviceRequest.ID, enums.ServiceRequestStatusResolved.String(), nil, &notes)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	updateData := map[string]interface{}{
		"booking_status": enums.Fulfilled,
	}

	err = u.Update.UpdateBooking(ctx, &domain.Booking{ID: bookingID}, updateData)
	if err != nil {
		return false, err
	}

	return true, nil
}
