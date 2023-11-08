package enums

import (
	"fmt"
	"io"
	"strconv"
)

// BookingCodeStatus is the allowed states of a booking's verification code.
type BookingCodeStatus string

const (
	Verified   BookingCodeStatus = "VERIFIED"
	UnVerified BookingCodeStatus = "UNVERIFIED"
)

// IsValid returns true if a booking code status is valid
func (m BookingCodeStatus) IsValid() bool {
	switch m {
	case Verified, UnVerified:
		return true
	}
	return false
}

func (m BookingCodeStatus) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a booking type.
func (m *BookingCodeStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = BookingCodeStatus(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid BookingCodeStatus", str)
	}
	return nil
}

// MarshalGQL writes the booking code status type to the supplied
func (m BookingCodeStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}

// BookingStatus is the allowed states of a booking.
type BookingStatus string

const (
	Pending   BookingStatus = "PENDING"
	Fulfilled BookingStatus = "FULFILLED"
)

// IsValid returns true if a booking status is valid
func (m BookingStatus) IsValid() bool {
	switch m {
	case Pending, Fulfilled:
		return true
	}
	return false
}

func (m BookingStatus) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a booking status.
func (m *BookingStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = BookingStatus(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid BookingStatus", str)
	}
	return nil
}

// MarshalGQL writes the booking status type to the supplied
func (m BookingStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
