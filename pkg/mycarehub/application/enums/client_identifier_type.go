package enums

import (
	"fmt"
	"io"
	"strconv"
)

// UserIdentifierType is a list of all the user identifier types.
type UserIdentifierType string

const (
	//UserIdentifierTypeCCC represents the CCC user identifier type
	UserIdentifierTypeCCC UserIdentifierType = "CCC"
	//UserIdentifierTypeNationalID represents the national id user identifier type
	UserIdentifierTypeNationalID UserIdentifierType = "NATIONAL_ID"
)

// IsValid returns true if a user identifier type is valid
func (f UserIdentifierType) IsValid() bool {
	switch f {
	case UserIdentifierTypeCCC, UserIdentifierTypeNationalID:
		return true
	}
	return false
}

func (f UserIdentifierType) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a user identifier type.
func (f *UserIdentifierType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = UserIdentifierType(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid UserIdentifierType", str)
	}
	return nil
}

// MarshalGQL writes the user identifier type to the supplied
func (f UserIdentifierType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}
