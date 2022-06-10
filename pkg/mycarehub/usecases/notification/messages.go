package notification

import (
	"fmt"
	"strings"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// StaffNotificationArgs is a collection of arguments required to compose a notification and the associated message
type StaffNotificationArgs struct {
	// Subject is the user responsible for creating/dispatching the notification
	// i.e "who it's about". Used to personalize the notification
	Subject *domain.User

	// Arguments for a service request notification
	ServiceRequestType *enums.ServiceRequestType

	// Arguments for a role assignment notification
	RoleTypes []enums.UserRoleType
}

// ComposeStaffNotification composes a staff notification which will be sent to the staff at a facility
func ComposeStaffNotification(notificationType enums.NotificationType, args StaffNotificationArgs) *domain.Notification {
	notification := &domain.Notification{
		Flavour: feedlib.FlavourPro,
		Type:    notificationType,
	}

	switch notificationType {
	case enums.NotificationTypeServiceRequest:
		notificationBody := fmt.Sprintf(
			"%s from %s requires your attention. Please follow up and resolve it.",
			ServiceRequestMessage(*args.ServiceRequestType),
			args.Subject.Name,
		)

		notification.Title = "A service request has been created"
		notification.Body = notificationBody

		return notification
	case enums.NotificationTypeRoleAssignment:
		notificationBody := "You have been assigned the following role(s): "
		for i, role := range args.RoleTypes {
			if i == 0 {
				notificationBody += role.Name()
			} else {
				notificationBody += fmt.Sprintf(", %s", role.Name())
			}
		}

		notification.Title = "You have been assigned a new role"
		notification.Body = notificationBody

		return notification
	default:
		return nil
	}
}

// ServiceRequestMessage determines the notification message based on the service request type
func ServiceRequestMessage(request enums.ServiceRequestType) string {
	switch request {
	case enums.ServiceRequestTypeRedFlag:
		return "A flagged health diary entry service request"
	case enums.ServiceRequestTypePinReset:
		return "A PIN reset service request"
	case enums.ServiceRequestTypeStaffPinReset:
		return "A staff PIN reset service request"
	case enums.ServiceRequestTypeHomePageHealthDiary:
		return "A shared health diary service request"
	case enums.ServiceRequestTypeAppointments:
		return "An appointment reschedule request"
	case enums.ServiceRequestTypeScreeningToolsRedFlag:
		return "A flagged screening tool response service request"
	default:
		return ""
	}
}

// ClientNotificationArgs is a collection of arguments required to compose a notification and the associated message
type ClientNotificationArgs struct {
	// Arguments to a community invite notification
	Community *domain.Community
	Inviter   *domain.User

	// Arguments to an appointment notification
	Appointment   *domain.Appointment
	IsRescheduled bool

	// Args to a survey notification
	Survey *domain.UserSurvey
}

// ComposeClientNotification composes a client notification which will be sent to the client at a facility
func ComposeClientNotification(notificationType enums.NotificationType, args ClientNotificationArgs) *domain.Notification {
	notification := &domain.Notification{
		Flavour: feedlib.FlavourConsumer,
		Type:    notificationType,
	}

	switch notificationType {
	case enums.NotificationTypeCommunities:
		notificationBody := fmt.Sprintf(
			"Invitation to join %s community by %s. To join, accept the invite.",
			args.Community.Name,
			args.Inviter.Name,
		)

		notification.Title = "You have been invited to join a conversation"
		notification.Body = notificationBody

		return notification

	case enums.NotificationTypeAppointment:
		reason := strings.ToLower(args.Appointment.Reason)
		date := args.Appointment.Date.AsTime().Format("January 02, 2006")

		if args.IsRescheduled {
			notificationBody := fmt.Sprintf(
				"Your %s appointment has been rescheduled to %s.",
				reason,
				date,
			)

			notification.Title = "An appointment has been rescheduled"
			notification.Body = notificationBody
		} else {
			notificationBody := fmt.Sprintf(
				"You have a new %s appointment scheduled for %s.",
				reason,
				date,
			)

			notification.Title = "You have a new scheduled appointment"
			notification.Body = notificationBody
		}

		return notification

	case enums.NotificationTypeSurveys:
		notification.Title = "You have a new survey"
		notification.Body = fmt.Sprintf("You have a new %s survey. Please navigate to the homepage and fill it.", args.Survey.Title)

		return notification

	default:
		return nil
	}
}
