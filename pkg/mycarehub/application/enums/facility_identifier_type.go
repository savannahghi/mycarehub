package enums

import (
	"fmt"
	"io"
	"strconv"
)

// FacilityIdentifierType is a list of all the facility identifier types.
type FacilityIdentifierType string

const (
	//FacilityIdentifierTypeMFLCode represents the mfl facility identifier type
	FacilityIdentifierTypeMFLCode FacilityIdentifierType = "MFL_CODE"

	//FacilityIdentifierTypeHealthCRM represents the health crm facility identifier type
	FacilityIdentifierTypeHealthCRM FacilityIdentifierType = "HEALTH_CRM"

	// FacilityIdentifierTypeSladeCode represents the health crm slade code identifier type
	FacilityIdentifierTypeSladeCode FacilityIdentifierType = "SLADE_CODE"
)

// IsValid returns true if a facility identifier type is valid
func (f FacilityIdentifierType) IsValid() bool {
	switch f {
	case FacilityIdentifierTypeMFLCode, FacilityIdentifierTypeHealthCRM, FacilityIdentifierTypeSladeCode:
		return true
	}
	return false
}

func (f FacilityIdentifierType) String() string {
	return string(f)
}

// UnmarshalGQL converts the supplied value to a facility identifier type.
func (f *FacilityIdentifierType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*f = FacilityIdentifierType(str)
	if !f.IsValid() {
		return fmt.Errorf("%s is not a valid FacilityIdentifierType", str)
	}
	return nil
}

// MarshalGQL writes the facility identifier type to the supplied
func (f FacilityIdentifierType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}
