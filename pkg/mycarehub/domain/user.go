package domain

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// User  holds details that both the client and staff have in common
//
// Client and Staff cannot exist without being a user
type User struct {
	ID *string `json:"userID"`

	Username string `json:"userName"`

	Name string `json:"name"`

	Gender enumutils.Gender `json:"gender"`
	Active bool

	Contacts *Contact `json:"primaryContact"`

	// for the preferred language list, order matters
	// Languages []enumutils.Language `json:"languages"`

	PushTokens []string `json:"pushTokens"`

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `json:"lastSuccessfulLogin"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `json:"lastFailedLogin"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount int `json:"failedLoginCount"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `json:"nextAllowedLogin"`

	PinChangeRequired bool `json:"pinChangeRequired"`

	HasSetPin              bool `json:"hasSetPin"`
	HasSetSecurityQuestion bool `json:"hasSetSecurityQuestions"`
	IsPhoneVerified        bool `json:"isPhoneVerified"`

	TermsAccepted       bool                   `json:"termsAccepted"`
	AcceptedTermsID     int                    `json:"acceptedTermsID"` // foreign key to version of terms they accepted
	Suspended           bool                   `json:"suspended"`
	Avatar              string                 `json:"avatar"`
	Roles               []*AuthorityRole       `json:"roles"`
	Permissions         []*AuthorityPermission `json:"permissions"`
	DateOfBirth         *time.Time             `json:"dateOfBirth"`
	FailedSecurityCount int                    `json:"failedSecurityCount"`
	PinUpdateRequired   bool                   `json:"pinUpdateRequired"`
	HasSetNickname      bool                   `json:"hasSetNickname"`

	CurrentOrganizationID string `json:"currentOrganizationID"`
	CurrentProgramID      string `json:"currentProgramID"`
}

// ClientProfile holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
// It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type ClientProfile struct {
	ID          *string            `json:"id"`
	User        *User              `json:"user"`
	Active      bool               `json:"active"`
	ClientTypes []enums.ClientType `json:"clientTypes"`
	UserID      string             `json:"userID"`

	TreatmentEnrollmentDate *time.Time `json:"treatmentEnrollmentDate"`

	FHIRPatientID *string `json:"fhirPatientID"`

	HealthRecordID *string `json:"healthRecordID"`

	TreatmentBuddy string `json:"treatmentBuddy"`

	ClientCounselled bool `json:"counselled"`

	OrganisationID string `json:"organisationID"`

	DefaultFacility *Facility `json:"defaultFacility"`

	CHVUserID   *string     `json:"chvUserID"`
	CHVUserName string      `json:"chvUserName"`
	CaregiverID *string     `json:"caregiverID"`
	CCCNumber   string      `json:"CCCNumber"`
	Facilities  []*Facility `json:"facilities"` // TODO: needs at least one
	ProgramID   string      `json:"programID"`
}

// StaffProfile represents the staff profile model
type StaffProfile struct {
	ID *string `json:"id"`

	User *User `json:"user"`

	UserID string `json:"user_id"` // foreign key to user

	Active bool `json:"active"`

	StaffNumber string `json:"staff_number"`

	Facilities []*Facility `json:"facilities"` // TODO: needs at least one

	DefaultFacility *Facility `json:"defaultFacility"`
	OrganisationID  string    `json:"organisationID"`
	ProgramID       string    `json:"programID"`
}

// UserPIN is used to store users' PINs and their entire change history.
type UserPIN struct {
	UserID    string    `json:"userID"`
	HashedPIN string    `json:"column:hashedPin"`
	ValidFrom time.Time `json:"column:validFrom"`
	ValidTo   time.Time `json:"column:validTo"`
	IsValid   bool      `json:"isValid"`
	Salt      string    `json:"salt"`
}

// Contact hold contact information/details for users
type Contact struct {
	ID *string `json:"id"`

	ContactType  string `json:"contactType"`
	ContactValue string `json:"contactValue"`

	Active bool `json:"active"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `json:"optedIn"`

	UserID         *string
	OrganisationID string
}

// Identifier is used to store a user's identifying details e.d ID number, CCC number
type Identifier struct {
	ID                  string    `json:"id"`
	IdentifierType      string    `json:"identifierType"`
	IdentifierValue     string    `json:"identifierValue"`
	IdentifierUse       string    `json:"identifierUse"`
	Description         string    `json:"description"`
	ValidFrom           time.Time `json:"validFrom"`
	ValidTo             time.Time `json:"validTo"`
	IsPrimaryIdentifier bool      `json:"isPrimaryIdentifier"`
	Active              bool      `json:"active"`
	ProgramID           string    `json:"programID"`
}

// ClientRegistrationPayload is the payload for a client registration
type ClientRegistrationPayload struct {
	UserProfile      User
	Phone            Contact
	ClientIdentifier Identifier
	Client           ClientProfile
}

// CaregiverRegistration is the input used for creating a caregiver
type CaregiverRegistration struct {
	User      *User      `json:"user"`
	Contact   *Contact   `json:"contact"`
	Caregiver *Caregiver `json:"caregiver"`
}

// StaffRegistrationPayload carries with it the staff registration details
type StaffRegistrationPayload struct {
	UserProfile     User
	Phone           Contact
	StaffIdentifier Identifier
	Staff           StaffProfile
}
