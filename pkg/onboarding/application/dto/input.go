package dto

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Active      bool   `json:"active"`
	County      string `json:"county"`
	Description string `json:"description"`
}
