package enums

import (
	"fmt"
	"io"
	"strconv"
)

// CountryType defines various countries available
type CountryType string

const (
	// CountryTypeKenya defines country type KENYA
	CountryTypeKenya CountryType = "KENYA"
	// Other countries
)

// AllCountries represents a slice of all countries available
var AllCountries = []CountryType{
	CountryTypeKenya,
}

// IsValid returns true if a country type is valid
func (e CountryType) IsValid() bool {
	switch e {
	case CountryTypeKenya:
		return true
	}
	return false
}

// String converts country type to string.
func (e CountryType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a country type.
func (e *CountryType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CountryType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CountryType", str)
	}
	return nil
}

// MarshalGQL writes the country type to the supplied writer
func (e CountryType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
