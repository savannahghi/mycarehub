package domain

import (
	"time"

	"github.com/savannahghi/feedlib"
)

type OTP struct {
	UserID      string          `json:"userID"`
	Valid       bool            `json:"valid"`
	GeneratedAt time.Time       `json:"generatedAt"`
	ValidUntil  time.Time       `json:"validUntil"`
	Channel     string          `json:"channel"`
	Flavour     feedlib.Flavour `json:"flavour"`
	PhoneNumber string          `json:"phoneNumber"`
}
