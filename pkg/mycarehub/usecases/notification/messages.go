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
}

// ComposeStaffNotification composes a staff notification which will be sent to the staff at a facility
func ComposeStaffNotification(notificationType enums.NotificationType, input StaffNotificationArgs) *domain.Notification {
	notification := &domain.Notification{
		Flavour: feedlib.FlavourPro,
		Type:    notificationType,
	}

	switch notificationType {
	case enums.NotificationTypeServiceRequest:
		notificationBody := fmt.Sprintf(
			"%s from %s requires your attention. Please follow up and resolve it.",
			ServiceRequestMessage(*input.ServiceRequestType),
			input.Subject.Name,
		)

		notification.Title = "A service request has been created"
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
	case enums.ServiceRequestTypeSurveyRedFlag:
		return "A flagged survey response service request"
	default:
		return ""
	}
}

// ClientNotificationInput is a collection of arguments required to compose a notification and the associated message
type ClientNotificationInput struct {
	// Arguments to a community invite notification
	Community *domain.Community
	Inviter   *domain.User
	Demoter   *domain.User
	Promoter  *domain.User

	// Arguments to an appointment notification
	Appointment   *domain.Appointment
	IsRescheduled bool

	// Args to a survey notification
	Survey *domain.UserSurvey
}

// ComposeClientNotification composes a client notification which will be sent to the client at a facility
func ComposeClientNotification(notificationType enums.NotificationType, input ClientNotificationInput) *domain.Notification {
	notification := &domain.Notification{
		Flavour: feedlib.FlavourConsumer,
		Type:    notificationType,
	}

	switch notificationType {
	case enums.NotificationTypeCommunities:
		notificationBody := fmt.Sprintf(
			"Invitation to join %s community by %s. To join, accept the invite.",
			input.Community.Name,
			input.Inviter.Name,
		)

		notification.Title = "You have been invited to join a conversation"
		notification.Body = notificationBody

		return notification

	case enums.NotificationTypeAppointment:
		reason := strings.ToLower(input.Appointment.Reason)
		date := input.Appointment.Date.AsTime().Format("January 02, 2006")

		if input.IsRescheduled {
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
		notification.Body = fmt.Sprintf("You have a new %s survey. Please navigate to the homepage and fill it.", input.Survey.Title)

		return notification

	case enums.NotificationTypeDemoteModerator:
		notification.Title = "You have been demoted to a regular user"
		notification.Body = fmt.Sprintf("You have been demoted to a regular user by %s in %s community.", input.Demoter.Username, input.Community.Name)

		return notification

	case enums.NotificationTypePromoteToModerator:
		notification.Title = "You have been promoted to a moderator"
		notification.Body = fmt.Sprintf("You have been promoted to a moderator by %s in %s community.", input.Promoter.Username, input.Community.Name)

		return notification

	default:
		return nil
	}
}
