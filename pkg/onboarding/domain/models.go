package domain

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	ID          *string // globally unique when set
	Name        string  // unique within this structure
	Code        string  // MFL Code for Kenyan facilities, globally unique
	Active      bool
	County      string // TODO: Controlled list of counties
	Description string
}

// FacilityPage ...
type FacilityPage struct {
	Facilities   []*Facility
	Count        int
	CurrentPage  int
	NextPage     *int
	PreviousPage *int
}

// FilterParam ...
type FilterParam struct {
	Name     string
	DataType string // TODO: Ideally a controlled list i.e enum
	Date     string // TODO: Clear spec on validation e.g dates must be ISO 8601
}
