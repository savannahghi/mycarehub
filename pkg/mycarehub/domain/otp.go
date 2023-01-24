package domain

import (
	"time"

	"github.com/savannahghi/feedlib"
)

// OTP model the OTP details of OTP data
type OTP struct {
	UserID      string          `json:"userID"`
	Valid       bool            `json:"valid"`
	GeneratedAt time.Time       `json:"generatedAt"`
	ValidUntil  time.Time       `json:"validUntil"`
	Channel     string          `json:"channel"`
	Flavour     feedlib.Flavour `json:"flavour"`
	PhoneNumber string          `json:"phoneNumber"`
	OTP         string          `json:"otp"`
}

// OTPResponse models the object to be returned to the client when an OTP is generated
type OTPResponse struct {
	OTP         string `json:"otp"`
	PhoneNumber string `json:"phoneNumber"`
}
