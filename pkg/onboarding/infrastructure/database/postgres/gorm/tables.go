package gorm

import (
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
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
	FacilityID *uuid.UUID `gorm:"primaryKey;unique;column:facility_id"`
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
	id := uuid.New()
	f.FacilityID = &id
	return
}

// TableName customizes how the table name is generated
func (Facility) TableName() string {
	return "facility"
}

func allTables() []interface{} {
	tables := []interface{}{
		&Facility{},
		&Metric{},
	}
	return tables
}

// Metric reprents the metrics data structure input
type Metric struct {
	Base

	// ensures we don't re-save the same metric; opaque; globally unique
	MetricID *uuid.UUID `gorm:"primaryKey;autoIncrement:true;unique;column:metric_id"`

	// TODO Metric types should be a controlled list i.e enum
	Type domain.MetricType `gorm:"column:metric_type"`

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
	id := uuid.New()
	m.MetricID = &id
	return
}

// TableName customizes how the table name is generated
func (Metric) TableName() string {
	return "metric"
}

type User struct {
	Base
	//globally unique when set
	UserID   *uuid.UUID `gorm:"primaryKey;unique;column:user_id"`
	Username string     `gorm:"column:username"` // @handle, also globally unique; nickname

	DisplayName string `gorm:"column:displayName"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string  `gorm:"column:firstName"` // given name
	MiddleName *string `gorm:"column:middleName"`
	LastName   string  `gorm:"column:lastName"`

	UserType domain.UserType `gorm:"column:userType"` // TODO enum; e.g client, health care worker

	Gender enumutils.Gender `gorm:"column:gender"`

	Active bool `gorm:"column:active"`

	Contacts []*string `gorm:"column:contacts"` // TODO: validate, ensure. This should come from contacts

	// for the preferred language list, order matters
	Languages []string `gorm:"column:languages"` // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens []string `gorm:"column:push_tokens"`

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `gorm:"column:last_successful_login"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `gorm:"column:last_failed_login"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount int `gorm:"column:failed_login_count"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `gorm:"column:next_allowed_login"`

	TermsAccepted   bool   `gorm:"column:terms_accepted"`
	AcceptedTermsID string `gorm:"column:accepted_terms_id"`
}
