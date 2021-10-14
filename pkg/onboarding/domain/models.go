package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global customer ID(GCID)
	ID uuid.UUID
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
// 	DataType FacilityIdentifiers // TODO: Ideally a controlled list i.e enum (MFL code, Active, County )
// 	Value    string              // TODO: Clear spec on validation e.g dates must be ISO 8601. This is the actual data being filtered
// }

// Metric reprents the metrics data structure input
type Metric struct {
	// ensures we don't re-save the same metric; opaque; globally unique
	MetricID uuid.UUID

	// TODO Metric types should be a controlled list i.e enum
	Type MetricType

	// this will vary by context
	// should not identify the user (there's a UID field)
	// focus on the actual event
	Payload datatypes.JSON `gorm:"column:payload"`

	Timestamp time.Time

	// a user identifier, can be hashed for anonymity
	// with a predictable one way hash
	UID string
}
