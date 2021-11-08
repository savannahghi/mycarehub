package enums

import (
	"fmt"
	"io"
	"strconv"
)

// UsersType is a list of all the user types.
type UsersType string

// contacts type constants
const (
	HealthcareWorkerUser UsersType = "HEALTHCAREWORKER"
	ClientUser           UsersType = "CLIENT"
)

// AllUsers is a set of a  valid and known user types.
var AllUsers = []UsersType{
	HealthcareWorkerUser,
	ClientUser,
}

// IsValid returns true if a user is valid
func (m UsersType) IsValid() bool {
	switch m {
	case HealthcareWorkerUser, ClientUser:
		return true
	}
	return false
}

func (m UsersType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a user type.
func (m *UsersType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = UsersType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid UsersType", str)
	}
	return nil
}

// MarshalGQL writes the user type to the supplied
func (m UsersType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
