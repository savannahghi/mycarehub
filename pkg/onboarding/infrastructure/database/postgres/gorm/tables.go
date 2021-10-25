package gorm

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Base model contains defines commin fields across tables
type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	Base
	//globally unique when set
	FacilityID *string `gorm:"primaryKey;unique;column:facility_id"`
	// unique within this structure
	Name string `gorm:"column:name"`
	// MFL Code for Kenyan facilities, globally unique
	Code        string `gorm:"unique;column:mfl_code"`
	Active      string `gorm:"column:active"`
	County      string `gorm:"column:county"` // TODO: Controlled list of counties
	Description string `gorm:"column:description"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *Facility) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.FacilityID = &id
	return
}

// TableName customizes how the table name is generated
func (Facility) TableName() string {
	return "facility"
}

// Metric reprents the metrics data structure input
type Metric struct {
	Base

	// ensures we don't re-save the same metric; opaque; globally unique
	MetricID *string `gorm:"primaryKey;autoIncrement:true;unique;column:metric_id"`

	// TODO Metric types should be a controlled list i.e enum
	Type enums.MetricType `gorm:"column:metric_type"`

	// this will vary by context
	// should not identify the user (there's a UID field)
	// focus on the actual event
	Payload datatypes.JSON `gorm:"column:payload"`

	Timestamp time.Time `gorm:"column:timestamp"`

	// a user identifier, can be hashed for anonymity
	// with a predictable one way hash
	UID string `gorm:"column:uid"`
}

// BeforeCreate is a hook run before creating a new facility
func (m *Metric) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	m.MetricID = &id
	return
}

// TableName customizes how the table name is generated
func (Metric) TableName() string {
	return "metric"
}

// User represents the table data structure for a user
type User struct {
	Base

	UserID *string `gorm:"primaryKey;unique;column:user_id"` // globally unique ID

	Username string `gorm:"column:username"` // @handle, also globally unique; nickname

	DisplayName string `gorm:"column:display_name"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string `gorm:"column:first_name"` // given name
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`

	UserType enums.UsersType `gorm:"column:user_type"` // TODO enum; e.g client, health care worker

	Gender enumutils.Gender `gorm:"column:gender"` // TODO enum; genders; keep it simple

	Active bool `gorm:"column:active"`

	Contacts []Contact `gorm:"ForeignKey:UserID"` // TODO: validate, ensure

	// // for the preferred language list, order matters
	Languages pq.StringArray `gorm:"type:text[];column:languages"` // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens []string `gorm:"type:text[];column:push_tokens"`

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `gorm:"type:time;column:last_successful_login"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `gorm:"type:time;column:last_failed_login"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount string `gorm:"column:failed_login_count"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `gorm:"type:time;column:next_allowed_login"`

	TermsAccepted   bool            `gorm:"type:bool;column:terms_accepted"`
	AcceptedTermsID string          `gorm:"column:accepted_terms_id"` // foreign key to version of terms they accepted
	Flavour         feedlib.Flavour `gorm:"column:flavour"`
}

// BeforeCreate is a hook run before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	u.UserID = &id
	return
}

// TableName customizes how the table name is generated
func (User) TableName() string {
	return "user"
}

// Contact hold contact information/details for users
type Contact struct {
	Base

	ContactID *string `gorm:"primaryKey;unique;column:contact_id"`

	Type enums.ContactType `gorm:"column:type"` // TODO enum

	Contact string `gorm:"unique;column:contact"` // TODO Validate: phones are E164, emails are valid

	Active bool `gorm:"column:active"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `gorm:"column:opted_in"`

	UserID *string `gorm:"column:user_id"` // Foreign key
}

// BeforeCreate is a hook run before creating a new contact
func (c *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ContactID = &id
	return
}

// TableName customizes how the table name is generated
func (Contact) TableName() string {
	return "contact"
}

// StaffProfile contains all the information a staff should have about themselves
type StaffProfile struct {
	StaffProfileID *string `gorm:"primaryKey;unique;column:staff_profile_id"`

	UserID *string `gorm:"unique;column:user_id"` // foreign key to user
	User   User    `gorm:"foreignKey:user_id;references:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	StaffNumber string `gorm:"unique;column:staff_number"`

	// Facilities []*Facility `gorm:"many2many:staffprofile_facility;not null;"` // TODO: needs at least one

	// A UI switcher optionally toggles the default
	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
	DefaultFacilityID *string  `gorm:"column:default_facility_id"` // TODO: required, FK to facility
	Facility          Facility `gorm:"foreignKey:default_facility_id;references:facility_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// // there is nothing special about super-admin; just the set of roles they have
	// Roles []string `gorm:"type:text[];column:roles"` // TODO: roles are an enum (controlled list), known to both FE and BE

	// Addresses []*UserAddress `gorm:"many2many:staffprofile_useraddress;"`
}

// BeforeCreate is a hook run before creating a new staff profile
func (s *StaffProfile) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	s.StaffProfileID = &id
	return
}

// TableName customizes how the table name is generated
func (StaffProfile) TableName() string {
	return "staffprofile"
}

// UserAddress are value objects for user address e.g postal code
type UserAddress struct {
	UserAddressID *string `gorm:"primaryKey;unique;column:useraddress_id"` // globally unique

	Type       string `gorm:"column:type"`    // TODO: enum; postal, physical or both
	Text       string `gorm:"column:text"`    // actual address, can be multi-line
	Country    string `gorm:"column:country"` // TODO: enum
	PostalCode string `gorm:"column:postal_code"`
	County     string `gorm:"column:county"` // TODO: counties belong to a country
	Active     bool   `gorm:"column:active"`
}

// BeforeCreate is a hook run before creating a new address
func (a *UserAddress) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	a.UserAddressID = &id
	return
}

// TableName customizes how the table name is generated
func (UserAddress) TableName() string {
	return "useraddress"
}

// StaffUserProfile combines user and staff profile
type StaffUserProfile struct {
	User  *User
	Staff *StaffProfile
}

// ClientProfile holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
// It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type ClientProfile struct {
	Base

	ID *string `gorm:"primaryKey;unique;column:id"`

	UserID *string `gorm:"unique;column:user_id"`
	User   User    `gorm:"foreignKey:user_id;references:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	TreatmentEnrollmentDate *time.Time

	ClientType enums.ClientType `gorm:"column:client_type"`

	Active bool `gorm:"column:active"`

	HealthRecordID *string `gorm:"column:health_record_id"`

	// TODO: a client can have many identifiers; an identifier belongs to a client
	// (implement reverse relation lookup)
	Identifiers []*Identifier `gorm:"foreignKey:ClientID"`

	// Addresses []*domain.UserAddress `gorm:"column:addresses"`

	// RelatedPersons []*domain.RelatedPerson `gorm:"column:related_persons"`

	// client's currently assigned facility
	FacilityID string `gorm:"column:facility_id"` // TODO: FK

	TreatmentBuddy string `gorm:"column:treatment_buddy"` // TODO: optional, free text OR FK to user?

	CHVUserID string `gorm:"column:chvuser_id"` // TODO: optional, FK to User

	ClientCounselled bool `gorm:"column:client_counselled"`
}

// BeforeCreate is a hook run before creating a client profile
func (c *ClientProfile) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = &id
	return
}

// TableName customizes how the table name is generated
func (ClientProfile) TableName() string {
	return "clientprofile"
}

// ClientUserProfile holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc together with their user profile
type ClientUserProfile struct {
	User   *User
	Client *ClientProfile
}

// PINData is the PIN's gorm data model.
type PINData struct {
	Base

	PINDataID *uuid.UUID      `gorm:"primaryKey;unique;column:pin_data_id"`
	UserID    string          `gorm:"unique;column:user_id"`
	HashedPIN string          `gorm:"column:hashed_pin"`
	ValidFrom time.Time       `gorm:"column:valid_from"`
	ValidTo   time.Time       `gorm:"column:valid_to"`
	IsValid   bool            `gorm:"column:is_valid"`
	Flavour   feedlib.Flavour `gorm:"column:flavour"`
	Salt      string          `gorm:"column:salt"`
}

// BeforeCreate is a hook run before creating a new facility
func (p *PINData) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New()
	p.PINDataID = &id
	return
}

// TableName customizes how the table name is generated
func (PINData) TableName() string {
	return "pindata"
}

// Identifier are specific/unique identifiers for a user
type Identifier struct {
	Base

	ID *string `gorm:"primaryKey;unique;column:id"`

	ClientID string `gorm:"column:client_id"`

	IdentifierType enums.IdentifierType `gorm:"identifier_type"`
	IdentifierUse  enums.IdentifierUse  `gorm:"identifier_use"`

	// TODO: Validate identifier value against type e.g format of CCC number
	// TODO: Unique together: identifier value & type i.e the same identifier can't be used for more than one client
	IdentifierValue     string     `gorm:"identifier_value"`
	Description         string     `gorm:"description"`
	ValidFrom           *time.Time `gorm:"valid_from"`
	ValidTo             *time.Time `gorm:"valid_to"`
	Active              bool       `gorm:"active"`
	IsPrimaryIdentifier bool       `gorm:"is_primary_identifier"`
}

// BeforeCreate is a hook run before creating a client profile
func (c *Identifier) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = &id
	return
}

// TableName customizes how the table name is generated
func (Identifier) TableName() string {
	return "client_clientidentifier"
}

func allTables() []interface{} {
	tables := []interface{}{
		&Facility{},
		&Metric{},
		&User{},
		&Contact{},
		&StaffProfile{},
		&ClientProfile{},
		&Identifier{},
		&UserAddress{},
		&PINData{},
	}
	return tables
}
