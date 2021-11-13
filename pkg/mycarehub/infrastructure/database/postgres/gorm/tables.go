package gorm

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// OrganizationID assign a default organisation to a type
var OrganizationID = serverutils.MustGetEnvVar(common.OrganizationID)

// Base model contains defines commin fields across tables
type Base struct {
	CreatedAt time.Time `gorm:"column:created"`
	UpdatedAt time.Time `gorm:"column:updated"`
	// OrganisationID string    `gorm:"column:organisation_id"`
	//DeletedAt      time.Time `gorm:"column:deleted_at"`
}

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	Base
	//globally unique when set
	FacilityID *string `gorm:"primaryKey;unique;column:id"`
	// unique within this structure
	Name string `gorm:"column:name;unique;not null"`
	// MFL Code for Kenyan facilities, globally unique
	Code           int    `gorm:"unique;column:mfl_code;not null"`
	Active         bool   `gorm:"column:active;not null"`
	County         string `gorm:"column:county;not null"` // TODO: Controlled list of counties
	Description    string `gorm:"column:description;not null"`
	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *Facility) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.FacilityID = &id
	f.OrganisationID = OrganizationID
	return
}

// TableName customizes how the table name is generated
func (Facility) TableName() string {
	return common.FacilityTableName
}

// User represents the table data structure for a user
type User struct {
	// Base

	UserID *string `gorm:"primaryKey;unique;column:id"` // globally unique ID

	Username string `gorm:"column:username;unique;not null"` // @handle, also globally unique; nickname

	// DisplayName string `gorm:"column:display_name"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string `gorm:"column:first_name;not null"` // given name
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name;not null"`

	UserType enums.UsersType `gorm:"column:user_type;not null"` // TODO enum; e.g client, health care worker

	Gender enumutils.Gender `gorm:"column:gender;not null"` // TODO enum; genders; keep it simple

	Active bool `gorm:"column:is_active;not null"`

	Contacts []Contact `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"` // TODO: validate, ensure

	// // for the preferred language list, order matters
	// Languages pq.StringArray `gorm:"type:text[];column:languages;not null"` // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens []string `gorm:"type:text[];column:push_tokens"`

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `gorm:"type:time;column:last_successful_login"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `gorm:"type:time;column:last_failed_login"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount int `gorm:"column:failed_login_count"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `gorm:"type:time;column:next_allowed_login"`

	TermsAccepted   bool            `gorm:"type:bool;column:terms_accepted;not null"`
	AcceptedTermsID *int            `gorm:"column:accepted_terms_of_service_id"` // foreign key to version of terms they accepted
	Flavour         feedlib.Flavour `gorm:"column:flavour;not null"`

	// Django required fields
	OrganisationID   string `gorm:"column:organisation_id"`
	Password         string `gorm:"column:password"`
	IsSuperuser      bool   `gorm:"column:is_superuser"`
	IsStaff          bool   `gorm:"column:is_staff"`
	Email            string `gorm:"column:email"`
	DateJoined       string `gorm:"column:date_joined"`
	Name             string `gorm:"column:name"`
	IsApproved       bool   `gorm:"column:is_approved"`
	ApprovalNotified bool   `gorm:"column:approval_notified"`
	Handle           string `gorm:"column:handle"`
}

// BeforeCreate is a hook run before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	u.UserID = &id
	u.OrganisationID = OrganizationID
	salt, _ := extension.NewExternalMethodsImpl().EncryptPIN(ksuid.New().String(), nil)
	bytePass := []byte(salt)
	u.Password = string(bytePass[0:127])
	u.IsSuperuser = false
	u.IsStaff = false
	u.Email = gofakeit.Email()
	u.DateJoined = time.Now().UTC().Format(time.RFC1123Z)
	u.Name = gofakeit.Name()
	u.IsApproved = false
	u.ApprovalNotified = false
	u.Handle = "@" + u.Username

	return
}

// TableName customizes how the table name is generated
func (User) TableName() string {
	return "users_user"
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

	PINDataID *int            `gorm:"primaryKey;unique;column:id;autoincrement"`
	UserID    string          `gorm:"column:user_id;not null"`
	HashedPIN string          `gorm:"column:hashed_pin;not null"`
	ValidFrom time.Time       `gorm:"column:valid_from;not null"`
	ValidTo   time.Time       `gorm:"column:valid_to;not null"`
	IsValid   bool            `gorm:"column:active;not null"`
	Flavour   feedlib.Flavour `gorm:"column:flavour;not null"`
	Salt      string          `gorm:"column:hashed_pin;not null"`
}

// TableName customizes how the table name is generated
func (PINData) TableName() string {
	return "users_userpin"
}

// TermsOfService is the gorms terms of service model
type TermsOfService struct {
	Base

	TermsID   *int       `gorm:"primaryKey;unique;column:id"`
	Text      *string    `gorm:"column:text;not null"`
	ValidFrom *time.Time `gorm:"column:valid_from;not null"`
	ValidTo   *time.Time `gorm:"column:valid_to;not null"`
	// Django reqired fields
	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating terms of service
func (t *TermsOfService) BeforeCreate(tx *gorm.DB) (err error) {
	t.OrganisationID = OrganizationID
	return
}

// TableName customizes how the table name is generated
func (TermsOfService) TableName() string {
	return "users_termsofservice"
}

// Organisation maps the organization table. This will be useful when running integration
// tests since many models have an organization ID as a foreign key.
type Organisation struct {
	Base

	OrganisationID   *string `gorm:"primaryKey;unique;column:id"`
	Active           bool    `gorm:"column:active;not null"`
	Deleted          bool    `gorm:"column:deleted;not null"`
	OrgCode          string  `gorm:"column:org_code"`
	Code             int     `gorm:"column:code"`
	OrganisationName string  `gorm:"column:organisation_name"`
	EmailAddress     string  `gorm:"column:email_address"`
	PhoneNumber      string  `gorm:"column:phone_number"`
	PostalAddress    string  `gorm:"column:postal_address"`
	PhysicalAddress  string  `gorm:"column:physical_address"`
	DefaultCountry   string  `gorm:"column:default_country"`
}

// BeforeCreate is a hook run before creating a new organisation
func (t *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	t.OrganisationID = &OrganizationID
	return
}

// TableName customizes how the table name is generated
func (Organisation) TableName() string {
	return "common_organisation"
}
