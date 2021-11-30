package domain

import "time"

// ClientHealthDiaryEntry models the health diary entry. It is used to capture the
// client's moods on a day-by-day basis
type ClientHealthDiaryEntry struct {
	Active                bool      `json:"active"`
	Mood                  string    `json:"mood"`
	Note                  string    `json:"note"`
	EntryType             string    `json:"entryType"`
	ShareWithHealthWorker bool      `json:"shareWithHealthWorker"`
	SharedAt              time.Time `json:"sharedAt"`
	ClientID              string    `json:"clientID"`
	CreatedAt             time.Time `json:"createdAt"`
}

// ClientServiceRequest models a service request created for the healthcare worker.
type ClientServiceRequest struct {
	Active         bool      `json:"active"`
	RequestType    string    `json:"requestType"`
	Request        string    `json:"request"`
	Status         string    `json:"status"`
	InProgressAt   time.Time `json:"inProgressAt"`
	ResolvedAt     time.Time `json:"resolvedAt"`
	ClientID       string    `json:"clientID"`
	InProgressByID string    `json:"inProgressByID"`
	ResolvedByID   string    `json:"resolvedByID"`
}
