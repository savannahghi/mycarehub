package domain

// Program defines the program structure
type Program struct {
	ID             string `json:"id"`
	Active         bool   `json:"active"`
	Name           string `json:"name"`
	OrganisationID string `json:"organisationID"`
}
