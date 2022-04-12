package fcm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
	fcmMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm/mock"
)

func TestServiceFCMImpl_SendNotification(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx     context.Context
		payload *firebasetools.SendNotificationPayload
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully send notification",
			args: args{
				ctx: ctx,
				payload: &firebasetools.SendNotificationPayload{
					RegistrationTokens: []string{uuid.New().String()},
					Data: map[string]string{
						"data": "user",
					},
					Notification: &firebasetools.FirebaseSimpleNotificationInput{},
					Android:      &firebasetools.FirebaseAndroidConfigInput{},
					Ios:          &firebasetools.FirebaseAPNSConfigInput{},
					Web:          &firebasetools.FirebaseWebpushConfigInput{},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcmService := fcm.NewService()
			fakeFcm := fcmMock.NewFCMServiceMock()
			got, err := fcmService.SendNotification(tt.args.ctx, tt.args.payload)

			if tt.name == "Sad Case - Fail to send notification" {
				fakeFcm.MockSendNotificationFn = func(ctx context.Context, payload *firebasetools.SendNotificationPayload) (bool, error) {
					return false, fmt.Errorf("failed to send notification")
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceFCMImpl.SendNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceFCMImpl.SendNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}
