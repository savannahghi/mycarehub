package domain

import "time"

// ClientHealthDiaryEntry models the health diary entry. It is used to capture the
// client's moods on a day-by-day basis
type ClientHealthDiaryEntry struct {
	ID                    *string   `json:"id"`
	Active                bool      `json:"active"`
	Mood                  string    `json:"mood"`
	Note                  string    `json:"note"`
	EntryType             string    `json:"entryType"`
	ShareWithHealthWorker bool      `json:"shareWithHealthWorker"`
	SharedAt              time.Time `json:"sharedAt"`
	ClientID              string    `json:"clientID"`
	CreatedAt             time.Time `json:"createdAt"`
	PhoneNumber           string    `json:"phoneNumber"`
	ClientName            string    `json:"clientName"`
	CCCNumber             string    `json:"cccNumber"`
}
