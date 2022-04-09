package domain

import "github.com/savannahghi/feedlib"

// Notification represents a notification
type Notification struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Type   string `json:"type"`
	IsRead bool   `json:"isRead"`

	UserID     *string
	FacilityID *string
	Flavour    feedlib.Flavour
}

// NotificationsPage response for fetching notifications
type NotificationsPage struct {
	Notifications []*Notification `json:"notifications"`
	Pagination    Pagination      `json:"pagination"`
}
