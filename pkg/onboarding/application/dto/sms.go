package dto

// AfricasTalkingMessage contains the SMS message data
type AfricasTalkingMessage struct {
	// A unique identifier attached to each incoming message.
	LinkID string `json:"linkId"`
	// The content of the message received.
	Text string `json:"text"`
	// Your registered short code that the sms was sent out to.
	To string `json:"to"`
	// The id of the message.
	ID string `json:"id"`
	// The date when the sms was sent.
	Date string `json:"date"`
	// The sender’s phone number.
	From string `json:"from"`
}
