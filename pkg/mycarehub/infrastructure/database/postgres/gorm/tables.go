package gorm

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"gorm.io/gorm"
)

// Base model contains defines commin fields across tables
type Base struct {
	CreatedAt time.Time `gorm:"column:created"`
	UpdatedAt time.Time `gorm:"column:updated"`
	// TODO: Add deleted column when all tables are updated to have it
	// DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	Base
	//globally unique when set
	FacilityID *string `gorm:"primaryKey;unique;column:facility_id"`
	// unique within this structure
	Name string `gorm:"column:name;unique;not null"`
	// MFL Code for Kenyan facilities, globally unique
	Code        string           `gorm:"unique;column:mfl_code;not null"`
	Active      string           `gorm:"column:active;not null"`
	County      enums.CountyType `gorm:"column:county;not null"` // TODO: Controlled list of counties
	Description string           `gorm:"column:description;not null"`
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

// User represents the table data structure for a user
type User struct {
	Base

	UserID *string `gorm:"primaryKey;unique;column:id"` // globally unique ID

	Username string `gorm:"column:username;unique;not null"` // @handle, also globally unique; nickname

	// Handle string `gorm:"column:handle;unique;not null"`

	DisplayName string `gorm:"column:name"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string `gorm:"column:first_name;not null"` // given name
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name;not null"`

	UserType enums.UsersType `gorm:"column:user_type;not null"` // TODO enum; e.g client, health care worker

	Gender enumutils.Gender `gorm:"column:gender;not null"` // TODO enum; genders; keep it simple

	// DateOfBirth string `gorm:"type:date;column:date_of_birth"`

	Active bool `gorm:"type:bool;column:is_active;not null"`

	// for the preferred language list, order matters
	Languages pq.StringArray `gorm:"type:text[];column:languages;not null"` // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens []string `gorm:"type:text[];column:push_tokens"`

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `gorm:"type:time;column:last_successful_login"`

	// This is a required field in DRF but we default to false
	// IsSUperUser bool `gorm:"type:bool;column:is_superuser; default:false"`

	// This is a required field in DRF but we default to a random hashed password
	// Password string `gorm:"column:password;not null"`

	// This is a required field in DRF but we default to a random email
	// Email string `gorm:"column:email;not null"`

	// This is a required field in DRF but we default to false
	// IsStaff string `gorm:"column:is_staff;not null"`

	// When true, the user is able to log in to the main website (and vice versa)
	// IsApproved string `gorm:"type:bool;column:is_approved"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `gorm:"type:time;column:last_failed_login"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount string `gorm:"column:failed_login_count"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `gorm:"type:time;column:next_allowed_login"`

	TermsAccepted   bool            `gorm:"type:bool;column:accepted_terms_of_service;not null"`
	AcceptedTermsID string          `gorm:"column:accepted_terms_id"` // foreign key to version of terms they accepted
	Flavour         feedlib.Flavour `gorm:"column:flavour;not null"`
}

// BeforeCreate is a hook run before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	u.UserID = &id
	return
}

// TableName customizes how the table name is generated
func (User) TableName() string {
	return "users_user"
}

// Contact hold contact information/details for users
type Contact struct {
	Base

	ContactID *string `gorm:"primaryKey;unique;column:id"`

	Type string `gorm:"column:contact_type;not null"` // TODO enum

	Contact string `gorm:"unique;column:contact_value;not null"` // TODO Validate: phones are E164, emails are valid

	Active bool `gorm:"column:active;not null"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `gorm:"column:opted_in;not null"`
}

// BeforeCreate is a hook run before creating a new contact
func (c *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ContactID = &id
	return
}

// TableName customizes how the table name is generated
func (Contact) TableName() string {
	return "common_contact"
}

// PINData is the PIN's gorm data model.
type PINData struct {
	Base

	PINDataID *int      `gorm:"primaryKey;unique;column:id;autoincrement"`
	UserID    string    `gorm:"column:user_id;not null"`
	HashedPIN string    `gorm:"column:hashed_pin;not null"`
	ValidFrom time.Time `gorm:"column:valid_from;not null"`
	ValidTo   time.Time `gorm:"column:valid_to;not null"`
	IsValid   bool      `gorm:"column:active;not null"`
	Salt      string    `gorm:"column:hashed_pin;not null"`
}

// BeforeCreate is a hook run before creating a new facility
// func (p *PINData) BeforeCreate(tx *gorm.DB) (err error) {
// 	id := uuid.New().String()
// 	p.PINDataID = &id
// 	return
// }

// TableName customizes how the table name is generated
func (PINData) TableName() string {
	return "users_userpin"
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

	TreatmentEnrollmentDate *time.Time `gorm:"type:time;column:enrollment_date"`

	ClientType string `gorm:"column:client_type"`

	Active bool `gorm:"column:active"`

	HealthRecordID *string `gorm:"column:health_record_id"`

	// TODO: a client can have many identifiers; an identifier belongs to a client
	// (implement reverse relation lookup)
	// Identifiers []*Identifier `gorm:"foreignKey:ClientID"`

	Contacts []*Contact `gorm:"many2many:common_contact;not null;"`

	// Addresses []*domain.Addresses `gorm:"column:addresses"`

	// RelatedPersons []*domain.RelatedPerson `gorm:"column:related_persons"`

	// client's currently assigned facility
	FacilityID string `gorm:"column:current_facility_id"` // TODO: FK

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
	return "clients_client"
}
