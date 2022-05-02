package notification

import (
	"fmt"

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
