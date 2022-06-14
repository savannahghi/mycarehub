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

	// NotificationTypeDemoteModerator represents a demote moderator notification
	NotificationTypeDemoteModerator NotificationType = "DEMOTE_MODERATOR"

	// NotificationTypePromoteToModerator represents a promote to moderator notification
	NotificationTypePromoteToModerator NotificationType = "PROMOTE_TO_MODERATOR"
)

// IsValid returns true if a notification type is valid
func (n NotificationType) IsValid() bool {
	switch n {
	case
		NotificationTypeAppointment,
		NotificationTypeServiceRequest,
		NotificationTypeCommunities,
		NotificationTypeRoleAssignment,
		NotificationTypeSurveys,
		NotificationTypeDemoteModerator,
		NotificationTypePromoteToModerator:
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
