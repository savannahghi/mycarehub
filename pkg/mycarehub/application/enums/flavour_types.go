package enums

import (
	"fmt"
	"io"
	"strconv"
)

// Flavour is a list of all the flavours.
type Flavour string

// contacts type constants
const (
	PRO      Flavour = "PRO"
	CONSUMER Flavour = "CONSUMER"
)

// AllFlavours is a set of a  valid and known user types.
var AllFlavours = []Flavour{
	PRO,
	CONSUMER,
}

// IsValid returns true if a user is valid
func (m Flavour) IsValid() bool {
	switch m {
	case PRO, CONSUMER:
		return true
	}
	return false
}

func (m Flavour) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a user type.
func (m *Flavour) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = Flavour(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid flavour", str)
	}
	return nil
}

// MarshalGQL writes the user type to the supplied
func (m Flavour) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
