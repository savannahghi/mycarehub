package enums

import (
	"fmt"
	"io"
	"strconv"
)

// ClientType defines various client types
type ClientType string

const (
	// ClientTypePmtct represents a Prevention of mother-to-child transmission client type
	ClientTypePmtct ClientType = "PMTCT"

	// ClientTypeOvc represents an Orphan and Vulnerable Children client type
	ClientTypeOvc ClientType = "OVC"
)

// AllClientType represents a slice of all possible `ClientType` values
var AllClientType = []ClientType{
	ClientTypePmtct,
	ClientTypeOvc,
}

// IsValid returns true if a client type is valid
func (e ClientType) IsValid() bool {
	switch e {
	case ClientTypePmtct, ClientTypeOvc:
		return true
	}
	return false
}

// String ...
func (e ClientType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a metric type.
func (e *ClientType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ClientType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ClientType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e ClientType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
