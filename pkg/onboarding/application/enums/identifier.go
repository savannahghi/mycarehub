package enums

import (
	"fmt"
	"io"
	"strconv"
)

// IdentifierType defines the various identifier types
type IdentifierType string

const (
	// IdentifierTypeCCC represents a Comprehensive Care Centre identifier type
	IdentifierTypeCCC IdentifierType = "CCC"

	// IdentifierTypeID represents the national ID identifier type
	IdentifierTypeID IdentifierType = "NATIONAL ID"

	// IdentifierTypePassport represents a passport identifier type
	IdentifierTypePassport IdentifierType = "PASSPORT"
)

// AllIdentifierType represents a slice of all possible `IdentifierTypes` values
var AllIdentifierType = []IdentifierType{
	IdentifierTypeCCC,
	IdentifierTypeID,
	IdentifierTypePassport,
}

// IsValid returns true if an identifier type is valid
func (e IdentifierType) IsValid() bool {
	switch e {
	case IdentifierTypeCCC, IdentifierTypeID, IdentifierTypePassport:
		return true
	}
	return false
}

// String ...
func (e IdentifierType) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a metric type.
func (e *IdentifierType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentifierType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentifierType", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e IdentifierType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// IdentifierUse defines the different kinds of identifiers use
type IdentifierUse string

const (
	// IdentifierUseOfficial represents an `official` identifier use
	IdentifierUseOfficial IdentifierUse = "OFFICIAL"

	// IdentifierUseTemporary represents a `temporary` identifier use
	IdentifierUseTemporary IdentifierUse = "TEMPORARY"

	// IdentifierUseOld represents an `old` identifier use
	IdentifierUseOld IdentifierUse = "OLD"
)

// AllIdentifierUse represents a slice of all possible `IdentifierUse` values
var AllIdentifierUse = []IdentifierUse{
	IdentifierUseOfficial,
	IdentifierUseTemporary,
	IdentifierUseOld,
}

// IsValid returns true if an identifier use is valid
func (e IdentifierUse) IsValid() bool {
	switch e {
	case IdentifierUseOfficial, IdentifierUseTemporary, IdentifierUseOld:
		return true
	}
	return false
}

// String ...
func (e IdentifierUse) String() string {
	return string(e)
}

// UnmarshalGQL converts the supplied value to a metric type.
func (e *IdentifierUse) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IdentifierUse(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IdentifierUse", str)
	}
	return nil
}

// MarshalGQL writes the metric type to the supplied writer
func (e IdentifierUse) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
