package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ClientIdentifierType is a list of all the client identifier types.
type ClientIdentifierType string

const (
	//ClientIdentifierTypeCCC represents the CCC client identifier type
	ClientIdentifierTypeCCC ClientIdentifierType = "CCC"
	//ClientIdentifierTypeNationalID represents the national id client identifier type
	ClientIdentifierTypeNationalID ClientIdentifierType = "NATIONAL_ID"
)

// IsValid returns true if a client identifier type is valid
func (f ClientIdentifierType) IsValid() bool {
	switch f {
	case ClientIdentifierTypeCCC, ClientIdentifierTypeNationalID:
		return true
	}
	return false
}

func (f ClientIdentifierType) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a client identifier type.
func (f *ClientIdentifierType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = ClientIdentifierType(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid ClientIdentifierType", str)
	}
	return nil
}

// MarshalGQL writes the client identifier type to the supplied
func (f ClientIdentifierType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}
