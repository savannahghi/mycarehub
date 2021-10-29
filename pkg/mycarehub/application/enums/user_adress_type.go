package enums

import (
	"fmt"
	"io"
	"strconv"
)

// AddressesType defines various addresses available
type AddressesType string

const (
	// AddressesTypePostal adderess country type POSTAL
	AddressesTypePostal AddressesType = "POSTAL"
	// AddressesTypePhysical adderess country type PHYSICAL
	AddressesTypePhysical AddressesType = "PHYSICAL"
	// AddressesTypePostalPhysical adderess country type POSTALPHYSICAL
	AddressesTypePostalPhysical AddressesType = "POSTALPHYSICAL"
)

// AllAddresses represents a slice of all addresses available
var AllAddresses = []AddressesType{
	AddressesTypePostal,
	AddressesTypePhysical,
	AddressesTypePostalPhysical,
}

// IsValid returns true if an address type is valid
func (e AddressesType) IsValid() bool {
	switch e {
	case AddressesTypePostal, AddressesTypePhysical, AddressesTypePostalPhysical:
		return true
	}
	return false
}

// String converts addresses type to string.
func (e AddressesType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a user address type.
func (e *AddressesType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AddressesType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AddressesType", str)
	}
	return nil
}

// MarshalGQL writes the user address type to the supplied writer
func (e AddressesType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
