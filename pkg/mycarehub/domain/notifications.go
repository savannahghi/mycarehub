package domain

import (
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Notification represents a notification
type Notification struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Type      enums.NotificationType `json:"type"`
	IsRead    bool                   `json:"isRead"`
	CreatedAt time.Time              `json:"createdAt"`

	UserID         *string         `json:"userID"`
	FacilityID     *string         `json:"facilityID"`
	Flavour        feedlib.Flavour `json:"flavour"`
	ProgramID      string          `json:"programID"`
	OrganisationID string          `json:"organisationID"`
}

// NotificationsPage response for fetching notifications
type NotificationsPage struct {
	Notifications []*Notification `json:"notifications"`
	Pagination    Pagination      `json:"pagination"`
}

// NotificationTypeFilter represents an enum and its name value
type NotificationTypeFilter struct {
	Enum enums.NotificationType
	Name string `json:"name"`
}

// NotificationFilters represents the filters used to fetch notifications
type NotificationFilters struct {
	IsRead            *bool                     `json:"isRead"`
	NotificationTypes []*enums.NotificationType `json:"notificationTypes"`
}
