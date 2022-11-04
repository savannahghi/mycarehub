package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ConsentState is a custom type that defines consentState that a client caregiver can be in
type ConsentState string

const (
	// ConsentStateAccepted represents a state where a client or caregiver has accepted a consentState
	ConsentStateAccepted ConsentState = "ACCEPTED"

	// ConsentStatePending represents a state where a client or caregiver has not yet acted upon consentState
	ConsentStatePending ConsentState = "PENDING"

	// ConsentStateRejected represents a state where a client or caregiver has rejected a consentState
	ConsentStateRejected ConsentState = "REJECTED"
)

// IsValid returns true if a consentState type is valid
func (c ConsentState) IsValid() bool {
	switch c {
	case ConsentStateAccepted, ConsentStatePending, ConsentStateRejected:
		return true
	}
	return false
}

// String converts the consentState type enum to a string
func (c ConsentState) String() string {
	return string(c)
}

// UnmarshalGQL converts the supplied value to a consentState type.
func (c *ConsentState) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*c = ConsentState(str)
	if !c.IsValid() {
		return fmt.Errorf("%s is not a valid consentState type", str)
	}
	return nil
}

// MarshalGQL writes the consentState type to the supplied writer
func (c ConsentState) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(c.String()))
}
