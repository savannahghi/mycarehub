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

	UserID     *string
	FacilityID *string
	Flavour    feedlib.Flavour
}

// NotificationsPage response for fetching notifications
type NotificationsPage struct {
	Notifications []*Notification `json:"notifications"`
	Pagination    Pagination      `json:"pagination"`
}

// NotificationTypeFilter represents an enum and its name value
type NotificationTypeFilter struct {
	Enum enums.NotificationType
	Name string
}

// NotificationFilters represents the filters used to fetch notifications
type NotificationFilters struct {
	IsRead            *bool
	NotificationTypes []*enums.NotificationType
}
