package domain

import "gorm.io/gorm"

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	gorm.Model

	//globally unique when set
	FacilityID *int64 `gorm:"primaryKey;autoIncrement:true;unique;column:id"`
	// unique within this structure
	Name string `gorm:"column:name"`
	// MFL Code for Kenyan facilities, globally unique
	Code        string `gorm:"unique;column:mfl_code"`
	Active      bool   `gorm:"column:active"`
	County      string `gorm:"column:county"` // TODO: Controlled list of counties
	Description string `gorm:"column:description"`
}

// TableName customizes how the table name is generated
func (Facility) TableName() string {
	return "facility"
}

// FacilityPage models the structure of all facilities including pagination
type FacilityPage struct {
	Facilities   []*Facility
	Count        int
	CurrentPage  int
	NextPage     *int
	PreviousPage *int
}

// FilterParam models the structure of the the filter parameters
type FilterParam struct {
	Name     string
	DataType string // TODO: Ideally a controlled list i.e enum
	Date     string // TODO: Clear spec on validation e.g dates must be ISO 8601
}

// AllTables is a collection of all tables in the domain
func AllTables() []interface{} {
	tables := []interface{}{
		Facility{},
	}
	return tables
}
