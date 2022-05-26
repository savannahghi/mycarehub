package enums

import (
	"fmt"
	"io"
	"strconv"
)

// CaregiverType is a list of all the sort types.
type CaregiverType string

const (
	// CaregiverTypeFather denotes the caregiver type father
	CaregiverTypeFather CaregiverType = "FATHER"

	// CaregiverTypeMother defines the value for the mother caregiver type
	CaregiverTypeMother CaregiverType = "MOTHER"

	// CaregiverTypeSibling defines the value for the sibling caregiver type
	CaregiverTypeSibling CaregiverType = "SIBLING"

	// CaregiverTypeHealthCareProfessional defines the value for the health care provider caregiver type
	CaregiverTypeHealthCareProfessional CaregiverType = "HEALTHCARE_PROFESSIONAL"
)

// IsValid returns true if a sort is valid
func (m CaregiverType) IsValid() bool {
	switch m {
	case CaregiverTypeFather,
		CaregiverTypeMother,
		CaregiverTypeSibling,
		CaregiverTypeHealthCareProfessional:
		return true
	}
	return false
}

func (m CaregiverType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *CaregiverType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = CaregiverType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid CaregiverType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m CaregiverType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
