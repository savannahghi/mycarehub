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
	// ID is the Global facility ID(GCID)
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
// 	DataType string // TODO: Ideally a controlled list i.e enum
// 	Date     string // TODO: Clear spec on validation e.g dates must be ISO 8601
// }

// User  holds details that both the client and staff have in common
//
// Client and Staff cannot exist without being a user
type User struct {
	ID uuid.UUID // globally unique ID

	Username string // @handle, also globally unique; nickname

	DisplayName string // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string // given name
	MiddleName *string
	LastName   string

	UserType string // TODO enum; e.g client, health care worker

	Gender string // TODO enum; genders; keep it simple

	Active bool

	Contacts []*Contact // TODO: validate, ensure

	// for the preferred language list, order matters
	Languages []string // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens []string

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount int

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time

	TermsAccepted   bool
	AcceptedTermsID string // foreign key to version of terms they accepted
}

// AuthCredentials is the authentication credentials for a given user
type AuthCredentials struct {
	User *User

	RefreshToken string
	IDToken      string
	ExpiresIn    time.Time
}

// UserPIN is used to store users' PINs and their entire change history.
type UserPIN struct {
	UserID string // TODO: At the DB, this should be indexed

	HashedPIN string
	ValidFrom time.Time
	ValidTo   time.Time

	// TODO: Compute this each time an operation involving the PIN is carried out
	// 	in order to make routine things e.g login via PIN fast
	IsValid bool // TODO: Consider a composite or partial DB index with UserID, IsValid, flavour
	Flavour string
}

// Identifier are specific/unique identifiers for a user
type Identifier struct {
	ID             *string // globally unique identifier
	ClientID       string  // TODO: FK to client
	IdentifierType string  // TODO: Enum; start with basics e.g CCC number, ID number
	IdentifierUse  string  // TODO: Enum; e.g official, temporary, old (see FHIR Person for enum)

	// TODO: Validate identifier value against type e.g format of CCC number
	// TODO: Unique together: identifier value & type i.e the same identifier can't be used for more than one client
	IdentifierValue     string // the actual identifier e.g CCC number
	Description         string
	ValidFrom           *time.Time
	ValidTo             *time.Time
	Active              bool
	IsPrimaryIdentifier bool
}

// ClientProfile holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
//It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type ClientProfile struct {
	ID string // globally unique identifier; synthetic i.e has no encoded meaning

	// every client is a user first
	// biodata is linked to the user record
	// the client record is for bridging to other identifiers e.g patient record IDs
	UserID string // TODO: Foreign key to User

	TreatmentEnrollmentDate *time.Time // use for date of treatment enrollment

	ClientType string // TODO: enum; e.g PMTCT, OVC

	Active bool

	HealthRecordID *string // optional link to a health record e.g FHIR Patient ID

	// TODO: a client can have many identifiers; an identifier belongs to a client
	// (implement reverse relation lookup)
	Identifiers []*Identifier

	Addresses []*Address

	RelatedPersons []*RelatedPerson // e.g next of kin

	// client's currently assigned facility
	FacilityID string // TODO: FK

	TreatmentBuddyUserID string // TODO: optional, FK to User

	CHVUserID string // TODO: optional, FK to User

	ClientCounselled bool
}

// Address are value objects for user address e.g postal code
type Address struct {
	ID string

	Type       string // TODO: enum; postal, physical or both
	Text       string // actual address, can be multi-line
	Country    string // TODO: enum
	PostalCode string
	County     string // TODO: counties belong to a country
	Active     bool
}

// RelatedPerson holds the details for person we consider relates to a Client
//
// It servers as Next of Kin details
type RelatedPerson struct {
	ID string

	Active           bool
	RelatedTo        string // TODO: FK to client
	RelationshipType string // TODO: enum
	FirstName        string
	LastName         string
	OtherName        string // TODO: optional
	Gender           string // TODO: enum

	DateOfBirth *time.Time // TODO: optional
	Addresses   []*Address // TODO: optional
	Contacts    []*Contact // TODO: optional
}

// ClientProfileRegistrationPayload holds the registration input we need to register a client
//
// into the system. Every Client us a user first
type ClientProfileRegistrationPayload struct {
	// every client is a user first
	// biodata is linked to the user record
	// the client record is for bridging to other identifiers e.g patient record IDs
	UserID string // TODO: Foreign key to User

	ClientType string // TODO: enum; e.g PMTCT, OVC

	PrimaryIdentifier *Identifier // TODO: optional, default set if not givemn

	Addresses []*Address

	FacilityID string

	TreatmentEnrollmentDate *time.Time

	ClientCounselled bool

	// TODO: when returning to UI, calculate length of treatment (return as days for ease of use in frontend)
}

// Contact hold contact information/details for users
type Contact struct {
	ID string

	Type string // TODO enum

	Contact string // TODO Validate: phones are E164, emails are valid

	Active bool

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool
}

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
