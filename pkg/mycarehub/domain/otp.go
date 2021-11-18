package domain

import (
	"time"

	"github.com/savannahghi/feedlib"
)

// OTP models the details of the OTP payload
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
