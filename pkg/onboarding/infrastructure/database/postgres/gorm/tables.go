package gorm

import "gorm.io/gorm"

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	gorm.Model

	//globally unique when set
	FacilityID *int64 `gorm:"primaryKey;autoIncrement:true;unique;column:facility_id"`
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

func allTables() []interface{} {
	tables := []interface{}{
		&Facility{},
	}
	return tables
}
