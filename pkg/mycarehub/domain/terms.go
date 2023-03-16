package domain

import (
	"time"

	"github.com/savannahghi/feedlib"
)

// TermsOfService contains the struct field to hold the required display data for the terms of service.
type TermsOfService struct {
	TermsID   int             `json:"termsID"`
	Text      *string         `json:"text"`
	Flavour   feedlib.Flavour `json:"flavour"`
	ValidFrom time.Time       `json:"validFrom"`
	ValidTo   time.Time       `json:"validTo"`
}
