package domain

// Organisation represents the DAO of an organisation
type Organisation struct {
	ID              string     `json:"id"`
	Active          bool       `json:"active"`
	Code            string     `json:"code"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	EmailAddress    string     `json:"emailAddress"`
	PhoneNumber     string     `json:"phoneNumber"`
	PostalAddress   string     `json:"postalAddress"`
	PhysicalAddress string     `json:"physicalAddress"`
	DefaultCountry  string     `json:"defaultCountry"`
	Programs        []*Program `json:"programs"`
}
