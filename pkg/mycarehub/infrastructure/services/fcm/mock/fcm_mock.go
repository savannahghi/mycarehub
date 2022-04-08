package mock

import (
	"context"

	"github.com/savannahghi/firebasetools"
)

// FCMMock mocks the fcm service implementations
type FCMMock struct {
	MockSendNotificationFn func(ctx context.Context, payload *firebasetools.SendNotificationPayload) (bool, error)
}

// NewFCMServiceMock initializes the fcm mock service
func NewFCMServiceMock() *FCMMock {
	return &FCMMock{
		MockSendNotificationFn: func(ctx context.Context, payload *firebasetools.SendNotificationPayload) (bool, error) {
			return true, nil
		},
	}
}

// SendNotification mocks the implementation of sending push notifications to the app
func (f FCMMock) SendNotification(ctx context.Context, payload *firebasetools.SendNotificationPayload) (bool, error) {
	return f.MockSendNotificationFn(ctx, payload)
}
