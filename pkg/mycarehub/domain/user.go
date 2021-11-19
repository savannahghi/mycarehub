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

	TermsAccepted   bool            `json:"termsAccepted"`
	AcceptedTermsID int             `json:"AcceptedTermsID"` // foreign key to version of terms they accepted
	Flavour         feedlib.Flavour `json:"flavour"`
	Suspended       bool            `json:"suspended"`
	Avatar          string          `json:"avatar"`
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

// Contact hold contact information/details for users
type Contact struct {
	ID *string

	Type    string //TODO: Make this an enum
	Contact string // TODO Validate: phones are E164, emails are valid

	Active bool

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool
}
