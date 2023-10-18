package enums

import (
	"fmt"
	"io"
	"strconv"
)

// Country is is a list of allowed countries
type Country string

const (
	//CountryKenya represents Kenya
	CountryKenya Country = "KE"
)

// IsValid returns true if a country is valid
func (f Country) IsValid() bool {
	switch f {
	case CountryKenya:
		return true
	}
	return false
}

func (f Country) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a country
func (f *Country) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = Country(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid Country", str)
	}
	return nil
}

// MarshalGQL writes the country to the supplied
func (f Country) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}
