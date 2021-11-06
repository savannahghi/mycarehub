package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global facility ID(GCID)
	ID *string `json:"id"`

	// unique within this structure
	Name string `json:"name"`

	// MFL Code for Kenyan facilities, globally unique
	Code        string           `json:"code"`
	Active      bool             `json:"active"`
	County      enums.CountyType `json:"county"` // TODO: Controlled list of counties
	Description string           `json:"description"`
}

//FacilityPage returns a list of paginates facilities
type FacilityPage struct {
	Pagination Pagination
	Facilities []Facility
}
