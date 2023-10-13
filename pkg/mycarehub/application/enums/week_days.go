package enums

import (
	"fmt"
	"io"
	"strconv"
)

// DayOfWeek is a list of all the 7 days in a week.
type DayOfWeek string

// DayOfWeek constants
const (
	DayOfWeekMonday    DayOfWeek = "MONDAY"
	DayOfWeekTuesday   DayOfWeek = "TUESDAY"
	DayOfWeekWednesday DayOfWeek = "WEDNESDAY"
	DayOfWeekThursday  DayOfWeek = "THURSDAY"
	DayOfWeekFriday    DayOfWeek = "FRIDAY"
	DayOfWeekSaturday  DayOfWeek = "SATURDAY"
	DayOfWeekSunday    DayOfWeek = "SUNDAY"
)

// IsValid returns true if a user is valid
func (m DayOfWeek) IsValid() bool {
	switch m {
	case DayOfWeekMonday, DayOfWeekTuesday, DayOfWeekWednesday, DayOfWeekThursday, DayOfWeekFriday, DayOfWeekSaturday, DayOfWeekSunday:
		return true
	}
	return false
}

func (m DayOfWeek) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a user type.
func (m *DayOfWeek) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = DayOfWeek(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid day of the week", str)
	}
	return nil
}

// MarshalGQL writes the user type to the supplied
func (m DayOfWeek) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
