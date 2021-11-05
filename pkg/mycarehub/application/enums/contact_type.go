package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ContactType is a list of all the contact types.
type ContactType string

// contacts type constants
const (
	PhoneContact ContactType = "PHONE"
	EmailContact ContactType = "EMAIL"
)

// AllContacts is a set of a  valid and known contact types.
var AllContacts = []ContactType{
	PhoneContact,
	EmailContact,
}

// IsValid returns true if a contact is valid
func (m ContactType) IsValid() bool {
	switch m {
	case PhoneContact, EmailContact:
		return true
	}
	return false
}

func (m ContactType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a contact type.
func (m *ContactType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = ContactType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid ContactType", str)
	}
	return nil
}

// MarshalGQL writes the contact type to the supplied
func (m ContactType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
