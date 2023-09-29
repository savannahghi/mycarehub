package domain

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global facility ID(GCID)
	ID *string `json:"id"`

	// unique within this structure
	Name string `json:"name"`

	Phone              string `json:"phone"`
	Active             bool   `json:"active"`
	Country            string `json:"country"`
	Description        string `json:"description"`
	FHIROrganisationID string `json:"fhirOrganisationId"`

	Identifier FacilityIdentifier `json:"identifier"`

	WorkStationDetails WorkStationDetails `json:"workStationDetails"`
	Branches           []Facility         `json:"facility"`
	Parent             string             `json:"parent"`
	Status             string             `json:"status"`
	Type               string             `json:"type"`
	County             string             `json:"county"`
	Address            Address            `json:"address"`
	Coordinates        Coordinates        `json:"coordinates"`
	AccreditationType  string             `json:"accreditationType"` // TODO: enums
}

type Address struct {
	Location        Coordinates `json:"location"`
	Country         string      `json:"country"` // TODO: Should be country enums
	County          string      `json:"county"`  // TODO: Should be county enums
	PhysicalAddress string      `json:"physicalAddress"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// FacilityPage returns a list of paginates facilities
type FacilityPage struct {
	Pagination Pagination  `json:"pagination"`
	Facilities []*Facility `json:"facilities"`
}

// UpdateFacilityPayload is the payload for updating facility(s) fhir organization ID
type UpdateFacilityPayload struct {
	FacilityID         string `json:"facilityID"`
	FHIROrganisationID string `json:"fhirOrganisationID"`
}

// FacilityIdentifier is the identifier of the facility
type FacilityIdentifier struct {
	ID     string                       `json:"id"`
	Active bool                         `json:"active"`
	Type   enums.FacilityIdentifierType `json:"type"`
	Value  string                       `json:"value"`
}

// Service models the details of 'services' that are available in a facility
type Service struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Category    string `json:"category"`
	Type        string `json:"type"` // TODO: This should be enum

	Price float64 `json:"price"`
}

// ServiceCategory models the details the category in which a facility's service belongs to.
type ServiceCategory struct {
	ID          string `json:"id"`
	Name        string `json:"categoryName"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	ServiceID   string `json:"serviceID"`
}

// FacilityBusinessHours models the facility's business hours
type FacilityBusinessHours struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Day       string    `json:"day"`
}

// FacilityContact models the contact details of a facility
type FacilityContact struct {
	ID           string   `json:"id"`
	ContactType  string   `json:"contactType"`
	ContactValue string   `json:"contactValue"`
	Active       bool     `json:"active"`
	Role         string   `json:"role"`
	Facility     Facility `json:"facility"`
}

// FacilityPhoto models the the structure of the data used to show facility's photo
type FacilityPhoto struct {
	ID        string   `json:"id"`
	Photo     string   `json:"photo"`
	PhotoType string   `json:"photoType"` // TODO: Should be an enum
	Facility  Facility `json:"facility"`
}

// Rating models facility/service ratings
type Rating struct {
	Value   int    `json:"value"`
	Comment string `json:"comment"`
}
