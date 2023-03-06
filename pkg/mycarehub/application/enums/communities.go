package enums

import (
	"fmt"
	"io"
	"strconv"
)

// Visibility indicates that the room will be shown in the published room list
type Visibility string

const (
	// PrivateVisibility hides the room from the published room list
	PrivateVisibility Visibility = "private"

	// PublicVisibility shows the room in the published room list
	PublicVisibility Visibility = "public"
)

// IsValid returns true if the Visibility is a valid value
func (e Visibility) IsValid() bool {
	switch e {
	case PrivateVisibility, PublicVisibility:
		return true
	}
	return false
}

// String is the string representation of the Visibility
func (e Visibility) String() string {
	return string(e)
}

// UnmarshalGQL parses the input string as a Visibility
func (e *Visibility) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Visibility(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Visibility", str)
	}
	return nil
}

// MarshalGQL serializes the Visibility
func (e Visibility) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// Preset is a convenience parameter for setting various default state events based on a preset.
type Preset string

const (
	// PresetPrivateChat gives additional information about the history visibility and guest access
	PresetPrivateChat Preset = "private_chat"

	// PresetPublicChat gives additional information about the history visibility and guest access
	PresetPublicChat Preset = "public_chat"

	// PresetTrustedPrivateChat gives all invitees the same power level as the room creator.
	PresetTrustedPrivateChat Preset = "trusted_private_chat"
)

// IsValid returns true if the Preset is a valid value
func (p Preset) IsValid() bool {
	switch p {
	case PresetPrivateChat, PresetPublicChat, PresetTrustedPrivateChat:
		return true
	}
	return false
}

// String is the string representation of the Preset
func (p Preset) String() string {
	return string(p)
}

// String is the string representation of the Preset
func (p Preset) Preset() string {
	return string(p)
}

// UnmarshalGQL parses the input string as a Preset
func (p *Preset) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*p = Preset(str)
	if !p.IsValid() {
		return fmt.Errorf("%s is not a valid Preset", str)
	}
	return nil
}

// MarshalGQL serializes the Preset
func (p Preset) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(p.String()))
}
