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

	DisplayName string `json:"displayName"`

	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`

	UserType enums.UsersType `json:"userType"`

	Gender enumutils.Gender `json:"gender"`
	Active bool

	Contacts *Contact `json:"primaryContact"`

	// for the preferred language list, order matters
	// Languages []enumutils.Language `json:"languages"`

	// PushTokens []string

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `json:"lastSuccessfulLogin"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `json:"lastFailedLogin"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount int `json:"failedLoginCount"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `json:"NextAllowedLogin"`

	PinChangeRequired bool `json:"pinChangeRequired"`

	HasSetPin              bool `json:"hasSetPin"`
	HasSetSecurityQuestion bool `json:"hasSetSecurityQuestions"`
	IsPhoneVerified        bool `json:"isPhoneVerified"`

	TermsAccepted   bool            `json:"termsAccepted"`
	AcceptedTermsID int             `json:"AcceptedTermsID"` // foreign key to version of terms they accepted
	Flavour         feedlib.Flavour `json:"flavour"`
	Suspended       bool            `json:"suspended"`
	Avatar          string          `json:"avatar"`
}

// ClientProfile holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
// It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type ClientProfile struct {
	ID         *string `json:"id"`
	User       *User   `json:"user"`
	Active     bool    `json:"Active"`
	ClientType string  `json:"ClientType"`
	UserID     string  `json:"userID"`

	TreatmentEnrollmentDate *time.Time `json:"treatmentEnrollmentDate"`

	FHIRPatientID string `json:"fhirPatientID"`

	HealthRecordID *string `json:"healthRecordID"`

	TreatmentBuddy string `json:"treatmentBuddy"`

	ClientCounselled bool `json:"counselled"`

	OrganisationID string `json:"organisationID"`

	FacilityID string `json:"facilityID"`

	CHVUserID string `json:"CHVUserID"`
}

// AuthCredentials is the authentication credentials for a given user
type AuthCredentials struct {
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

// Contact hold contact information/details for users
type Contact struct {
	ID *string `json:"id"`

	ContactType  string `json:"contactType"`
	ContactValue string `json:"contactValue"`

	Active bool `json:"active"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `json:"optedIn"`
}

// LoginResponse models the response that will be returned after a user logs in
type LoginResponse struct {
	Client          *ClientProfile  `json:"clientProfile"`
	AuthCredentials AuthCredentials `json:"credentials"`
	Code            int             `json:"code"`
	Message         string          `json:"message"`
}
