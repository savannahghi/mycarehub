package dto

//SessionDetails input for ussd
type SessionDetails struct {
	SessionID   string  `json:"sessionID"`
	PhoneNumber *string `json:"phoneNumber"`
	Level       int     `json:"level"`
	Text        string  `json:"text"`
}

//EndSessionDetails store end of session details
type EndSessionDetails struct {
	SessionID    string
	ServiceCode  string
	NetworkCode  string
	PhoneNumber  *string
	Status       string
	Input        string
	ErrorMessage string
}
