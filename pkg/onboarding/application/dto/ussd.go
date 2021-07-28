package dto

//SessionDetails input for ussd
type SessionDetails struct {
	SessionID   string  `json:"sessionID"`
	PhoneNumber *string `json:"phoneNumber"`
	Level       int     `json:"level"`
	Text        string  `json:"text"`
}
