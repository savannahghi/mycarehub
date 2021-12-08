package domain

import "github.com/savannahghi/feedlib"

// FAQ domain entity contains the FAQ information
type FAQ struct {
	ID          *string         `json:"id"`
	Active      bool            `json:"active"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Body        string          `json:"body"`
	Flavour     feedlib.Flavour `json:"flavour"`
}
