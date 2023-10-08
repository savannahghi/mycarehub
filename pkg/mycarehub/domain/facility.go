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

	Phone              string `json:"phone"`
	Active             bool   `json:"active"`
	Country            string `json:"country"`
	County             string `json:"county"`
	Address            string `json:"address"`
	Description        string `json:"description"`
	FHIROrganisationID string `json:"fhirOrganisationId"`

	Identifier FacilityIdentifier `json:"identifier"`

	WorkStationDetails WorkStationDetails `json:"workStationDetails"`

	Coordinates Coordinates `json:"coordinates"`
}

// Coordinates is used to show geographical locations
type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// FacilityPage returns a list of paginates facilities
type FacilityPage struct {
	Pagination Pagination  `json:"pagination"`
	Facilities []*Facility `json:"facilities"`
}

// UpdateFacilityPayload is the payload for updating faacility(s) fhir organization ID
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
