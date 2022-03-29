package domain

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global facility ID(GCID)
	ID *string `json:"id"`

	// unique within this structure
	Name string `json:"name"`

	// MFL Code for Kenyan facilities, globally unique
	Code               int    `json:"code"`
	Phone              string `json:"phone"`
	Active             bool   `json:"active"`
	County             string `json:"county"` // TODO: Controlled list of counties
	Description        string `json:"description"`
	FHIROrganisationID string `json:"fhirOrganisationId"`
}

//FacilityPage returns a list of paginates facilities
type FacilityPage struct {
	Pagination Pagination
	Facilities []Facility
}

// UpdateFacilityPayload is the payload for updating faacility(s) fhir organization ID
type UpdateFacilityPayload struct {
	FacilityID         string `json:"facilityID"`
	FHIROrganisationID string `json:"fhirOrganisationID"`
}
