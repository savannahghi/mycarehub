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

	// ClientTypeOtz represents an Operation Triple Zero client type
	ClientTypeOtz ClientType = "OTZ"

	// ClientTypeOtzPlus represents all women aged between 10-24 years who are pregnant or breastfeeding.
	// Exit when baby turns 2 years
	ClientTypeOtzPlus ClientType = "OTZ_PLUS"

	// ClientTypeHvl represents HIV-positive clients with a viral load >1,000 copies per ml
	ClientTypeHvl ClientType = "HVL"

	// ClientTypeDreams represents Determined, Resilient, Empowered, AIDS-free, Mentored and Safe  client types
	ClientTypeDreams ClientType = "DREAMS"

	// ClientTypeHighRisk represents all pediatric patients 0-4 yrs, all 0-4-year-olds with an HIV negative guardian, all clients with low viremia (50-999 copies/ml).
	ClientTypeHighRisk ClientType = "HIGH_RISK"

	// ClientTypeSpouses represents the spouses of pmtct mothers who have disclosed
	ClientTypeSpouses ClientType = "SPOUSES"

	// ClientTypeYouth represents 20-24 year olds, both male and female
	ClientTypeYouth ClientType = "YOUTH"
)

// AllClientType represents a slice of all possible `ClientType` values
var AllClientType = []ClientType{
	ClientTypePmtct,
	ClientTypeOvc,
	ClientTypeOtz,
	ClientTypeOtzPlus,
	ClientTypeHvl,
	ClientTypeDreams,
	ClientTypeHighRisk,
	ClientTypeSpouses,
	ClientTypeYouth,
}

// IsValid returns true if a client type is valid
func (e ClientType) IsValid() bool {
	switch e {
	case
		ClientTypePmtct,
		ClientTypeOvc,
		ClientTypeOtz,
		ClientTypeOtzPlus,
		ClientTypeHvl,
		ClientTypeDreams,
		ClientTypeHighRisk,
		ClientTypeSpouses,
		ClientTypeYouth:
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
