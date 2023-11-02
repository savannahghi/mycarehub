package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global facility ID(GCID)
	ID *string `json:"id,omitempty"`

	// unique within this structure
	Name string `json:"name,omitempty"`

	Phone              string                `json:"phone,omitempty"`
	Active             bool                  `json:"active,omitempty"`
	Country            string                `json:"country,omitempty"`
	County             string                `json:"county,omitempty"`
	Address            string                `json:"address,omitempty"`
	Description        string                `json:"description,omitempty"`
	FHIROrganisationID string                `json:"fhirOrganisationId,omitempty"`
	Distance           float64               `json:"distance,omitempty"`
	Identifiers        []*FacilityIdentifier `json:"identifiers,omitempty"`
	WorkStationDetails WorkStationDetails    `json:"workStationDetails,omitempty"`
	Coordinates        *Coordinates          `json:"coordinates,omitempty"`
	Services           []FacilityService     `json:"services,omitempty"`
	BusinessHours      []BusinessHours       `json:"businessHours,omitempty"`
}

// BusinessHours models data that show facility's operational hours
type BusinessHours struct {
	ID          string `json:"id"`
	Day         string `json:"day"`
	OpeningTime string `json:"openingTime"`
	ClosingTime string `json:"closingTime"`
	FacilityID  string `json:"facilityID"`
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

// FacilityServicePage models the services offered in a facility
type FacilityServicePage struct {
	Results     []FacilityService `json:"results"`
	Count       int               `json:"count"`
	Next        string            `json:"next"`
	Previous    string            `json:"previous"`
	PageSize    int               `json:"page_size"`
	CurrentPage int               `json:"current_page"`
	TotalPages  int               `json:"total_pages"`
	StartIndex  int               `json:"start_index"`
	EndIndex    int               `json:"end_index"`
}

// FacilityService models the data class that is used to show facility services
type FacilityService struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Identifiers []ServiceIdentifier `json:"identifiers"`
}

// ServiceIdentifier models the structure of facility's service identifiers
type ServiceIdentifier struct {
	ID              string `json:"id"`
	IdentifierType  string `json:"identifierType"`
	IdentifierValue string `json:"identifierValue"`
	ServiceID       string `json:"serviceID"`
}
