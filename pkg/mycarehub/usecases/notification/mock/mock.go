package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// NotificationUseCaseMock mocks the notifications usecase methods
type NotificationUseCaseMock struct {
	MockNotifyUserFn                 func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error
	MockNotifyFacilityStaffsFn       func(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error
	MockFetchNotificationsFn         func(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput, filters *domain.NotificationFilters) (*domain.NotificationsPage, error)
	MockFetchNotificationTypeFilters func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.NotificationTypeFilter, error)
	MockSendNotificationFn           func(
		ctx context.Context,
		registrationTokens []string,
		data map[string]interface{},
		notification *firebasetools.FirebaseSimpleNotificationInput,
	) (bool, error)
	MockReadNotificationsFn func(ctx context.Context, ids []string) (bool, error)
}

// NewServiceNotificationMock initializes a new notification mock instance
func NewServiceNotificationMock() *NotificationUseCaseMock {
	return &NotificationUseCaseMock{
		MockFetchNotificationTypeFilters: func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.NotificationTypeFilter, error) {
			return []*domain.NotificationTypeFilter{{Enum: enums.NotificationTypeAppointment, Name: enums.NotificationTypeAppointment.String()}}, nil
		},
		MockNotifyUserFn: func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
			return nil
		},
		MockReadNotificationsFn: func(ctx context.Context, ids []string) (bool, error) {
			return true, nil
		},
		MockNotifyFacilityStaffsFn: func(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error {
			return nil
		},
		MockFetchNotificationsFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput, filters *domain.NotificationFilters) (*domain.NotificationsPage, error) {
			UUID := uuid.New().String()
			return &domain.NotificationsPage{
				Notifications: []*domain.Notification{
					{
						ID:         UUID,
						Title:      "New Teleconsult",
						Body:       "Teleconsult with Doctor Who at the Tardis",
						Type:       "TELECONSULT",
						IsRead:     false,
						UserID:     &UUID,
						FacilityID: &UUID,
						Flavour:    flavour,
					},
				},
				Pagination: domain.Pagination{},
			}, nil
		},
		MockSendNotificationFn: func(
			ctx context.Context,
			registrationTokens []string,
			data map[string]interface{},
			notification *firebasetools.FirebaseSimpleNotificationInput,
		) (bool, error) {
			return true, nil
		},
	}
}

// NotifyUser mocks the implementation of sending a fcm notification to a user
func (n NotificationUseCaseMock) NotifyUser(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
	return n.MockNotifyUserFn(ctx, userProfile, notificationPayload)
}

// FetchNotifications mocks the implementation of fetching notifications from the database
func (n NotificationUseCaseMock) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput, filters *domain.NotificationFilters) (*domain.NotificationsPage, error) {
	return n.MockFetchNotificationsFn(ctx, userID, flavour, paginationInput, filters)
}

// NotifyFacilityStaffs is used to save and send a FCM notification to a user/ staff at a facility
func (n NotificationUseCaseMock) NotifyFacilityStaffs(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error {
	return n.MockNotifyFacilityStaffsFn(ctx, facility, notificationPayload)
}

// SendNotification is used to send a FCM notification
func (n NotificationUseCaseMock) SendNotification(
	ctx context.Context,
	registrationTokens []string,
	data map[string]interface{},
	notification *firebasetools.FirebaseSimpleNotificationInput,
) (bool, error) {
	return n.MockSendNotificationFn(ctx, registrationTokens, data, notification)
}

//ReadNotifications indicates that the notification as bee
func (n NotificationUseCaseMock) ReadNotifications(ctx context.Context, ids []string) (bool, error) {
	return n.MockReadNotificationsFn(ctx, ids)
}

//FetchNotificationTypeFilters fetches the available notification types for a user
func (n NotificationUseCaseMock) FetchNotificationTypeFilters(ctx context.Context, flavour feedlib.Flavour) ([]*domain.NotificationTypeFilter, error) {
	return n.MockFetchNotificationTypeFilters(ctx, flavour)
}
