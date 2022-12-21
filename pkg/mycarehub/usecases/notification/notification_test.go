package notification_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	fakeFCM "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
)

func TestUseCaseNotificationImpl_NotifyUser(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                 context.Context
		userProfile         *domain.User
		notificationPayload *domain.Notification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully notify user",
			args: args{
				ctx: ctx,
				userProfile: &domain.User{
					PushTokens: []string{uuid.New().String()},
				},
				notificationPayload: &domain.Notification{},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to notify user",
			args: args{
				ctx: ctx,
				userProfile: &domain.User{
					Name: gofakeit.Name(),
				},
				notificationPayload: &domain.Notification{},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save notification",
			args: args{
				ctx: ctx,
				userProfile: &domain.User{
					Name: gofakeit.Name(),
				},
				notificationPayload: &domain.Notification{
					Title: "Test title",
					Body:  "Test Body",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFCMService := fakeFCM.NewFCMServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			n := notification.NewNotificationUseCaseImpl(fakeFCMService, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad Case - Fail to notify user" {
				fakeFCMService.MockSendNotificationFn = func(ctx context.Context, payload *firebasetools.SendNotificationPayload) (bool, error) {
					return false, fmt.Errorf("failed to send notification")
				}
			}

			if tt.name == "Sad Case - Fail to save notification" {
				fakeDB.MockSaveNotificationFn = func(ctx context.Context, payload *domain.Notification) error {
					return fmt.Errorf("failed to save notification")
				}
			}

			if err := n.NotifyUser(tt.args.ctx, tt.args.userProfile, tt.args.notificationPayload); (err != nil) != tt.wantErr {
				t.Errorf("UseCaseNotificationImpl.NotifyUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCaseNotificationImpl_FetchNotifications(t *testing.T) {
	type args struct {
		ctx             context.Context
		userID          string
		flavour         feedlib.Flavour
		paginationInput dto.PaginationsInput
		filters         *domain.NotificationFilters
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.NotificationsPage
		wantErr bool
	}{
		{
			name: "happy case: list client notifications",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourConsumer,
				paginationInput: dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       5,
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: list staff notifications",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
				paginationInput: dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       5,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: cannot retrieve staff profile",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
				paginationInput: dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       5,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: cannot list notifications",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
				paginationInput: dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       5,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: fail validation",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
				paginationInput: dto.PaginationsInput{
					Limit: 5,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - fail to user profile by user id",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
				paginationInput: dto.PaginationsInput{
					Limit: 5,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFCMService := fakeFCM.NewFCMServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			n := notification.NewNotificationUseCaseImpl(fakeFCMService, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case - fail to user profile by user id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: cannot list notifications" {
				fakeDB.MockListNotificationsFn = func(ctx context.Context, params *domain.Notification, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Notification, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("cannot list notifications")
				}
			}

			if tt.name == "sad case: cannot retrieve staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get a staff profile")
				}
			}

			got, err := n.FetchNotifications(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.paginationInput, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseNotificationImpl.FetchNotifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got = %v", got)
				return
			}
		})
	}
}

func TestUseCaseNotificationImpl_NotifyFacilityStaffs(t *testing.T) {
	id := gofakeit.UUID()

	type args struct {
		ctx                 context.Context
		facility            *domain.Facility
		notificationPayload *domain.Notification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad case: cannot save notification",
			args: args{
				ctx: context.Background(),
				facility: &domain.Facility{
					ID: &id,
				},
				notificationPayload: &domain.Notification{
					Title: "Test notification title",
					Body:  "Test notification body",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: cannot retrieve facility staff",
			args: args{
				ctx: context.Background(),
				facility: &domain.Facility{
					ID: &id,
				},
				notificationPayload: &domain.Notification{
					Title: "Test notification title",
					Body:  "Test notification body",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: cannot send notification",
			args: args{
				ctx: context.Background(),
				facility: &domain.Facility{
					ID: &id,
				},
				notificationPayload: &domain.Notification{
					Title: "Test notification title",
					Body:  "Test notification body",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFCMService := fakeFCM.NewFCMServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			n := notification.NewNotificationUseCaseImpl(fakeFCMService, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: cannot save notification" {
				fakeDB.MockSaveNotificationFn = func(ctx context.Context, payload *domain.Notification) error {
					return fmt.Errorf("cannot save notification")
				}
			}

			if tt.name == "sad case: cannot retrieve facility staff" {
				fakeDB.MockGetFacilityStaffsFn = func(ctx context.Context, facilityID string) ([]*domain.StaffProfile, error) {
					return nil, fmt.Errorf("cannot get facility staffs")
				}
			}

			if tt.name == "sad case: cannot send notification" {
				fakeFCMService.MockSendNotificationFn = func(ctx context.Context, payload *firebasetools.SendNotificationPayload) (bool, error) {
					return false, fmt.Errorf("cannot send notification")
				}
			}

			if err := n.NotifyFacilityStaffs(tt.args.ctx, tt.args.facility, tt.args.notificationPayload); (err != nil) != tt.wantErr {
				t.Errorf("UseCaseNotificationImpl.NotifyFacilityStaffs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCaseNotificationImpl_SendNotification(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                context.Context
		registrationTokens []string
		data               map[string]interface{}
		notification       *firebasetools.FirebaseSimpleNotificationInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully send a notification",
			args: args{
				ctx:                ctx,
				registrationTokens: []string{},
				data:               map[string]interface{}{},
				notification:       &firebasetools.FirebaseSimpleNotificationInput{},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFCMService := fakeFCM.NewFCMServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			n := notification.NewNotificationUseCaseImpl(fakeFCMService, fakeDB, fakeDB, fakeDB, fakeExtension)

			got, err := n.SendNotification(tt.args.ctx, tt.args.registrationTokens, tt.args.data, tt.args.notification)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseNotificationImpl.SendNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseNotificationImpl.SendNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseNotificationImpl_ReadNotifications(t *testing.T) {

	type args struct {
		ctx context.Context
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: mark notification as read",
			args: args{
				ctx: context.Background(),
				ids: []string{gofakeit.UUID()},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: non existent notification",
			args: args{
				ctx: context.Background(),
				ids: []string{gofakeit.UUID()},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: error updating notification",
			args: args{
				ctx: context.Background(),
				ids: []string{gofakeit.UUID()},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFCMService := fakeFCM.NewFCMServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			n := notification.NewNotificationUseCaseImpl(fakeFCMService, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: non existent notification" {
				fakeDB.MockGetNotificationFn = func(ctx context.Context, notificationID string) (*domain.Notification, error) {
					return nil, fmt.Errorf("fail to update a notification")
				}
			}

			if tt.name == "sad case: error updating notification" {
				fakeDB.MockUpdateNotificationFn = func(ctx context.Context, notification *domain.Notification, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update notification")
				}
			}

			got, err := n.ReadNotifications(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseNotificationImpl.ReadNotifications() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseNotificationImpl.ReadNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseNotificationImpl_FetchNotificationTypeFilters(t *testing.T) {

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.NotificationTypeFilter
		wantErr bool
	}{
		{
			name: "happy case: list available filters",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			want: []*domain.NotificationTypeFilter{
				{
					Enum: enums.NotificationTypeAppointment,
					Name: enums.NotificationTypeAppointment.String(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to get logged in user",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: fail to get staff profile",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: fail to list available notifications",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Sad case - fail to user profile by user id",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeFCMService := fakeFCM.NewFCMServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			n := notification.NewNotificationUseCaseImpl(fakeFCMService, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "Sad case - fail to user profile by user id" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "sad case: fail to get staff profile" {
				fakeDB.MockGetStaffProfileFn = func(ctx context.Context, userID string, programID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile")
				}
			}

			if tt.name == "sad case: fail to list available notifications" {
				fakeDB.MockListAvailableNotificationTypesFn = func(ctx context.Context, params *domain.Notification) ([]enums.NotificationType, error) {
					return []enums.NotificationType{}, fmt.Errorf("fail to fetch notification types")
				}
			}

			got, err := n.FetchNotificationTypeFilters(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseNotificationImpl.FetchNotificationTypeFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCaseNotificationImpl.FetchNotificationTypeFilters() expected a response but got = %v", got)
				return
			}
		})
	}
}
