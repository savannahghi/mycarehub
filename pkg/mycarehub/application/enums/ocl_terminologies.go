package enums

import (
	"fmt"
	"io"
	"strconv"
)

// Terminologies is a list of all the clinical terminologies available
type Terminologies string

// Terminologies constants
const (
	TerminologiesCIEL Terminologies = "CIEL"
)

// IsValid returns true if a user is valid
func (m Terminologies) IsValid() bool {
	switch m {
	case TerminologiesCIEL:
		return true
	}
	return false
}

func (m Terminologies) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a user type.
func (m *Terminologies) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = Terminologies(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid terminology", str)
	}
	return nil
}

// MarshalGQL writes the user type to the supplied
func (m Terminologies) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
