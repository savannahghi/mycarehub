package notification

import (
	"context"
	"fmt"
	"log"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
)

// IServiceNotify specifies a set of method signatures that are used to send notifications to client, staffs or facilities
type IServiceNotify interface {
	NotifyUser(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error
	NotifyFacilityStaffs(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error
	FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput, filters *domain.NotificationFilters) (*domain.NotificationsPage, error)
	FetchNotificationTypeFilters(ctx context.Context, flavour feedlib.Flavour) ([]*domain.NotificationTypeFilter, error)
	ReadNotifications(ctx context.Context, ids []string) (bool, error)

	SendNotification(
		ctx context.Context,
		registrationTokens []string,
		data map[string]interface{},
		notification *firebasetools.FirebaseSimpleNotificationInput,
	) (bool, error)
}

// UseCaseNotification holds the method signatures that are implemented in the notification usecase
type UseCaseNotification interface {
	IServiceNotify
}

// UseCaseNotificationImpl embeds the notifications logic
type UseCaseNotificationImpl struct {
	FCM         fcm.ServiceFCM
	ExternalExt extension.ExternalMethodsExtension
	Query       infrastructure.Query
	Create      infrastructure.Create
	Update      infrastructure.Update
}

// NewNotificationUseCaseImpl initialized a new notifications service implementation
func NewNotificationUseCaseImpl(
	fcm fcm.ServiceFCM,
	query infrastructure.Query,
	create infrastructure.Create,
	update infrastructure.Update,
	ext extension.ExternalMethodsExtension,
) UseCaseNotification {
	return &UseCaseNotificationImpl{
		FCM:         fcm,
		Query:       query,
		Create:      create,
		Update:      update,
		ExternalExt: ext,
	}
}

// NotifyUser is used to save and send a FCM notification to a user
func (n UseCaseNotificationImpl) NotifyUser(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
	notificationPayload.UserID = userProfile.ID
	notificationPayload.ProgramID = userProfile.CurrentProgramID
	if notificationPayload.Body != "" {
		err := n.Create.SaveNotification(ctx, notificationPayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return err
		}
	}

	notificationData := &dto.FCMNotificationMessage{
		Title: notificationPayload.Title,
	}

	payload := helpers.ComposeNotificationPayload(userProfile, *notificationData)
	_, err := n.FCM.SendNotification(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		log.Printf("failed to send notification: %v", err)
	}
	return nil
}

// NotifyFacilityStaffs is used to save and send a FCM notification to a user/ staff at a facility
func (n UseCaseNotificationImpl) NotifyFacilityStaffs(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error {
	notificationPayload.FacilityID = facility.ID
	if notificationPayload.Body != "" {
		err := n.Create.SaveNotification(ctx, notificationPayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return err
		}
	}

	notificationData := &dto.FCMNotificationMessage{
		Title: notificationPayload.Title,
		Body:  notificationPayload.Body,
	}

	staffs, err := n.Query.GetFacilityStaffs(ctx, *facility.ID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	for _, staff := range staffs {
		payload := helpers.ComposeNotificationPayload(staff.User, *notificationData)
		_, err = n.FCM.SendNotification(ctx, payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("failed to send notification: %v", err)
		}
	}

	return nil
}

// FetchNotifications retrieves a users notifications
func (n UseCaseNotificationImpl) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput, filters *domain.NotificationFilters) (*domain.NotificationsPage, error) {
	// if user did not provide current page, throw an error
	if err := paginationInput.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("pagination input validation failed: %v", err)
	}

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	parameters := &domain.Notification{UserID: &userID, Flavour: flavour}

	switch flavour {
	case feedlib.FlavourPro:
		staff, err := n.Query.GetStaffProfileByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		parameters.FacilityID = staff.DefaultFacility.ID
	}

	var notificationFilters []*firebasetools.FilterParam
	if filters != nil {
		if filters.IsRead != nil {
			filter := &firebasetools.FilterParam{
				FieldName:           "is_read",
				FieldType:           enumutils.FieldTypeBoolean,
				ComparisonOperation: enumutils.OperationEqual,
				FieldValue:          *filters.IsRead,
			}
			notificationFilters = append(notificationFilters, filter)
		}

		if filters.NotificationTypes != nil && len(filters.NotificationTypes) > 0 {
			filter := &firebasetools.FilterParam{
				FieldName:           "notification_type",
				FieldType:           enumutils.FieldTypeString,
				ComparisonOperation: enumutils.OperationIn,
				FieldValue:          filters.NotificationTypes,
			}
			notificationFilters = append(notificationFilters, filter)
		}
	}

	notifications, pageInfo, err := n.Query.ListNotifications(ctx, parameters, notificationFilters, page)
	if err != nil {
		return nil, err
	}

	response := &domain.NotificationsPage{
		Notifications: notifications,
		Pagination:    *pageInfo,
	}

	return response, nil
}

// SendNotification is used to send an FCM notification to a registered push token. This API will mainly
// be used for test purposes i.e. to check whether FCMs are being sent to the passed push token
func (n UseCaseNotificationImpl) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]interface{},
	notification *firebasetools.FirebaseSimpleNotificationInput,
) (bool, error) {
	notificationData, err := converterandformatter.MapInterfaceToMapString(data)
	if err != nil {
		return false, err
	}

	payload := &firebasetools.SendNotificationPayload{
		RegistrationTokens: registrationTokens,
		Data:               notificationData,
		Notification:       notification,
	}
	return n.FCM.SendNotification(ctx, payload)
}

// ReadNotifications indicates that the notification as bee
func (n UseCaseNotificationImpl) ReadNotifications(ctx context.Context, ids []string) (bool, error) {

	for _, id := range ids {
		notification, err := n.Query.GetNotification(ctx, id)
		if err != nil {
			return false, err
		}

		update := map[string]interface{}{
			"is_read": true,
		}
		err = n.Update.UpdateNotification(ctx, notification, update)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// FetchNotificationTypeFilters fetches the available notification types for a user
func (n UseCaseNotificationImpl) FetchNotificationTypeFilters(ctx context.Context, flavour feedlib.Flavour) ([]*domain.NotificationTypeFilter, error) {
	userID, err := n.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	parameters := &domain.Notification{UserID: &userID, Flavour: flavour}

	switch flavour {
	case feedlib.FlavourPro:
		staff, err := n.Query.GetStaffProfileByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		parameters.FacilityID = staff.DefaultFacility.ID
	}

	notificationTypes, err := n.Query.ListAvailableNotificationTypes(ctx, parameters)
	if err != nil {
		return nil, err
	}

	var filters []*domain.NotificationTypeFilter

	for _, notificationType := range notificationTypes {
		filter := &domain.NotificationTypeFilter{
			Enum: notificationType,
			Name: notificationType.Name(),
		}
		filters = append(filters, filter)
	}

	return filters, nil
}
