package enums

import (
	"fmt"
	"io"
	"strconv"
)

// NotificationType represents a type of notification
type NotificationType string

const (
	// NotificationTypeAppointment represents notifications from appointments
	NotificationTypeAppointment NotificationType = "APPOINTMENT"

	// NotificationTypeServiceRequest represents notifications from service requests
	NotificationTypeServiceRequest NotificationType = "SERVICE_REQUEST"

	// NotificationTypeCommunities represents notifications from communities
	NotificationTypeCommunities NotificationType = "COMMUNITIES"

	// NotificationTypeSurveys represents notifications from surveys
	NotificationTypeSurveys NotificationType = "SURVEYS"

	// NotificationTypeRoleAssignment represents a role assignment notification
	NotificationTypeRoleAssignment NotificationType = "ROLE_ASSIGNMENT"

	// NotificationTypeRoleRevocation represents a role assignment notification
	NotificationTypeRoleRevocation NotificationType = "ROLE_REVOCATION"

	// NotificationTypeDemoteModerator represents a demote moderator notification
	NotificationTypeDemoteModerator NotificationType = "DEMOTE_MODERATOR"

	// NotificationTypePromoteToModerator represents a promote to moderator notification
	NotificationTypePromoteToModerator NotificationType = "PROMOTE_TO_MODERATOR"

	// NotificationTypeBooking represents a booking notification
	NotificationTypeBooking NotificationType = "BOOKING"
)

// AllNotificationTypes holds all types of notification
var AllNotificationTypes = []NotificationType{
	NotificationTypeAppointment,
	NotificationTypeServiceRequest,
	NotificationTypeCommunities,
	NotificationTypeRoleRevocation,
	NotificationTypeRoleAssignment,
	NotificationTypeSurveys,
	NotificationTypeDemoteModerator,
	NotificationTypePromoteToModerator,
	NotificationTypeBooking,
}

// IsValid returns true if a notification type is valid
func (n NotificationType) IsValid() bool {
	switch n {
	case
		NotificationTypeAppointment,
		NotificationTypeServiceRequest,
		NotificationTypeCommunities,
		NotificationTypeRoleRevocation,
		NotificationTypeRoleAssignment,
		NotificationTypeSurveys,
		NotificationTypeDemoteModerator,
		NotificationTypePromoteToModerator,
		NotificationTypeBooking:
		return true
	}
	return false
}

// String returns a string representation of the enum
func (n NotificationType) String() string {
	return string(n)
}

// UnmarshalGQL converts the supplied value to a metric type.
func (n *NotificationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*n = NotificationType(str)
	if !n.IsValid() {
		return fmt.Errorf("%s is not a valid NotificationType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (n NotificationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(n.String()))
}

// Name returns a human readable format for an enum
func (n NotificationType) Name() string {
	switch n {
	case NotificationTypeAppointment:
		return "Appointments"
	case NotificationTypeServiceRequest:
		return "Service Requests"
	case NotificationTypeCommunities:
		return "Communities"
	case NotificationTypeRoleRevocation:
		return "Role Addition"
	case NotificationTypeRoleAssignment:
		return "Role Removal"
	case NotificationTypeSurveys:
		return "Surveys"
	case NotificationTypeDemoteModerator:
		return "Moderator Demotion"
	case NotificationTypePromoteToModerator:
		return "Moderator Promotion"
	case NotificationTypeBooking:
		return "Booking"
	}
	return "UNKNOWN"
}
