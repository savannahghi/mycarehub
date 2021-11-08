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

	UserID *string `gorm:"primaryKey;unique;column:user_id"` // globally unique ID

	Username string `gorm:"column:username;unique;not null"` // @handle, also globally unique; nickname

	DisplayName string `gorm:"column:display_name"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string `gorm:"column:first_name;not null"` // given name
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name;not null"`

	UserType enums.UsersType `gorm:"column:user_type;not null"` // TODO enum; e.g client, health care worker

	Gender enumutils.Gender `gorm:"column:gender;not null"` // TODO enum; genders; keep it simple

	Active bool `gorm:"column:active;not null"`

	Contacts []Contact `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"` // TODO: validate, ensure

	// // for the preferred language list, order matters
	Languages pq.StringArray `gorm:"type:text[];column:languages;not null"` // TODO: turn this into a slice of enums, start small (en, sw)

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

	TermsAccepted   bool            `gorm:"type:bool;column:terms_accepted;not null"`
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
	return "user_users"
}

// Contact hold contact information/details for users
type Contact struct {
	Base

	ContactID *string `gorm:"primaryKey;unique;column:contact_id"`

	Type string `gorm:"column:type;not null"` // TODO enum

	Contact string `gorm:"unique;column:contact;not null"` // TODO Validate: phones are E164, emails are valid

	Active bool `gorm:"column:active;not null"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `gorm:"column:opted_in;not null"`

	UserID *string `gorm:"column:user_id;not null"` // Foreign key
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

// PINData is the PIN's gorm data model.
type PINData struct {
	Base

	PINDataID *string         `gorm:"primaryKey;unique;column:pin_data_id"`
	UserID    string          `gorm:"column:user_id;not null"`
	HashedPIN string          `gorm:"column:hashed_pin;not null"`
	ValidFrom time.Time       `gorm:"column:valid_from;not null"`
	ValidTo   time.Time       `gorm:"column:valid_to;not null"`
	IsValid   bool            `gorm:"column:is_valid;not null"`
	Flavour   feedlib.Flavour `gorm:"column:flavour;not null"`
	Salt      string          `gorm:"column:salt;not null"`
}

// BeforeCreate is a hook run before creating a new facility
func (p *PINData) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	p.PINDataID = &id
	return
}

// TableName customizes how the table name is generated
func (PINData) TableName() string {
	return "pindata"
}

func allTables() []interface{} {
	tables := []interface{}{
		&Facility{},
		&User{},
		&Contact{},
		&PINData{},
	}
	return tables
}
