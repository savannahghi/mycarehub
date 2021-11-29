package enums

import (
	"fmt"
	"io"
	"strconv"
)

// Mood is a custom type that defines moods that a client can go through
type Mood string

const (
	// MoodVerySad represents a state where a client is in a terrible state of mind
	MoodVerySad Mood = "VERY_SAD"

	// MoodSad represents a state where a client is in a bad state of mind
	MoodSad Mood = "SAD"

	// MoodHappy represents a contented mood
	MoodHappy Mood = "HAPPY"

	// MoodVeryHappy represents a delighted and joyous mood
	MoodVeryHappy Mood = "VERY_HAPPY"

	// MoodNeutral is a state of mind where a user is impartial
	MoodNeutral Mood = "NEUTRAL"
)

// AllMoods represents a slice of all available moods
var AllMoods = []Mood{
	MoodHappy, MoodNeutral, MoodSad, MoodVeryHappy, MoodVerySad,
}

// IsValid returns true if a mood type is valid
func (m Mood) IsValid() bool {
	switch m {
	case MoodHappy, MoodNeutral, MoodSad, MoodVeryHappy, MoodVerySad:
		return true
	}
	return false
}

// String converts the mood type enum to a string
func (m Mood) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a mood type.
func (m *Mood) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = Mood(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid mood type", str)
	}
	return nil
}

// MarshalGQL writes the mood type to the supplied writer
func (m Mood) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
