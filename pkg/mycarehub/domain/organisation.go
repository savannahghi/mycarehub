package domain

// Organisation represents the DAO of an organisation
type Organisation struct {
	ID               string `json:"id"`
	Active           bool   `json:"active"`
	OrganisationCode string `json:"org_code"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	EmailAddress     string `json:"email_address"`
	PhoneNumber      string `json:"phone_number"`
	PostalAddress    string `json:"postal_address"`
	PhysicalAddress  string `json:"physical_address"`
	DefaultCountry   string `json:"default_country"`
}
