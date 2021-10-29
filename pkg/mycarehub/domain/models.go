package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
	"gorm.io/datatypes"
)

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	// ID is the Global facility ID(GCID)
	ID *string `json:"id"`
	// unique within this structure
	Name string `json:"name"`
	// MFL Code for Kenyan facilities, globally unique
	Code        string `json:"code"`
	Active      bool   `json:"active"`
	County      string `json:"county"` // TODO: Controlled list of counties
	Description string `json:"description"`
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

// User  holds details that both the client and staff have in common
//
// Client and Staff cannot exist without being a user
type User struct {
	ID *string `json:"id"` // globally unique ID

	Username string `json:"username"` // @handle, also globally unique; nickname

	DisplayName string `json:"displayName"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string `json:"firstName"` // given name
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`

	UserType enums.UsersType `json:"userType"`

	Gender enumutils.Gender `json:"gender"`
	Active bool

	Contacts []*Contact `json:"contact"` // TODO: validate, ensure

	// for the preferred language list, order matters
	Languages []enumutils.Language `json:"languages"`

	// PushTokens []string

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `json:"lastSuccessfulLogin"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `json:"lastFailedLogin"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount string `json:"failedLoginCount"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `json:"NextAllowedLogin"`

	TermsAccepted   bool            `json:"termsAccepted"`
	AcceptedTermsID string          `json:"AcceptedTermsID"` // foreign key to version of terms they accepted
	Flavour         feedlib.Flavour `json:"flavour"`
}

// AuthCredentials is the authentication credentials for a given user
type AuthCredentials struct {
	User *User `json:"user"`

	RefreshToken string `json:"refreshToken"`
	IDToken      string `json:"idToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// UserPIN is used to store users' PINs and their entire change history.
type UserPIN struct {
	UserID    string          `json:"userID"`
	HashedPIN string          `json:"column:hashedPin"`
	ValidFrom time.Time       `json:"column:validFrom"`
	ValidTo   time.Time       `json:"column:validTo"`
	Flavour   feedlib.Flavour `json:"flavour"`
	IsValid   bool            `json:"isValid"`
	Salt      string          `json:"salt"`
}

// Identifier are specific/unique identifiers for a user
type Identifier struct {
	ID             *string              `json:"id"`
	ClientID       string               `json:"clientID"`
	IdentifierType enums.IdentifierType `json:"identifierType"`
	IdentifierUse  enums.IdentifierUse  `json:"identifierUse"`

	// TODO: Validate identifier value against type e.g format of CCC number
	// TODO: Unique together: identifier value & type i.e the same identifier can't be used for more than one client
	IdentifierValue     string     `json:"identifierValue"`
	Description         string     `json:"description"`
	ValidFrom           *time.Time `json:"validFrom"`
	ValidTo             *time.Time `json:"validTo"`
	Active              bool       `json:"active"`
	IsPrimaryIdentifier bool       `json:"isPrimaryIdentifier"`
}

// ClientProfile holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
//It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type ClientProfile struct {
	ID *string `json:"id"` // globally unique identifier; synthetic i.e has no encoded meaning

	// every client is a user first
	// biodata is linked to the user record
	// the client record is for bridging to other identifiers e.g patient record IDs
	UserID *string `json:"userID"`

	TreatmentEnrollmentDate *time.Time `json:"treatmentEnrollmentDate"` // use for date of treatment enrollment

	ClientType enums.ClientType `json:"ClientType"`

	Active bool `json:"Active"`

	HealthRecordID *string `json:"healthRecordID"` // optional link to a health record e.g FHIR Patient ID

	// TODO: a client can have many identifiers; an identifier belongs to a client
	// (implement reverse relation lookup)
	Identifiers []*Identifier `json:"identifiers"`

	Addresses []*Addresses `json:"addresses"`

	RelatedPersons []*RelatedPerson `json:"relatedPersons"` // e.g next of kin

	// client's currently assigned facility
	FacilityID string `json:"facilityID"` // TODO: FK

	TreatmentBuddyUserID string `json:"treatmentBuddyUserID"` // TODO: optional, FK to User

	CHVUserID string `json:"CHVUserID"` // TODO: optional, FK to User

	ClientCounselled bool `json:"clientCounselled"`
}

// Addresses are value objects for user address e.g postal code
type Addresses struct {
	ID         string              `json:"id"`
	Type       enums.AddressesType `json:"type"`
	Text       string              `json:"text"` // actual address, can be multi-line
	Country    enums.CountryType   `json:"country"`
	PostalCode string              `json:"postalCode"`
	County     enums.CountyType    `json:"county"`
	Active     bool                `json:"active"`
}

// RelatedPerson holds the details for person we consider relates to a Client
//
// It servers as Next of Kin details
type RelatedPerson struct {
	ID *string

	Active           bool
	RelatedTo        string // TODO: FK to client
	RelationshipType string // TODO: enum
	FirstName        string
	LastName         string
	OtherName        string // TODO: optional
	Gender           string // TODO: enum

	DateOfBirth *scalarutils.Date // TODO: optional
	Addresses   []*Addresses      // TODO: optional
	Contacts    []*Contact        // TODO: optional
}

// ClientProfileRegistrationPayload holds the registration input we need to register a client
//
// into the system. Every Client us a user first
type ClientProfileRegistrationPayload struct {
	// every client is a user first
	// biodata is linked to the user record
	// the client record is for bridging to other identifiers e.g patient record IDs
	UserID *string // TODO: Foreign key to User

	ClientType string // TODO: enum; e.g PMTCT, OVC

	PrimaryIdentifier *Identifier // TODO: optional, default set if not givemn

	Addresses []*Addresses

	FacilityID uuid.UUID

	TreatmentEnrollmentDate *time.Time

	ClientCounselled bool

	// TODO: when returning to UI, calculate length of treatment (return as days for ease of use in frontend)
}

// Contact hold contact information/details for users
type Contact struct {
	ID *string

	Type    enums.ContactType
	Contact string // TODO Validate: phones are E164, emails are valid

	Active bool

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool
}

// Metric reprents the metrics data structure input
type Metric struct {
	// ensures we don't re-save the same metric; opaque; globally unique
	MetricID *string

	Type enums.MetricType

	// this will vary by context
	// should not identify the user (there's a UID field)
	// focus on the actual event
	Payload datatypes.JSON `gorm:"column:payload"`

	Timestamp time.Time

	// a user identifier, can be hashed for anonymity
	// with a predictable one way hash
	UID string
}

// StaffProfile contains all the information a staff should have about themselves
type StaffProfile struct {
	ID *string

	UserID *string // foreign key to user

	StaffNumber string

	Facilities []*Facility // TODO: needs at least one

	// A UI switcher optionally toggles the default
	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
	DefaultFacilityID *string

	// there is nothing special about super-admin; just the set of roles they have
	Roles []enums.RolesType // TODO: roles are an enum (controlled list), known to both FE and BE

	Addresses []*Addresses
}

// Pin model contain the information about given PIN
type Pin struct {
	UserID       string
	PIN          string
	ConfirmedPin string
	Flavour      feedlib.Flavour
}

// StaffUserProfile combines user and staff profile
type StaffUserProfile struct {
	User  *User
	Staff *StaffProfile
}

// ClientUserProfile represents the clients profile and associated user profile
type ClientUserProfile struct {
	User   *User          `json:"user"`
	Client *ClientProfile `json:"client"`
}
