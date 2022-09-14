package domain

import (
	"time"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

// User  holds details that both the client and staff have in common
//
// Client and Staff cannot exist without being a user
type User struct {
	ID *string `json:"userID"`

	Username string `json:"userName"`

	UserType enums.UsersType `json:"userType"`
	Name     string          `json:"name"`

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
	Flavour             feedlib.Flavour        `json:"flavour"`
	Suspended           bool                   `json:"suspended"`
	Avatar              string                 `json:"avatar"`
	Roles               []*AuthorityRole       `json:"roles"`
	Permissions         []*AuthorityPermission `json:"permissions"`
	DateOfBirth         *time.Time             `json:"dateOfBirth"`
	FailedSecurityCount int                    `json:"failedSecurityCount"`
	PinUpdateRequired   bool                   `json:"pinUpdateRequired"`
	HasSetNickname      bool                   `json:"hasSetNickname"`
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

	FacilityID   string `json:"facilityID"`
	FacilityName string `json:"facilityName"`

	CHVUserID   *string    `json:"chvUserID"`
	CHVUserName string     `json:"chvUserName"`
	CaregiverID *string    `json:"caregiverID"`
	CCCNumber   string     `json:"CCCNumber"`
	Facilities  []Facility `json:"facilities"` // TODO: needs at least one
}

// StaffProfile represents the staff profile model
type StaffProfile struct {
	ID *string `json:"id"`

	User *User `json:"user"`

	UserID string `json:"user_id"` // foreign key to user

	Active bool `json:"active"`

	StaffNumber string `json:"staff_number"`

	Facilities []Facility `json:"facilities"` // TODO: needs at least one

	// A UI switcher optionally toggles the default
	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
	DefaultFacilityID   string `json:"default_facility"` // TODO: required, FK to facility
	DefaultFacilityName string `json:"defaultFacilityName"`
}

// AuthCredentials is the authentication credentials for a given user
type AuthCredentials struct {
	RefreshToken string `json:"refreshToken"`
	IDToken      string `json:"idToken"`
	ExpiresIn    string `json:"expiresIn"`
}

// GetStreamToken models the response received when generating a getstream token
type GetStreamToken struct {
	Token string `json:"getStreamToken"`
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

// Contact hold contact information/details for users
type Contact struct {
	ID *string `json:"id"`

	ContactType  string `json:"contactType"`
	ContactValue string `json:"contactValue"`

	Active bool `json:"active"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `json:"optedIn"`

	UserID  *string
	Flavour feedlib.Flavour
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
}

// ClientRegistrationPayload is the payload for a client registration
type ClientRegistrationPayload struct {
	UserProfile      User
	Phone            Contact
	ClientIdentifier Identifier
	Client           ClientProfile
}

// StaffRegistrationPayload carries with it the staff registration details
type StaffRegistrationPayload struct {
	UserProfile     User
	Phone           Contact
	StaffIdentifier Identifier
	Staff           StaffProfile
}
