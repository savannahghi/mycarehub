package mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// NotificationUseCaseMock mocks the notifications usecase methods
type NotificationUseCaseMock struct {
	MockNotifyUserFn           func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error
	MockNotifyFacilityStaffsFn func(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error
	MockFetchNotificationsFn   func(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error)
}

// NewServiceNotificationMock initializes a new notification mock instance
func NewServiceNotificationMock() *NotificationUseCaseMock {
	return &NotificationUseCaseMock{
		MockNotifyUserFn: func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
			return nil
		},
		MockNotifyFacilityStaffsFn: func(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error {
			return nil
		},
		MockFetchNotificationsFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error) {
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
	}
}

// NotifyUser mocks the implementation of sending a fcm notification to a user
func (n NotificationUseCaseMock) NotifyUser(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
	return n.MockNotifyUserFn(ctx, userProfile, notificationPayload)
}

// FetchNotifications mocks the implementation of fetching notifications from the database
func (n NotificationUseCaseMock) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error) {
	return n.MockFetchNotificationsFn(ctx, userID, flavour, paginationInput)
}

// NotifyFacilityStaffs is used to save and send a FCM notification to a user/ staff at a facility
func (n NotificationUseCaseMock) NotifyFacilityStaffs(ctx context.Context, facility *domain.Facility, notificationPayload *domain.Notification) error {
	return n.MockNotifyFacilityStaffsFn(ctx, facility, notificationPayload)
}
