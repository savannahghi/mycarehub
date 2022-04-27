package notification

import (
	"context"
	"fmt"
	"log"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
)

// IServiceNotify specifies a set of method signatures that are used to send notifications to client, staffs or facilities
type IServiceNotify interface {
	NotifyUser(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error
	NotifyFacilityStaffs(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error
	FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error)
}

// UseCaseNotification holds the method signatures that are implemented in the notification usecase
type UseCaseNotification interface {
	IServiceNotify
}

// UseCaseNotificationImpl embeds the notifications logic
type UseCaseNotificationImpl struct {
	FCM    fcm.ServiceFCM
	Query  infrastructure.Query
	Create infrastructure.Create
}

// NewNotificationUseCaseImpl initialized a new notifications service implementation
func NewNotificationUseCaseImpl(
	fcm fcm.ServiceFCM,
	query infrastructure.Query,
	create infrastructure.Create,
) UseCaseNotification {
	return &UseCaseNotificationImpl{
		FCM:    fcm,
		Query:  query,
		Create: create,
	}
}

// NotifyUser is used to save and send a FCM notification to a user
func (n UseCaseNotificationImpl) NotifyUser(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
	notificationPayload.UserID = userProfile.ID
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
func (n UseCaseNotificationImpl) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error) {
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
		parameters.FacilityID = &staff.DefaultFacilityID
	}

	notifications, pageInfo, err := n.Query.ListNotifications(ctx, parameters, page)
	if err != nil {
		return nil, err
	}

	response := &domain.NotificationsPage{
		Notifications: notifications,
		Pagination:    *pageInfo,
	}

	return response, nil
}
