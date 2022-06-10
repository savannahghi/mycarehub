package notification

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

func TestServiceRequestMessage(t *testing.T) {
	type args struct {
		request enums.ServiceRequestType
	}

	type test struct {
		name    string
		args    args
		wantErr bool
	}

	tests := []test{
		{
			name: "sad case: unknown service request type",
			args: args{
				request: "UNKNOWN",
			},
			wantErr: true,
		},
	}

	// tests/ensures that every defined service request has an associated notification message
	for _, requestType := range enums.AllServiceRequestType {
		t := test{
			name: fmt.Sprintf("happy case: %s service request type", requestType),
			args: args{
				request: requestType,
			},
			wantErr: false,
		}
		tests = append(tests, t)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ServiceRequestMessage(tt.args.request)
			if !tt.wantErr && got == "" {
				t.Errorf("ServiceRequestMessage() expected a notification message for: %v", tt.args.request)
			}
		})
	}
}

func TestComposeStaffNotification(t *testing.T) {
	redFlag := enums.ServiceRequestTypeRedFlag
	type args struct {
		notificationType enums.NotificationType
		args             StaffNotificationArgs
	}
	tests := []struct {
		name string
		args args
		want *domain.Notification
	}{
		{
			name: "service request notification",
			args: args{
				notificationType: enums.NotificationTypeServiceRequest,
				args: StaffNotificationArgs{
					Subject: &domain.User{
						Name: "John Doe",
					},
					ServiceRequestType: &redFlag,
				},
			},
			want: &domain.Notification{
				Title:   "A service request has been created",
				Body:    "A flagged health diary entry service request from John Doe requires your attention. Please follow up and resolve it.",
				Type:    enums.NotificationTypeServiceRequest,
				Flavour: feedlib.FlavourPro,
			},
		},
		{
			name: "unknown notification type",
			args: args{
				notificationType: "UNKNOWN",
				args: StaffNotificationArgs{
					Subject: &domain.User{
						Name: "John Doe",
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeStaffNotification(tt.args.notificationType, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComposeStaffNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComposeClientNotification(t *testing.T) {
	type args struct {
		notificationType enums.NotificationType
		args             ClientNotificationArgs
	}
	tests := []struct {
		name string
		args args
		want *domain.Notification
	}{
		{
			name: "community invite notification",
			args: args{
				notificationType: enums.NotificationTypeCommunities,
				args: ClientNotificationArgs{
					Community: &domain.Community{
						Name: "Good Life",
					},
					Inviter: &domain.User{
						Name: "John Doe",
					},
				},
			},
			want: &domain.Notification{
				Title:   "You have been invited to join a conversation",
				Body:    "Invitation to join Good Life community by John Doe. To join, accept the invite.",
				Type:    enums.NotificationTypeCommunities,
				Flavour: feedlib.FlavourConsumer,
			},
		},
		{
			name: "new appointment notification",
			args: args{
				notificationType: enums.NotificationTypeAppointment,
				args: ClientNotificationArgs{
					Appointment: &domain.Appointment{
						Reason: "Dental Check",
						Date: scalarutils.Date{
							Year:  2022,
							Month: 1,
							Day:   1,
						},
					},
				},
			},
			want: &domain.Notification{
				Title:   "You have a new scheduled appointment",
				Body:    "You have a new dental check appointment scheduled for January 01, 2022.",
				Type:    enums.NotificationTypeAppointment,
				Flavour: feedlib.FlavourConsumer,
			},
		},
		{
			name: "appointment reschedule notification",
			args: args{
				notificationType: enums.NotificationTypeAppointment,
				args: ClientNotificationArgs{
					Appointment: &domain.Appointment{
						Reason: "Dental Check",
						Date: scalarutils.Date{
							Year:  2022,
							Month: 2,
							Day:   1,
						},
					},
					IsRescheduled: true,
				},
			},
			want: &domain.Notification{
				Title:   "An appointment has been rescheduled",
				Body:    "Your dental check appointment has been rescheduled to February 01, 2022.",
				Type:    enums.NotificationTypeAppointment,
				Flavour: feedlib.FlavourConsumer,
			},
		},
		{
			name: "unknown notification type",
			args: args{
				notificationType: "UNKNOWN",
				args: ClientNotificationArgs{
					Inviter: &domain.User{
						Name: "John Doe",
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposeClientNotification(tt.args.notificationType, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComposeClientNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}
