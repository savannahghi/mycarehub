package domain

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global customer ID(GCID)
	ID int64
	// unique within this structure
	Name string
	// MFL Code for Kenyan facilities, globally unique
	Code        string
	Active      bool
	County      string // TODO: Controlled list of counties
	Description string
}

// // FacilityPage models the structure of all facilities including pagination
// type FacilityPage struct {
// 	Facilities   []*Facility
// 	Count        int
// 	CurrentPage  int
// 	NextPage     *int
// 	PreviousPage *int
// }

// // FilterParam models the structure of the the filter parameters
// type FilterParam struct {
// 	Name     string
// 	DataType string // TODO: Ideally a controlled list i.e enum
// 	Date     string // TODO: Clear spec on validation e.g dates must be ISO 8601
// }
