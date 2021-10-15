package dto

import (
	"net/url"
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"gorm.io/datatypes"
)

// FacilityInput describes the facility input
type FacilityInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Active      bool   `json:"active"`
	County      string `json:"county"`
	Description string `json:"description"`
}

// FacilityFilterInput is used to supply filter parameters for healthcare facility filter inputs
type FacilityFilterInput struct {
	Search  *string `json:"search"`
	Name    *string `json:"name"`
	MFLCode *string `json:"code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *FacilityFilterInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.Name != nil {
		vals.Add("name", *i.Name)
	}
	if i.MFLCode != nil {
		vals.Add("code", *i.MFLCode)
	}
	return vals
}

// FacilitySortInput is used to supply sort input for healthcare facility list queries
type FacilitySortInput struct {
	Name    *enumutils.SortOrder `json:"name"`
	MFLCode *enumutils.SortOrder `json:"code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *FacilitySortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == enumutils.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.MFLCode != nil {
		if *i.Name == enumutils.SortOrderAsc {
			vals.Add("code", "number")
		} else {
			vals.Add("code", "-number")
		}
	}
	return vals
}

// MetricInput reprents the metrics data structure input
type MetricInput struct {

	// TODO Metric types should be a controlled list i.e enum
	Type domain.MetricType `json:"metric_type"`

	// this will vary by context
	// should not identify the user (there's a UID field)
	// focus on the actual event
	Payload datatypes.JSON `gorm:"column:payload"`

	Timestamp time.Time `json:"time"`

	// a user identifier, can be hashed for anonymity
	// with a predictable one way hash
	UID string `json:"uid"`
}

// UserInput contains the input for the User
type UserInput struct {
	Username string // @handle, also globally unique; nickname

	DisplayName string // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string // given name
	MiddleName *string
	LastName   string

	UserType string // TODO enum; e.g client, health care worker

	Gender string // TODO enum; genders; keep it simple

	Contacts []*domain.Contact // TODO: validate, ensure

	// for the preferred language list, order matters
	Languages []string // TODO: turn this into a slice of enums, start small (en, sw)
}

// StaffProfileInput contains input required to register a staff
type StaffProfileInput struct {
	StaffNumber string

	Facilities []*domain.Facility // TODO: needs at least one

	// A UI switcher optionally toggles the default
	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
	DefaultFacilityID *string // TODO: required, FK to facility

	// there is nothing special about super-admin; just the set of roles they have
	Roles []domain.RoleType `json:"roles"` // TODO: roles are an enum (controlled list), known to both FE and BE

	Addresses []*domain.UserAddress
}

// ContactInput hold contact information/details for users
type ContactInput struct {
	Type string // TODO enum

	Contact string // TODO Validate: phones are E164, emails are valid

	Active bool

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool
}

// ContactAddressInput are value objects for user address e.g postal code
type ContactAddressInput struct {
	Type       string // TODO: enum; postal, physical or both
	Text       string // actual address, can be multi-line
	Country    string // TODO: enum
	PostalCode string
	County     string // TODO: counties belong to a country
	Active     bool
}
